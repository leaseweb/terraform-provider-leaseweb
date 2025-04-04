package dedicatedserver

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedserver/v2"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &installationResource{}
	_ resource.ResourceWithImportState = &installationResource{}
)

func NewInstallationResource() resource.Resource {
	return &installationResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "dedicated_server_installation",
		},
	}
}

type installationResource struct {
	utils.ResourceAPI
}

type installationResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	DedicatedServerID types.String   `tfsdk:"dedicated_server_id"`
	CallbackURL       types.String   `tfsdk:"callback_url"`
	ControlPanelID    types.String   `tfsdk:"control_panel_id"`
	Device            types.String   `tfsdk:"device"`
	Hostname          types.String   `tfsdk:"hostname"`
	OperatingSystemID types.String   `tfsdk:"operating_system_id"`
	Partitions        types.List     `tfsdk:"partitions"`
	Password          types.String   `tfsdk:"password"`
	PostInstallScript types.String   `tfsdk:"post_install_script"`
	PowerCycle        types.Bool     `tfsdk:"power_cycle"`
	Raid              types.Object   `tfsdk:"raid"`
	SSHKeys           []types.String `tfsdk:"ssh_keys"`
	Timezone          types.String   `tfsdk:"timezone"`
}

type raidResourceModel struct {
	Level         types.Int32  `tfsdk:"level"`
	NumberOfDisks types.Int32  `tfsdk:"number_of_disks"`
	Type          types.String `tfsdk:"type"`
}

type partitionsResourceModel struct {
	Filesystem types.String `tfsdk:"filesystem"`
	Mountpoint types.String `tfsdk:"mountpoint"`
	Size       types.String `tfsdk:"size"`
}

