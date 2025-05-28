package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestListDeploymentAccessTokens(t *testing.T) {
	// Create a sample response
	tokens := AccessTokensList{
		{
			ID:          "token-id-1",
			Secret:      "****",
			Type:        AccessModeRead,
			Description: "Test token 1",
			CreatedBy:   "test-user",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "token-id-2",
			Secret:      "****",
			Type:        AccessModeReadWrite,
			Description: "Test token 2",
			CreatedBy:   "test-user",
			CreatedAt:   time.Now(),
		},
	}

	deploymentID := "123e4567-e89b-12d3-a456-426614174000"

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(tokens)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/deployments", deploymentID, "access_tokens")
	defer server.Close()

	// Call the method
	result, err := client.ListDeploymentAccessTokens(context.Background(), deploymentID)
	if err != nil {
		t.Fatalf("ListDeploymentAccessTokens() error = %v", err)
	}

	// Check the result
	if len(result) != len(tokens) {
		t.Errorf("ListDeploymentAccessTokens() returned %d tokens, want %d", len(result), len(tokens))
	}

	// Check the first token
	if result[0].ID != tokens[0].ID {
		t.Errorf("ListDeploymentAccessTokens() first token ID = %s, want %s", result[0].ID, tokens[0].ID)
	}
	if result[0].Type != tokens[0].Type {
		t.Errorf("ListDeploymentAccessTokens() first token Type = %s, want %s", result[0].Type, tokens[0].Type)
	}
	if result[0].Description != tokens[0].Description {
		t.Errorf("ListDeploymentAccessTokens() first token Description = %s, want %s", result[0].Description, tokens[0].Description)
	}
}

func TestCreateDeploymentAccessToken(t *testing.T) {
	// Create a sample request
	request := AccessTokenCreateRequest{
		Type:        AccessModeReadWrite,
		Description: "Test token description",
	}

	// Create a sample response
	response := AccessToken{
		ID:          "token-id-1",
		Secret:      "token-value-123",
		Type:        request.Type,
		Description: request.Description,
		CreatedBy:   "test-user",
		CreatedAt:   time.Now(),
	}

	deploymentID := "123e4567-e89b-12d3-a456-426614174000"

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/deployments", deploymentID, "access_tokens")
	defer server.Close()

	// Call the method
	result, err := client.CreateDeploymentAccessToken(context.Background(), deploymentID, request)
	if err != nil {
		t.Fatalf("CreateDeploymentAccessToken() error = %v", err)
	}

	// Check the result
	if result.Type != request.Type {
		t.Errorf("CreateDeploymentAccessToken() Type = %s, want %s", result.Type, request.Type)
	}
	if result.Description != request.Description {
		t.Errorf("CreateDeploymentAccessToken() Description = %s, want %s", result.Description, request.Description)
	}
	if result.Secret != response.Secret {
		t.Errorf("CreateDeploymentAccessToken() Secret = %s, want %s", result.Secret, response.Secret)
	}
}

func TestRevealDeploymentAccessToken(t *testing.T) {
	// Create a sample response
	response := AccessToken{
		ID:          "token-id-1",
		Secret:      "token-value-123",
		Type:        AccessModeReadWrite,
		Description: "Test token description",
		CreatedBy:   "test-user",
		CreatedAt:   time.Now(),
	}

	deploymentID := "123e4567-e89b-12d3-a456-426614174000"
	tokenID := "token-id-1"

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/deployments", deploymentID, "access_tokens", tokenID)
	defer server.Close()

	// Call the method

	result, err := client.RevealDeploymentAccessToken(context.Background(), deploymentID, tokenID)
	if err != nil {
		t.Fatalf("RevealDeploymentAccessToken() error = %v", err)
	}

	// Check the result
	if result.ID != response.ID {
		t.Errorf("RevealDeploymentAccessToken() ID = %s, want %s", result.ID, response.ID)
	}
	if result.Secret != response.Secret {
		t.Errorf("RevealDeploymentAccessToken() Secret = %s, want %s", result.Secret, response.Secret)
	}
}

func TestDeleteDeploymentAccessToken(t *testing.T) {
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"
	tokenID := "token-id-1"

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/deployments", deploymentID, "access_tokens", tokenID)
	defer server.Close()

	// Call the method
	err := client.DeleteDeploymentAccessToken(context.Background(), deploymentID, tokenID)
	if err != nil {
		t.Fatalf("DeleteDeploymentAccessToken() error = %v", err)
	}
}

func TestDeleteDeploymentAccessToken_InvalidDeploymentID(t *testing.T) {
	// Create a client
	client, err := New("test-api-key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Call the method with an invalid deployment ID
	err = client.DeleteDeploymentAccessToken(context.Background(), "invalid-id", "token-id-1")
	if err == nil {
		t.Fatalf("DeleteDeploymentAccessToken() with invalid deployment ID should return an error")
	}
}
