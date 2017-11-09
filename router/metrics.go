package router

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(routerBacklog, routerDroppedEvents, routerProcessingDuration)
}

var routerDroppedEvents = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "router",
		Name:      "dropped_events_total",
		Help:      "Dropped events due to insufficient processing power.",
	})

var routerBacklog = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Namespace: "gateway",
		Subsystem: "router",
		Name:      "backlog_events",
		Help:      "Gauge of events count waiting to be processed by the router.",
	})

var routerProcessingDuration = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Namespace: "gateway",
		Subsystem: "router",
		Name:      "event_processing_seconds",
		Help: "Bucketed histogram of processing duration of an event in the router. " +
			"From receiving the event to calling a function.",
		Buckets: prometheus.ExponentialBuckets(0.00001, 2, 20),
	})

var receivedEventsMutex = sync.Mutex{}
var receivedEvents = map[string]time.Time{}

func reportReceivedEvent(id string) {
	receivedEventsMutex.Lock()
	defer receivedEventsMutex.Unlock()
	receivedEvents[id] = time.Now()
}

func reportProceededEvent(id string) {
	receivedEventsMutex.Lock()
	defer receivedEventsMutex.Unlock()
	if startTime, ok := receivedEvents[id]; ok {
		routerProcessingDuration.Observe(time.Since(startTime).Seconds())
		delete(receivedEvents, id)
	}
}
