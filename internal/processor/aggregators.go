package processor

import (
	"time"
	"valery-datadog-datastream-demo/internal/data"
)

type Aggregator func(time.Time, []*data.MetricRecord) *data.TimeDataPoint

func CountAggregator(timestamp time.Time, metrics []*data.MetricRecord) *data.TimeDataPoint {
	return data.NewTimeDataPoint(timestamp, float64(len(metrics)))
}

func SumAggregator(timestamp time.Time, metrics []*data.MetricRecord) *data.TimeDataPoint {
	var sum = 0.0
	for _, m := range metrics {
		sum += m.MetricValue()
	}
	return data.NewTimeDataPoint(timestamp, roundTo2DecimalPoints(sum))
}

func AvgAggregator(timestamp time.Time, metrics []*data.MetricRecord) *data.TimeDataPoint {
	var sum = 0.0
	for _, m := range metrics {
		sum += m.MetricValue()
	}
	return data.NewTimeDataPoint(timestamp, roundTo2DecimalPoints(sum/float64(len(metrics))))
}

func roundTo2DecimalPoints(value float64) float64 {
	return float64(int(value*100)) / 100
}
