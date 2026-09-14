package main

import (
	"context"
	"fmt"
	"log"

	"github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"
)

func main() {
	// Create a new client with your API key
	client, err := v1.New("your-api-key")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a context for API requests
	ctx := context.Background()

	// Example 1: List available cloud providers
	fmt.Println("=== Available Cloud Providers ===")
	providers, err := client.ListCloudProviders(ctx)
	if err != nil {
		log.Fatalf("Failed to list cloud providers: %v", err)
	}

	for _, provider := range providers {
		fmt.Printf("Provider: %s (ID: %s)\n", provider.ID, provider.ID)
	}
	fmt.Println()

	// Example 2: List available regions
	fmt.Println("=== Available Regions ===")
	regions, err := client.ListRegions(ctx)
	if err != nil {
		log.Fatalf("Failed to list regions: %v", err)
	}

	for _, region := range regions {
		fmt.Printf("Region: %s (Provider: %s)\n",
			region.Name,
			region.CloudProvider,
		)
	}
	fmt.Println()

	// Example 3: List available tiers
	// Without a filter the API returns tiers of every deployment type: VictoriaMetrics
	// single-node and cluster, VictoriaLogs, and VictoriaTraces where it is enabled.
	fmt.Println("=== Available Tiers ===")
	tiers, err := client.ListTiers(ctx)
	if err != nil {
		log.Fatalf("Failed to list tiers: %v", err)
	}

	for _, tier := range tiers {
		fmt.Printf("Tier: %s (ID: %d)\n", tier.Name, tier.ID)
		fmt.Printf("  Type: %s, Cloud Provider: %s\n", tier.Type, tier.CloudProvider)
		fmt.Printf("  Cost per hour: $%.4f\n", tier.ComputeCostPerHour)
		// the limits a tier reports depend on its type, so exactly one of the
		// Metrics, Logs and Traces structs is filled in
		switch {
		case tier.Metrics != nil:
			fmt.Printf("  Ingestion Rate: %d, Active Time Series: %d\n", tier.Metrics.IngestionRate, tier.Metrics.ActiveTimeSeries)
			fmt.Printf("  Read Rate: %d, Series per Query: %d\n", tier.Metrics.DatapointsReadRate, tier.Metrics.SeriesReadPerQuery)
		case tier.Logs != nil:
			fmt.Printf("  Ingestion Rate: %d bytes/s, Active Streams: %d\n", tier.Logs.IngestionRateBytes, tier.Logs.ActiveLogStreams)
			fmt.Printf("  Read Rate: %d bytes/s, Bytes per Query: %d\n", tier.Logs.DataReadRate, tier.Logs.BytesPerQuery)
		case tier.Traces != nil:
			fmt.Printf("  Ingestion Rate: %d bytes/s, Active Streams: %d\n", tier.Traces.IngestionRateBytes, tier.Traces.ActiveLogStreams)
			fmt.Printf("  Read Rate: %d bytes/s, Bytes per Query: %d\n", tier.Traces.DataReadRate, tier.Traces.BytesPerQuery)
		}
		fmt.Println()
	}

	// Example 3b: List only the VictoriaLogs tiers
	fmt.Println("=== VictoriaLogs Tiers ===")
	vlogsTiers, err := client.ListTiers(ctx, v1.WithDeploymentType(v1.DeploymentTypeVLogs))
	if err != nil {
		log.Fatalf("Failed to list VictoriaLogs tiers: %v", err)
	}

	for _, tier := range vlogsTiers {
		fmt.Printf("Tier: %s (ID: %d), cost per hour: $%.4f\n", tier.Name, tier.ID, tier.ComputeCostPerHour)
	}
	fmt.Println()

	// Example 4: Find a specific tier by ID
	fmt.Println("=== Finding a Specific Tier ===")
	const targetTierID = 21 // Example tier ID (s.small.a)
	var foundTier *v1.TierInfo

	for _, tier := range tiers {
		if tier.ID == targetTierID {
			foundTier = &tier
			break
		}
	}

	if foundTier != nil {
		fmt.Printf("Found tier: %s (ID: %d)\n", foundTier.Name, foundTier.ID)
		fmt.Printf("  Type: %s, Cloud Provider: %s\n", foundTier.Type, foundTier.CloudProvider)
		fmt.Printf("  Cost per hour: $%.4f\n", foundTier.ComputeCostPerHour)
		if m := foundTier.Metrics; m != nil {
			fmt.Printf("  Ingestion Rate: %d, Active Time Series: %d\n", m.IngestionRate, m.ActiveTimeSeries)
		}
	} else {
		fmt.Printf("Tier with ID %d not found\n", targetTierID)
	}
}