func (i *installationResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	raid := func() schema.SingleNestedAttribute {
		return schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"level": schema.Int32Attribute{
					Description: "RAID level to apply to your installation, this value is only required if you specify a type HW or SW. Valid options are \n  - *0*\n  - *1*\n  - *5*\n  - *10*\n",
					Optional:    true,
					Validators: []validator.Int32{
						int32validator.OneOf([]int32{0, 1, 5, 10}...),
					},
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.RequiresReplace(),
					},
				},
				"number_of_disks": schema.Int32Attribute{
					Description: "The number of disks you want to apply RAID on. If not specified all disks are used",
					Optional:    true,
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.RequiresReplace(),
					},
				},
				"type": schema.StringAttribute{
					Description: "RAID type to apply to your installation. NONE is the equivalent of pass-through mode on HW RAID equipped servers. Valid options are \n  - *HW*\n  - *SW*\n  - *NONE*\n",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"HW", "SW", "NONE"}...),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
		}
	}

	partitions := func() schema.ListNestedAttribute {
		return schema.ListNestedAttribute{
			Optional: true,
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"filesystem": schema.StringAttribute{
						Description: "File system in which partition would be mounted",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"mountpoint": schema.StringAttribute{
						Description: "The partition mount point (eg /home). Mandatory for the root partition (/) and not intended to be used in swap partition",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"size": schema.StringAttribute{
						Description: "Size of the partition (Normally in MB, but this is OS-specific)",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
		}
	}

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the installation job",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dedicated_server_id": schema.StringAttribute{
				Description: "The ID of a server",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"callback_url": schema.StringAttribute{
				Description: "Url which will receive callbacks when the installation is finished or failed",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"control_panel_id": schema.StringAttribute{
				Description: "Control panel identifier",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device": schema.StringAttribute{
				Description: `Block devices in a disk set in which the partitions will be installed. Supported values are any disk set id, ` + "`SATA_SAS`" + ` or ` + "`NVME`" + `.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hostname": schema.StringAttribute{
				Description: "Hostname to be used in your installation",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"operating_system_id": schema.StringAttribute{
				Description: "Operating system identifier",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"partitions": partitions(),
			"password": schema.StringAttribute{
				Description: "Server root password. If not provided, it would be automatically generated",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"post_install_script": schema.StringAttribute{
				Description: "A valid bash script to run right after the installation.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"power_cycle": schema.BoolAttribute{
				Description: "If true, allows system reboots to happen automatically within the process. Otherwise, you should do them manually",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"raid": raid(),
			"ssh_keys": schema.SetAttribute{
				Description: "List of public sshKeys to be setup in your installation",
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"timezone": schema.StringAttribute{
				Description: "Timezone represented as Geographical_Area/City",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}

	utils.AddUnsupportedActionsNotation(
		resp,
		[]utils.Action{utils.ReadAction, utils.UpdateAction, utils.DeleteAction},
	)
}

func (i *installationResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan installationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	// Extract the Raid configuration from the plan
	var raidPlan raidResourceModel
	plan.Raid.As(ctx, &raidPlan, basetypes.ObjectAsOptions{})

	// Extract Partitions configuration from the plan into a Go slice
	partitionsPlan := make([]partitionsResourceModel, 0, len(plan.Partitions.Elements()))
	plan.Partitions.ElementsAs(ctx, &partitionsPlan, false)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preparing partitions for the installation options.
	var partitions []dedicatedserver.Partition
	if !plan.Partitions.IsNull() && !plan.Partitions.IsUnknown() {

		for _, p := range partitionsPlan {
			// Check if all fields are either null or unknown, if so, skip
			if utils.AdaptStringPointerValueToNullableString(p.Mountpoint) == nil &&
				utils.AdaptStringPointerValueToNullableString(p.Size) == nil &&
				utils.AdaptStringPointerValueToNullableString(p.Filesystem) == nil {
				continue
			}

			partitions = append(partitions, dedicatedserver.Partition{
				Filesystem: utils.AdaptStringPointerValueToNullableString(p.Filesystem),
				Size:       utils.AdaptStringPointerValueToNullableString(p.Size),
				Mountpoint: utils.AdaptStringPointerValueToNullableString(p.Mountpoint),
			})
		}

	}

	// Preparing RAID configuration for the installation options
	var raid *dedicatedserver.Raid
	// Check that at least one RAID field is set before initializing the RAID struct.
	if !plan.Raid.IsNull() && !plan.Raid.IsUnknown() &&
		(utils.AdaptInt32PointerValueToNullableInt32(raidPlan.Level) != nil ||
			utils.AdaptInt32PointerValueToNullableInt32(raidPlan.NumberOfDisks) != nil ||
			utils.AdaptStringPointerValueToNullableString(raidPlan.Type) != nil) {

		raid = &dedicatedserver.Raid{
			Level:         (*dedicatedserver.RaidLevel)(utils.AdaptInt32PointerValueToNullableInt32(raidPlan.Level)),
			NumberOfDisks: utils.AdaptInt32PointerValueToNullableInt32(raidPlan.NumberOfDisks),
			Type:          (*dedicatedserver.RaidType)(utils.AdaptStringPointerValueToNullableString(raidPlan.Type)),
		}
	}

	// Preparing SSH keys for the installation options, combining them into a single string
	var SSHKeysList []string
	for _, k := range plan.SSHKeys {
		if utils.AdaptStringPointerValueToNullableString(k) != nil {
			SSHKeysList = append(SSHKeysList, k.ValueString())
		}
	}
	SSHKeys := strings.Join(SSHKeysList, "\n")

	opts := dedicatedserver.NewInstallOperatingSystemOpts(plan.OperatingSystemID.ValueString())
	opts.CallbackUrl = utils.AdaptStringPointerValueToNullableString(plan.CallbackURL)
	opts.ControlPanelId = utils.AdaptStringPointerValueToNullableString(plan.ControlPanelID)
	opts.Device = utils.AdaptStringPointerValueToNullableString(plan.Device)
	opts.Hostname = utils.AdaptStringPointerValueToNullableString(plan.Hostname)
	opts.Partitions = partitions
	opts.Password = utils.AdaptStringPointerValueToNullableString(plan.Password)
	opts.PostInstallScript = utils.AdaptStringValueToNullableString(base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(plan.PostInstallScript.ValueString()))))
	opts.PowerCycle = utils.AdaptBoolPointerValueToNullableBool(plan.PowerCycle)
	opts.Raid = raid
	opts.Timezone = utils.AdaptStringPointerValueToNullableString(plan.Timezone)
	if len(SSHKeysList) > 0 {
		opts.SshKeys = &SSHKeys
	}

	serverID := plan.DedicatedServerID.ValueString()
	installationJob, response, err := i.DedicatedserverAPI.InstallOperatingSystem(ctx, serverID).
		InstallOperatingSystemOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	job, err := i.waitForJobAndRetrieveUntilFinished(serverID, installationJob.GetUuid(), ctx, resp)
	if err != nil {
		utils.ReportError(err.Error(), &resp.Diagnostics)
		return
	}
	plan.ID = types.StringValue(job.GetUuid())

	diags := i.syncResourceModelWithSDK(&plan, job.GetPayload(), ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (i *installationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("dedicated_server_id"), req, resp)
}

func (i *installationResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state installationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := state.DedicatedServerID.ValueString()

	result, response, err := i.DedicatedserverAPI.GetJobList(ctx, serverID).
		Offset(0).Limit(1).Type_("install").Status("FINISHED").Execute()

	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	jobs := result.GetJobs()

	if len(jobs) == 0 {
		utils.ReportError(fmt.Sprintf("No installation jobs found for server %s", serverID), &resp.Diagnostics)
		return
	}
	job := jobs[0]
	state.ID = types.StringValue(job.GetUuid())

	diags := i.syncResourceModelWithSDK(&state, job.GetPayload(), ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (i *installationResource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	_ *resource.UpdateResponse,
) {
}

func (i *installationResource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
}

func (i *installationResource) getJob(serverID, jobID string, ctx context.Context, resp *resource.CreateResponse) *dedicatedserver.CurrentJob {
	// Fetch the job
	result, response, err := i.DedicatedserverAPI.GetJob(ctx, serverID, jobID).Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
	}

	return result
}

// waitForJobAndRetrieveUntilFinished handles polling with retry and timeout.
func (i *installationResource) waitForJobAndRetrieveUntilFinished(serverID, jobID string, ctx context.Context, resp *resource.CreateResponse) (*dedicatedserver.CurrentJob, error) {
	var job *dedicatedserver.CurrentJob
	// Create a constant backoff with a 30-second retry interval
	bo := backoff.NewConstantBackOff(30 * time.Second)

	// Set the retry limit to 120 retries (60 minutes total)
	retryCount := 0
	maxRetries := 120

	// Start polling and retrying
	for {
		if retryCount >= maxRetries {
			return nil, errors.New("timed out waiting for job to finish after 60 minutes")
		}

		// Call the function to get job
		job = i.getJob(serverID, jobID, ctx, resp)

		// check if the job is finished
		if job.GetStatus() == "FINISHED" {
			// Job is finished, exit the loop
			return job, nil
		}

		// Sleep for the backoff interval before retrying
		time.Sleep(bo.NextBackOff())
		retryCount++
	}
}

func (i *installationResource) syncResourceModelWithSDK(
	state *installationResourceModel,
	payload dedicatedserver.ServerJobPayload,
	ctx context.Context,
) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Device = types.StringValue(payload.GetDevice())
	state.OperatingSystemID = types.StringValue(payload.GetOperatingSystemId())
	state.PowerCycle = types.BoolValue(payload.GetPowerCycle())
	state.Timezone = types.StringValue(payload.GetTimezone())
	state.Hostname = types.StringValue(payload.GetHostname())

	partitionAttributeTypes := map[string]attr.Type{
		"filesystem": types.StringType,
		"mountpoint": types.StringType,
		"size":       types.StringType,
	}

	// Preparing and converting partitions into types.Object to store in the state
	var partitionsObjects []attr.Value
	for _, p := range payload.GetPartitions() {
		partition := partitionsResourceModel{
			Filesystem: types.StringValue(p.GetFilesystem()),
			Mountpoint: types.StringValue(p.GetMountpoint()),
			Size:       types.StringValue(p.GetSize()),
		}

		partitionObj, partitionDiags := types.ObjectValueFrom(ctx, partitionAttributeTypes, partition)
		diags.Append(partitionDiags...)
		partitionsObjects = append(partitionsObjects, partitionObj)
	}

	// Convert the slice of partition objects to a types.List and store it in the plan
	partitionsList, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: partitionAttributeTypes}, partitionsObjects)
	diags.Append(listDiags...)
	state.Partitions = partitionsList

	if payload.HasRaid() {
		raidAttributeTypes := map[string]attr.Type{
			"level":           types.Int32Type,
			"number_of_disks": types.Int32Type,
			"type":            types.StringType,
		}

		payloadRaid := payload.GetRaid()

		raid := raidResourceModel{
			Level:         types.Int32Value(payloadRaid.GetLevel()),
			NumberOfDisks: types.Int32Value(payloadRaid.GetNumberOfDisks()),
			Type:          types.StringValue(payloadRaid.GetType()),
		}

		raidObj, _ := types.ObjectValueFrom(ctx, raidAttributeTypes, raid)
		state.Raid = raidObj
	}

	return diags
}
