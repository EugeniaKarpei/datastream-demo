package config

// Configuration metadate about CSV-dataset which is used in this demo.

const CsvDataSetFilePath = "./data/dataset.csv"

// Metadata constants
var (
	MetricName                 = "online.spent"
	MetricIdColumnIndex        = 0  // Avg_Price column
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

type CsvRecordMetaData struct {
	ColumnIndex int // 1-based
	Name        string
}
