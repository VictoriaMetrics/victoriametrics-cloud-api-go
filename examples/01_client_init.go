package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/vmcloud/v1"
)

func main() {
	// Example 1: Basic client initialization with API key
	client, err := vmcloud.New("your-api-key")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	fmt.Println("Successfully created client with default settings")
	fmt.Printf("Base URL: %s\n", client.BaseURL())

	// Example 2: Client with custom HTTP client (with timeout)
	customHTTPClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	client, err = vmcloud.New("your-api-key", vmcloud.WithHTTPClient(customHTTPClient))
	if err != nil {
		log.Fatalf("Failed to create client with custom HTTP client: %v", err)
	}

	// ...
}
