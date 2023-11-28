package processor

// Main entry point into metric data processing.
//
// MetricProcessor has 2 roles (implement 2 interfaces):
//
// 1. StreamProcessor - accepts data from input DataStream and partitions it using tags internally
// We could potentially have things like
// * flexible time intervals and
// * metric names
// here as parameters and part of the system, but for this demo we took it out of scope because filtering
// by time is more complex.
// At this point we keep time interval constant and stick to single metric.
//
// 2. MetricDataProvider - accepts API calls and provides data points for API clients.
//
// The design of MetricProcessor is based on nested map of
//   map[tagName] -> map[tagValue] -> Metrics - which is a collection of metrics with another map inside
// It allows us to have O(1) time for retrieval of metric records for 0 or 1 filters scenarios.
// If number of filters > 1, we are merging metric data sets by iterating on the smallest of filtered data-sets. The complexity of this step is O(min(ni)) where ni - number of metric records within i-th filter partition.
// It is possible to achieve O(n) performance for this step but it would cost substantial memory profile increase as we would need to pre-compute data sets for combined tags i.e. tag1:value1;tag2:value2;etc.
//
// After metrics retrieved we apply partition by time (using one of our static time partitioners) and aggregation.
//
// For filter search (/getFilters) we are using Trie data structure to be able to quickly retrieve all availble tag:value pairs. The complexity of this step is O(sn + tn) where sn - length of search term and tn - combined length of all tag:value strings that exist in our dataset.

import (
	"sort"
	"valery-datadog-datastream-demo/internal/data"
)

// Provide data to external users (for ex. - API handlers)
type MetricDataProvider interface {
	GetMetricDataPoints(filters []*data.Tag, timePartition TimePartitioner, aggregator Aggregator) []data.TimeDataPoint
	GetMetricTagFilters(searchTerm string) []string
}

var _ MetricDataProvider = (*InMemoryMetricStreamProcessor)(nil)

var _ data.StreamProcessor = (*InMemoryMetricStreamProcessor)(nil)

func NewInMemoryMetricStreamProcessor() *InMemoryMetricStreamProcessor {
	return &InMemoryMetricStreamProcessor{
		allMetrics:    data.NewMetrics(),
		taggedMetrics: make(map[string]map[string]*data.Metrics),
		tagFilters:    NewTrieNode(),
	}
}

// Stateful metrics processor.
// * Works as a singleton
// * Uses metadata to build indexes on data stream
// * Uses aggregators to aggregate incoming metrics into displayable data points
type InMemoryMetricStreamProcessor struct {

	// all metrics time series for general layout
	allMetrics *data.Metrics

	// nested map for tagged metrics: tagName -> tagValue -> TaggedMetrics
	taggedMetrics map[string]map[string]*data.Metrics

	tagFilters *TrieNode
}

// Process incoming data stream, build indices based on tags
func (mp *InMemoryMetricStreamProcessor) Process(dataRecord []string) error {
	metricRecord, tags, err := data.FromCsvDataRecord(dataRecord)
	if err != nil {
		return err
	}

	// based on tag names and values specified for the data record - populate nested metric data-storage
	for tagName, tag := range tags {
		tagValueMap, found := mp.taggedMetrics[tagName]
		if !found {
			tagValueMap = make(map[string]*data.Metrics)
		}

		taggedMetrics, found := tagValueMap[tag.Value()]
		if !found {
			taggedMetrics = data.NewMetrics()
		}
		taggedMetrics.AddRecord(metricRecord)
		tagValueMap[tag.Value()] = taggedMetrics
		mp.taggedMetrics[tagName] = tagValueMap

		// and also update metric tag-filters storage
		filterStr := tag.AsFilter()
		mp.tagFilters.AddWord(filterStr)
	}

	// add metric to the total collection
	mp.allMetrics.AddRecord(metricRecord)

	return nil // no errors, we are done
}

// Returns key-value pairs of tagName:tagValue - available for filtering in the current data-set
func (mp *InMemoryMetricStreamProcessor) GetMetricTagFilters(searchTerm string) []string {
	filters := mp.tagFilters.GetWordsInSubtrie(searchTerm)
	// todo: we might also add remaining tag:value pairs here
	return filters
}

// Fetch data from the internal data structures, use indices to filter and aggregator to aggregate and prepare data points
// we only implement filtering for now.
func (mp *InMemoryMetricStreamProcessor) GetMetricDataPoints(filters []*data.Tag, timePartition TimePartitioner, aggregate Aggregator) []data.TimeDataPoint {
	// 1. We need to choose data to partition or aggregate
	metrics := mp.getInputMetrics(filters)

	// 2. Partition by time
	partitionedMetrics := timePartition(metrics.MetricRecords())

	// 3. aggregate using aggregator function
	dataPoints := make([]data.TimeDataPoint, len(partitionedMetrics))
	i := 0
	for pKey, partition := range partitionedMetrics {
		dataPoints[i] = aggregate(pKey, partition)
		i++
	}

	// 4. sort result data points by timestamp
	sort.Slice(dataPoints, func(i, j int) bool {
		return dataPoints[i].Timestamp < dataPoints[j].Timestamp
	})
	return dataPoints
}

func (mp *InMemoryMetricStreamProcessor) getInputMetrics(filters []*data.Tag) *data.Metrics {
	// by default we take an empty set of metrics
	metrics := data.NewMetrics()

	if len(filters) == 0 {
		// if no filters specified - we use all metrics for aggregation
		metrics = mp.allMetrics
	} else if len(filters) == 1 {
		// if only one filter is specified - we can simply use pre-partitioned time series
		filterTag := filters[0]
		tagValueMap, found := mp.taggedMetrics[filterTag.Name()]
		if found {
			metrics, _ = tagValueMap[filterTag.Value()]
		}
	} else {
		// there are multiple filters specified which means we have to merge filtering results
		// gather time series to be merged
		filterTagMetrics := make([]*data.Metrics, len(filters))
		for i, filterTag := range filters {
			// if at least one of tags (name->value->...) does not have any data - filtering is not necessary
			// as we'll get empty result in the end anyway
			tagValueMap, found := mp.taggedMetrics[filterTag.Name()]
			if !found {
				return metrics
			}
			tagMetrics, found := tagValueMap[filterTag.Value()]
			if !found {
				return metrics
			}
			filterTagMetrics[i] = tagMetrics
		}
		// we pick the smallest set of metrics
		minLenMetrics := pickWithMinLength(filterTagMetrics)

		// and we merge
		merged := data.NewMetrics()
		for _, m := range minLenMetrics.MetricRecords() {
			addToMerged := true
			for _, fm := range filterTagMetrics {
				if !fm.IsRecordPresent(m.Id()) {
					addToMerged = false
					break
				}
			}
			if addToMerged {
				merged.AddRecord(m)
			}
		}
		metrics = merged
	}
	return metrics
}

func pickWithMinLength(metrics []*data.Metrics) *data.Metrics {
	minLen := metrics[0]
	for _, m := range metrics {
		if m.Len() < minLen.Len() {
			minLen = m
		}
	}
	return minLen
}
