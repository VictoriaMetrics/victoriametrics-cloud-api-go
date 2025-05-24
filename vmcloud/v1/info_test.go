package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListCloudProviders(t *testing.T) {
	// Create a sample response
	providers := CloudProviderInfoList{
		{
			ID:  DeploymentCloudProviderAWS,
			URL: "https://aws.amazon.com",
		},
	}

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(providers)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/cloud_providers")
	defer server.Close()

	// Call the method
	result, err := client.ListCloudProviders(context.Background())
	if err != nil {
		t.Fatalf("ListCloudProviders() error = %v", err)
	}

	// Check the result
	if len(result) != len(providers) {
		t.Errorf("ListCloudProviders() returned %d providers, want %d", len(result), len(providers))
	}

	// Check the first provider
	if result[0].ID != providers[0].ID {
		t.Errorf("ListCloudProviders() first provider ID = %s, want %s", result[0].ID, providers[0].ID)
	}
	if result[0].URL != providers[0].URL {
		t.Errorf("ListCloudProviders() first provider URL = %s, want %s", result[0].URL, providers[0].URL)
	}
}

func TestListRegions(t *testing.T) {
	// Create a sample response
	regions := RegionInfoList{
		{
			CloudProvider: DeploymentCloudProviderAWS,
			Name:          "us-east-1",
		},
		{
			CloudProvider: DeploymentCloudProviderAWS,
			Name:          "us-west-1",
		},
	}

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(regions)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/regions")
	defer server.Close()

	// Call the method
	result, err := client.ListRegions(context.Background())
	if err != nil {
		t.Fatalf("ListRegions() error = %v", err)
	}

	// Check the result
	if len(result) != len(regions) {
		t.Errorf("ListRegions() returned %d regions, want %d", len(result), len(regions))
	}

	// Check the first region
	if result[0].CloudProvider != regions[0].CloudProvider {
		t.Errorf("ListRegions() first region CloudProvider = %s, want %s", result[0].CloudProvider, regions[0].CloudProvider)
	}
	if result[0].Name != regions[0].Name {
		t.Errorf("ListRegions() first region Name = %s, want %s", result[0].Name, regions[0].Name)
	}
}

func TestListTiers(t *testing.T) {
	// Create a sample response
	tiers := TierInfoList{
		{
			ID:                 21,
			Type:               DeploymentTypeSingleNode,
			CloudProvider:      DeploymentCloudProviderAWS,
			Name:               "s.small.a",
			ComputeCostPerHour: 0.1,
			IngestionRate:      10000,
			ActiveTimeSeries:   100000,
		},
		{
			ID:                 22,
			Type:               DeploymentTypeSingleNode,
			CloudProvider:      DeploymentCloudProviderAWS,
			Name:               "s.medium.a",
			ComputeCostPerHour: 0.2,
			IngestionRate:      20000,
			ActiveTimeSeries:   200000,
		},
	}

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(tiers)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Setup test server
	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/tiers")
	defer server.Close()

	// Call the method
	result, err := client.ListTiers(context.Background())
	if err != nil {
		t.Fatalf("ListTiers() error = %v", err)
	}

	// Check the result
	if len(result) != len(tiers) {
		t.Errorf("ListTiers() returned %d tiers, want %d", len(result), len(tiers))
	}

	// Check the first tier
	if result[0].ID != tiers[0].ID {
		t.Errorf("ListTiers() first tier ID = %d, want %d", result[0].ID, tiers[0].ID)
	}
	if result[0].Type != tiers[0].Type {
		t.Errorf("ListTiers() first tier Type = %s, want %s", result[0].Type, tiers[0].Type)
	}
	if result[0].CloudProvider != tiers[0].CloudProvider {
		t.Errorf("ListTiers() first tier CloudProvider = %s, want %s", result[0].CloudProvider, tiers[0].CloudProvider)
	}
	if result[0].Name != tiers[0].Name {
		t.Errorf("ListTiers() first tier Name = %s, want %s", result[0].Name, tiers[0].Name)
	}
	if result[0].ComputeCostPerHour != tiers[0].ComputeCostPerHour {
		t.Errorf("ListTiers() first tier ComputeCostPerHour = %f, want %f", result[0].ComputeCostPerHour, tiers[0].ComputeCostPerHour)
	}
}
