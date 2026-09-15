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
		deduplication     *uint32
		deduplicationUnit *DurationUnit
		wantErr           bool
		errContains       string
	}{
		{
			name:              "Seconds for single node",
			deploymentType:    DeploymentTypeSingleNode,
			deduplication:     new(uint32(30)),
			deduplicationUnit: new(DurationUnitSecond),
			wantErr:           false,
		},
		{
			name:              "Milliseconds for cluster",
			deploymentType:    DeploymentTypeCluster,
			deduplication:     new(uint32(500)),
			deduplicationUnit: new(DurationUnitMillisecond),
			wantErr:           false,
		},
		{
			// a zero window is a real setting, distinct from having no window at all
			name:              "Zero window for single node",
			deploymentType:    DeploymentTypeSingleNode,
			deduplication:     new(uint32(0)),
			deduplicationUnit: new(DurationUnitSecond),
			wantErr:           false,
		},
		{
			name:              "Invalid unit for cluster",
			deploymentType:    DeploymentTypeCluster,
			deduplication:     new(uint32(30)),
			deduplicationUnit: new(DurationUnit("invalid-unit")),
			wantErr:           true,
			errContains:       "invalid deduplication unit",
		},
		{
			name:              "Days are not a deduplication unit",
			deploymentType:    DeploymentTypeSingleNode,
			deduplication:     new(uint32(30)),
			deduplicationUnit: new(DurationUnitDay),
			wantErr:           true,
			errContains:       "invalid deduplication unit",
		},
		{
			name:              "Window without a unit for single node",
			deploymentType:    DeploymentTypeSingleNode,
			deduplication:     new(uint32(30)),
			deduplicationUnit: nil,
			wantErr:           true,
			errContains:       "must be set together",
		},
		{
			name:              "Unit without a window for single node",
			deploymentType:    DeploymentTypeSingleNode,
			deduplication:     nil,
			deduplicationUnit: new(DurationUnitSecond),
			wantErr:           true,
			errContains:       "must be set together",
		},
		{
			name:           "Nothing set for single node",
			deploymentType: DeploymentTypeSingleNode,
			wantErr:        true,
			errContains:    "must be set together",
		},
		{
			// VictoriaLogs has no deduplication window, so leaving it unset is the norm
			name:           "Nothing set for VictoriaLogs",
			deploymentType: DeploymentTypeVLogs,
			wantErr:        false,
		},
		{
			name:           "Nothing set for VictoriaTraces",
			deploymentType: DeploymentTypeVTraces,
			wantErr:        false,
		},
		{
			// setting a window on a type that has none is a caller mistake, and saying so
			// beats dropping it silently on the way to the API
			name:              "Window set for VictoriaLogs",
			deploymentType:    DeploymentTypeVLogs,
			deduplication:     new(uint32(10)),
			deduplicationUnit: new(DurationUnitSecond),
			wantErr:           true,
			errContains:       "not supported for vlogs_single deployments",
		},
		{
			name:              "Unit alone set for VictoriaTraces",
			deploymentType:    DeploymentTypeVTraces,
			deduplicationUnit: new(DurationUnitSecond),
			wantErr:           true,
			errContains:       "not supported for vtraces_single deployments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDeduplicationForCreate(tt.deploymentType, tt.deduplication, tt.deduplicationUnit)
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
		deduplication     *uint32
		deduplicationUnit *DurationUnit
		wantErr           bool
		errContains       string
	}{
		{
			name:              "Valid window",
			deduplication:     new(uint32(30)),
			deduplicationUnit: new(DurationUnitSecond),
			wantErr:           false,
		},
		{
			// the case a plain uint32 could not express: zero is a real window, and with a
			// pointer it stays distinct from "not changing the window"
			name:              "Zero window",
			deduplication:     new(uint32(0)),
			deduplicationUnit: new(DurationUnitMillisecond),
			wantErr:           false,
		},
		{
			// an update that is not touching deduplication, including every vlogs_single
			// and vtraces_single update
			name:              "Both unset leaves the window alone",
			deduplication:     nil,
			deduplicationUnit: nil,
			wantErr:           false,
		},
		{
			name:              "Window without a unit",
			deduplication:     new(uint32(30)),
			deduplicationUnit: nil,
			wantErr:           true,
			errContains:       "must be set together",
		},
		{
			name:              "Unit without a window",
			deduplication:     nil,
			deduplicationUnit: new(DurationUnitSecond),
			wantErr:           true,
			errContains:       "must be set together",
		},
		{
			name:              "Invalid unit",
			deduplication:     new(uint32(30)),
			deduplicationUnit: new(DurationUnit("invalid-unit")),
			wantErr:           true,
			errContains:       "invalid deduplication unit",
		},
		{
			name:              "Days are not a deduplication unit",
			deduplication:     new(uint32(30)),
			deduplicationUnit: new(DurationUnitDay),
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
