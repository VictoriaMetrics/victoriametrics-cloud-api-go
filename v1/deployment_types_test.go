package v1

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

// TestListWithUnknownDeploymentTypeIsForwarded checks that a list filter this build does
// not recognise still reaches the API. Filtering is the API's decision, so a type added
// after this SDK was released must remain usable rather than being rejected locally.
func TestListWithUnknownDeploymentTypeIsForwarded(t *testing.T) {
	server, client, queries := setupQueryCapturingServer(t, "[]")
	defer server.Close()

	if _, err := client.ListTiers(context.Background(), WithDeploymentType("vlogs_cluster")); err != nil {
		t.Errorf("ListTiers() with an unknown type returned error = %v", err)
	}
	if _, err := client.ListDeployments(context.Background(), WithDeploymentType("vlogs_cluster")); err != nil {
		t.Errorf("ListDeployments() with an unknown type returned error = %v", err)
	}

	want := []string{"type=vlogs_cluster", "type=vlogs_cluster"}
	if len(*queries) != len(want) {
		t.Fatalf("captured queries = %q, want %q", *queries, want)
	}
	for i, got := range *queries {
		if got != want[i] {
			t.Errorf("captured query %d = %q, want %q", i, got, want[i])
		}
	}
}

// TestCreateRejectsUnknownDeploymentType checks that create still refuses a type this
// build cannot construct a valid body for - the check that was removed from the list path
// is a create-time rule, not a general one.
func TestCreateRejectsUnknownDeploymentType(t *testing.T) {
	server, client := setupTestServer(t, http.StatusOK, `{}`, "/api/v1/deployments")
	defer server.Close()

	request := DeploymentCreationRequest{
		Name: "example", Type: "vlogs_cluster", Provider: DeploymentCloudProviderAWS,
		Region: "us-east-1", Tier: 101, StorageSize: 20, StorageSizeUnit: StorageUnitGB,
		Retention: 30, RetentionUnit: DurationUnitDay,
	}
	if _, err := client.CreateDeployment(context.Background(), request); err == nil {
		t.Error("CreateDeployment() with an unknown type returned no error")
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
	if tier.AccessTokenLimit != 10 {
		t.Errorf("tier AccessTokenLimit = %d, want 10", tier.AccessTokenLimit)
	}

	logs := tier.Logs
	if logs == nil {
		t.Fatalf("tier Logs = nil, want the vlogs limits")
	}
	if logs.IngestionRateBytes != 10485760 {
		t.Errorf("tier IngestionRateBytes = %d, want 10485760", logs.IngestionRateBytes)
	}
	if logs.ActiveLogStreams != 10000 {
		t.Errorf("tier ActiveLogStreams = %d, want 10000", logs.ActiveLogStreams)
	}
	if logs.NewStreamsOver24h != 20000 {
		t.Errorf("tier NewStreamsOver24h = %d, want 20000", logs.NewStreamsOver24h)
	}
	if logs.DataReadRate != 10485760 {
		t.Errorf("tier DataReadRate = %d, want 10485760", logs.DataReadRate)
	}
	if logs.BytesPerQuery != 1073741824 {
		t.Errorf("tier BytesPerQuery = %d, want 1073741824", logs.BytesPerQuery)
	}
	// a logs tier reports no metrics or traces limits
	if tier.Metrics != nil || tier.Traces != nil {
		t.Errorf("tier reported metrics/traces limits, want neither")
	}
}

func TestTierInfoUnmarshalVTracesLimits(t *testing.T) {
	const body = `[{
		"id": 201,
		"type": "vtraces_single",
		"cloud_provider": "aws",
		"name": "t.small.a",
		"compute_cost_per_hour": 0.35,
		"ingestion_rate_bytes": 20971520,
		"active_log_streams": 15000,
		"new_streams_over_24h": 30000,
		"data_read_rate": 20971520,
		"bytes_per_query": 2147483648,
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
	if tier.Type != DeploymentTypeVTraces {
		t.Errorf("tier Type = %s, want %s", tier.Type, DeploymentTypeVTraces)
	}

	traces := tier.Traces
	if traces == nil {
		t.Fatalf("tier Traces = nil, want the vtraces limits")
	}
	if traces.IngestionRateBytes != 20971520 {
		t.Errorf("tier IngestionRateBytes = %d, want 20971520", traces.IngestionRateBytes)
	}
	if traces.ActiveLogStreams != 15000 {
		t.Errorf("tier ActiveLogStreams = %d, want 15000", traces.ActiveLogStreams)
	}
	if traces.NewStreamsOver24h != 30000 {
		t.Errorf("tier NewStreamsOver24h = %d, want 30000", traces.NewStreamsOver24h)
	}
	if traces.DataReadRate != 20971520 {
		t.Errorf("tier DataReadRate = %d, want 20971520", traces.DataReadRate)
	}
	if traces.BytesPerQuery != 2147483648 {
		t.Errorf("tier BytesPerQuery = %d, want 2147483648", traces.BytesPerQuery)
	}
	// a traces tier reports no metrics or logs limits
	if tier.Metrics != nil || tier.Logs != nil {
		t.Errorf("tier reported metrics/logs limits, want neither")
	}
}

// TestTierInfoMarshalValue checks that a tier marshals its limits whether it is passed
// by value, by pointer or inside a list. encoding/json reaches a pointer-receiver
// MarshalJSON only for addressable values, so a pointer receiver on TierInfo.MarshalJSON
// would silently drop the limits of a tier passed by value.
func TestTierInfoMarshalValue(t *testing.T) {
	tier := TierInfo{
		TierInfoCommon: TierInfoCommon{
			ID:   21,
			Type: DeploymentTypeSingleNode,
			Name: "s.small.a",
		},
		Metrics: &MetricsTierInfo{IngestionRate: 10000},
	}

	f := func(name string, v any) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			encoded, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			if !strings.Contains(string(encoded), `"ingestion_rate":10000`) {
				t.Errorf("Marshal() = %s, want it to carry the metrics limits", encoded)
			}
		})
	}

	f("value", tier)
	f("pointer", &tier)
	f("list", TierInfoList{tier})
}

// TestTierInfoRoundTrip checks that a tier survives a marshal/unmarshal cycle with
// its limits landing back in the struct its type calls for.
func TestTierInfoRoundTrip(t *testing.T) {
	f := func(name string, tier TierInfo) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			encoded, err := json.Marshal(tier)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			var decoded TierInfo
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if !reflect.DeepEqual(tier, decoded) {
				t.Errorf("round trip returned %+v, want %+v", decoded, tier)
			}
		})
	}

	f("metrics", TierInfo{
		TierInfoCommon: TierInfoCommon{
			ID:                            21,
			Type:                          DeploymentTypeSingleNode,
			CloudProvider:                 DeploymentCloudProviderAWS,
			Name:                          "s.small.a",
			ComputeCostPerHour:            0.1,
			AccessTokenConcurrentRequests: 5,
			AccessTokenLimit:              10,
		},
		Metrics: &MetricsTierInfo{
			IngestionRate:      10000,
			ActiveTimeSeries:   100000,
			NewSeriesOver24h:   200000,
			DatapointsReadRate: 30000,
			SeriesReadPerQuery: 1500,
		},
	})
	f("logs", TierInfo{
		TierInfoCommon: TierInfoCommon{
			ID:                 101,
			Type:               DeploymentTypeVLogs,
			CloudProvider:      DeploymentCloudProviderAWS,
			Name:               "l.small.a",
			ComputeCostPerHour: 0.25,
		},
		Logs: &LogsTierInfo{
			IngestionRateBytes: 10485760,
			ActiveLogStreams:   10000,
			NewStreamsOver24h:  20000,
			DataReadRate:       10485760,
			BytesPerQuery:      1073741824,
		},
	})
	f("traces", TierInfo{
		TierInfoCommon: TierInfoCommon{
			ID:                 201,
			Type:               DeploymentTypeVTraces,
			CloudProvider:      DeploymentCloudProviderAWS,
			Name:               "t.small.a",
			ComputeCostPerHour: 0.35,
		},
		Traces: &TracesTierInfo{
			IngestionRateBytes: 20971520,
			ActiveLogStreams:   15000,
			NewStreamsOver24h:  30000,
			DataReadRate:       20971520,
			BytesPerQuery:      2147483648,
		},
	})
}

func TestTierInfoMarshalFollowsType(t *testing.T) {
	tier := TierInfo{
		TierInfoCommon: TierInfoCommon{
			ID:                 201,
			Type:               DeploymentTypeVTraces,
			CloudProvider:      DeploymentCloudProviderAWS,
			Name:               "t.small.a",
			ComputeCostPerHour: 0.35,
		},
		Traces: &TracesTierInfo{
			IngestionRateBytes: 20971520,
			ActiveLogStreams:   15000,
			NewStreamsOver24h:  30000,
			DataReadRate:       20971520,
			BytesPerQuery:      1235689,
		},
		// neither of these matches Type, so neither may reach the wire
		Metrics: &MetricsTierInfo{
			IngestionRate:      10000,
			ActiveTimeSeries:   2000000,
			NewSeriesOver24h:   100000000,
			DatapointsReadRate: 345678,
			SeriesReadPerQuery: 5678934,
		},
		Logs: &LogsTierInfo{
			IngestionRateBytes: 230000000,
			ActiveLogStreams:   123214,
			NewStreamsOver24h:  800,
			DataReadRate:       4504945868,
			BytesPerQuery:      1235689,
		},
	}

	encoded, err := json.Marshal(tier)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if got := string(encoded); strings.Contains(got, "ingestion_rate\"") || strings.Contains(got, "active_time_series") {
		t.Errorf("Marshal() = %s, want no metrics limits on a vtraces tier", got)
	}

	var decoded TierInfo
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if decoded.Metrics != nil || decoded.Logs != nil {
		t.Errorf("decoded tier carries metrics/logs limits, want neither")
	}
	if decoded.Traces == nil || *decoded.Traces != *tier.Traces {
		t.Errorf("decoded Traces = %+v, want %+v", decoded.Traces, tier.Traces)
	}
}

func TestListTiersAllDeploymentTypes(t *testing.T) {
	const body = `[
		{
			"id": 21, "type": "single_node", "cloud_provider": "aws", "name": "s.small.a",
			"compute_cost_per_hour": 0.1,
			"ingestion_rate": 10000, "active_time_series": 100000, "new_series_over_24h": 200000,
			"datapoints_read_rate": 30000, "series_read_per_query": 1500,
			"access_token_concurrent_requests": 5, "access_token_limit": 10
		},
		{
			"id": 41, "type": "cluster", "cloud_provider": "aws", "name": "c.large.a",
			"compute_cost_per_hour": 1.5,
			"ingestion_rate": 500000, "active_time_series": 5000000, "new_series_over_24h": 1000000,
			"datapoints_read_rate": 900000, "series_read_per_query": 50000,
			"access_token_concurrent_requests": 20, "access_token_limit": 40
		},
		{
			"id": 101, "type": "vlogs_single", "cloud_provider": "aws", "name": "l.small.a",
			"compute_cost_per_hour": 0.25,
			"ingestion_rate_bytes": 10485760, "active_log_streams": 10000, "new_streams_over_24h": 20000,
			"data_read_rate": 10485760, "bytes_per_query": 1073741824,
			"access_token_concurrent_requests": 5, "access_token_limit": 10
		},
		{
			"id": 201, "type": "vtraces_single", "cloud_provider": "aws", "name": "t.small.a",
			"compute_cost_per_hour": 0.35,
			"ingestion_rate_bytes": 20971520, "active_log_streams": 15000, "new_streams_over_24h": 30000,
			"data_read_rate": 20971520, "bytes_per_query": 2147483648,
			"access_token_concurrent_requests": 5, "access_token_limit": 10
		}
	]`

	server, client := setupTestServer(t, http.StatusOK, body, "/api/v1/tiers")
	defer server.Close()

	result, err := client.ListTiers(context.Background())
	if err != nil {
		t.Fatalf("ListTiers() error = %v", err)
	}
	if len(result) != 4 {
		t.Fatalf("ListTiers() returned %d tiers, want 4", len(result))
	}

	// the response order is preserved, and every tier keeps its common fields
	wantTypes := []DeploymentType{
		DeploymentTypeSingleNode, DeploymentTypeCluster, DeploymentTypeVLogs, DeploymentTypeVTraces,
	}
	wantIDs := []uint32{21, 41, 101, 201}
	for i, tier := range result {
		if tier.Type != wantTypes[i] {
			t.Errorf("tier %d Type = %s, want %s", i, tier.Type, wantTypes[i])
		}
		if tier.ID != wantIDs[i] {
			t.Errorf("tier %d ID = %d, want %d", i, tier.ID, wantIDs[i])
		}
		if tier.CloudProvider != DeploymentCloudProviderAWS {
			t.Errorf("tier %d CloudProvider = %s, want %s", i, tier.CloudProvider, DeploymentCloudProviderAWS)
		}
		if tier.AccessTokenLimit == 0 {
			t.Errorf("tier %d AccessTokenLimit = 0, want it decoded", i)
		}
	}

	// each tier carries exactly the limits of its own type
	singleNode, cluster, logs, traces := result[0], result[1], result[2], result[3]

	if singleNode.Metrics == nil || singleNode.Logs != nil || singleNode.Traces != nil {
		t.Fatalf("single_node tier limits = %+v/%+v/%+v, want metrics only",
			singleNode.Metrics, singleNode.Logs, singleNode.Traces)
	}
	if singleNode.Metrics.IngestionRate != 10000 || singleNode.Metrics.ActiveTimeSeries != 100000 {
		t.Errorf("single_node limits = %+v, want IngestionRate 10000 and ActiveTimeSeries 100000", singleNode.Metrics)
	}

	if cluster.Metrics == nil || cluster.Logs != nil || cluster.Traces != nil {
		t.Fatalf("cluster tier limits = %+v/%+v/%+v, want metrics only",
			cluster.Metrics, cluster.Logs, cluster.Traces)
	}
	if cluster.Metrics.SeriesReadPerQuery != 50000 {
		t.Errorf("cluster SeriesReadPerQuery = %d, want 50000", cluster.Metrics.SeriesReadPerQuery)
	}

	if logs.Logs == nil || logs.Metrics != nil || logs.Traces != nil {
		t.Fatalf("vlogs tier limits = %+v/%+v/%+v, want logs only", logs.Metrics, logs.Logs, logs.Traces)
	}
	if logs.Logs.ActiveLogStreams != 10000 || logs.Logs.BytesPerQuery != 1073741824 {
		t.Errorf("vlogs limits = %+v, want ActiveLogStreams 10000 and BytesPerQuery 1073741824", logs.Logs)
	}

	if traces.Traces == nil || traces.Metrics != nil || traces.Logs != nil {
		t.Fatalf("vtraces tier limits = %+v/%+v/%+v, want traces only", traces.Metrics, traces.Logs, traces.Traces)
	}
	if traces.Traces.ActiveLogStreams != 15000 || traces.Traces.BytesPerQuery != 2147483648 {
		t.Errorf("vtraces limits = %+v, want ActiveLogStreams 15000 and BytesPerQuery 2147483648", traces.Traces)
	}
}

func TestTierInfoUnknownType(t *testing.T) {
	const body = `{
		"id": 9,
		"type": "vtraces_cluster",
		"cloud_provider": "aws",
		"name": "t.big.a",
		"ingestion_rate": 5
	}`

	var tier TierInfo
	if err := json.Unmarshal([]byte(body), &tier); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if tier.ID != 9 || tier.Name != "t.big.a" {
		t.Errorf("Unmarshal() = %+v, want the common fields decoded", tier.TierInfoCommon)
	}
	if tier.Metrics != nil || tier.Logs != nil || tier.Traces != nil {
		t.Errorf("unknown type decoded limits %+v/%+v/%+v, want all nil", tier.Metrics, tier.Logs, tier.Traces)
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
	// a server that spells the absent window out rather than omitting it must not be read
	// as a deployment that has a zero-second window
	f("VictoriaLogs deployment with the window spelled out",
		`{"type":"vlogs_single","deduplication_value":0,"deduplication_unit":""}`, 0, "", false)
	f("VictoriaTraces deployment with the window spelled out",
		`{"type":"vtraces_single","deduplication_value":0,"deduplication_unit":"s"}`, 0, "", false)
	// an absent or unknown type falls back to the fields, so a newer type still reads
	f("deployment without a type", `{"deduplication_value":30,"deduplication_unit":"s"}`, 30, DurationUnitSecond, true)
	f("deployment of an unknown type", `{"type":"vlogs_cluster","deduplication_value":30,"deduplication_unit":"s"}`, 30, DurationUnitSecond, true)
	// a metrics deployment missing half the window is not "no window", but there is no
	// third state to report, so it stays false
	f("metrics deployment missing the unit", `{"type":"single_node","deduplication_value":30}`, 0, "", false)
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
		Deduplication:     new(uint32(10)),
		DeduplicationUnit: new(DurationUnitSecond),
		MaintenanceWindow: MaintenanceWindowWeekendDays,
	}

	server, client := setupTestServer(t, http.StatusOK, `{"id":"123e4567-e89b-12d3-a456-426614174000","type":"single_node"}`, "/api/v1/deployments")
	defer server.Close()

	if _, err := client.CreateDeployment(context.Background(), request); err != nil {
		t.Fatalf("CreateDeployment() error = %v", err)
	}
}

// TestUpdateVLogsDeploymentWithoutDeduplication covers the reason the update request
// carries pointers: an update body has no deployment type, so leaving both deduplication
// fields nil is what lets a VictoriaLogs or VictoriaTraces deployment be updated without
// a window, and it keeps both fields out of the body entirely.
func TestUpdateVLogsDeploymentWithoutDeduplication(t *testing.T) {
	const deploymentID = "123e4567-e89b-12d3-a456-426614174000"

	var body []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"` + deploymentID + `","type":"vlogs_single"}`))
	}))
	defer server.Close()

	client, err := New("test-api-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	request := DeploymentUpdateRequest{
		Name:              "example-logs",
		Tier:              101,
		StorageSize:       20,
		StorageSizeUnit:   StorageUnitGB,
		Retention:         30,
		RetentionUnit:     DurationUnitDay,
		MaintenanceWindow: MaintenanceWindowWeekendDays,
	}
	if _, err := client.UpdateDeployment(context.Background(), deploymentID, request); err != nil {
		t.Fatalf("UpdateDeployment() error = %v", err)
	}

	var sent map[string]json.RawMessage
	if err := json.Unmarshal(body, &sent); err != nil {
		t.Fatalf("Failed to unmarshal the request body: %v", err)
	}
	if _, ok := sent["deduplication"]; ok {
		t.Errorf("update body = %s, want no deduplication field", body)
	}
	if _, ok := sent["deduplication_unit"]; ok {
		t.Errorf("update body = %s, want no deduplication_unit field", body)
	}
}

