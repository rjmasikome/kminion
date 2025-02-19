package prometheus

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

func (e *Exporter) collectEndToEnd(ctx context.Context, ch chan<- prometheus.Metric) bool {

	// Producer is up, if produce duration is within the configured SLA
	// TODO: discuss the used semantic, perhaps add a dedicated SLA metric
	produceMs, exists := e.minionSvc.ProduceDurationMs(ctx)
	if exists && produceMs <= e.minionSvc.Cfg.EndToEnd.Producer.LatencySla.Milliseconds() {
		ch <- prometheus.MustNewConstMetric(e.endToEndProducerUp, prometheus.GaugeValue, 1.0)
	} else {
		ch <- prometheus.MustNewConstMetric(e.endToEndProducerUp, prometheus.GaugeValue, 0.0)
	}

	// Consumer is up, if consume duration is within the configured SLA
	// TODO: discuss the used semantic, perhaps add a dedicated SLA metric
	consumeMs, exists := e.minionSvc.ConsumeDurationMs(ctx)
	if exists && consumeMs <= e.minionSvc.Cfg.EndToEnd.Consumer.LatencySla.Milliseconds() {
		ch <- prometheus.MustNewConstMetric(e.endToEndConsumerUp, prometheus.GaugeValue, 1.0)
	} else {
		ch <- prometheus.MustNewConstMetric(e.endToEndConsumerUp, prometheus.GaugeValue, 0.0)
	}

	// Check if the consumer can commit offset (manually), if yes the value is 1
	commitOffsetsOk := e.minionSvc.OffsetCommitAvailability(ctx)
	if commitOffsetsOk {
		ch <- prometheus.MustNewConstMetric(e.endToEndOffsetCommitAvailability, prometheus.GaugeValue, 1.0)
	} else {
		ch <- prometheus.MustNewConstMetric(e.endToEndOffsetCommitAvailability, prometheus.GaugeValue, 0.0)
	}

	messageProduced := e.minionSvc.ProduceCounts(ctx)
	messageProducedAcks := e.minionSvc.ProduceAcks(ctx)
	if e.minionSvc.Cfg.EndToEnd.Enabled && e.minionSvc.Cfg.EndToEnd.Producer.RequiredAcks != 0 {
		if messageProduced == 0 {
			// Avoid division by zero
			ch <- prometheus.MustNewConstMetric(e.endToEndMessageLossRate, prometheus.GaugeValue, 0)
		} else {
			// lossRate is 1 - successRate
			lossRate := 1 - (float64(messageProducedAcks) / float64(messageProduced))
			ch <- prometheus.MustNewConstMetric(e.endToEndMessageLossRate, prometheus.GaugeValue, lossRate)
		}
	}

	return true
}
