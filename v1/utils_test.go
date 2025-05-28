package v1

import (
	"testing"
)

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name string
		uuid string
		want bool
	}{
		{
			name: "valid UUID",
			uuid: "123e4567-e89b-12d3-a456-426614174000",
			want: true,
		},
		{
			name: "invalid UUID - wrong format",
			uuid: "123e4567-e89b-12d3-a456-42661417400",
			want: false,
		},
		{
			name: "invalid UUID - empty string",
			uuid: "",
			want: false,
		},
		{
			name: "invalid UUID - not a UUID",
			uuid: "not-a-uuid",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidUUID(tt.uuid); got != tt.want {
				t.Errorf("isValidUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDeploymentID(t *testing.T) {
	tests := []struct {
		name         string
		deploymentID string
		wantErr      bool
	}{
		{
			name:         "valid deployment ID",
			deploymentID: "123e4567-e89b-12d3-a456-426614174000",
			wantErr:      false,
		},
		{
			name:         "invalid deployment ID - empty string",
			deploymentID: "",
			wantErr:      true,
		},
		{
			name:         "invalid deployment ID - not a UUID",
			deploymentID: "not-a-uuid",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkDeploymentID(tt.deploymentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkDeploymentID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidTenantID(t *testing.T) {
	tests := []struct {
		name     string
		tenantID string
		want     bool
	}{
		{
			name:     "valid tenant ID - simple number",
			tenantID: "12345",
			want:     true,
		},
		{
			name:     "valid tenant ID - with colon and number",
			tenantID: "12345:67890",
			want:     true,
		},
		{
			name:     "invalid tenant ID - empty string",
			tenantID: "",
			want:     false,
		},
		{
			name:     "invalid tenant ID - not a number",
			tenantID: "not-a-number",
			want:     false,
		},
		{
			name:     "invalid tenant ID - wrong format",
			tenantID: "12345:abcde",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidTenantID(tt.tenantID); got != tt.want {
				t.Errorf("isValidTenantID() = %v, want %v", got, tt.want)
			}
		})
	}
}
