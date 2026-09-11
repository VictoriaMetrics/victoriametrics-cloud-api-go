package v1

import (
	"fmt"
	"regexp"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func isValidUUID(uuid string) bool {
	return uuidRegex.MatchString(uuid)
}

func checkDeploymentID(deploymentID string) error {
	if deploymentID == "" {
		return fmt.Errorf("deployment ID cannot be empty")
	}
	if !isValidUUID(deploymentID) {
		return fmt.Errorf("invalid deployment ID format: %s", deploymentID)
	}
	return nil
}

var tenantIDRegex = regexp.MustCompile(`^(\d+)(:\d+)?$`)

func isValidTenantID(tenantID string) bool {
	return tenantIDRegex.MatchString(tenantID)
}

// validateCommonDeploymentParams validates parameters common to both create and update operations.
//
// Storage size is only checked for a unit and a non-zero value. The valid range and the
// step it must sit on depend on the storage type, the deployment topology and the
// per-installation configuration, none of which the client can see, so the API is the
// authority on them.
func validateCommonDeploymentParams(
	name string,
	tier uint32,
	maintenanceWindow MaintenanceWindow,
	storageSize uint64,
	storageSizeUnit StorageUnit,
	retention uint32,
	retentionUnit DurationUnit,
) error {
	if name == "" {
		return fmt.Errorf("deployment name cannot be empty")
	}
	if tier == 0 {
		return fmt.Errorf("deployment tier cannot be empty")
	}
	if maintenanceWindow != MaintenanceWindowWeekendDays &&
		maintenanceWindow != MaintenanceWindowBusinessDays {
		return fmt.Errorf("invalid maintenance window: %s", maintenanceWindow)
	}
	if storageSize == 0 {
		return fmt.Errorf("deployment storage size cannot be zero")
	}
	if storageSizeUnit != StorageUnitGB && storageSizeUnit != StorageUnitTB {
		return fmt.Errorf("invalid storage size unit: %s", storageSizeUnit)
	}
	if retention == 0 {
		return fmt.Errorf("deployment retention cannot be zero")
	}
	if retentionUnit != DurationUnitDay && retentionUnit != DurationUnitMonth {
		return fmt.Errorf("invalid retention unit: %s, only days and months are supported", retentionUnit)
	}
	return nil
}

// isValidDeduplicationUnit reports whether unit is a unit the API accepts for a
// deduplication window.
func isValidDeduplicationUnit(unit DurationUnit) bool {
	return unit == DurationUnitSecond || unit == DurationUnitMillisecond
}

// validateDeduplicationForCreate checks the deduplication window of a create request.
// The deployment type is known here, so the rule is exact: metrics deployments must
// carry a unit, and for the other types the API ignores both fields.
func validateDeduplicationForCreate(deploymentType DeploymentType, deduplicationUnit DurationUnit) error {
	if !deploymentType.SupportsDeduplication() {
		return nil
	}
	if !isValidDeduplicationUnit(deduplicationUnit) {
		return fmt.Errorf("invalid deduplication unit: %s, only seconds and milliseconds are supported", deduplicationUnit)
	}
	return nil
}

// validateDeduplicationForUpdate checks the deduplication window of an update request.
// An update request does not carry the deployment type, so a request that leaves both
// deduplication fields unset is taken as one for a deployment that has no deduplication
// window, and the API rejects it if the deployment is a metrics one. A request that sets
// either field is validated as a metrics one.
func validateDeduplicationForUpdate(deduplication uint32, deduplicationUnit DurationUnit) error {
	if deduplication == 0 && deduplicationUnit == "" {
		return nil
	}
	if !isValidDeduplicationUnit(deduplicationUnit) {
		return fmt.Errorf("invalid deduplication unit: %s, only seconds and milliseconds are supported", deduplicationUnit)
	}
	return nil
}

// isValidDeploymentType reports whether deploymentType is a type the API can create.
func isValidDeploymentType(deploymentType DeploymentType) bool {
	switch deploymentType {
	case DeploymentTypeSingleNode, DeploymentTypeCluster, DeploymentTypeVLogs, DeploymentTypeVTraces:
		return true
	default:
		return false
	}
}

// validateCreateDeploymentParams validates parameters specific to deployment creation
func validateCreateDeploymentParams(
	deploymentType DeploymentType,
	region string,
	provider DeploymentCloudProvider,
) error {
	if !isValidDeploymentType(deploymentType) {
		return fmt.Errorf("invalid deployment type: %s", deploymentType)
	}
	if region == "" {
		return fmt.Errorf("deployment region cannot be empty")
	}
	if provider != DeploymentCloudProviderAWS {
		return fmt.Errorf("unsupported deployment cloud provider: %s", provider)
	}
	return nil
}
