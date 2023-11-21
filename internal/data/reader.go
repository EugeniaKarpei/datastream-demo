package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type DataSourceReader interface {
	// Returns next data source record and true if it exists, false - if there is no more data
	GetNextDataItem() (*MetricRecord, bool)
}

func NewFileDataSourceReader(filePath string) DataSourceReader {
	fileDataSource := &FileDataSourceReader{
		filePath:      filePath,
		currentOffset: 0,
	}
	fileDataSource.readFileIntoMemory()
	return fileDataSource
}

var _ DataSourceReader = (*FileDataSourceReader)(nil)

type FileDataSourceReader struct {
	filePath      string
	dataInMemory  []*MetricRecord
	currentOffset int
}

func (fileDataSource *FileDataSourceReader) GetNextDataItem() (*MetricRecord, bool) {
	// We have reached the end of the data set
	if fileDataSource.currentOffset == len(fileDataSource.dataInMemory) {
		fmt.Println("Reached the end of the dataset")
		return nil, false
	}
	// return data record at the current offset
	dataRec := fileDataSource.dataInMemory[fileDataSource.currentOffset]
	fileDataSource.currentOffset++
	return dataRec, true
}

func (fileDataSource *FileDataSourceReader) readFileIntoMemory() {
	file, err := os.Open(fileDataSource.filePath)
	if err != nil {
		log.Fatalf("Unable to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var records []*MetricRecord

	// Assuming the first row is a header
	if _, err := reader.Read(); err != nil {
		log.Fatalf("Unable to read header from CSV: %v", err)
	}

	// Iterate through the records
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unable to read CSV record: %v", err)
		}

		metricRecord, err := FromCsvDataRecord(record)

		records = append(records, metricRecord)
	}

	fileDataSource.dataInMemory = records
}
