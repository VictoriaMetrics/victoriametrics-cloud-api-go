package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListDeploymentRuleFileNames(t *testing.T) {
	// Create a sample response
	ruleFiles := []string{
		"alert1.yml",
		"alert2.yml",
		"recording.yml",
	}
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(ruleFiles)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/deployments", deploymentID, "rule-sets", "files")
	defer server.Close()

	// Call the method
	result, err := client.ListDeploymentRuleFileNames(context.Background(), deploymentID)
	if err != nil {
		t.Fatalf("ListDeploymentRuleFileNames() error = %v", err)
	}

	// Check the result
	if len(result) != len(ruleFiles) {
		t.Errorf("ListDeploymentRuleFileNames() returned %d rule files, want %d", len(result), len(ruleFiles))
	}

	// Check the rule file names
	for i, name := range result {
		if name != ruleFiles[i] {
			t.Errorf("ListDeploymentRuleFileNames() rule file %d = %s, want %s", i, name, ruleFiles[i])
		}
	}
}

func TestGetDeploymentRuleFileContent(t *testing.T) {
	// Create a sample response
	ruleContent := `
groups:
- name: example
  rules:
  - alert: HighRequestLatency
    expr: job:request_latency_seconds:mean5m{job="myjob"} > 0.5
    for: 10m
    labels:
      severity: page
    annotations:
      summary: High request latency
`
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"
	ruleFileName := "alert1.yml"

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, ruleContent, "/api/v1/deployments", deploymentID, "rule-sets", "files", ruleFileName)
	defer server.Close()

	// Call the method
	result, err := client.GetDeploymentRuleFileContent(context.Background(), deploymentID, ruleFileName)
	if err != nil {
		t.Fatalf("GetDeploymentRuleFileContent() error = %v", err)
	}

	// Check the result
	if result != ruleContent {
		t.Errorf("GetDeploymentRuleFileContent() = %s, want %s", result, ruleContent)
	}
}

func TestUpdateDeploymentRuleFileContent(t *testing.T) {
	// Create a sample rule content
	ruleContent := `
groups:
- name: example
  rules:
  - alert: HighRequestLatency
    expr: job:request_latency_seconds:mean5m{job="myjob"} > 0.5
    for: 10m
    labels:
      severity: page
    annotations:
      summary: High request latency
`
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"
	ruleFileName := "alert1.yml"

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/deployments", deploymentID, "rule-sets", "files", ruleFileName)
	defer server.Close()

	// Call the method
	err := client.UpdateDeploymentRuleFileContent(context.Background(), deploymentID, ruleFileName, ruleContent)
	if err != nil {
		t.Fatalf("UpdateDeploymentRuleFileContent() error = %v", err)
	}
}

func TestCreateDeploymentRuleFileContent(t *testing.T) {
	// Create a sample rule content
	ruleContent := `
groups:
- name: example
  rules:
  - alert: HighRequestLatency
    expr: job:request_latency_seconds:mean5m{job="myjob"} > 0.5
    for: 10m
    labels:
      severity: page
    annotations:
      summary: High request latency
`
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"
	ruleFileName := "new-alert.yml"

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/deployments", deploymentID, "rule-sets", "files", ruleFileName)
	defer server.Close()

	// Call the method
	err := client.CreateDeploymentRuleFileContent(context.Background(), deploymentID, ruleFileName, ruleContent)
	if err != nil {
		t.Fatalf("CreateDeploymentRuleFileContent() error = %v", err)
	}
}

func TestDeleteDeploymentRuleFile(t *testing.T) {
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"
	ruleFileName := "alert1.yml"

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/deployments", deploymentID, "rule-sets", "files", ruleFileName)
	defer server.Close()

	// Call the method
	err := client.DeleteDeploymentRuleFile(context.Background(), deploymentID, ruleFileName)
	if err != nil {
		t.Fatalf("DeleteDeploymentRuleFile() error = %v", err)
	}
}

