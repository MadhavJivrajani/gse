#!/bin/bash

cd /sys/kernel/debug/tracing
echo 'sig==23' > events/signal/signal_generate/filter 
echo 1 > events/signal/signal_generate/enable
: > trace
echo 1 > tracing_on
