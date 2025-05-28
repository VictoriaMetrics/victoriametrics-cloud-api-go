package v1

import (
	"strings"
	"testing"
)

func TestValidateCommonDeploymentParams(t *testing.T) {
	tests := []struct {
		name              string
		deploymentName    string
		tier              uint32
		maintenanceWindow MaintenanceWindow
		storageSize       uint64
		storageSizeUnit   StorageUnit
		retention         uint32
		retentionUnit     DurationUnit
		deduplicationUnit DurationUnit
		wantErr           bool
		errContains       string
	}{
		{
			name:              "Valid parameters",
			deploymentName:    "test-deployment",
			tier:              21,
			maintenanceWindow: MaintenanceWindowWeekendDays,
			storageSize:       10,
			storageSizeUnit:   StorageUnitGB,
			retention:         30,
			retentionUnit:     DurationUnitDay,
			deduplicationUnit: DurationUnitSecond,
			wantErr:           false,
		},
		{
			name:              "Empty name",
			deploymentName:    "",
			tier:              21,
			maintenanceWindow: MaintenanceWindowWeekendDays,
			storageSize:       10,
			storageSizeUnit:   StorageUnitGB,
			retention:         30,
			retentionUnit:     DurationUnitDay,
			deduplicationUnit: DurationUnitSecond,
			wantErr:           true,
			errContains:       "name cannot be empty",
		},
		{
			name:              "Zero tier",
			deploymentName:    "test-deployment",
			tier:              0,
			maintenanceWindow: MaintenanceWindowWeekendDays,
			storageSize:       10,
			storageSizeUnit:   StorageUnitGB,
			retention:         30,
			retentionUnit:     DurationUnitDay,
			deduplicationUnit: DurationUnitSecond,
			wantErr:           true,
			errContains:       "tier cannot be empty",
		},
		{
			name:              "Invalid maintenance window",
			deploymentName:    "test-deployment",
			tier:              21,
			maintenanceWindow: "invalid-window",
			storageSize:       10,
			storageSizeUnit:   StorageUnitGB,
			retention:         30,
			retentionUnit:     DurationUnitDay,
			deduplicationUnit: DurationUnitSecond,
			wantErr:           true,
			errContains:       "invalid maintenance window",
		},
		{
			name:              "Zero storage size",
			deploymentName:    "test-deployment",
			tier:              21,
			maintenanceWindow: MaintenanceWindowWeekendDays,
			storageSize:       0,
			storageSizeUnit:   StorageUnitGB,
			retention:         30,
			retentionUnit:     DurationUnitDay,
			deduplicationUnit: DurationUnitSecond,
			wantErr:           true,
			errContains:       "storage size cannot be zero",
		},
		{
			name:              "Invalid storage size unit",
			deploymentName:    "test-deployment",
			tier:              21,
			maintenanceWindow: MaintenanceWindowWeekendDays,
			storageSize:       10,
			storageSizeUnit:   "invalid-unit",
			retention:         30,
			retentionUnit:     DurationUnitDay,
			deduplicationUnit: DurationUnitSecond,
			wantErr:           true,
			errContains:       "invalid storage size unit",
		},
		{
			name:              "Storage size too small",
			deploymentName:    "test-deployment",
			tier:              21,
			maintenanceWindow: MaintenanceWindowWeekendDays,
			storageSize:       5,
			storageSizeUnit:   StorageUnitGB,
			retention:         30,
			retentionUnit:     DurationUnitDay,
			deduplicationUnit: DurationUnitSecond,
			wantErr:           true,
			errContains:       "must be at least 10 GB",
		},
		{
			name:              "Zero retention",
			deploymentName:    "test-deployment",
			tier:              21,
			maintenanceWindow: MaintenanceWindowWeekendDays,
			storageSize:       10,
			storageSizeUnit:   StorageUnitGB,
			retention:         0,
			retentionUnit:     DurationUnitDay,
			deduplicationUnit: DurationUnitSecond,
			wantErr:           true,
			errContains:       "retention cannot be zero",
		},
		{
			name:              "Invalid retention unit",
			deploymentName:    "test-deployment",
			tier:              21,
			maintenanceWindow: MaintenanceWindowWeekendDays,
			storageSize:       10,
			storageSizeUnit:   StorageUnitGB,
			retention:         30,
			retentionUnit:     "invalid-unit",
			deduplicationUnit: DurationUnitSecond,
			wantErr:           true,
			errContains:       "invalid retention unit",
		},
		{
			name:              "Invalid deduplication unit",
			deploymentName:    "test-deployment",
			tier:              21,
			maintenanceWindow: MaintenanceWindowWeekendDays,
			storageSize:       10,
			storageSizeUnit:   StorageUnitGB,
			retention:         30,
			retentionUnit:     DurationUnitDay,
			deduplicationUnit: "invalid-unit",
			wantErr:           true,
			errContains:       "invalid deduplication unit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommonDeploymentParams(
				tt.deploymentName,
				tt.tier,
				tt.maintenanceWindow,
				tt.storageSize,
				tt.storageSizeUnit,
				tt.retention,
				tt.retentionUnit,
				tt.deduplicationUnit,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateCommonDeploymentParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("validateCommonDeploymentParams() error = %v, want it to contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestValidateCreateDeploymentParams(t *testing.T) {
	tests := []struct {
		name              string
		deploymentType    DeploymentType
		region            string
		provider          DeploymentCloudProvider
		deploymentStorage uint64
		storageSizeUnit   StorageUnit
		wantErr           bool
		errContains       string
	}{
		{
			name:              "Valid parameters",
			deploymentType:    DeploymentTypeSingleNode,
			region:            "us-east-1",
			provider:          DeploymentCloudProviderAWS,
			deploymentStorage: 10,
			storageSizeUnit:   StorageUnitGB,
			wantErr:           false,
		},
		{
			name:              "Invalid deployment type",
			deploymentType:    "invalid-type",
			region:            "us-east-1",
			provider:          DeploymentCloudProviderAWS,
			deploymentStorage: 10,
			storageSizeUnit:   StorageUnitGB,
			wantErr:           true,
			errContains:       "invalid deployment type",
		},
		{
			name:              "Empty region",
			deploymentType:    DeploymentTypeSingleNode,
			region:            "",
			provider:          DeploymentCloudProviderAWS,
			deploymentStorage: 10,
			storageSizeUnit:   StorageUnitGB,
			wantErr:           true,
			errContains:       "region cannot be empty",
		},
		{
			name:              "Invalid provider",
			deploymentType:    DeploymentTypeSingleNode,
			region:            "us-east-1",
			provider:          "invalid-provider",
			deploymentStorage: 10,
			storageSizeUnit:   StorageUnitGB,
			wantErr:           true,
			errContains:       "unsupported deployment cloud provider",
		},
		{
			name:              "Storage too large for single-node",
			deploymentType:    DeploymentTypeSingleNode,
			region:            "us-east-1",
			provider:          DeploymentCloudProviderAWS,
			deploymentStorage: 20,
			storageSizeUnit:   StorageUnitTB,
			wantErr:           true,
			errContains:       "cannot have more than 16 TB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateDeploymentParams(
				tt.deploymentType,
				tt.region,
				tt.provider,
				tt.deploymentStorage,
				tt.storageSizeUnit,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateCreateDeploymentParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("validateCreateDeploymentParams() error = %v, want it to contain %v", err, tt.errContains)
				}
			}
		})
	}
}
