package provider

import (
	"context"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"kubiya-control-plane/internal/clients"
	kubiyasentry "kubiya-control-plane/internal/sentry"
)

type kubiyaControlPlaneProvider struct {
	version string
}

var _ provider.Provider = (*kubiyaControlPlaneProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &kubiyaControlPlaneProvider{version: version}
	}
}

func (p *kubiyaControlPlaneProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAgentResource,
		NewTeamResource,
		NewProjectResource,
		NewEnvironmentResource,
		NewToolSetResource,
		NewWorkerResource,
		NewPolicyResource,
	}
}

func (p *kubiyaControlPlaneProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAgentDataSource,
		NewTeamDataSource,
		NewProjectDataSource,
		NewEnvironmentDataSource,
	}
}

func (p *kubiyaControlPlaneProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *kubiyaControlPlaneProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kubiya_control_plane"
}

func (p *kubiyaControlPlaneProvider) Configure(ctx context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Initialize Sentry when provider is configured
	if err := kubiyasentry.Initialize(); err != nil {
		// Log error but don't fail provider configuration
		// Sentry is optional for functionality
	}

	// Get logger and add to context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	// Start a transaction for provider configuration
	ctx, span := kubiyasentry.StartTransaction(ctx, "provider.configure", "Configuring Kubiya Control Plane provider")
	defer kubiyasentry.FinishSpan(span)

	// Log provider configuration start
	logger.Info("Configuring Kubiya Control Plane provider", "version", p.version)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("provider", "Configuring Kubiya Control Plane provider", sentry.LevelInfo, nil)

	const (
		apiKeyEnvVar         = "KUBIYA_CONTROL_PLANE_API_KEY"
		orgIDEnvVar          = "KUBIYA_CONTROL_PLANE_ORG_ID"
		envKeyEnvVar         = "KUBIYA_CONTROL_PLANE_ENV"
		missingAPIKey        = "Kubiya Control Plane API Key Not Configured"
		missingAPIKeyDetails = "Please set the Kubiya Control Plane API Key using the environment variable 'KUBIYA_CONTROL_PLANE_API_KEY'. " +
			"Use the command below:\n> export KUBIYA_CONTROL_PLANE_API_KEY=YOUR_API_KEY"
		missingOrgID        = "Kubiya Control Plane Organization ID Not Configured"
		missingOrgIDDetails = "Please set the Organization ID using the environment variable 'KUBIYA_CONTROL_PLANE_ORG_ID'. " +
			"Use the command below:\n> export KUBIYA_CONTROL_PLANE_ORG_ID=YOUR_ORG_ID"
	)

	apiKey := os.Getenv(apiKeyEnvVar)
	if apiKey == "" {
		logger.Error("API key not configured", "env_var", apiKeyEnvVar)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		resp.Diagnostics.AddError(missingAPIKey, missingAPIKeyDetails)
		return
	}

	orgID := os.Getenv(orgIDEnvVar)
	if orgID == "" {
		logger.Error("Organization ID not configured", "env_var", orgIDEnvVar)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		resp.Diagnostics.AddError(missingOrgID, missingOrgIDDetails)
		return
	}

	// Fetch the environment or set to default
	env := os.Getenv(envKeyEnvVar)
	if env == "" {
		env = "development"
		logger.Debug("Using default environment", "environment", env)
	} else {
		logger.Info("Using configured environment", "environment", env)
	}

	// Set Sentry environment tag
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("environment", env)
		scope.SetTag("provider.version", p.version)
		scope.SetTag("organization_id", orgID)
	})

	// Log client creation attempt
	logger.Debug("Creating Kubiya Control Plane client", "environment", env)

	// Create a new Kubiya Control Plane client using the API key and environment
	client, err := clients.New(apiKey, env)
	if err != nil {
		logger.Error("Failed to create Kubiya Control Plane client", "error", err, "environment", env)
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		resp.Diagnostics.AddError("Failed to Create Kubiya Control Plane Client", "An error occurred while creating the Kubiya Control Plane client: "+err.Error())
		return
	}

	// Success
	logger.Info("Successfully configured Kubiya Control Plane provider", "environment", env, "version", p.version)
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)

	// Attach the client to be used by resources and data sources
	resp.ResourceData = client
	resp.DataSourceData = client
}
