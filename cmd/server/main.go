package main

import (
	"fmt"
	"log"
	"net/http"

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

	dataPoints := metricProcessor.GetMetricDataPoints(
		[]*data.Tag{
			data.NewFilterTag("location", "Chicago"),
			data.NewFilterTag("gender", "F"),
		},
		processor.MonthlyTimePartitioner,
		processor.AvgAggregator,
	)

	for _, dp := range dataPoints {
		fmt.Println("Prepared data-point > time:", dp.Timestamp(), " value:", dp.Value())
	}

	router.GET("/getData", func(c *gin.Context) {
		handleWebSocket(c.Writer, c.Request)
	})

	router.POST("/setStreams", func(c *gin.Context) {
		// Implement logic to handle stream preferences
		c.JSON(http.StatusOK, gin.H{"status": "stream preferences updated"})
	})

	log.Fatal(router.Run(":8080"))
}

// handleWebSocket handles the WebSocket connection and data streaming
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	defer ws.Close()

	// Stream data in a loop or based on some event
	for {
		// Example data, replace with actual time series data fetching logic
		message := "example time series data"
		if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("Error writing websocket message:", err)
			break
		}
		// Implement the logic to periodically send data
	}
}
