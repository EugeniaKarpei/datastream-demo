package data

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"valery-datadog-datastream-demo/internal/config"
)

// Generate metric from CSV data record
func FromCsvDataRecord(csvDataRecord []string) (*MetricRecord, error) {
	id, err := parseInt(csvDataRecord[config.MetricIdColumnIndex])
	if err != nil {
		return nil, err
	}

	timestamp, err := parseDate(csvDataRecord[config.MetricTimestampColumnIndex])
	if err != nil {
		return nil, err
	}

	metricValue, err := parseFloat64(csvDataRecord[config.MetricValueColumnIndex])
	if err != nil {
		return nil, err
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
		tags:      tags,
	}, nil
}

// Main metric data structures

type MetricRecord struct {
	id        int
	timestamp time.Time
	name      string
	value     float64
	tags      map[string]*Tag
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

func (metric *MetricRecord) Tags() []string {
	tagNames := make([]string, len(metric.tags))
	i := 0
	for tagName := range metric.tags {
		tagNames[i] = tagName
		i++
	}
	return tagNames
}

func (metric *MetricRecord) HasTagWithAnyValue(tagName string) bool {
	_, hasTag := metric.tags[tagName]
	return hasTag
}

func (metric *MetricRecord) HasTagWithValue(tagName string, tagValue string) bool {
	tag, hasTag := metric.tags[tagName]
	if !hasTag {
		return false
	}
	return tag.Value() == tagValue
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

type MetricValue struct {
	id        int
	timestamp time.Time
}

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
