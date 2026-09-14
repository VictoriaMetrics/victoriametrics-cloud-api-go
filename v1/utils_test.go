package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestContextWithDynamicAPIKey(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
	}{
		{
			name:   "non-empty key",
			apiKey: "my-secret-key",
		},
		{
			name:   "empty key",
			apiKey: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithDynamicAPIKey(context.Background(), tt.apiKey)
			got, ok := ctx.Value(apiKeyContextKey).(string)
			if !ok {
				t.Fatal("ContextWithDynamicAPIKey() did not store a string value in context")
			}
			if got != tt.apiKey {
				t.Errorf("ContextWithDynamicAPIKey() stored %q, want %q", got, tt.apiKey)
			}
		})
	}
}

func TestRequestAPIDynamicAPIKey(t *testing.T) {
	tests := []struct {
		name           string
		clientKey      string
		ctxKey         *string // nil means don't set context key
		expectedHeader string
	}{
		{
			name:           "static key, no ctx",
			clientKey:      "static-key",
			ctxKey:         nil,
			expectedHeader: "static-key",
		},
		{
			name:           "static key + ctx key, client takes precedence",
			clientKey:      "static-key",
			ctxKey:         strPtr("dynamic-key"),
			expectedHeader: "static-key",
		},
		{
			name:           "dynamic + ctx key",
			clientKey:      DynamicAPIKey,
			ctxKey:         strPtr("per-request-key"),
			expectedHeader: "per-request-key",
		},
		{
			name:           "dynamic, no ctx key",
			clientKey:      DynamicAPIKey,
			ctxKey:         nil,
			expectedHeader: "",
		},
		{
			name:           "dynamic + empty ctx key",
			clientKey:      DynamicAPIKey,
			ctxKey:         strPtr(""),
			expectedHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedHeader string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedHeader = r.Header.Get(AccessTokenHeader)
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode([]CloudProviderInfo{})
			}))
			defer server.Close()

			client, err := New(tt.clientKey, WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			ctx := context.Background()
			if tt.ctxKey != nil {
				ctx = ContextWithDynamicAPIKey(ctx, *tt.ctxKey)
			}

			_, _ = client.ListCloudProviders(ctx)

			if capturedHeader != tt.expectedHeader {
				t.Errorf("request header %q = %q, want %q", AccessTokenHeader, capturedHeader, tt.expectedHeader)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
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
