package main

import (
	"context"
	"fmt"
	"log"

	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/vmcloud/v1"
)

func main() {
	// Create a new client with your API key
	client, err := vmcloud.New("your-api-key")
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
	fmt.Println("=== Available Tiers ===")
	tiers, err := client.ListTiers(ctx)
	if err != nil {
		log.Fatalf("Failed to list tiers: %v", err)
	}

	for _, tier := range tiers {
		fmt.Printf("Tier: %s (ID: %d)\n", tier.Name, tier.ID)
		fmt.Printf("  Type: %s, Cloud Provider: %s\n", tier.Type, tier.CloudProvider)
		fmt.Printf("  Cost per hour: $%.4f\n", tier.ComputeCostPerHour)
		fmt.Printf("  Ingestion Rate: %d, Active Time Series: %d\n", tier.IngestionRate, tier.ActiveTimeSeries)
		fmt.Println()
	}

	// Example 4: Find a specific tier by ID
	fmt.Println("=== Finding a Specific Tier ===")
	const targetTierID = 21 // Example tier ID (s.small.a)
	var foundTier *vmcloud.TierInfo

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
		fmt.Printf("  Ingestion Rate: %d, Active Time Series: %d\n", foundTier.IngestionRate, foundTier.ActiveTimeSeries)
	} else {
		fmt.Printf("Tier with ID %d not found\n", targetTierID)
	}
}
