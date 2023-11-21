package processor

import "valery-datadog-datastream-demo/internal/data/model"

type Aggregator func([]*model.MetricRecord) *DataPoint

func CountAggregator(metrics []*model.MetricRecord) *DataPoint {

}
