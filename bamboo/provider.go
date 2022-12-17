package bamboo

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	bamboo "github.com/rcarmstrong/go-bamboo"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &bambooProvider{}
)

// bambooProviderModel maps provider schema data to a Go type.
type bambooProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &bambooProvider{}
}

// bambooProvider is the provider implementation.
type bambooProvider struct{}

// Metadata returns the provider type name.
func (p *bambooProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bamboo"
}

// Schema defines the provider-level schema for configuration data.
func (p *bambooProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server": schema.StringAttribute{
				Optional: false,
			},
			"username": schema.StringAttribute{
				Optional: false,
			},
			"password": schema.StringAttribute{
				Optional:  false,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a Bamboo API client for data sources and resources.
func (p *bambooProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config bambooProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Bamboo API Host",
			"The provider cannot create the Bamboo API client as there is an unknown configuration value for the Bamboo API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BAMBOO_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Bamboo API Username",
			"The provider cannot create the Bamboo API client as there is an unknown configuration value for the Bamboo API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BAMBOO_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Bamboo API Password",
			"The provider cannot create the Bamboo API client as there is an unknown configuration value for the Bamboo API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BAMBOO_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("BAMBOO_HOST")
	username := os.Getenv("BAMBOO_USERNAME")
	password := os.Getenv("BAMBOO_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Bamboo API Host",
			"The provider cannot create the Bamboo API client as there is a missing or empty value for the Bamboo API host. "+
				"Set the host value in the configuration or use the BAMBOO_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Bamboo API Username",
			"The provider cannot create the Bamboo API client as there is a missing or empty value for the Bamboo API username. "+
				"Set the username value in the configuration or use the BAMBOO_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Bamboo API Password",
			"The provider cannot create the Bamboo API client as there is a missing or empty value for the Bamboo API password. "+
				"Set the password value in the configuration or use the BAMBOO_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Bamboo client using the configuration values
	client := bamboo.NewSimpleClient(nil, username, password)
	client.SetURL(host)

	// Make the Bamboo client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *bambooProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *bambooProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
