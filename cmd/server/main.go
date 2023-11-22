package main

// Main service startup entry point.

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"valery-datadog-datastream-demo/internal/api"
	"valery-datadog-datastream-demo/internal/config"
	"valery-datadog-datastream-demo/internal/data"
	"valery-datadog-datastream-demo/internal/processor"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func main() {
	router := gin.Default()

	// Create data stream from csv file and a metric processor instance. For this demo it handles a single metric
	dataStream := data.NewFileDataStream(config.CsvDataSetFilePath)
	metricProcessor := processor.NewInMemoryMetricStreamProcessor()

	// Stream data into the metric processor
	dataStream.Stream(metricProcessor)

	// Register API endpoints
	// getData - main flow - to fetch metrics using filters, partitioners and aggregate them
	router.GET("/getData", func(c *gin.Context) {
		api.HandleGetDataWebSocket(metricProcessor, upgrader, c.Request, c.Writer)
	})

	// getFilters - secondary flow - get available tag:value pairs by given prefix string
	router.GET("/getFilters", func(c *gin.Context) {
		api.HandleGetFiltersWebSocket(metricProcessor, upgrader, c.Request, c.Writer)
	})

	// Start the server
	log.Fatal(router.Run(":8080"))
}
