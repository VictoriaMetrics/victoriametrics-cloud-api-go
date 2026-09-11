package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupQueryCapturingServer serves response and records the raw query of every request,
// so that a test can assert on the filters the client sent.
func setupQueryCapturingServer(t *testing.T, response string) (*httptest.Server, *VMCloudAPIClient, *[]string) {
	t.Helper()
	var queries []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.RawQuery)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	client, err := New("test-api-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	return server, client, &queries
}

func TestDeploymentTypeCapabilities(t *testing.T) {
	f := func(dType DeploymentType, wantDeduplication, wantRules bool) {
		t.Helper()
		if got := dType.SupportsDeduplication(); got != wantDeduplication {
			t.Errorf("%s.SupportsDeduplication() = %v, want %v", dType, got, wantDeduplication)
		}
		if got := dType.SupportsAlertingRules(); got != wantRules {
			t.Errorf("%s.SupportsAlertingRules() = %v, want %v", dType, got, wantRules)
		}
	}
	f(DeploymentTypeSingleNode, true, true)
	f(DeploymentTypeCluster, true, true)
	f(DeploymentTypeVLogs, false, false)
	f(DeploymentTypeVTraces, false, false)
}

func TestListDeploymentsWithDeploymentType(t *testing.T) {
	server, client, queries := setupQueryCapturingServer(t, "[]")
	defer server.Close()

	if _, err := client.ListDeployments(context.Background()); err != nil {
		t.Fatalf("ListDeployments() error = %v", err)
	}
	if _, err := client.ListDeployments(context.Background(), WithDeploymentType(DeploymentTypeVTraces)); err != nil {
		t.Fatalf("ListDeployments(WithDeploymentType) error = %v", err)
	}

	want := []string{"", "type=vtraces_single"}
	if len(*queries) != len(want) {
		t.Fatalf("captured %d requests, want %d", len(*queries), len(want))
	}
	for i, q := range want {
		if (*queries)[i] != q {
			t.Errorf("request %d query = %q, want %q", i, (*queries)[i], q)
		}
	}
}

func TestListTiersWithDeploymentType(t *testing.T) {
	server, client, queries := setupQueryCapturingServer(t, "[]")
	defer server.Close()

	if _, err := client.ListTiers(context.Background(), WithDeploymentType(DeploymentTypeVLogs)); err != nil {
		t.Fatalf("ListTiers(WithDeploymentType) error = %v", err)
	}
	if len(*queries) != 1 || (*queries)[0] != "type=vlogs_single" {
		t.Errorf("captured queries = %q, want [type=vlogs_single]", *queries)
	}
}

func TestListWithUnknownDeploymentTypeIsRejected(t *testing.T) {
	server, client, queries := setupQueryCapturingServer(t, "[]")
	defer server.Close()

	if _, err := client.ListTiers(context.Background(), WithDeploymentType("vmetrics")); err == nil {
		t.Error("ListTiers() with an unknown type returned no error")
	}
	if _, err := client.ListDeployments(context.Background(), WithDeploymentType("vmetrics")); err == nil {
		t.Error("ListDeployments() with an unknown type returned no error")
	}
	if len(*queries) != 0 {
		t.Errorf("an unknown type reached the API as %q, want no request at all", *queries)
	}
}

