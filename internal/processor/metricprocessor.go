package processor

import (
	"valery-datadog-datastream-demo/internal/data/model"
)

type MetricProcessor interface {
	ProcessMetric(metricRec *model.MetricRecord)
}

type MetricDataProvider interface {
	// We could potentially have time interval here as a parameter, but filtering
	// by time is more complex, so for now we keep time interval constant
	GetMetricDataPoints(filters []*model.Tag, partitions []*model.Tag, aggregator Aggregator) []*DataPoint
}

var _ MetricProcessor = (*InMemoryMetricProcessor)(nil)

// Stateful metrics processor.
// * Works as a singleton
// * Uses metadata to build indexes on data stream
// * Uses aggregators to aggregate incoming metrics into displayable data points
type InMemoryMetricProcessor struct {
	allMetrics map[int]MetricValue
}

func (mp *InMemoryMetricProcessor) ProcessMetric(metricRec *model.MetricRecord) {

}
