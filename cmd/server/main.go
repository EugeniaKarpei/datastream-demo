package main

import (
	"log"
	"net/http"
	"valery-datadog-datastream-demo/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

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

	dataSource := data.NewFileDataStream(config.CsvDataSetFilePath)
	metricProcessor := processor.NewInMemoryMetricStreamProcessor()

	// stream csv data into the metric processor
	dataSource.Stream(metricProcessor)

	router.GET("/getData", func(c *gin.Context) {
		api.HandleGetDataWebSocket(metricProcessor, upgrader, c.Request, c.Writer)
	})

	router.GET("/getFilters", func(c *gin.Context) {
		api.HandleGetFiltersWebSocket(metricProcessor, upgrader, c.Request, c.Writer)
	})

	log.Fatal(router.Run(":8080"))
}
