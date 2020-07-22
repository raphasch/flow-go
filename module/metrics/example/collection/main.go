package main

import (
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/dapperlabs/flow-go/engine"
	"github.com/dapperlabs/flow-go/module/metrics"
	"github.com/dapperlabs/flow-go/module/metrics/example"
	"github.com/dapperlabs/flow-go/module/trace"
	"github.com/dapperlabs/flow-go/utils/unittest"
)

func main() {
	example.WithMetricsServer(func(logger zerolog.Logger) {
		tracer, err := trace.NewTracer(logger, "collection")
		if err != nil {
			panic(err)
		}
		collector := struct {
			*metrics.HotstuffCollector
			*metrics.CollectionCollector
			*metrics.NetworkCollector
		}{
			HotstuffCollector:   metrics.NewHotstuffCollector("some_chain_id", prometheus.DefaultRegisterer),
			CollectionCollector: metrics.NewCollectionCollector(tracer, prometheus.DefaultRegisterer),
			NetworkCollector:    metrics.NewNetworkCollector(prometheus.DefaultRegisterer),
		}

		topic1 := engine.ChannelName(engine.TestNetwork)
		topic2 := engine.ChannelName(engine.TestMetrics)

		for i := 0; i < 100; i++ {
			collector.TransactionIngested(unittest.IdentifierFixture())
			collector.HotStuffBusyDuration(10, metrics.HotstuffEventTypeTimeout)
			collector.HotStuffWaitDuration(10, metrics.HotstuffEventTypeTimeout)
			collector.HotStuffIdleDuration(10)
			collector.SetCurView(uint64(i))
			collector.SetQCView(uint64(i))

			collector.NetworkMessageSent(rand.Intn(1000), topic1)
			collector.NetworkMessageSent(rand.Intn(1000), topic2)

			collector.NetworkMessageReceived(rand.Intn(1000), topic1)
			collector.NetworkMessageReceived(rand.Intn(1000), topic2)

			time.Sleep(1 * time.Second)
		}
	})
}
