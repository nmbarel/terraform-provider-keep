package provider

import (
	"context"
	"fmt"
	"os"

	"terraform-provider-keep/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ provider.Provider = &keepProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
    return func() provider.Provider {
        return &keepProvider{
            version: version,
        }
    }
}

// keepProvider is the provider implementation.
type keepProvider struct {
    // version is set to the provider version on release, "dev" when the
    // provider is built and ran locally, and "test" when running acceptance
    // testing.
    version string
}

// Metadata returns the provider type name.
func (p *keepProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "keep"
    resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *keepProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
			Attributes: map[string]schema.Attribute{
					"api_key": schema.StringAttribute{
							Required: true,
					},
					"host_url": schema.StringAttribute{
							Required: true,
					},
			},
	}
}

// keepProviderModel maps provider schema data to a Go type.
type keepProviderModel struct {
	HostURL		 types.String `tfsdk:"host_url"`
	ApiKey     types.String `tfsdk:"api_key"`
}


// Configure prepares a keep API client for data sources and resources.
func (p *keepProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config keepProviderModel
	fmt.Print("configing")
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
			return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ApiKey.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
					path.Root("api_key"),
					"Unknown Keep API Key",
					"The provider cannot create the Keep API client as there is an unknown configuration value for the Keep API Key. "+
							"Either target apply the source of the value first, set the value statically in the configuration, or use the KEEP_API_KEY environment variable.",
			)
	}

	if config.HostURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
				path.Root("host_url"),
				"Unknown Keep API host url",
				"The provider cannot create the Keep API client as there is an unknown configuration value for the Keep API host url. "+
						"Either target apply the source of the value first, set the value statically in the configuration, or use the KEEP_API_HOST_URL environment variable.",
		)
}

	if resp.Diagnostics.HasError() {
			return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	api_key := os.Getenv("KEEP_API_KEY")
	host_url := os.Getenv("KEEP_API_HOST_URL")

	if !config.ApiKey.IsNull() {
			api_key = config.ApiKey.ValueString()

	if !config.HostURL.IsNull() {
			host_url = config.HostURL.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if api_key == "" {
			resp.Diagnostics.AddAttributeError(
					path.Root("api_key"),
					"Missing Keep API Key",
					"The provider cannot create the Keep API client as there is a missing or empty value for the Keep API Key. "+
							"Set the host value in the configuration or use the KEEP_API_KEY environment variable. "+
							"If either is already set, ensure the value is not empty.",
			)
	}

	if host_url == "" {
			resp.Diagnostics.AddAttributeError(
					path.Root("host_url"),
					"Missing Keep API Host URL",
					"The provider cannot create the Keep API client as there is a missing or empty value for the Keep API Host URL. "+
							"Set the host value in the configuration or use the KEEP_API_HOST_URL environment variable. "+
							"If either is already set, ensure the value is not empty.",
			)
	}

	if resp.Diagnostics.HasError() {
			return
	}

	// Create a new Keep client using the configuration values
	client, err := client.NewClient(host_url, api_key)
	if err != nil {
			resp.Diagnostics.AddError(
					"Unable to Create Keep API Client",
					"An unexpected error occurred when creating the Keep API client. "+
							"If the error is not clear, please contact the provider developers.\n\n"+
							"Keep Client Error: "+err.Error(),
			)
			return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}
}

// DataSources defines the data sources implemented in the provider.
func (p *keepProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource {
			NewWorkflowsDataSource,
		}
}

// Resources defines the resources implemented in the provider.
func (p *keepProvider) Resources(_ context.Context) []func() resource.Resource {
    return nil
}