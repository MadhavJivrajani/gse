/*
Copyright © 2021 Madhav Jivrajani madhav.jiv@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import (
	"runtime"
	"strconv"
	"strings"
)

// SchedTrace holds information about the extracted
// scheduler trace.
type SchedTrace struct {
	Gomaxprocs           int
	Idleprocs            int
	Threads              int
	SpinningThreads      int
	IdleThreads          int
	GlobalRunQueueLength int
	// the length of this slice will
	// be equal to the value of GOMAXPORCS.
	LocalRunQueueLengths []int
}

// NewSchedTrace is a constructor for SchedTrace
// and initializes LocalRunQueueLengths to a slice
// of length equal to GOMAXPROCS.
func NewSchedTrace() *SchedTrace {
	n := runtime.GOMAXPROCS(-1)
	return &SchedTrace{
		LocalRunQueueLengths: make([]int, n),
	}
}

// UpdateSchedTraceFromRawTrace updates the SchedTrace in-place from a string based
// scheduler trace and returns this updated SchedTrace.
func (st *SchedTrace) UpdateSchedTraceFromRawTrace(rawTrace string) *SchedTrace {
	split := strings.Split(rawTrace, " ")

	// this is a lot of hacky code to get the information out of a scheduler trace
	// line that looks as follows with GOMAXPROCS being 8:
	// SCHED 0ms: gomaxprocs=8 idleprocs=6 threads=5 spinningthreads=1 idlethreads=2 runqueue=0 [0 0 0 0 0 0 0 0]
	st.Gomaxprocs = getValueFromKVPair(split[2])
	st.Idleprocs = getValueFromKVPair(split[3])
	st.Threads = getValueFromKVPair(split[4])
	st.SpinningThreads = getValueFromKVPair(split[5])
	st.IdleThreads = getValueFromKVPair(split[6])
	st.GlobalRunQueueLength = getValueFromKVPair(split[7])

	st.LocalRunQueueLengths[0] = toInt(split[8][1:])
	for i := 1; i < runtime.GOMAXPROCS(-1); i++ {
		st.LocalRunQueueLengths[i] = toInt(split[8+i])
	}

	return st
}

// pair is of the form key=value, and value here is assumed to be an int.
func getValueFromKVPair(pair string) int {
	equalSignIndex := strings.Index(pair, "=")
	return toInt(pair[equalSignIndex+1:])
}

func toInt(s string) int {
	lenS := len(s)
	if s[lenS-1] == ']' {
		s = s[:lenS-1]
	}
	// we shall assume that error cannot occur in this parse because
	// of laziness.
	val, _ := strconv.ParseInt(s, 10, 0)
	return int(val)
}
