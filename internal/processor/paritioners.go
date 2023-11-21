package processor

import (
	"time"
	"valery-datadog-datastream-demo/internal/data"
)

type TimePartitioner func([]*data.MetricRecord) map[time.Time][]*data.MetricRecord

func MonthlyTimePartitioner(inputMetrics []*data.MetricRecord) map[time.Time][]*data.MetricRecord {
	return partitionByTime(inputMetrics, startOfTheMonth)
}

func DailyTimePartitioner(inputMetrics []*data.MetricRecord) map[time.Time][]*data.MetricRecord {
	return partitionByTime(inputMetrics, startOfTheDay)
}

func partitionByTime(inputMetrics []*data.MetricRecord, partitionKey timePartitionKey) map[time.Time][]*data.MetricRecord {
	partitioned := make(map[time.Time][]*data.MetricRecord)
	for _, m := range inputMetrics {
		pKey := partitionKey(m.Timestamp())

		partition, found := partitioned[pKey]
		if !found {
			partition = make([]*data.MetricRecord, 0)
		}
		partition = append(partition, m)
		partitioned[pKey] = partition
	}
	return partitioned
}

type timePartitionKey func(time.Time) time.Time

func startOfTheMonth(timestamp time.Time) time.Time {
	year, month, _ := timestamp.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func startOfTheDay(timestamp time.Time) time.Time {
	year, month, day := timestamp.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
