package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestListDeployments(t *testing.T) {
	// Create a sample response
	deployments := DeploymentSummaryList{
		{
			ID:            "123e4567-e89b-12d3-a456-426614174000",
			Name:          "test-deployment-1",
			Type:          DeploymentTypeSingleNode,
			CloudProvider: DeploymentCloudProviderAWS,
			Region:        "us-east-1",
			Status:        DeploymentStatusRunning,
			Version:       "1.0.0",
			CreatedAt:     time.Now(),
		},
		{
			ID:            "223e4567-e89b-12d3-a456-426614174001",
			Name:          "test-deployment-2",
			Type:          DeploymentTypeCluster,
			CloudProvider: DeploymentCloudProviderAWS,
			Region:        "us-west-1",
			Status:        DeploymentStatusProvisioning,
			Version:       "1.0.0",
			CreatedAt:     time.Now(),
		},
	}

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(deployments)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/deployments")
	defer server.Close()

	// Call the method
	result, err := client.ListDeployments(context.Background())
	if err != nil {
		t.Fatalf("ListDeployments() error = %v", err)
	}

	// Check the result
	if len(result) != len(deployments) {
		t.Errorf("ListDeployments() returned %d deployments, want %d", len(result), len(deployments))
	}

	// Check the first deployment
	if result[0].ID != deployments[0].ID {
		t.Errorf("ListDeployments() first deployment ID = %s, want %s", result[0].ID, deployments[0].ID)
	}
	if result[0].Name != deployments[0].Name {
		t.Errorf("ListDeployments() first deployment Name = %s, want %s", result[0].Name, deployments[0].Name)
	}
	if result[0].Type != deployments[0].Type {
		t.Errorf("ListDeployments() first deployment Type = %s, want %s", result[0].Type, deployments[0].Type)
	}
}

func TestGetDeploymentDetails(t *testing.T) {
	// Create a sample response
	deployment := DeploymentInfo{
		ID:                 "123e4567-e89b-12d3-a456-426614174000",
		Name:               "test-deployment",
		Type:               DeploymentTypeSingleNode,
		CloudProvider:      DeploymentCloudProviderAWS,
		Region:             "us-east-1",
		Status:             DeploymentStatusRunning,
		Version:            "1.0.0",
		CreatedAt:          time.Now(),
		Tier:               21,
		StorageSizeGb:      10,
		RetentionValue:     30,
		RetentionUnit:      DurationUnitDay,
		DeduplicationValue: 10,
		DeduplicationUnit:  DurationUnitSecond,
		MaintenanceWindow:  MaintenanceWindowWeekendDays,
		AccessEndpoint:     "https://test-deployment.victoriametrics.com",
	}

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(deployment)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/deployments", deployment.ID)
	defer server.Close()

	// Call the method
	result, err := client.GetDeploymentDetails(context.Background(), deployment.ID)
	if err != nil {
		t.Fatalf("GetDeploymentDetails() error = %v", err)
	}

	// Check the result
	if result.ID != deployment.ID {
		t.Errorf("GetDeploymentDetails() ID = %s, want %s", result.ID, deployment.ID)
	}
	if result.Name != deployment.Name {
		t.Errorf("GetDeploymentDetails() Name = %s, want %s", result.Name, deployment.Name)
	}
	if result.Type != deployment.Type {
		t.Errorf("GetDeploymentDetails() Type = %s, want %s", result.Type, deployment.Type)
	}
	if result.CloudProvider != deployment.CloudProvider {
		t.Errorf("GetDeploymentDetails() CloudProvider = %s, want %s", result.CloudProvider, deployment.CloudProvider)
	}
	if result.StorageSizeGb != deployment.StorageSizeGb {
		t.Errorf("GetDeploymentDetails() StorageSizeGb = %d, want %d", result.StorageSizeGb, deployment.StorageSizeGb)
	}
}