func TestTierInfoUnmarshalVLogsLimits(t *testing.T) {
	const body = `[{
		"id": 101,
		"type": "vlogs_single",
		"cloud_provider": "aws",
		"name": "l.small.a",
		"compute_cost_per_hour": 0.25,
		"ingestion_rate_bytes": 10485760,
		"active_log_streams": 10000,
		"new_streams_over_24h": 20000,
		"data_read_rate": 10485760,
		"bytes_per_query": 1073741824,
		"access_token_concurrent_requests": 5,
		"access_token_limit": 10
	}]`

	server, client := setupTestServer(t, http.StatusOK, body, "/api/v1/tiers")
	defer server.Close()

	result, err := client.ListTiers(context.Background())
	if err != nil {
		t.Fatalf("ListTiers() error = %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("ListTiers() returned %d tiers, want 1", len(result))
	}
	tier := result[0]
	if tier.Type != DeploymentTypeVLogs {
		t.Errorf("tier Type = %s, want %s", tier.Type, DeploymentTypeVLogs)
	}
	if tier.IngestionRateBytes != 10485760 {
		t.Errorf("tier IngestionRateBytes = %d, want 10485760", tier.IngestionRateBytes)
	}
	if tier.ActiveLogStreams != 10000 {
		t.Errorf("tier ActiveLogStreams = %d, want 10000", tier.ActiveLogStreams)
	}
	if tier.NewStreamsOver24h != 20000 {
		t.Errorf("tier NewStreamsOver24h = %d, want 20000", tier.NewStreamsOver24h)
	}
	if tier.DataReadRate != 10485760 {
		t.Errorf("tier DataReadRate = %d, want 10485760", tier.DataReadRate)
	}
	if tier.BytesPerQuery != 1073741824 {
		t.Errorf("tier BytesPerQuery = %d, want 1073741824", tier.BytesPerQuery)
	}
	if tier.AccessTokenLimit != 10 {
		t.Errorf("tier AccessTokenLimit = %d, want 10", tier.AccessTokenLimit)
	}
	// the metrics limits are absent from a logs tier
	if tier.IngestionRate != 0 || tier.ActiveTimeSeries != 0 {
		t.Errorf("tier reported metrics limits %d/%d, want 0/0", tier.IngestionRate, tier.ActiveTimeSeries)
	}
}

func TestDeploymentInfoDeduplication(t *testing.T) {
	f := func(name, body string, wantValue uint32, wantUnit DurationUnit, wantOK bool) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			var info DeploymentInfo
			if err := json.Unmarshal([]byte(body), &info); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			value, unit, ok := info.Deduplication()
			if ok != wantOK {
				t.Fatalf("Deduplication() ok = %v, want %v", ok, wantOK)
			}
			if value != wantValue || unit != wantUnit {
				t.Errorf("Deduplication() = (%d, %s), want (%d, %s)", value, unit, wantValue, wantUnit)
			}
		})
	}
	f("metrics deployment", `{"type":"single_node","deduplication_value":30,"deduplication_unit":"s"}`, 30, DurationUnitSecond, true)
	// a zero window is a valid setting, and stays distinct from having no window at all
	f("metrics deployment with a zero window", `{"type":"cluster","deduplication_value":0,"deduplication_unit":"ms"}`, 0, DurationUnitMillisecond, true)
	f("VictoriaLogs deployment", `{"type":"vlogs_single"}`, 0, "", false)
	f("VictoriaTraces deployment", `{"type":"vtraces_single"}`, 0, "", false)
}

func TestCreateVLogsDeploymentOmitsDeduplication(t *testing.T) {
	request := DeploymentCreationRequest{
		Name:              "example-logs",
		Type:              DeploymentTypeVLogs,
		Provider:          DeploymentCloudProviderAWS,
		Region:            "us-east-1",
		Tier:              101,
		StorageSize:       20,
		StorageSizeUnit:   StorageUnitGB,
		Retention:         30,
		RetentionUnit:     DurationUnitDay,
		MaintenanceWindow: MaintenanceWindowWeekendDays,
	}

	server, client := setupTestServer(t, http.StatusOK, `{"id":"123e4567-e89b-12d3-a456-426614174000","type":"vlogs_single"}`, "/api/v1/deployments")
	defer server.Close()

	result, err := client.CreateDeployment(context.Background(), request)
	if err != nil {
		t.Fatalf("CreateDeployment() error = %v", err)
	}
	if result.Type != DeploymentTypeVLogs {
		t.Errorf("CreateDeployment() Type = %s, want %s", result.Type, DeploymentTypeVLogs)
	}
	if _, _, ok := result.Deduplication(); ok {
		t.Error("CreateDeployment() reported a deduplication window for a VictoriaLogs deployment")
	}
}

func TestCreateSingleNodeDeploymentRequiresDeduplicationUnit(t *testing.T) {
	request := DeploymentCreationRequest{
		Name:              "example-metrics",
		Type:              DeploymentTypeSingleNode,
		Provider:          DeploymentCloudProviderAWS,
		Region:            "us-east-1",
		Tier:              21,
		StorageSize:       20,
		StorageSizeUnit:   StorageUnitGB,
		Retention:         30,
		RetentionUnit:     DurationUnitDay,
		MaintenanceWindow: MaintenanceWindowWeekendDays,
	}

	client, err := New("test-api-key", WithBaseURL("https://api.victoriametrics.cloud"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	if _, err := client.CreateDeployment(context.Background(), request); err == nil {
		t.Error("CreateDeployment() without a deduplication unit returned no error")
	}
}

func TestCreateDeploymentAllowsLargeStorage(t *testing.T) {
	// the valid storage range and step depend on the storage type, the topology and the
	// installation config, so the client no longer second-guesses the API on them
	request := DeploymentCreationRequest{
		Name:              "example-metrics",
		Type:              DeploymentTypeSingleNode,
		Provider:          DeploymentCloudProviderAWS,
		Region:            "us-east-1",
		Tier:              21,
		StorageSize:       32,
		StorageSizeUnit:   StorageUnitTB,
		Retention:         30,
		RetentionUnit:     DurationUnitDay,
		Deduplication:     10,
		DeduplicationUnit: DurationUnitSecond,
		MaintenanceWindow: MaintenanceWindowWeekendDays,
	}

	server, client := setupTestServer(t, http.StatusOK, `{"id":"123e4567-e89b-12d3-a456-426614174000","type":"single_node"}`, "/api/v1/deployments")
	defer server.Close()

	if _, err := client.CreateDeployment(context.Background(), request); err != nil {
		t.Fatalf("CreateDeployment() error = %v", err)
	}
}
