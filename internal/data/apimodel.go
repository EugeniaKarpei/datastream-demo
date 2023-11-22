package data

import "time"

// Request / reponse protocol entities

// /getData request
type GetDataRequest struct {
	Filters    []string `json:"filters"` // in the format of "tagName:tagValue"
	Scale      string   `json:"scale"`
	Aggregator string   `json:"aggregator"`
}

// /getFilters request
type GetFiltersRequest struct {
	Query string `json:"query"`
}

// Data points (response)

func NewTimeDataPoint(timestamp time.Time, value float64) TimeDataPoint {
	return TimeDataPoint{
		Timestamp: timestamp.UnixMilli(),
		Value:     value,
	}
}

type TimeDataPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}