func TestCreateDeployment(t *testing.T) {
	// Create a sample request
	request := DeploymentCreationRequest{
		Name:              "test-deployment",
		Type:              DeploymentTypeSingleNode,
		Provider:          DeploymentCloudProviderAWS,
		Region:            "us-east-1",
		Tier:              21,
		StorageSize:       10,
		StorageSizeUnit:   StorageUnitGB,
		Retention:         30,
		RetentionUnit:     DurationUnitDay,
		Deduplication:     10,
		DeduplicationUnit: DurationUnitSecond,
		MaintenanceWindow: MaintenanceWindowWeekendDays,
	}

	// Create a sample response
	response := DeploymentInfo{
		ID:                 "123e4567-e89b-12d3-a456-426614174000",
		Name:               request.Name,
		Type:               request.Type,
		CloudProvider:      request.Provider,
		Region:             request.Region,
		Status:             DeploymentStatusProvisioning,
		Version:            "1.0.0",
		CreatedAt:          time.Now(),
		Tier:               request.Tier,
		StorageSizeGb:      request.StorageSize,
		RetentionValue:     request.Retention,
		RetentionUnit:      request.RetentionUnit,
		DeduplicationValue: request.Deduplication,
		DeduplicationUnit:  request.DeduplicationUnit,
		MaintenanceWindow:  request.MaintenanceWindow,
		AccessEndpoint:     "https://test-deployment.victoriametrics.com",
	}

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/deployments")
	defer server.Close()

	// Call the method
	result, err := client.CreateDeployment(context.Background(), request)
	if err != nil {
		t.Fatalf("CreateDeployment() error = %v", err)
	}

	// Check the result
	if result.Name != request.Name {
		t.Errorf("CreateDeployment() Name = %s, want %s", result.Name, request.Name)
	}
	if result.Type != request.Type {
		t.Errorf("CreateDeployment() Type = %s, want %s", result.Type, request.Type)
	}
	if result.CloudProvider != request.Provider {
		t.Errorf("CreateDeployment() CloudProvider = %s, want %s", result.CloudProvider, request.Provider)
	}
	if result.StorageSizeGb != request.StorageSize {
		t.Errorf("CreateDeployment() StorageSizeGb = %d, want %d", result.StorageSizeGb, request.StorageSize)
	}
}

func TestUpdateDeployment(t *testing.T) {
	// Create a sample request
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"
	request := DeploymentUpdateRequest{
		Name:              "updated-deployment",
		Tier:              22,
		StorageSize:       20,
		StorageSizeUnit:   StorageUnitGB,
		Retention:         60,
		RetentionUnit:     DurationUnitDay,
		Deduplication:     15,
		DeduplicationUnit: DurationUnitSecond,
		MaintenanceWindow: MaintenanceWindowBusinessDays,
	}

	// Create a sample response
	response := DeploymentInfo{
		ID:                 deploymentID,
		Name:               request.Name,
		Type:               DeploymentTypeSingleNode,
		CloudProvider:      DeploymentCloudProviderAWS,
		Region:             "us-east-1",
		Status:             DeploymentStatusRunning,
		Version:            "1.0.0",
		CreatedAt:          time.Now(),
		Tier:               request.Tier,
		StorageSizeGb:      request.StorageSize,
		RetentionValue:     request.Retention,
		RetentionUnit:      request.RetentionUnit,
		DeduplicationValue: request.Deduplication,
		DeduplicationUnit:  request.DeduplicationUnit,
		MaintenanceWindow:  request.MaintenanceWindow,
		AccessEndpoint:     "https://test-deployment.victoriametrics.com",
	}

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/deployments", deploymentID)
	defer server.Close()

	// Call the method
	result, err := client.UpdateDeployment(context.Background(), deploymentID, request)
	if err != nil {
		t.Fatalf("UpdateDeployment() error = %v", err)
	}

	// Check the result
	if result.ID != deploymentID {
		t.Errorf("UpdateDeployment() ID = %s, want %s", result.ID, deploymentID)
	}
	if result.Name != request.Name {
		t.Errorf("UpdateDeployment() Name = %s, want %s", result.Name, request.Name)
	}
	if result.Tier != request.Tier {
		t.Errorf("UpdateDeployment() Tier = %d, want %d", result.Tier, request.Tier)
	}
	if result.StorageSizeGb != request.StorageSize {
		t.Errorf("UpdateDeployment() StorageSizeGb = %d, want %d", result.StorageSizeGb, request.StorageSize)
	}
}

func TestDeleteDeployment(t *testing.T) {
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/deployments", deploymentID)
	defer server.Close()

	// Call the method
	err := client.DeleteDeployment(context.Background(), deploymentID)
	if err != nil {
		t.Fatalf("DeleteDeployment() error = %v", err)
	}
}

func TestDeleteDeployment_InvalidID(t *testing.T) {
	// Create a client
	client, err := New("test-api-key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Call the method with an invalid ID
	err = client.DeleteDeployment(context.Background(), "invalid-id")
	if err == nil {
		t.Fatalf("DeleteDeployment() with invalid ID should return an error")
	}
}
