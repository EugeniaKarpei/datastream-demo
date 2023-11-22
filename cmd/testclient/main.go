package main

// Test client, simulates requests to service endpoints and prints responses to standard output.
// Used for testing locally.

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"valery-datadog-datastream-demo/internal/data"
)

const SERVICE_BASE_URL = "ws://localhost:8080"

func main() {
	// Test /getData API
	jsonReq, err := json.Marshal(data.GetDataRequest{
		Filters: []string{
			"location:Chicago",
			"gender:F",
		},
		Scale: "Weekly",
	})
	if err != nil {
		log.Fatal("error marshaling JSON:", err)
	}
	testApi(SERVICE_BASE_URL+"/getData", jsonReq)

	// Test /getFilters API
	jsonReq, err = json.Marshal(data.GetFiltersRequest{
		Query: "gen",
	})
	testApi(SERVICE_BASE_URL+"/getFilters", jsonReq)
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
