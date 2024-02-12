package provider

import (
	"context"
	"fmt"

	"terraform-provider-keep/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource = &workflowsResource{}
		_ resource.ResourceWithConfigure = &workflowsResource{}
)

// NewWorkflowsResource is a helper function to simplify the provider implementation.
func NewWorkflowsResource() resource.Resource {
    return &workflowsResource{}
}

// workflowsResource is the resource implementation.
type workflowsResource struct{
	client *client.Client
}

type workflowsResourceModel struct{
	Workflows []workflowsModel `tfsdk:"workflows"`
}


// Configure adds the provider configured client to the resource.
func (r *workflowsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
			return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
			resp.Diagnostics.AddError(
					"Unexpected Data Source Configure Type",
					fmt.Sprintf("Expected *Keep.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
			return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *workflowsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workflows"
}

// Schema defines the schema for the resource.
func (d *workflowsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
				"workflows": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
										"yaml": schema.StringAttribute{
												Required: true,
										},
								},
						},
				},
		},
}
}

// Create creates the resource and sets the initial Terraform state.
func (r *workflowsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Got to create")
	var plan workflowsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API Request body from plan
	for _, workflow := range plan.Workflows {
		workflowYaml := workflow.Yaml

		_, err := r.client.PostWorkflow(workflowYaml)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Workflow",
				"Could not create workflow, unexpected  error: "+err.Error(),
			)
			return
		
		}
		//plan.Workflows[index].Yaml = receive_workflow.Yaml
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *workflowsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workflowsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get workflows from keep
	workflows, err := r.client.GetWorkflows()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workflows",
			"Error is " + err.Error(),
		)
		return
	}
	state.Workflows = []workflowsModel{}
	for _, workflow := range workflows {
		state.Workflows = append(state.Workflows, workflowsModel{workflow.Yaml})

	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *workflowsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan workflowsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API Request from plan
	var workflows []workflowsModel
	for _, workflow := range plan.Workflows {
		workflows = append(workflows, workflow)
	}

	// Update existing Workflows
	for _, workflow := range workflows {
		_, err := r.client.PostWorkflow(workflow.Yaml)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Workflow",
				"Could not update workflow, unexpected order: "+err.Error(),
			)
			return
		}
	}
	fetchedWorkflows, err := r.client.GetWorkflows()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workflows",
			"Could not read workflows, error is: "+err.Error(),
		)
		return
	}

	plan.Workflows = []workflowsModel{}
	for _, workflow := range fetchedWorkflows {
		plan.Workflows = append(plan.Workflows, workflowsModel{workflow.Yaml})
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *workflowsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Got to delete")
	var state workflowsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	for _, workflow := range state.Workflows {
		tflog.Info(ctx, "workflow is "+workflow.Yaml)
		id, err := r.client.DeleteWorkflow(workflow.Yaml)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Workflows",
				"Could not delete workflows, err is "+ err.Error(),
			)
			return
		}
		tflog.Info(ctx, id)
	}
}
