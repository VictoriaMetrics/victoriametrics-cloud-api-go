package main

import (
	"context"
	"fmt"
	"log"
	"time"

	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"
)

func main() {
	// Create a new client with your API key
	client, err := vmcloud.New("your-api-key")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a context for API requests
	ctx := context.Background()

	// Example 1: List all deployments
	deployments, err := client.ListDeployments(ctx)
	if err != nil {
		log.Fatalf("Failed to list deployments: %v", err)
	}

	for _, deployment := range deployments {
		fmt.Printf("Deployment: %s (ID: %s)\n", deployment.Name, deployment.ID)
		fmt.Printf("  Type: %s, Status: %s\n", deployment.Type, deployment.Status)
		fmt.Printf("  Provider: %s, Region: %s\n", deployment.CloudProvider, deployment.Region)
		fmt.Printf("  VM Version: %s GB\n", deployment.Version)
		fmt.Printf("  Created: %s\n", deployment.CreatedAt.Format(time.RFC3339))
		fmt.Println()
	}

	// Example 2: Create a new deployment
	// Note: This is a long-running operation and will create actual resources in your account
	// Uncomment and modify as needed
	/*
		newDeployment := vmcloud.DeploymentCreationRequest{
			Name:              "example-deployment",
			Type:              vmcloud.DeploymentTypeSingleNode,
			Provider:          vmcloud.DeploymentCloudProviderAWS,
			Region:            "us-east-2", // US East (Ohio)
			Tier:              21,          // s.small.a
			StorageSize:       10,          // 10 GB storage
			StorageSizeUnit:   vmcloud.StorageUnitGB,
			Retention:         30, // 30 days retention
			RetentionUnit:     vmcloud.DurationUnitDay,
			Deduplication:     10, // 10 seconds deduplication
			DeduplicationUnit: vmcloud.DurationUnitSecond,
			MaintenanceWindow: vmcloud.MaintenanceWindowWeekendDays, // Maintenance on weekends
		}

		createdDeployment, err := client.CreateDeployment(ctx, newDeployment)
		if err != nil {
			log.Fatalf("Failed to create deployment: %v", err)
		}

		fmt.Printf("Created deployment: %s (ID: %s)\n", createdDeployment.Name, createdDeployment.ID)
		fmt.Printf("  Type: %s, Status: %s\n", createdDeployment.Type, createdDeployment.Status)
		fmt.Printf("  Provider: %s, Region: %s\n", createdDeployment.CloudProvider, createdDeployment.Region)
		fmt.Printf("  VM Version: %s GB\n", createdDeployment.Version)
		fmt.Printf("  Created: %s\n", createdDeployment.CreatedAt.Format(time.RFC3339))
		fmt.Printf("  Access Endpoint: %s\n", createdDeployment.AccessEndpoint)
		fmt.Println()

		// Store the deployment ID for later examples
		deploymentID := createdDeployment.ID
	*/

	// For demonstration purposes, use a placeholder deployment ID
	// In a real application, you would use the ID from a created deployment
	deploymentID := "example-deployment-id"

	// Example 3: Get deployment details
	fmt.Printf("Getting details for deployment ID: %s\n", deploymentID)
	deploymentDetails, err := client.GetDeploymentDetails(ctx, deploymentID)
	if err != nil {
		log.Fatalf("Failed to get deployment details: %v", err)
	}

	fmt.Printf("Deployment details: %s (ID: %s)\n", deploymentDetails.Name, deploymentDetails.ID)
	fmt.Printf("  Type: %s, Status: %s\n", deploymentDetails.Type, deploymentDetails.Status)
	fmt.Printf("  Provider: %s, Region: %s\n", deploymentDetails.CloudProvider, deploymentDetails.Region)
	fmt.Printf("  Storage Size: %d GB\n", deploymentDetails.StorageSizeGb)
	fmt.Printf("  Retention: %d %s\n", deploymentDetails.RetentionValue, deploymentDetails.RetentionUnit)
	fmt.Printf("  Deduplication: %d %s\n", deploymentDetails.DeduplicationValue, deploymentDetails.DeduplicationUnit)
	fmt.Printf("  Access Endpoint: %s\n", deploymentDetails.AccessEndpoint)
	fmt.Println()

	// Example 4: Update a deployment
	fmt.Printf("Updating deployment ID: %s\n", deploymentID)
	// Uncomment in a real application
	/*
		updateRequest := vmcloud.DeploymentUpdateRequest{
			Name:              "updated-example-deployment",
			Tier:              22, // Upgrade to a larger tier (s.micro.a)
			StorageSize:       20, // Increase storage
			StorageSizeUnit:   vmcloud.StorageUnitGB,
			Retention:         60, // Increase retention
			RetentionUnit:     vmcloud.DurationUnitDay,
			Deduplication:     10, // Keep deduplication the same
			DeduplicationUnit: vmcloud.DurationUnitSecond,
			MaintenanceWindow: vmcloud.MaintenanceWindowBusinessDays, // Change maintenance window to business days
		}

		updatedDeployment, err := client.UpdateDeployment(ctx, deploymentID, updateRequest)
		if err != nil {
			log.Fatalf("Failed to update deployment: %v", err)
		}

		fmt.Printf("Updated deployment: %s (ID: %s)\n", updatedDeployment.Name, updatedDeployment.ID)
		fmt.Printf("  New Tier: %d\n", updatedDeployment.Tier)
		fmt.Printf("  New Storage Size: %d GB\n", updatedDeployment.StorageSizeGb)
		fmt.Printf("  New Retention: %d %s\n", updatedDeployment.RetentionValue, updatedDeployment.RetentionUnit)
		fmt.Printf("  New Deduplication: %d %s\n", updatedDeployment.DeduplicationValue, updatedDeployment.DeduplicationUnit)
		fmt.Printf("  New Maintenance Window: %s\n", updatedDeployment.MaintenanceWindow)
		fmt.Println()
	*/

	// Example 5: Delete a deployment
	fmt.Printf("Deleting deployment ID: %s\n", deploymentID)
	// Uncomment in a real application
	/*
		err = client.DeleteDeployment(ctx, deploymentID)
		if err != nil {
			log.Fatalf("Failed to delete deployment: %v", err)
		}
		fmt.Printf("Successfully deleted deployment ID: %s\n", deploymentID)
	*/
	fmt.Println("Note: Deletion code is commented out to prevent accidental deletion")
}
