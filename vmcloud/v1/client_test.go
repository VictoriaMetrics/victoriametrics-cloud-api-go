package v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		options []VMCloudAPIClientOption
		wantErr bool
	}{
		{
			name:    "valid client with API key",
			apiKey:  "test-api-key",
			options: nil,
			wantErr: false,
		},
		{
			name:    "empty API key",
			apiKey:  "",
			options: nil,
			wantErr: true,
		},
		{
			name:   "valid client with custom HTTP client",
			apiKey: "test-api-key",
			options: []VMCloudAPIClientOption{
				WithHTTPClient(&http.Client{
					Timeout: 30 * time.Second,
				}),
			},
			wantErr: false,
		},
		{
			name:   "valid client with custom base URL",
			apiKey: "test-api-key",
			options: []VMCloudAPIClientOption{
				WithBaseURL("https://custom-api.victoriametrics.com"),
			},
			wantErr: false,
		},
		{
			name:   "invalid base URL",
			apiKey: "test-api-key",
			options: []VMCloudAPIClientOption{
				WithBaseURL("://invalid-url"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.apiKey, tt.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if client == nil {
					t.Errorf("New() returned nil client without error")
				} else {
					if client.apiKey != tt.apiKey {
						t.Errorf("New() client.apiKey = %v, want %v", client.apiKey, tt.apiKey)
					}
				}
			}
		})
	}
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	client, err := New("test-api-key", WithHTTPClient(customClient))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if client.c != customClient {
		t.Errorf("WithHTTPClient() did not set the custom HTTP client")
	}
}

func TestWithBaseURL(t *testing.T) {
	customBaseURL := "https://custom-api.victoriametrics.com"

	client, err := New("test-api-key", WithBaseURL(customBaseURL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if client.BaseURL() != customBaseURL {
		t.Errorf("WithBaseURL() client.BaseURL() = %v, want %v", client.BaseURL(), customBaseURL)
	}
}

func TestBaseURL(t *testing.T) {
	defaultClient, err := New("test-api-key")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if defaultClient.BaseURL() != DefaultBaseURL {
		t.Errorf("BaseURL() = %v, want %v", defaultClient.BaseURL(), DefaultBaseURL)
	}

	customBaseURL := "https://custom-api.victoriametrics.com"
	customClient, err := New("test-api-key", WithBaseURL(customBaseURL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if customClient.BaseURL() != customBaseURL {
		t.Errorf("BaseURL() = %v, want %v", customClient.BaseURL(), customBaseURL)
	}
}

// setupTestServer creates a test HTTP server that returns the given response
func setupTestServer(t *testing.T, statusCode int, response string, path ...string) (*httptest.Server, *VMCloudAPIClient) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the API key is set in the request header
		if r.Header.Get(AccessTokenHeader) != "test-api-key" {
			t.Errorf("Request missing API key header")
		}
		expectedPath := strings.Join(path, "/")
		if len(path) > 0 && r.URL.Path != expectedPath {
			t.Errorf("Expected request path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))

	client, err := New("test-api-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return server, client
}
