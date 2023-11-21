package data

// todo: serialize into json
type TimeDataPoint struct {
	timestamp int
	points    []*DataPoint
}

type DataPoint struct {
	name  string
	value float64
}