func TestRuleFileOperations_InvalidDeploymentID(t *testing.T) {
	// Create a client
	client, err := New("test-api-key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Call the methods with an invalid deployment ID
	invalidID := "invalid-id"
	ruleFileName := "alert1.yml"
	ruleContent := "test content"

	_, err = client.ListDeploymentRuleFileNames(context.Background(), invalidID)
	if err == nil {
		t.Errorf("ListDeploymentRuleFileNames() with invalid deployment ID should return an error")
	}

	_, err = client.GetDeploymentRuleFileContent(context.Background(), invalidID, ruleFileName)
	if err == nil {
		t.Errorf("GetDeploymentRuleFileContent() with invalid deployment ID should return an error")
	}

	err = client.UpdateDeploymentRuleFileContent(context.Background(), invalidID, ruleFileName, ruleContent)
	if err == nil {
		t.Errorf("UpdateDeploymentRuleFileContent() with invalid deployment ID should return an error")
	}

	err = client.CreateDeploymentRuleFileContent(context.Background(), invalidID, ruleFileName, ruleContent)
	if err == nil {
		t.Errorf("CreateDeploymentRuleFileContent() with invalid deployment ID should return an error")
	}

	err = client.DeleteDeploymentRuleFile(context.Background(), invalidID, ruleFileName)
	if err == nil {
		t.Errorf("DeleteDeploymentRuleFile() with invalid deployment ID should return an error")
	}
}

const conditionalDeploymentID = "123e4567-e89b-12d3-a456-426614174000"

// recordingServer answers every request with status, body and the given
// headers, and keeps the headers of the last request it received.
func recordingServer(t *testing.T, status int, body string, respHeaders map[string]string) (*VMCloudAPIClient, *http.Header) {
	t.Helper()
	var got http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
		for k, v := range respHeaders {
			w.Header().Set(k, v)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	client, err := New("test-api-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	return client, &got
}

func TestRuleFileWriteOptions(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		call            func(c *VMCloudAPIClient) error
		wantIfMatch     string
		wantIfNoneMatch string
		wantActor       string
	}{
		"update at a version": {
			call: func(c *VMCloudAPIClient) error {
				return c.UpdateDeploymentRuleFileContentWithOptions(ctx, conditionalDeploymentID, "a.yml", "groups: []", WithIfMatch("abc"), WithActor("ops@example.com"))
			},
			wantIfMatch: `"abc"`,
			wantActor:   "ops@example.com",
		},
		"already quoted etag is kept": {
			call: func(c *VMCloudAPIClient) error {
				return c.UpdateDeploymentRuleFileContentWithOptions(ctx, conditionalDeploymentID, "a.yml", "groups: []", WithIfMatch(`"abc"`))
			},
			wantIfMatch: `"abc"`,
		},
		"create only if absent": {
			call: func(c *VMCloudAPIClient) error {
				return c.UpdateDeploymentRuleFileContentWithOptions(ctx, conditionalDeploymentID, "a.yml", "groups: []", WithIfNoneMatchAny())
			},
			wantIfNoneMatch: "*",
		},
		"delete at a version": {
			call: func(c *VMCloudAPIClient) error {
				return c.DeleteDeploymentRuleFileWithOptions(ctx, conditionalDeploymentID, "a.yml", WithIfMatch("abc"))
			},
			wantIfMatch: `"abc"`,
		},
		"weak etag is kept": {
			call: func(c *VMCloudAPIClient) error {
				return c.UpdateDeploymentRuleFileContentWithOptions(ctx, conditionalDeploymentID, "a.yml", "groups: []", WithIfMatch(`W/"abc"`))
			},
			wantIfMatch: `W/"abc"`,
		},
		"no options sends no conditions": {
			call: func(c *VMCloudAPIClient) error {
				return c.UpdateDeploymentRuleFileContent(ctx, conditionalDeploymentID, "a.yml", "groups: []")
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client, got := recordingServer(t, http.StatusOK, "", nil)
			if err := tc.call(client); err != nil {
				t.Fatalf("call: %v", err)
			}
			h := *got
			if h.Get("If-Match") != tc.wantIfMatch || h.Get("If-None-Match") != tc.wantIfNoneMatch || h.Get("X-VM-Actor") != tc.wantActor {
				t.Fatalf("headers If-Match=%q If-None-Match=%q X-VM-Actor=%q, want %q %q %q",
					h.Get("If-Match"), h.Get("If-None-Match"), h.Get("X-VM-Actor"),
					tc.wantIfMatch, tc.wantIfNoneMatch, tc.wantActor)
			}
		})
	}
}

func TestRuleFileEmptyETagNotSent(t *testing.T) {
	ctx := context.Background()
	calls := map[string]func(c *VMCloudAPIClient) error{
		"update": func(c *VMCloudAPIClient) error {
			return c.UpdateDeploymentRuleFileContentWithOptions(ctx, conditionalDeploymentID, "a.yml", "groups: []", WithIfMatch(""))
		},
		"delete": func(c *VMCloudAPIClient) error {
			return c.DeleteDeploymentRuleFileWithOptions(ctx, conditionalDeploymentID, "a.yml", WithIfMatch(""))
		},
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			client, got := recordingServer(t, http.StatusOK, "", nil)
			if err := call(client); !errors.Is(err, ErrEmptyETag) {
				t.Fatalf("err = %v, want ErrEmptyETag", err)
			}
			if *got != nil {
				t.Fatalf("request was sent with headers %v", *got)
			}
		})
	}
}

func TestCreateRuleFileErrorText(t *testing.T) {
	client, _ := recordingServer(t, http.StatusConflict, `{"error": "exists"}`, nil)
	err := client.CreateDeploymentRuleFileContent(context.Background(), conditionalDeploymentID, "a.yml", "x")
	if err == nil || !strings.Contains(err.Error(), "failed to create rule file") {
		t.Fatalf("err = %v, want a create error", err)
	}
}

func TestGetDeploymentRuleFile(t *testing.T) {
	for header, want := range map[string]string{`"abc"`: `"abc"`, `W/"abc"`: `W/"abc"`, "": ""} {
		client, _ := recordingServer(t, http.StatusOK, "groups: []", map[string]string{"ETag": header})
		file, err := client.GetDeploymentRuleFile(context.Background(), conditionalDeploymentID, "a.yml")
		if err != nil {
			t.Fatalf("GetDeploymentRuleFile: %v", err)
		}
		if file.Content != "groups: []" || file.ETag != want {
			t.Fatalf("ETag header %q: got %+v, want content and ETag %q", header, file, want)
		}
	}
}

func TestRuleFileAPIError(t *testing.T) {
	client, _ := recordingServer(t, http.StatusConflict, `{"error": "rule file \"a.yml\" was changed since it was read"}`, nil)
	err := client.UpdateDeploymentRuleFileContentWithOptions(context.Background(), conditionalDeploymentID, "a.yml", "x", WithIfMatch("old"))

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %v, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusConflict || apiErr.Message != `rule file "a.yml" was changed since it was read` {
		t.Fatalf("APIError = %+v", apiErr)
	}
	// The text is unchanged, for callers that match on it.
	if !strings.Contains(err.Error(), "unexpected status code: 409, body: ") {
		t.Fatalf("err = %q, want the existing message format", err)
	}
}
