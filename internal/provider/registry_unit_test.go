package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitSplitInternalID(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedPart1 string
		expectedPart2 string
		expectErr     bool
	}{
		{
			name:          "valid composite ID",
			input:         "project123/warehouse456",
			expectedPart1: "project123",
			expectedPart2: "warehouse456",
			expectErr:     false,
		},
		{
			name:          "valid UUID composite ID",
			input:         "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111",
			expectedPart1: "00000000-0000-0000-0000-000000000000",
			expectedPart2: "11111111-1111-1111-1111-111111111111",
			expectErr:     false,
		},
		{
			name:      "missing slash separator",
			input:     "warehouse456",
			expectErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			expectErr: true,
		},
		{
			name:      "leading slash only",
			input:     "/warehouse456",
			expectErr: true,
		},
		{
			name:      "trailing slash only",
			input:     "project123/",
			expectErr: true,
		},
		{
			name:      "too many parts",
			input:     "project123/warehouse456/extra",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := splitInternalID(types.StringValue(tc.input))
			if tc.expectErr {
				require.Error(t, err)
				assert.Empty(t, part1)
				assert.Empty(t, part2)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedPart1, part1)
				assert.Equal(t, tc.expectedPart2, part2)
			}
		})
	}
}

func TestUnitRoleResource_ImportState(t *testing.T) {
	r := &lakekeeperRoleResource{}
	schemaResp := resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())

	testCases := []struct {
		name            string
		importID        string
		expectError     bool
		expectedProject string
		expectedRole    string
	}{
		{
			name:            "valid composite ID",
			importID:        "project-123/role-456",
			expectError:     false,
			expectedProject: "project-123",
			expectedRole:    "role-456",
		},
		{
			name:        "missing slash separator",
			importID:    "role-456",
			expectError: true,
		},
		{
			name:        "empty string",
			importID:    "",
			expectError: true,
		},
		{
			name:        "leading slash only",
			importID:    "/role-456",
			expectError: true,
		},
		{
			name:        "trailing slash only",
			importID:    "project-123/",
			expectError: true,
		},
		{
			name:        "too many parts",
			importID:    "project-123/role-456/extra",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := &resource.ImportStateResponse{
				State: tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(context.Background()), nil),
				},
			}
			r.ImportState(context.Background(), resource.ImportStateRequest{ID: tc.importID}, resp)
			if tc.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
				var proj, role, id types.String
				resp.State.GetAttribute(context.Background(), path.Root("project_id"), &proj)
				resp.State.GetAttribute(context.Background(), path.Root("role_id"), &role)
				resp.State.GetAttribute(context.Background(), path.Root("id"), &id)
				assert.Equal(t, tc.expectedProject, proj.ValueString())
				assert.Equal(t, tc.expectedRole, role.ValueString())
				assert.Equal(t, tc.importID, id.ValueString())
			}
		})
	}
}

func TestUnitWarehouseResource_ImportState(t *testing.T) {
	r := &lakekeeperWarehouseResource{}
	schemaResp := resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())

	testCases := []struct {
		name              string
		importID          string
		expectError       bool
		expectedProject   string
		expectedWarehouse string
	}{
		{
			name:              "valid composite ID",
			importID:          "project-123/warehouse-456",
			expectError:       false,
			expectedProject:   "project-123",
			expectedWarehouse: "warehouse-456",
		},
		{
			name:        "missing slash separator",
			importID:    "warehouse-456",
			expectError: true,
		},
		{
			name:        "empty string",
			importID:    "",
			expectError: true,
		},
		{
			name:        "leading slash only",
			importID:    "/warehouse-456",
			expectError: true,
		},
		{
			name:        "trailing slash only",
			importID:    "project-123/",
			expectError: true,
		},
		{
			name:        "too many parts",
			importID:    "project-123/warehouse-456/extra",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := &resource.ImportStateResponse{
				State: tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(context.Background()), nil),
				},
			}
			r.ImportState(context.Background(), resource.ImportStateRequest{ID: tc.importID}, resp)
			if tc.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
				var proj, wh, id types.String
				resp.State.GetAttribute(context.Background(), path.Root("project_id"), &proj)
				resp.State.GetAttribute(context.Background(), path.Root("warehouse_id"), &wh)
				resp.State.GetAttribute(context.Background(), path.Root("id"), &id)
				assert.Equal(t, tc.expectedProject, proj.ValueString())
				assert.Equal(t, tc.expectedWarehouse, wh.ValueString())
				assert.Equal(t, tc.importID, id.ValueString())
			}
		})
	}
}
