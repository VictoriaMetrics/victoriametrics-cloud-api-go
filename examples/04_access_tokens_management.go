package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	// For demonstration purposes, use a placeholder deployment ID
	// In a real application, you would use the ID from an actual deployment
	deploymentID := "example-deployment-id"

	// Example 1: List access tokens for a deployment
	fmt.Printf("Listing access tokens for deployment ID: %s\n", deploymentID)
	tokens, err := client.ListDeploymentAccessTokens(ctx, deploymentID)
	if err != nil {
		log.Fatalf("Failed to list access tokens: %v", err)
	}

	if len(tokens) == 0 {
		fmt.Println("No access tokens found for this deployment")
	} else {
		for _, token := range tokens {
			fmt.Printf("Token: %s (ID: %s)\n", token.Description, token.ID)
			fmt.Printf("  Type: %s\n", token.Type)
			fmt.Printf("  Created by: %s at %s\n", token.CreatedBy, token.CreatedAt.Format(time.RFC3339))
			fmt.Printf("  Secret: %s (first 4 chars)\n", token.Secret)
			fmt.Println()
		}
	}

	// Example 2: Create a new access token
	fmt.Printf("Creating access token for deployment ID: %s\n", deploymentID)
	// Uncomment in a real application
	/*
		tokenRequest := vmcloud.AccessTokenCreateRequest{
			Description: "Example API token for monitoring",
			Type:        vmcloud.AccessModeRead, // Read-only token
		}

		createdToken, err := client.CreateDeploymentAccessToken(ctx, deploymentID, tokenRequest)
		if err != nil {
			log.Fatalf("Failed to create access token: %v", err)
		}

		fmt.Printf("Created token: %s (ID: %s)\n", createdToken.Description, createdToken.ID)
		fmt.Printf("  Type: %s\n", createdToken.Type)
		fmt.Printf("  Created at: %s\n", createdToken.CreatedAt.Format(time.RFC3339))
		fmt.Printf("  Secret: %s\n", createdToken.Secret)
		fmt.Println("  IMPORTANT: Store this secret securely!")
		fmt.Println()

		// Store the token ID for later examples
		tokenID := createdToken.ID
	*/

	// For demonstration purposes, use a placeholder token ID
	tokenID := "example-token-id"

	// Example 3: Reveal a token's secret
	fmt.Printf("Revealing secret for token ID: %s in deployment ID: %s\n", tokenID, deploymentID)
	// Uncomment in a real application
	/*
		revealedToken, err := client.RevealDeploymentAccessToken(ctx, deploymentID, tokenID)
		if err != nil {
			log.Fatalf("Failed to reveal token secret: %v", err)
		}

		fmt.Printf("Token: %s (ID: %s)\n", revealedToken.Description, revealedToken.ID)
		fmt.Printf("  Full Secret: %s\n", revealedToken.Secret)
		fmt.Println("  IMPORTANT: Store this secret securely!")
		fmt.Println()
	*/

	// Example 4: Delete an access token
	fmt.Printf("Deleting token ID: %s from deployment ID: %s\n", tokenID, deploymentID)
	// Uncomment in a real application
	/*
		err = client.DeleteDeploymentAccessToken(ctx, deploymentID, tokenID)
		if err != nil {
			log.Fatalf("Failed to delete access token: %v", err)
		}
		fmt.Printf("Successfully deleted token ID: %s\n", tokenID)
	*/
	fmt.Println("Note: Token deletion code is commented out to prevent accidental deletion")
}
