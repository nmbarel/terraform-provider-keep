package provider

import (
	"context"
	"fmt"
	"terraform-provider-keep/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ datasource.DataSource = &workflowsDataSource{}
    _ datasource.DataSourceWithConfigure = &workflowsDataSource{}
  )
  
  // NewworkflowsDataSource is a helper function to simplify the provider implementation.
  func NewWorkflowsDataSource() datasource.DataSource {
    return &workflowsDataSource{}
  }
  
  // workflowsDataSource is the data source implementation.
  type workflowsDataSource struct{
    client *client.Client
  }

  // coffeesDataSourceModel maps the data source schema data.
type workflowsDataSourceModel struct {
    Workflows []workflowsModel `tfsdk:"workflows"`
  }
  
  // coffeesModel maps coffees schema data.
  type workflowsModel struct {
    Yaml          types.String               `tfsdk:"yaml"`
  }
  
  // Metadata returns the data source type name.
  func (d *workflowsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workflows"
  }
  
  // Schema defines the schema for the data source.
  func (d *workflowsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "workflows": schema.ListNestedAttribute{
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "yaml": schema.StringAttribute{
                            Computed: true,
                        },
                    },
                },
            },
        },
    }
  }
  
  // Read refreshes the Terraform state with the latest data.
  func (d *workflowsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var state workflowsDataSourceModel

    workflows, err := d.client.GetWorkflows()
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to get Workflows",
            err.Error(),
        )
        return
    }

    for _, workflow := range workflows {
        workflowState := workflowsModel{
            Yaml: types.StringValue(workflow.Yaml),
        }
        state.Workflows = append(state.Workflows, workflowState)

        // set state
        diags := resp.State.Set(ctx, &state)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
  }

  // Configure adds the provider configured client to the data source.
func (d *workflowsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    if req.ProviderData == nil {
      return
    }
  
    client, ok := req.ProviderData.(*client.Client)
    if !ok {
      resp.Diagnostics.AddError(
        "Unexpected Data Source Configure Type",
        fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
      )
  
      return
    }
  
    d.client = client
  }
