package station

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// GetLatestPodResponse Response represents the structure of the JSON response from the API
type GetLatestPodResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

func QueryLatestPodNumber() uint64 {

	// API URL
	url := "http://0.0.0.0:26667/tracks_pod_count"

	// Create a JSON request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "tracks_pod_count",
		"id":      1,
		"params":  []interface{}{},
	})
	if err != nil {
		log.Fatalf("Failed to marshal request body: %v", err)
	}

	// Send POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Unmarshal response into Response struct
	var response GetLatestPodResponse
	if err = json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Convert the result to an integer
	resultInt, err := strconv.Atoi(response.Result)
	if err != nil {
		log.Fatalf("Failed to convert result to integer: %v", err)
	}

	return uint64(resultInt)
}
