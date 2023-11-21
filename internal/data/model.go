package data

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"valery-datadog-datastream-demo/internal/config"
)

func FromRequestFilters(requestFilters []string) []*Tag {
	tags := make([]*Tag, len(requestFilters))
	for i, filter := range requestFilters {
		tags[i] = NewTagForFiltering(filter)
	}
	return tags
}

// Generate metric from CSV data record
func FromCsvDataRecord(csvDataRecord []string) (*MetricRecord, Tags, error) {
	id, err := parseInt(csvDataRecord[config.MetricIdColumnIndex])
	if err != nil {
		return nil, nil, err
	}

	timestamp, err := parseDate(csvDataRecord[config.MetricTimestampColumnIndex])
	if err != nil {
		return nil, nil, err
	}

	metricValue, err := parseFloat64(csvDataRecord[config.MetricValueColumnIndex])
	if err != nil {
		return nil, nil, err
	}

	tags := make(map[string]*Tag)
	for _, tagMetaData := range config.MetricTagsMetaData {
		tagValue := csvDataRecord[tagMetaData.ColumnIndex]
		// if no value in the tag column - don't apply the tag
		if len(tagValue) == 0 {
			continue
		}
		tag := &Tag{
			name:  tagMetaData.Name,
			value: tagValue,
		}
		tags[tag.name] = tag
	}

	return &MetricRecord{
		id:        id,
		timestamp: timestamp,
		name:      config.MetricName,
		value:     metricValue,
	}, tags, nil
}

// Main metric data structures

type MetricRecord struct {
	id        int
	timestamp time.Time
	name      string
	value     float64
}

func (metric *MetricRecord) Id() int {
	return metric.id
}

func (metric *MetricRecord) Timestamp() time.Time {
	return metric.timestamp
}

func (metric *MetricRecord) MetricName() string {
	return metric.name
}

func (metric *MetricRecord) MetricValue() float64 {
	return metric.value
}

// Collection of metrics with additional matching facilities

func NewMetrics() *Metrics {
	return &Metrics{
		presenceData: make(map[int]interface{}),
		records:      []*MetricRecord{},
	}
}

type Metrics struct {
	presenceData map[int]interface{}
	records      []*MetricRecord
}

func (metrics *Metrics) AddRecord(metricRecord *MetricRecord) {
	metrics.presenceData[metricRecord.id] = nil
	metrics.records = append(metrics.records, metricRecord)
}

func (metrics *Metrics) IsRecordPresent(metricRecordId int) bool {
	_, present := metrics.presenceData[metricRecordId]
	return present
}

func (metrics *Metrics) MetricRecords() []*MetricRecord {
	return metrics.records
}

func (metrics *Metrics) Len() int {
	return len(metrics.presenceData)
}

// Tags

func NewTagForFiltering(keyValueStr string) *Tag {
	keyValuePair := strings.Split(keyValueStr, ":")
	return &Tag{
		name:  keyValuePair[0],
		value: keyValuePair[1],
	}
}

type Tag struct {
	name  string
	value string
}

func (t *Tag) Name() string {
	return t.name
}

func (t *Tag) Value() string {
	return t.value
}

func (t *Tag) AsFilter() string {
	return t.name + ":" + t.value
}

type Tags map[string]*Tag

// helper functions

// returns error if strField is empty
func parseFloat64(strField string) (float64, error) {
	if len(strField) == 0 {
		return float64(0.0), errors.New("Metric value CSV field is empty")
	}
	return strconv.ParseFloat(strField, 2)
}

// if the value is with floating point, rounds to int
func parseInt(strField string) (int, error) {
	if len(strField) == 0 {
		return int(0.0), errors.New("Metric value CSV field is empty")
	}
	dotIndex := strings.LastIndex(strField, ".")
	if dotIndex > 0 {
		strField = strField[:dotIndex]
	}
	return strconv.Atoi(strField)
}

// returns error if there is no data to parse
func parseDate(strField string) (time.Time, error) {
	if len(strField) == 0 {
		return time.Now(), errors.New("Metric timestamp field is empty")
	}
	return time.Parse("2006-01-02", strField)
}
