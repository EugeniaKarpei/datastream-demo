package data

// Main interfaces for data access:
// * DataStream
// * StreamProcessor
// Main goal is to decouple stream source from the processor so that they do not know anything about each other.

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type StreamProcessor interface {
	Process(dataRecord []string) error
}

type DataStream interface {
	Stream(streamProcessor StreamProcessor)
}

func NewFileDataStream(filePath string) DataStream {
	fileDataSource := &FileDataStream{
		filePath: filePath,
	}
	return fileDataSource
}

var _ DataStream = (*FileDataStream)(nil)

type FileDataStream struct {
	filePath string
}

// Streams data into the processor
func (fds *FileDataStream) Stream(processor StreamProcessor) {
	file, err := os.Open(fds.filePath)
	if err != nil {
		log.Fatalf("Unable to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

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

		// Process records 1 by 1, if any processing issues - log error
		err = processor.Process(record)

		if err != nil {
			log.Printf("Failed to process CSV record: %v", err)
		}
	}
}
