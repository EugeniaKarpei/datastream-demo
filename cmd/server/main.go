package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader specifies parameters for upgrading an HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func main() {
	// Initialize Gin router
	router := gin.Default()

	// WebSocket endpoint for streaming data
	router.GET("/getData", func(c *gin.Context) {
		handleWebSocket(c.Writer, c.Request)
	})

	// API endpoint to set stream preferences
	router.POST("/setStreams", func(c *gin.Context) {
		// Implement logic to handle stream preferences
		c.JSON(http.StatusOK, gin.H{"status": "stream preferences updated"})
	})

	// Start the server
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
