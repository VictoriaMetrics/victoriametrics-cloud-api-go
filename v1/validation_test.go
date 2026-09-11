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
			wantErr:           true,
			errContains:       "invalid storage size unit",
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
			wantErr:           true,
			errContains:       "invalid retention unit",
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
		name           string
		deploymentType DeploymentType
		region         string
		provider       DeploymentCloudProvider
		wantErr        bool
		errContains    string
	}{
		{
			name:           "Valid parameters",
			deploymentType: DeploymentTypeSingleNode,
			region:         "us-east-1",
			provider:       DeploymentCloudProviderAWS,
			wantErr:        false,
		},
		{
			name:           "Invalid deployment type",
			deploymentType: "invalid-type",
			region:         "us-east-1",
			provider:       DeploymentCloudProviderAWS,
			wantErr:        true,
			errContains:    "invalid deployment type",
		},
		{
			name:           "Empty region",
			deploymentType: DeploymentTypeSingleNode,
			region:         "",
			provider:       DeploymentCloudProviderAWS,
			wantErr:        true,
			errContains:    "region cannot be empty",
		},
		{
			name:           "Invalid provider",
			deploymentType: DeploymentTypeSingleNode,
			region:         "us-east-1",
			provider:       "invalid-provider",
			wantErr:        true,
			errContains:    "unsupported deployment cloud provider",
		},
		{
			name:           "Valid cluster",
			deploymentType: DeploymentTypeCluster,
			region:         "us-east-1",
			provider:       DeploymentCloudProviderAWS,
			wantErr:        false,
		},
		{
			name:           "Valid VictoriaLogs",
			deploymentType: DeploymentTypeVLogs,
			region:         "us-east-1",
			provider:       DeploymentCloudProviderAWS,
			wantErr:        false,
		},
		{
			name:           "Valid VictoriaTraces",
			deploymentType: DeploymentTypeVTraces,
			region:         "us-east-1",
			provider:       DeploymentCloudProviderAWS,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateDeploymentParams(
				tt.deploymentType,
				tt.region,
				tt.provider,
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

func TestValidateDeduplicationForCreate(t *testing.T) {
	tests := []struct {
		name              string
		deploymentType    DeploymentType
		deduplicationUnit DurationUnit
		wantErr           bool
		errContains       string
	}{
		{
			name:              "Seconds for single node",
			deploymentType:    DeploymentTypeSingleNode,
			deduplicationUnit: DurationUnitSecond,
			wantErr:           false,
		},
		{
			name:              "Milliseconds for cluster",
			deploymentType:    DeploymentTypeCluster,
			deduplicationUnit: DurationUnitMillisecond,
			wantErr:           false,
		},
		{
			name:              "Invalid unit for cluster",
			deploymentType:    DeploymentTypeCluster,
			deduplicationUnit: "invalid-unit",
			wantErr:           true,
			errContains:       "invalid deduplication unit",
		},
		{
			name:              "Days are not a deduplication unit",
			deploymentType:    DeploymentTypeSingleNode,
			deduplicationUnit: DurationUnitDay,
			wantErr:           true,
			errContains:       "invalid deduplication unit",
		},
		{
			name:              "Missing unit for single node",
			deploymentType:    DeploymentTypeSingleNode,
			deduplicationUnit: "",
			wantErr:           true,
			errContains:       "invalid deduplication unit",
		},
		{
			// VictoriaLogs has no deduplication window, and the API ignores the field
			name:              "Missing unit for VictoriaLogs",
			deploymentType:    DeploymentTypeVLogs,
			deduplicationUnit: "",
			wantErr:           false,
		},
		{
			name:              "Missing unit for VictoriaTraces",
			deploymentType:    DeploymentTypeVTraces,
			deduplicationUnit: "",
			wantErr:           false,
		},
		{
			// the API ignores it rather than rejecting it, so neither does the client
			name:              "Invalid unit for VictoriaLogs is ignored",
			deploymentType:    DeploymentTypeVLogs,
			deduplicationUnit: "invalid-unit",
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDeduplicationForCreate(tt.deploymentType, tt.deduplicationUnit)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDeduplicationForCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("validateDeduplicationForCreate() error = %v, want it to contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestValidateDeduplicationForUpdate(t *testing.T) {
	tests := []struct {
		name              string
		deduplication     uint32
		deduplicationUnit DurationUnit
		wantErr           bool
		errContains       string
	}{
		{
			name:              "Valid window",
			deduplication:     30,
			deduplicationUnit: DurationUnitSecond,
			wantErr:           false,
		},
		{
			// a metrics deployment may legitimately deduplicate over a zero window
			name:              "Zero window with a unit",
			deduplication:     0,
			deduplicationUnit: DurationUnitSecond,
			wantErr:           false,
		},
		{
			// an update of a VictoriaLogs or VictoriaTraces deployment leaves both unset
			name:              "Both fields unset",
			deduplication:     0,
			deduplicationUnit: "",
			wantErr:           false,
		},
		{
			name:              "Invalid unit",
			deduplication:     30,
			deduplicationUnit: "invalid-unit",
			wantErr:           true,
			errContains:       "invalid deduplication unit",
		},
		{
			name:              "Window without a unit",
			deduplication:     30,
			deduplicationUnit: "",
			wantErr:           true,
			errContains:       "invalid deduplication unit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDeduplicationForUpdate(tt.deduplication, tt.deduplicationUnit)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDeduplicationForUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("validateDeduplicationForUpdate() error = %v, want it to contain %v", err, tt.errContains)
				}
			}
		})
	}
}
