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
deployments, err := client.ListDeployments(context.Background())
if err != nil {
	log.Fatalf("Failed to list deployments: %v", err)
}

for _, deployment := range deployments {
	fmt.Printf("Deployment: %s (ID: %s)\n", deployment.Name, deployment.ID)
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
	Deduplication:     10, // deduplication period in seconds
	DeduplicationUnit: vmcloud.DurationUnitSecond,
	MaintenanceWindow: vmcloud.MaintenanceWindowWeekendDays, // maintenance window on weekends
}

createdDeployment, err := client.CreateDeployment(context.Background(), deployment)
if err != nil {
	log.Fatalf("Failed to create deployment: %v", err)
}

fmt.Printf("Created deployment: %s (ID: %s)\n", createdDeployment.Name, createdDeployment.ID)
```

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

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
