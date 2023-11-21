package processor

import (
	"valery-datadog-datastream-demo/internal/data"
)

type MetricDataProvider interface {
	// We could potentially have time interval here as a parameter, but filtering
	// by time is more complex, so for now we keep time interval constant
	GetMetricDataPoints(filters []*data.Tag, timePartition TimePartitioner, aggregator Aggregator) []data.TimeDataPoint
	// GetMetricNames() []string
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
	filters := mp.tagFilters.GetWordsInOrder(searchTerm)
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
