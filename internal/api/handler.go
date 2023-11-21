package api

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"valery-datadog-datastream-demo/internal/data"
	"valery-datadog-datastream-demo/internal/processor"
)

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

		var getDataReq data.GetDataRequest
		if err := json.Unmarshal(message, &getDataReq); err != nil {
			log.Println("Error unmarshaling request:", err)
			continue
		}

		// Parse filtering, partitioning and aggregator from the request
		filters := data.FromRequestFilters(getDataReq.Filters)
		partitioner := processor.FromRequestScale(getDataReq.Scale)
		aggregator := processor.FromRequestAggregator(getDataReq.Aggregator)

		// Fetch data points
		dataPoints := metricDataProvider.GetMetricDataPoints(filters, partitioner, aggregator)

		// Send data points back to the client
		if err = ws.WriteJSON(dataPoints); err != nil {
			log.Println("Error sending data points:", err)
		}
	}
}

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

		var getFiltersReq data.GetFiltersRequest
		if err := json.Unmarshal(message, &getFiltersReq); err != nil {
			log.Println("Error unmarshaling request:", err)
			continue
		}

		// Fetch filters (tag name:value pairs)
		filters := metricDataProvider.GetMetricTagFilters(getFiltersReq.Query)

		// Send data points back to the client
		if err = ws.WriteJSON(filters); err != nil {
			log.Println("Error sending data points:", err)
		}
	}
}