// TestUpdateSendsZeroDeduplicationWindow pins the case a plain uint32 could not express:
// a metrics deployment whose deduplication window is set to zero. omitempty drops a nil
// pointer but keeps a pointer to zero, so the value reaches the API instead of being
// silently replaced by whatever default the API applies to a missing field.
func TestUpdateSendsZeroDeduplicationWindow(t *testing.T) {
	const deploymentID = "123e4567-e89b-12d3-a456-426614174000"

	var body []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"` + deploymentID + `","type":"single_node"}`))
	}))
	defer server.Close()

	client, err := New("test-api-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	request := DeploymentUpdateRequest{
		Name:              "example-metrics",
		Tier:              21,
		StorageSize:       20,
		StorageSizeUnit:   StorageUnitGB,
		Retention:         30,
		RetentionUnit:     DurationUnitDay,
		Deduplication:     new(uint32(0)),
		DeduplicationUnit: new(DurationUnitSecond),
		MaintenanceWindow: MaintenanceWindowWeekendDays,
	}
	if _, err := client.UpdateDeployment(context.Background(), deploymentID, request); err != nil {
		t.Fatalf("UpdateDeployment() error = %v", err)
	}

	var sent map[string]json.RawMessage
	if err := json.Unmarshal(body, &sent); err != nil {
		t.Fatalf("Failed to unmarshal the request body: %v", err)
	}
	if got, ok := sent["deduplication"]; !ok || string(got) != "0" {
		t.Errorf("update body = %s, want a deduplication field of 0", body)
	}
}

// TestUpdateWithHalfSetDeduplicationIsRejected checks the one combination the pointers
// make expressible but that is always a mistake: a window without a unit, or the reverse.
func TestUpdateWithHalfSetDeduplicationIsRejected(t *testing.T) {
	client, err := New("test-api-key", WithBaseURL("https://api.victoriametrics.cloud"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	base := DeploymentUpdateRequest{
		Name:              "example-metrics",
		Tier:              21,
		StorageSize:       20,
		StorageSizeUnit:   StorageUnitGB,
		Retention:         30,
		RetentionUnit:     DurationUnitDay,
		MaintenanceWindow: MaintenanceWindowWeekendDays,
	}

	windowOnly := base
	windowOnly.Deduplication = new(uint32(30))
	if _, err := client.UpdateDeployment(context.Background(), "123e4567-e89b-12d3-a456-426614174000", windowOnly); err == nil {
		t.Error("UpdateDeployment() with a window and no unit returned no error")
	}

	unitOnly := base
	unitOnly.DeduplicationUnit = new(DurationUnitSecond)
	if _, err := client.UpdateDeployment(context.Background(), "123e4567-e89b-12d3-a456-426614174000", unitOnly); err == nil {
		t.Error("UpdateDeployment() with a unit and no window returned no error")
	}
}
