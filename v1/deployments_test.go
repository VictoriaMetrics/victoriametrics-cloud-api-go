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
		DeduplicationValue: new(uint32(10)),
		DeduplicationUnit:  new(DurationUnitSecond),
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
		Deduplication:     new(uint32(10)),
		DeduplicationUnit: new(DurationUnitSecond),
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
		Deduplication:     new(uint32(15)),
		DeduplicationUnit: new(DurationUnitSecond),
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

// TestDeduplicationOmittedForLogsAndTraces checks that a create or update request for a
// deployment type without a deduplication window leaves both deduplication fields out of
// the body, rather than sending a zero window with an empty unit for the API to reject.
// TestCreateRequestDeduplicationBody pins down what a create request puts on the wire.
// nil leaves both fields out; a set window is sent as-is, zero included, so that a zero
// window stays distinct from having no window at all.
func TestCreateRequestDeduplicationBody(t *testing.T) {
	f := func(name string, deduplication *uint32, unit *DurationUnit, wantWindow string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			request := DeploymentCreationRequest{
				Name: "example", Type: DeploymentTypeSingleNode, Provider: DeploymentCloudProviderAWS,
				Region: "us-east-1", Tier: 21, StorageSize: 20, StorageSizeUnit: StorageUnitGB,
				Deduplication: deduplication, DeduplicationUnit: unit,
				Retention: 30, RetentionUnit: DurationUnitDay,
			}
			encoded, err := json.Marshal(request)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			var body map[string]json.RawMessage
			if err := json.Unmarshal(encoded, &body); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			window, hasWindow := body["deduplication"]
			_, hasUnit := body["deduplication_unit"]
			if hasWindow != hasUnit {
				t.Errorf("Marshal() = %s, sent a half-specified window", encoded)
			}
			if wantWindow == "" {
				if hasWindow {
					t.Errorf("Marshal() = %s, want no deduplication window", encoded)
				}
				return
			}
			if got := string(window); got != wantWindow {
				t.Errorf("Marshal() deduplication = %s, want %s", got, wantWindow)
			}
		})
	}

	f("window set", new(uint32(30)), new(DurationUnitSecond), "30")
	f("zero window", new(uint32(0)), new(DurationUnitSecond), "0")
	f("nothing set", nil, nil, "")
}

// byte-valued tier limits must stay 64-bit: a tier can report more than 2 GB, which does
// not fit in a 32-bit int and fails to decode on 32-bit targets. These assignments stop
// compiling if any of them is narrowed back to int.
var (
	_ int64 = LogsTierInfo{}.IngestionRateBytes
	_ int64 = LogsTierInfo{}.DataReadRate
	_ int64 = LogsTierInfo{}.BytesPerQuery
	_ int64 = TracesTierInfo{}.IngestionRateBytes
	_ int64 = TracesTierInfo{}.DataReadRate
	_ int64 = TracesTierInfo{}.BytesPerQuery
)

func TestDeduplicationOmittedForLogsAndTraces(t *testing.T) {
	f := func(name string, request any, wantDeduplication bool) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			encoded, err := json.Marshal(request)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			var body map[string]json.RawMessage
			if err := json.Unmarshal(encoded, &body); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			_, hasUnit := body["deduplication_unit"]
			if hasUnit != wantDeduplication {
				t.Errorf("Marshal() = %s, deduplication_unit present = %v, want %v", encoded, hasUnit, wantDeduplication)
			}
		})
	}

	f("create vlogs", DeploymentCreationRequest{
		Name:          "l1",
		Type:          DeploymentTypeVLogs,
		Tier:          101,
		Retention:     30,
		RetentionUnit: DurationUnitDay,
	}, false)
	f("create vtraces", DeploymentCreationRequest{
		Name:          "t1",
		Type:          DeploymentTypeVTraces,
		Tier:          201,
		Retention:     30,
		RetentionUnit: DurationUnitDay,
	}, false)
	f("create metrics", DeploymentCreationRequest{
		Name:              "m1",
		Type:              DeploymentTypeSingleNode,
		Tier:              21,
		Deduplication:     new(uint32(30)),
		DeduplicationUnit: new(DurationUnitSecond),
		Retention:         30,
		RetentionUnit:     DurationUnitDay,
	}, true)
	f("update without deduplication", DeploymentUpdateRequest{
		Name:          "l1",
		Tier:          101,
		Retention:     30,
		RetentionUnit: DurationUnitDay,
	}, false)
	f("update with deduplication", DeploymentUpdateRequest{
		Name:              "m1",
		Tier:              21,
		Deduplication:     new(uint32(30)),
		DeduplicationUnit: new(DurationUnitSecond),
		Retention:         30,
		RetentionUnit:     DurationUnitDay,
	}, true)
	f("update with a zero deduplication window", DeploymentUpdateRequest{
		Name:              "m1",
		Tier:              21,
		Deduplication:     new(uint32(0)),
		DeduplicationUnit: new(DurationUnitSecond),
		Retention:         30,
		RetentionUnit:     DurationUnitDay,
	}, true)
}
