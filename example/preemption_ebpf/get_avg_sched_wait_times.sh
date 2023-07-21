#!/bin/bash

# get avg sched wait times per thread.
while read i; do sudo perf sched timehist -t $i | tail -n +4 | awk '{ sum += $4; n++ } END { if (n > 0) print sum / n; }' | tail -n +1; done < tids.txt > wait_times.txt

# calculate avg wait time across threads.
cat wait_times.txt | awk '{ sum += $1; n++ } END { if (n > 0) print sum / n; }'
