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

	// For demonstration purposes, use a placeholder deployment ID
	// In a real application, you would use the ID from an actual deployment
	deploymentID := "example-deployment-id"

	// Example 1: List rule files for a deployment
	fmt.Printf("Listing rule files for deployment ID: %s\n", deploymentID)
	ruleFiles, err := client.ListDeploymentRuleFileNames(ctx, deploymentID)
	if err != nil {
		log.Fatalf("Failed to list rule files: %v", err)
	}

	if len(ruleFiles) == 0 {
		fmt.Println("No rule files found for this deployment")
	} else {
		fmt.Println("Rule files:")
		for _, fileName := range ruleFiles {
			fmt.Printf("  - %s\n", fileName)
		}
	}
	fmt.Println()

	// Example 2: Create a new rule file
	// Example Prometheus alerting rule
	ruleFileName := "high_cpu_alert.yml"
	// Example rule content - we're printing this for demonstration purposes
	exampleRuleContent := `groups:
- name: cpu_alerts
  rules:
  - alert: HighCPUUsage
    expr: cpu_usage_total{instance=~".*"} > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: High CPU usage detected
      description: "CPU usage is above 80% for 5 minutes on {{ $labels.instance }}"
`
	fmt.Println("Example rule content:")
	fmt.Println("---")
	fmt.Println(exampleRuleContent)
	fmt.Println("---")

	fmt.Printf("Creating rule file '%s' for deployment ID: %s\n", ruleFileName, deploymentID)

	// Comment out in a real application if you don't want to actually create the rule
	/*
		err = client.CreateDeploymentRuleFileContent(ctx, deploymentID, ruleFileName, ruleContent)
		if err != nil {
			log.Fatalf("Failed to create rule file: %v", err)
		}
		fmt.Printf("Successfully created rule file: %s\n", ruleFileName)
	*/
	fmt.Println("Note: Rule file creation code is commented out to prevent actual creation")
	fmt.Println()

	// Example 3: Get rule file content
	fmt.Printf("Getting content of rule file '%s' for deployment ID: %s\n", ruleFileName, deploymentID)
	// Comment out in a real application if you don't want to actually get the rule content
	/*
		content, err := client.GetDeploymentRuleFileContent(ctx, deploymentID, ruleFileName)
		if err != nil {
			log.Fatalf("Failed to get rule file content: %v", err)
		}

		fmt.Println("Rule file content:")
		fmt.Println("---")
		fmt.Println(content)
		fmt.Println("---")
	*/
	fmt.Println("Note: Rule file content retrieval code is commented out")
	fmt.Println()

	// Example 4: Update a rule file
	// Updated rule with additional alert - we're printing this for demonstration purposes
	updatedRuleContent := `groups:
- name: cpu_alerts
  rules:
  - alert: HighCPUUsage
    expr: cpu_usage_total{instance=~".*"} > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: High CPU usage detected
      description: "CPU usage is above 80% for 5 minutes on {{ $labels.instance }}"
  - alert: CriticalCPUUsage
    expr: cpu_usage_total{instance=~".*"} > 95
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: Critical CPU usage detected
      description: "CPU usage is above 95% for 2 minutes on {{ $labels.instance }}"
`
	fmt.Println("Example updated rule content (with additional alert):")
	fmt.Println("---")
	fmt.Println(updatedRuleContent)
	fmt.Println("---")

	fmt.Printf("Updating rule file '%s' for deployment ID: %s\n", ruleFileName, deploymentID)
	// Comment out in a real application if you don't want to actually update the rule
	/*
		err = client.UpdateDeploymentRuleFileContent(ctx, deploymentID, ruleFileName, updatedRuleContent)
		if err != nil {
			log.Fatalf("Failed to update rule file: %v", err)
		}
		fmt.Printf("Successfully updated rule file: %s\n", ruleFileName)
	*/
	fmt.Println("Note: Rule file update code is commented out to prevent actual update")
	fmt.Println()

	// Example 5: Delete a rule file
	fmt.Printf("Deleting rule file '%s' from deployment ID: %s\n", ruleFileName, deploymentID)

	// Comment out in a real application if you don't want to actually delete the rule
	/*
		err = client.DeleteDeploymentRuleFile(ctx, deploymentID, ruleFileName)
		if err != nil {
			log.Fatalf("Failed to delete rule file: %v", err)
		}
		fmt.Printf("Successfully deleted rule file: %s\n", ruleFileName)
	*/
	fmt.Println("Note: Rule file deletion code is commented out to prevent actual deletion")
}
