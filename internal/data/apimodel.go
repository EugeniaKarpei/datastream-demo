package data

import "time"

func NewTimeDataPoint(timestamp time.Time, value float64) *TimeDataPoint {
	return &TimeDataPoint{
		timestamp: timestamp.UnixMilli(),
		value:     value,
	}
}

// todo: serialize into json
type TimeDataPoint struct {
	timestamp int64
	value     float64
}

func (timeDataPoint *TimeDataPoint) Timestamp() int64 {
	return timeDataPoint.timestamp
}

func (timeDataPoint *TimeDataPoint) Value() float64 {
	return timeDataPoint.value
}
