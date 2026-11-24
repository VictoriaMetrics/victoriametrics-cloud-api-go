package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestListCustomAlertmanagers(t *testing.T) {
	ams := CustomAlertmanagersList{
		{
			ID:        "am-id-1",
			Name:      "External AM",
			URL:       "https://alertmanager.example.com",
			AuthType:  NotifierAuthTypeBasic,
			Username:  "demo",
			Status:    NotifierStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	responseJSON, err := json.Marshal(ams)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/integrations/custom-alertmanagers")
	defer server.Close()

	result, err := client.ListCustomAlertmanagers(context.Background())
	if err != nil {
		t.Fatalf("ListCustomAlertmanagers() error = %v", err)
	}

	if len(result) != len(ams) {
		t.Errorf("ListCustomAlertmanagers() returned %d notifiers, want %d", len(result), len(ams))
	}
	if result[0].ID != ams[0].ID || result[0].Name != ams[0].Name {
		t.Errorf("ListCustomAlertmanagers() = %+v, want %+v", result[0], ams[0])
	}
}

func TestGetCustomAlertmanager(t *testing.T) {
	am := CustomAlertmanager{
		ID:        "am-id-1",
		Name:      "External AM",
		URL:       "https://alertmanager.example.com",
		AuthType:  NotifierAuthTypeBearer,
		Status:    NotifierStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	responseJSON, err := json.Marshal(am)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/integrations/custom-alertmanagers", am.ID)
	defer server.Close()

	result, err := client.GetCustomAlertmanager(context.Background(), am.ID)
	if err != nil {
		t.Fatalf("GetCustomAlertmanager() error = %v", err)
	}
	if result.ID != am.ID || result.URL != am.URL {
		t.Errorf("GetCustomAlertmanager() = %+v, want %+v", result, am)
	}
}

func TestCreateCustomAlertmanager(t *testing.T) {
	request := CustomAlertmanagerRequest{
		Name:     "External AM",
		URL:      "https://alertmanager.example.com",
		AuthType: NotifierAuthTypeBasic,
		Username: "demo",
		Password: "demo-password",
	}

	response := CustomAlertmanager{
		ID:        "am-id-1",
		Name:      request.Name,
		URL:       request.URL,
		AuthType:  request.AuthType,
		Username:  request.Username,
		Status:    NotifierStatusUnknown,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/integrations/custom-alertmanagers")
	defer server.Close()

	result, err := client.CreateCustomAlertmanager(context.Background(), request)
	if err != nil {
		t.Fatalf("CreateCustomAlertmanager() error = %v", err)
	}
	if result.Name != request.Name || result.URL != request.URL {
		t.Errorf("CreateCustomAlertmanager() = %+v, want name/url from %+v", result, request)
	}
}

func TestUpdateCustomAlertmanager(t *testing.T) {
	notifierID := "am-id-1"
	request := CustomAlertmanagerRequest{
		Name:     "Renamed AM",
		URL:      "https://alertmanager.example.com",
		AuthType: NotifierAuthTypeBasic,
	}

	response := CustomAlertmanager{
		ID:        notifierID,
		Name:      request.Name,
		URL:       request.URL,
		AuthType:  request.AuthType,
		Status:    NotifierStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/integrations/custom-alertmanagers", notifierID)
	defer server.Close()

	result, err := client.UpdateCustomAlertmanager(context.Background(), notifierID, request)
	if err != nil {
		t.Fatalf("UpdateCustomAlertmanager() error = %v", err)
	}
	if result.Name != request.Name {
		t.Errorf("UpdateCustomAlertmanager() Name = %s, want %s", result.Name, request.Name)
	}
}

func TestDeleteCustomAlertmanager(t *testing.T) {
	notifierID := "am-id-1"

	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/integrations/custom-alertmanagers", notifierID)
	defer server.Close()

	if err := client.DeleteCustomAlertmanager(context.Background(), notifierID); err != nil {
		t.Fatalf("DeleteCustomAlertmanager() error = %v", err)
	}
}

func TestListCloudAlertmanagers(t *testing.T) {
	ams := CloudAlertmanagersList{
		{
			ID:        "am-id-2",
			Name:      "Hosted AM",
			Config:    "route:\n  receiver: default\n",
			Status:    NotifierStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	responseJSON, err := json.Marshal(ams)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/integrations/cloud-alertmanagers")
	defer server.Close()

	result, err := client.ListCloudAlertmanagers(context.Background())
	if err != nil {
		t.Fatalf("ListCloudAlertmanagers() error = %v", err)
	}
	if len(result) != len(ams) {
		t.Errorf("ListCloudAlertmanagers() returned %d notifiers, want %d", len(result), len(ams))
	}
	if result[0].ID != ams[0].ID || result[0].Config != ams[0].Config {
		t.Errorf("ListCloudAlertmanagers() = %+v, want %+v", result[0], ams[0])
	}
}

func TestGetCloudAlertmanager(t *testing.T) {
	am := CloudAlertmanager{
		ID:        "am-id-2",
		Name:      "Hosted AM",
		Config:    "route:\n  receiver: default\n",
		Status:    NotifierStatusError,
		Message:   "last notification failed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	responseJSON, err := json.Marshal(am)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/integrations/cloud-alertmanagers", am.ID)
	defer server.Close()

	result, err := client.GetCloudAlertmanager(context.Background(), am.ID)
	if err != nil {
		t.Fatalf("GetCloudAlertmanager() error = %v", err)
	}
	if result.ID != am.ID || result.Status != am.Status || result.Message != am.Message {
		t.Errorf("GetCloudAlertmanager() = %+v, want %+v", result, am)
	}
}

func TestCreateCloudAlertmanager(t *testing.T) {
	request := CloudAlertmanagerRequest{
		Name:   "Hosted AM",
		Config: "route:\n  receiver: default\n",
	}

	response := CloudAlertmanager{
		ID:        "am-id-2",
		Name:      request.Name,
		Config:    request.Config,
		Status:    NotifierStatusUnknown,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/integrations/cloud-alertmanagers")
	defer server.Close()

	result, err := client.CreateCloudAlertmanager(context.Background(), request)
	if err != nil {
		t.Fatalf("CreateCloudAlertmanager() error = %v", err)
	}
	if result.Name != request.Name || result.Config != request.Config {
		t.Errorf("CreateCloudAlertmanager() = %+v, want name/config from %+v", result, request)
	}
}

func TestUpdateCloudAlertmanager(t *testing.T) {
	notifierID := "am-id-2"
	request := CloudAlertmanagerRequest{
		Name:   "Renamed hosted AM",
		Config: "route:\n  receiver: default\n",
	}

	response := CloudAlertmanager{
		ID:        notifierID,
		Name:      request.Name,
		Config:    request.Config,
		Status:    NotifierStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/integrations/cloud-alertmanagers", notifierID)
	defer server.Close()

	result, err := client.UpdateCloudAlertmanager(context.Background(), notifierID, request)
	if err != nil {
		t.Fatalf("UpdateCloudAlertmanager() error = %v", err)
	}
	if result.Name != request.Name {
		t.Errorf("UpdateCloudAlertmanager() Name = %s, want %s", result.Name, request.Name)
	}
}

func TestDeleteCloudAlertmanager(t *testing.T) {
	notifierID := "am-id-2"

	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/integrations/cloud-alertmanagers", notifierID)
	defer server.Close()

	if err := client.DeleteCloudAlertmanager(context.Background(), notifierID); err != nil {
		t.Fatalf("DeleteCloudAlertmanager() error = %v", err)
	}
}

func TestNotifierOperations_InvalidNotifierID(t *testing.T) {
	client, err := New("test-api-key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	request := CustomAlertmanagerRequest{Name: "n", URL: "https://example.com", AuthType: NotifierAuthTypeBasic}
	cloudRequest := CloudAlertmanagerRequest{Name: "n", Config: "route: {}"}

	if _, err := client.GetCustomAlertmanager(context.Background(), ""); err == nil {
		t.Errorf("GetCustomAlertmanager() with empty notifier ID should return an error")
	}
	if _, err := client.UpdateCustomAlertmanager(context.Background(), "", request); err == nil {
		t.Errorf("UpdateCustomAlertmanager() with empty notifier ID should return an error")
	}
	if err := client.DeleteCustomAlertmanager(context.Background(), ""); err == nil {
		t.Errorf("DeleteCustomAlertmanager() with empty notifier ID should return an error")
	}

	if _, err := client.GetCloudAlertmanager(context.Background(), ""); err == nil {
		t.Errorf("GetCloudAlertmanager() with empty notifier ID should return an error")
	}
	if _, err := client.UpdateCloudAlertmanager(context.Background(), "", cloudRequest); err == nil {
		t.Errorf("UpdateCloudAlertmanager() with empty notifier ID should return an error")
	}
	if err := client.DeleteCloudAlertmanager(context.Background(), ""); err == nil {
		t.Errorf("DeleteCloudAlertmanager() with empty notifier ID should return an error")
	}
}

func TestGetDeploymentNotifier(t *testing.T) {
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"
	notifier := DeploymentNotifier{
		NotifierID:   "am-id-2",
		NotifierType: NotifierTypeCloudAlertmanager,
		NotifierName: "Hosted AM",
	}

	responseJSON, err := json.Marshal(notifier)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/deployments", deploymentID, "alerting", "notifier")
	defer server.Close()

	result, err := client.GetDeploymentNotifier(context.Background(), deploymentID)
	if err != nil {
		t.Fatalf("GetDeploymentNotifier() error = %v", err)
	}
	if result != notifier {
		t.Errorf("GetDeploymentNotifier() = %+v, want %+v", result, notifier)
	}
}

func TestSetDeploymentNotifier(t *testing.T) {
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"
	request := DeploymentNotifierRequest{
		NotifierID:   "am-id-2",
		NotifierType: NotifierTypeCloudAlertmanager,
	}
	response := DeploymentNotifier{
		NotifierID:   request.NotifierID,
		NotifierType: request.NotifierType,
		NotifierName: "Hosted AM",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	server, client := setupTestServer(t, http.StatusOK, string(responseJSON), "/api/v1/deployments", deploymentID, "alerting", "notifier")
	defer server.Close()

	result, err := client.SetDeploymentNotifier(context.Background(), deploymentID, request)
	if err != nil {
		t.Fatalf("SetDeploymentNotifier() error = %v", err)
	}
	if result != response {
		t.Errorf("SetDeploymentNotifier() = %+v, want %+v", result, response)
	}
}

func TestSetDeploymentNotifier_EmptyNotifierID(t *testing.T) {
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"

	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/deployments", deploymentID, "alerting", "notifier")
	defer server.Close()

	_, err := client.SetDeploymentNotifier(context.Background(), deploymentID, DeploymentNotifierRequest{NotifierType: NotifierTypeCloudAlertmanager})
	if err == nil {
		t.Fatalf("SetDeploymentNotifier() with empty notifier ID should return an error")
	}
}

func TestDeleteDeploymentNotifier(t *testing.T) {
	deploymentID := "123e4567-e89b-12d3-a456-426614174000"

	server, client := setupTestServer(t, http.StatusOK, "", "/api/v1/deployments", deploymentID, "alerting", "notifier")
	defer server.Close()

	if err := client.DeleteDeploymentNotifier(context.Background(), deploymentID); err != nil {
		t.Fatalf("DeleteDeploymentNotifier() error = %v", err)
	}
}

func TestDeploymentNotifierOperations_InvalidDeploymentID(t *testing.T) {
	client, err := New("test-api-key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	invalidID := "invalid-id"
	request := DeploymentNotifierRequest{NotifierID: "am-id-2", NotifierType: NotifierTypeCloudAlertmanager}

	if _, err := client.GetDeploymentNotifier(context.Background(), invalidID); err == nil {
		t.Errorf("GetDeploymentNotifier() with invalid deployment ID should return an error")
	}
	if _, err := client.SetDeploymentNotifier(context.Background(), invalidID, request); err == nil {
		t.Errorf("SetDeploymentNotifier() with invalid deployment ID should return an error")
	}
	if err := client.DeleteDeploymentNotifier(context.Background(), invalidID); err == nil {
		t.Errorf("DeleteDeploymentNotifier() with invalid deployment ID should return an error")
	}
}
