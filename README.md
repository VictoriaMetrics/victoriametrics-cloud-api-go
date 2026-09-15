# Client library for VictoriaMetrics Cloud API

![Latest Release](https://img.shields.io/github/v/release/VictoriaMetrics/victoriametrics-cloud-api-go?sort=semver&label=&logo=github&labelColor=gray&color=gray&link=https%3A%2F%2Fgithub.com%2FVictoriaMetrics%2Fvictoriametrics-cloud-api-go%2Freleases%2Flatest)
[![Go Reference](https://pkg.go.dev/badge/github.com/VictoriaMetrics/victoriametrics-cloud-api-go.svg)](https://pkg.go.dev/github.com/VictoriaMetrics/victoriametrics-cloud-api-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
![Slack](https://img.shields.io/badge/Join-4A154B?logo=slack&link=https%3A%2F%2Fslack.victoriametrics.com)
![X](https://img.shields.io/twitter/follow/VictoriaMetrics?style=flat&label=Follow&color=black&logo=x&labelColor=black&link=https%3A%2F%2Fx.com%2FVictoriaMetrics)
![Reddit](https://img.shields.io/reddit/subreddit-subscribers/VictoriaMetrics?style=flat&label=Join&labelColor=red&logoColor=white&logo=reddit&link=https%3A%2F%2Fwww.reddit.com%2Fr%2FVictoriaMetrics)

Go client library for interacting with the [VictoriaMetrics Cloud](https://victoriametrics.com/products/cloud/) API. 
This library provides a simple and idiomatic way to manage VictoriaMetrics Cloud resources programmatically:

- Full library documentation you can find in [package docs](https://pkg.go.dev/github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1).
- More info about VictoriaMetrics Cloud can be found in the [official documentation](https://docs.victoriametrics.com/victoriametrics-cloud/).
- More information about the API can be found in the [API documentation](https://docs.victoriametrics.com/victoriametrics-cloud/api/).

Just sign up for a [free trial](https://victoriametrics.cloud) to get started with VictoriaMetrics Cloud.

## Features

- Manage deployments (list, create, update, delete, get details)
- Manage access tokens for deployments (list, create, delete, reveal secret, revoke)
- Manage alerting/recording rule files for deployments (list, create, update, delete, get content)
- Retrieve information about cloud providers, regions and tiers

### Supported deployment types

| Type | Constant | Deduplication | Alerting/recording rules |
|---|---|---|---|
| VictoriaMetrics single-node | `DeploymentTypeSingleNode` | yes | yes |
| VictoriaMetrics cluster | `DeploymentTypeCluster` | yes | yes |
| VictoriaLogs | `DeploymentTypeVLogs` | no | no |
| VictoriaTraces | `DeploymentTypeVTraces` | no | no |

VictoriaTraces deployments are available only for accounts where VictoriaTraces is
enabled; for other accounts the API hides its tiers and rejects its deployments.

## Installation

```bash
go get github.com/VictoriaMetrics/victoriametrics-cloud-api-go
```

## Usage examples

For detailed examples, see the [examples](examples) directory:

- [Client initialization](examples/01_client_init/01_client_init.go) - Different ways to initialize the client
- [Listing cloud providers, regions, and tiers](examples/02_providers_regions_tiers/02_providers_regions_tiers.go) - How to retrieve information about available cloud providers, regions, and tiers
- [Deployments management](examples/03_deployments_management/03_deployments_management.go) - How to list, create, update, and delete deployments
- [Access tokens management](examples/04_access_tokens_management/04_access_tokens_management.go) - How to list, create, reveal, and delete access tokens
- [Rule files management](examples/05_rule_files_management/05_rule_files_management.go) - How to list, create, update, and delete alerting/recording rule files

### Creating a client

```go
package main

import (
	"log"

	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"
)

func main() {
	// Create a new client with your API key
	client, err := vmcloud.New("your-api-key")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Use the client to interact with the API
	// ...
}
```

### Listing deployments

```go
// Deployments of every type are returned unless the request is narrowed down.
deployments, err := client.ListDeployments(context.Background())
if err != nil {
	log.Fatalf("Failed to list deployments: %v", err)
}

for _, deployment := range deployments {
	fmt.Printf("Deployment: %s (ID: %s)\n", deployment.Name, deployment.ID)
}

// Only VictoriaLogs deployments. The same option narrows down ListTiers.
vlogsDeployments, err := client.ListDeployments(context.Background(),
	vmcloud.WithDeploymentType(vmcloud.DeploymentTypeVLogs))
if err != nil {
	log.Fatalf("Failed to list VictoriaLogs deployments: %v", err)
}
```

### Creating a deployment

```go
deployment := vmcloud.DeploymentCreationRequest{
	Name:              "my-deployment", // name of the deployment
	Type:              vmcloud.DeploymentTypeSingleNode, // single-node deployment
	Provider:          vmcloud.DeploymentCloudProviderAWS, // AWS as cloud provider
	Region:            "us-east-2", // US East (Ohio)
	Tier:              21, // s.starter.a
	StorageSize:       20, // storage size in GB
	StorageSizeUnit:   vmcloud.StorageUnitGB,
	Retention:         30, // data retention period in days
	RetentionUnit:     vmcloud.DurationUnitDay, 
	Deduplication:     new(uint32(10)), // deduplication period in seconds
	DeduplicationUnit: new(vmcloud.DurationUnitSecond),
	MaintenanceWindow: vmcloud.MaintenanceWindowWeekendDays, // maintenance window on weekends
}

createdDeployment, err := client.CreateDeployment(context.Background(), deployment)
if err != nil {
	log.Fatalf("Failed to create deployment: %v", err)
}

fmt.Printf("Created deployment: %s (ID: %s)\n", createdDeployment.Name, createdDeployment.ID)
```

VictoriaLogs and VictoriaTraces deployments are created the same way, with `Type` set to
`vmcloud.DeploymentTypeVLogs` or `vmcloud.DeploymentTypeVTraces`. They have no
deduplication window, so `Deduplication` and `DeduplicationUnit` are left nil. Setting
either of them on those types is rejected by the client rather than silently dropped:

```go
deployment := vmcloud.DeploymentCreationRequest{
	Name:              "my-logs",
	Type:              vmcloud.DeploymentTypeVLogs,
	Provider:          vmcloud.DeploymentCloudProviderAWS,
	Region:            "us-east-2",
	Tier:              101,
	StorageSize:       20,
	StorageSizeUnit:   vmcloud.StorageUnitGB,
	Retention:         30,
	RetentionUnit:     vmcloud.DurationUnitDay,
	MaintenanceWindow: vmcloud.MaintenanceWindowWeekendDays,
}
```

Storage size is validated by the API, not by this library: the valid range and the step a
size must sit on depend on the storage type, the deployment topology and the
per-installation configuration, none of which the client can see.

### Reading the deduplication window

`DeploymentInfo.DeduplicationValue` and `DeploymentInfo.DeduplicationUnit` are absent for
VictoriaLogs and VictoriaTraces deployments, which have no deduplication window. Use the
`Deduplication` accessor rather than dereferencing them:

```go
if value, unit, ok := deploymentDetails.Deduplication(); ok {
	fmt.Printf("Deduplication: %d %s\n", value, unit)
}
```

### Updating a deployment

`DeploymentUpdateRequest.Deduplication` and `.DeduplicationUnit` are pointers. An update
body carries no deployment type, so nil is what states the intent: it leaves the
deployment's current window untouched.

```go
// a metrics deployment, changing the tier but not the deduplication window
updateRequest := vmcloud.DeploymentUpdateRequest{
	Name:            "my-deployment",
	Tier:            22,
	StorageSize:     30,
	StorageSizeUnit: vmcloud.StorageUnitGB,
	Retention:       60,
	RetentionUnit:   vmcloud.DurationUnitDay,
	// Deduplication and DeduplicationUnit left nil: the window is not changed
	MaintenanceWindow: vmcloud.MaintenanceWindowBusinessDays,
}

// changing the window, zero included
updateRequest.Deduplication = new(uint32(0))
updateRequest.DeduplicationUnit = new(vmcloud.DurationUnitSecond)
```

VictoriaLogs and VictoriaTraces deployments have no deduplication window, so their updates
always leave both fields nil. Set the two together - a window without a unit, or a unit
without a window, is rejected before the request is sent.

### Managing access tokens

```go
// Create a new access token for a deployment
tokenRequest := vmcloud.AccessTokenCreateRequest{
	Description: "My API token",
	Type:        vmcloud.AccessModeReadWrite,
}

createdToken, err := client.CreateDeploymentAccessToken(context.Background(), "deployment-id", tokenRequest)
if err != nil {
	log.Fatalf("Failed to create access token: %v", err)
}

fmt.Printf("Created token: %s (ID: %s)\n", createdToken.Description, createdToken.ID)

// List access tokens for a deployment
tokens, err := client.ListDeploymentAccessTokens(context.Background(), "deployment-id")
if err != nil {
	log.Fatalf("Failed to list access tokens: %v", err)
}

for _, token := range tokens {
	fmt.Printf("Token: %s (ID: %s)\n", token.Description, token.ID)
}
```

### Managing alerting/recording rules

Only VictoriaMetrics single-node and cluster deployments run `vmalert`. The rule file
endpoints return an error for VictoriaLogs and VictoriaTraces deployments, which
`DeploymentType.SupportsAlertingRules` reports up front.

```go
// Create a new rule file
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

err := client.CreateDeploymentRuleFileContent(context.Background(), "deployment-id", "high-latency-alert.yml", ruleContent)
if err != nil {
	log.Fatalf("Failed to create rule file: %v", err)
}

// List rule files
ruleFiles, err := client.ListDeploymentRuleFileNames(context.Background(), "deployment-id")
if err != nil {
	log.Fatalf("Failed to list rule files: %v", err)
}

for _, fileName := range ruleFiles {
	fmt.Printf("Rule file: %s\n", fileName)
}
```

## Documentation

For more information about the VictoriaMetrics Cloud API, please refer to the [VictoriaMetrics Cloud documentation](https://docs.victoriametrics.com/victoriametrics-cloud/api/).

## Testing

The library includes a comprehensive test suite. To run the tests:

```bash
make test
```

The tests use mocked HTTP responses and don't require actual API credentials.

## Upgrading

### To v0.2.0

`DeploymentInfo.DeduplicationValue` and `DeploymentInfo.DeduplicationUnit` changed from
`uint32` and `DurationUnit` to `*uint32` and `*DurationUnit`, following the API, which
now omits both for deployment types that have no deduplication window. Read them through
the `Deduplication` accessor:

```go
// before
fmt.Printf("Deduplication: %d %s\n", info.DeduplicationValue, info.DeduplicationUnit)

// after
if value, unit, ok := info.Deduplication(); ok {
	fmt.Printf("Deduplication: %d %s\n", value, unit)
}
```

`ListDeployments` and `ListTiers` return entries of every deployment type, where they
previously returned only VictoriaMetrics single-node and cluster ones. Pass
`vmcloud.WithDeploymentType(...)` to keep the previous result of either call.

`TierInfo` no longer carries its limits as flat fields. The common fields moved into an
embedded `TierInfoCommon`, and the limits into a struct chosen by the tier's type -
`Metrics` for `single_node` and `cluster`, `Logs` for `vlogs_single`, `Traces` for
`vtraces_single`, with the other two nil:

```go
// before
fmt.Printf("Ingestion rate: %d\n", tier.IngestionRate)

// after
if m := tier.Metrics; m != nil {
	fmt.Printf("Ingestion rate: %d\n", m.IngestionRate)
}
```

`DeploymentCreationRequest.Deduplication` / `.DeduplicationUnit` and
`DeploymentUpdateRequest.Deduplication` / `.DeduplicationUnit` changed from `uint32` and
`DurationUnit` to `*uint32` and `*DurationUnit`, so that a zero window can actually be
sent and stays distinct from having no window at all. On an update, nil additionally means
"leave the window alone"; on a create, nil is required for the types that have no window:

```go
// before
Deduplication:     10,
DeduplicationUnit: vmcloud.DurationUnitSecond,

// after
Deduplication:     new(uint32(10)),
DeduplicationUnit: new(vmcloud.DurationUnitSecond),
```

`TierInfo` no longer carries the limits as flat fields. The fields every tier reports moved
into an embedded `TierInfoCommon`, and the limits moved into one struct per deployment type,
of which exactly one is non-nil:

```go
// before
fmt.Printf("%d %d\n", tier.IngestionRate, tier.ActiveTimeSeries)

// after
if m := tier.Metrics; m != nil {
	fmt.Printf("%d %d\n", m.IngestionRate, m.ActiveTimeSeries)
}
if l := tier.Logs; l != nil {
	fmt.Printf("%d %d\n", l.IngestionRateBytes, l.ActiveLogStreams)
}
if t := tier.Traces; t != nil {
	fmt.Printf("%d %d\n", t.IngestionRateBytes, t.ActiveLogStreams)
}
```

Reading `tier.ID`, `tier.Name` and `tier.Type` still works, since `TierInfoCommon` is
embedded. The byte-valued limits (`IngestionRateBytes`, `DataReadRate`, `BytesPerQuery`) are
`int64` rather than `int`, so that tiers reporting more than 2 GB decode on 32-bit targets.

`DeploymentUpdateRequest` also changed from a type alias to a defined type. Existing code
that builds it with field names keeps compiling; only code that passed a structurally
identical anonymous struct needs updating.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
