package v1

import (
	"context"
	"encoding/json"
	"net/http"
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
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/deployments", deploymentID, "rule-sets")
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
	server, client := setupTestServer(t, http.StatusOK, ruleContent, "/api/v1/deployments", deploymentID, "rule-sets", ruleFileName)
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
	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/deployments", deploymentID, "rule-sets", ruleFileName)
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
	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/deployments", deploymentID, "rule-sets", ruleFileName)
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
	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/deployments", deploymentID, "rule-sets", ruleFileName)
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
