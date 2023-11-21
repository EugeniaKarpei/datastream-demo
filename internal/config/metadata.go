package config

const CsvDataSetFilePath = "./data/dataset.csv"

var (
	MetricName                 = "online.spent"
	MetricIdColumnIndex        = 0  // Avg_Price
	MetricValueColumnIndex     = 11 // Avg_Price
	MetricTimestampColumnIndex = 6  // Transation_Date

	MetricTagsMetaData = map[string]*CsvRecordMetaData{
		"gender": {
			ColumnIndex: 2,
			Name:        "gender",
		},
		"location": {
			ColumnIndex: 3,
			Name:        "location",
		},
		"product_category": {
			ColumnIndex: 9,
			Name:        "product_category",
		},
		"coupon_status": {
			ColumnIndex: 13,
			Name:        "coupon_status",
		},
		"coupon_code": {
			ColumnIndex: 19,
			Name:        "coupon_code",
		},
	}
)

func TagNames() []string {
	tagNames := make([]string, len(MetricTagsMetaData))
	i := 0
	for tagName := range MetricTagsMetaData {
		tagNames[i] = tagName
		i++
	}
	return tagNames
}

type CsvRecordMetaData struct {
	ColumnIndex int // 1-based
	Name        string
}
