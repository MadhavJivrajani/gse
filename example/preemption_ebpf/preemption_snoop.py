#!/usr/bin/python

# Heavily inspired from: https://github.com/iovisor/bcc/blob/master/tools/killsnoop.py

from __future__ import print_function
from bcc import BPF
from bcc.utils import printb
from datetime import datetime
import sys

debug = 0

# define BPF program
bpf_text = """
#include <uapi/linux/ptrace.h>
#include <linux/sched.h>
struct val_t {
   u32 pid;
   int tid;
   int sig;
   int tgid;
   char comm[TASK_COMM_LEN];
};
struct data_t {
   u32 pid;
   int tgid;
   int tid;
   int sig;
   int ret;
   char comm[TASK_COMM_LEN];
};
BPF_HASH(infotmp, u32, struct val_t);
BPF_PERF_OUTPUT(events);
int syscall__tgkill(struct pt_regs *ctx, int tgid, int tid, int sig)
{
    u64 pid_tgid = bpf_get_current_pid_tgid();
    u32 pid = pid_tgid >> 32;
    u32 tid_key = (u32)pid_tgid;

    if (sig != 23) { return 0; }
    struct val_t val = {.pid = pid};
    if (bpf_get_current_comm(&val.comm, sizeof(val.comm)) == 0) {
        val.tgid = tgid;
        val.sig = sig;
        val.tid = tid;
        infotmp.update(&tid_key, &val);
    }
    return 0;
};
int do_ret_sys_tgkill(struct pt_regs *ctx)
{
    struct data_t data = {};
    struct val_t *valp;
    u64 pid_tgid = bpf_get_current_pid_tgid();
    u32 pid = pid_tgid >> 32;
    u32 tid_key = (u32)pid_tgid;
    valp = infotmp.lookup(&tid_key);
    if (valp == 0) {
        // missed entry
        return 0;
    }
    bpf_probe_read_kernel(&data.comm, sizeof(data.comm), valp->comm);
    data.pid = pid;
    data.tgid = valp->tgid;
    data.ret = PT_REGS_RC(ctx);
    data.sig = valp->sig;
    data.tid = valp->tid;
    events.perf_submit(ctx, &data, sizeof(data));
    infotmp.delete(&tid_key);
    return 0;
}
"""

# initialize BPF
b = BPF(text=bpf_text)
kill_fnname = b.get_syscall_fnname("tgkill")
b.attach_kprobe(event=kill_fnname, fn_name="syscall__tgkill")
b.attach_kretprobe(event=kill_fnname, fn_name="do_ret_sys_tgkill")

if len(sys.argv) > 1:
    # header
    print("%-13s %-6s %-16s %-4s %-6s" % (
        "TIME", "PID", "COMM", "SIG", "TID"))

# process event
def print_event(cpu, data, size):
    event = b["events"].event(data)
    if event.comm.decode('utf-8') != "lol":
        return
    if len(sys.argv) > 1:
        printb(b"%-13s %-6d %-16s %-4d %-6d" % (datetime.utcnow().strftime('%H:%M:%S.%f')[:-3].encode("ascii"),
            event.pid, event.comm, event.sig, event.tid))
    else:
        printb(b"%-6d" % event.tid)

# loop with callback to print_event
b["events"].open_perf_buffer(print_event)
while 1:
    try:
        b.perf_buffer_poll()
    except KeyboardInterrupt:
        exit()