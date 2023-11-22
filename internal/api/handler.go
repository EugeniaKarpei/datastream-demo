package api

// API request handler functions. Handle websocket request/response parsing and
// calls MetricProcessor for the main logic.

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"valery-datadog-datastream-demo/internal/data"
	"valery-datadog-datastream-demo/internal/processor"
)

// Handles /getData API call
func HandleGetDataWebSocket(
	metricDataProvider processor.MetricDataProvider,
	upgrader websocket.Upgrader,
	request *http.Request,
	responseWriter http.ResponseWriter,
) {
	ws, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Fatal("Error upgrading to WebSocket:", err)
		return
	}
	defer ws.Close()

	for {
		// Read message from client (getData request)
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Parse json
		var getDataReq data.GetDataRequest
		if err := json.Unmarshal(message, &getDataReq); err != nil {
			log.Println("Error unmarshaling request:", err)
			continue
		}

		// Convert apimodel -> common model entities to use them as parameters for the MetricProcessor
		filters := data.FromRequestFilters(getDataReq.Filters)
		partitioner := processor.FromRequestScale(getDataReq.Scale)
		aggregator := processor.FromRequestAggregator(getDataReq.Aggregator)

		// Fetch data points from MetricProcessor
		dataPoints := metricDataProvider.GetMetricDataPoints(filters, partitioner, aggregator)

		// Send data points to the client
		if err = ws.WriteJSON(dataPoints); err != nil {
			log.Println("Error sending data points:", err)
		}
	}
}

// Handles /getFilters API call
func HandleGetFiltersWebSocket(
	metricDataProvider processor.MetricDataProvider,
	upgrader websocket.Upgrader,
	request *http.Request,
	responseWriter http.ResponseWriter,
) {
	ws, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Fatal("Error upgrading to WebSocket:", err)
		return
	}
	defer ws.Close()

	for {
		// Read message from client (getData request)
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Parse json
		var getFiltersReq data.GetFiltersRequest
		if err := json.Unmarshal(message, &getFiltersReq); err != nil {
			log.Println("Error unmarshaling request:", err)
			continue
		}

		// Fetch filters (tag name:value pairs) from MetricProcessor
		filters := metricDataProvider.GetMetricTagFilters(getFiltersReq.Query)

		// Send data points to the client
		if err = ws.WriteJSON(filters); err != nil {
			log.Println("Error sending data points:", err)
		}
	}
}
