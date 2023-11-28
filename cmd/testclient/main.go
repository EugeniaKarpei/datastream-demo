package main

// Test client, simulates requests to service endpoints and prints responses to standard output.
// Used for testing locally.

import (
	"encoding/json"
	"fmt"
	"log"
	"valery-datadog-datastream-demo/internal/data"

	"github.com/gorilla/websocket"
)

const SERVICE_BASE_URL = "ws://localhost:8080"

var TestGetDataRequests = []data.GetDataRequest{
	{
		Filters: []string{
			"location:Chicago",
			"gender:F",
		},
		Scale:      "Weekly",
		Aggregator: "Sum",
	},
	{
		Filters: []string{
			"location:California",
		},
		Scale:      "Daily",
		Aggregator: "Count",
	},
	{
		Filters: []string{
			"coupon_status:Not Used",
		},
		Scale:      "Monthly",
		Aggregator: "Avg",
	},
}

var TestGetFiltersRequest = []data.GetFiltersRequest{
	{
		Query: "loc", // should return all locations
	},
	{
		Query: "coup", // should return both all coupon codes and statuses
	},
}

func main() {
	// Test /getData API
	for _, req := range TestGetDataRequests {
		jsonReq, _ := json.Marshal(req)
		testApi(SERVICE_BASE_URL+"/getData", jsonReq)
	}

	// Test /getFilters API
	for _, req := range TestGetFiltersRequest {
		jsonReq, _ := json.Marshal(req)
		testApi(SERVICE_BASE_URL+"/getFilters", jsonReq)
	}
}

func testApi(apiUrl string, jsonReq []byte) {
	c, _, err := websocket.DefaultDialer.Dial(apiUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Send the JSON request
	err = c.WriteMessage(websocket.TextMessage, jsonReq)
	if err != nil {
		log.Fatal("write:", err)
	}

	// Read and print the response
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Fatal("read:", err)
	}
	fmt.Printf("recv: %s\n", message)
}
