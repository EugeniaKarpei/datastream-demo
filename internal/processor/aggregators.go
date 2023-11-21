package processor

import (
	"time"
	"valery-datadog-datastream-demo/internal/data"
)

const (
	SUM_AGGREGATOR   = "Sum"
	COUNT_AGGREGATOR = "Count"
	AVG_AGGREGATOR   = "Avg"
)

func FromRequestAggregator(aggregator string) Aggregator {
	switch aggregator {
	case SUM_AGGREGATOR:
		return SumAggregator
	case AVG_AGGREGATOR:
		return AvgAggregator
	case COUNT_AGGREGATOR:
		return CountAggregator
	default:
		return CountAggregator
	}
}

type Aggregator func(time.Time, []*data.MetricRecord) data.TimeDataPoint

func CountAggregator(timestamp time.Time, metrics []*data.MetricRecord) data.TimeDataPoint {
	return data.NewTimeDataPoint(timestamp, float64(len(metrics)))
}

func SumAggregator(timestamp time.Time, metrics []*data.MetricRecord) data.TimeDataPoint {
	var sum = 0.0
	for _, m := range metrics {
		sum += m.MetricValue()
	}
	return data.NewTimeDataPoint(timestamp, roundTo2DecimalPoints(sum))
}

func AvgAggregator(timestamp time.Time, metrics []*data.MetricRecord) data.TimeDataPoint {
	var sum = 0.0
	for _, m := range metrics {
		sum += m.MetricValue()
	}
	return data.NewTimeDataPoint(timestamp, roundTo2DecimalPoints(sum/float64(len(metrics))))
}

func roundTo2DecimalPoints(value float64) float64 {
	return float64(int(value*100)) / 100
}
