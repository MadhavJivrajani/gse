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

package exporter

import (
	"fmt"
	"runtime"

	"github.com/MadhavJivrajani/gse/pkg/sched"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// SchedulreMetrics holds prometheus servable metrics.
type SchedulerMetrics struct {
	Idleprocs            prometheus.Gauge
	Threads              prometheus.Gauge
	SpinningThreads      prometheus.Gauge
	IdleThreads          prometheus.Gauge
	GlobalRunQueueLength prometheus.Gauge
	NeedsSpinningThreads prometheus.Gauge
	// the length of this slice will
	// be equal to the value of GOMAXPORCS.
	LocalRunQueueLengths []prometheus.Gauge
}

func NewSchedulerMetrics() *SchedulerMetrics {
	return &SchedulerMetrics{
		Idleprocs: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "idle_procs",
				Help: "Idle Procs",
			},
		),
		Threads: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "threads",
				Help: "Threads",
			},
		),
		SpinningThreads: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "spinning_threads",
				Help: "Spinning Threads",
			},
		),
		IdleThreads: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "idle_threads",
				Help: "Idle Threads",
			},
		),
		NeedsSpinningThreads: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "needs_spinning_threads",
				Help: "Needs Spinning Threads",
			},
		),
		GlobalRunQueueLength: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "global_runqueue_length",
				Help: "Global Run Queue Length",
			},
		),
		LocalRunQueueLengths: func() []prometheus.Gauge {
			res := make([]prometheus.Gauge, runtime.GOMAXPROCS(-1))
			for i := 0; i < len(res); i++ {
				res[i] = promauto.NewGauge(
					prometheus.GaugeOpts{
						Name: fmt.Sprintf("local_run_queue_%d", i),
						Help: fmt.Sprintf("Local Run Queue %d", i),
					},
				)
			}
			return res
		}(),
	}
}

// UpdateMetricsFromTrace updates the metric gauges based on the scheduler trace.
func (sm *SchedulerMetrics) UpdateMetricsFromTrace(st *sched.SchedTrace) {
	sm.Idleprocs.Set(float64(st.Idleprocs))
	sm.Threads.Set(float64(st.Threads))
	sm.SpinningThreads.Set(float64(st.SpinningThreads))
	sm.IdleThreads.Set(float64(st.IdleThreads))
	sm.NeedsSpinningThreads.Set(float64(st.NeedsSpinningThreads))
	sm.GlobalRunQueueLength.Set(float64(st.GlobalRunQueueLength))

	for i := 0; i < len(sm.LocalRunQueueLengths); i++ {
		sm.LocalRunQueueLengths[i].Set(float64(st.LocalRunQueueLengths[i]))
	}
}
