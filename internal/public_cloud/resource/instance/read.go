package instance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
	"terraform-provider-leaseweb/internal/utils"
)

func (i *instanceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state model.Instance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance, sdkResponse, err := i.client.PublicCloudClient.PublicCloudAPI.GetInstance(
		i.client.AuthContext(ctx),
		state.Id.ValueString(),
	).Execute()
	if err != nil {
		utils.HandleError(
			ctx,
			sdkResponse,
			&resp.Diagnostics,
			"Error Reading Public Cloud Instance "+state.Id.ValueString(),
			err.Error(),
		)
		return
	}

	var autoScalingGroupDetails *publicCloud.AutoScalingGroupDetails
	var autoScalingGroupDetailsResponse *http.Response

	// Get autoScalingGroup details for each Instance as the instanceDetails
	// endpoint is missing loadBalancer data.
	sdkAutoScalingGroup, _ := instance.GetAutoScalingGroupOk()
	if sdkAutoScalingGroup != nil {
		autoScalingGroupDetails, autoScalingGroupDetailsResponse, err = i.client.PublicCloudClient.PublicCloudAPI.GetAutoScalingGroup(
			i.client.AuthContext(ctx),
			sdkAutoScalingGroup.GetId(),
		).Execute()
		if err != nil {
			utils.HandleError(
				ctx,
				autoScalingGroupDetailsResponse,
				&resp.Diagnostics,
				fmt.Sprintf("Unable to Read Leaseweb Public Cloud Instance %v", instance.GetId()),
				err.Error(),
			)
			return
		}
	}

	state.Populate(instance, autoScalingGroupDetails, ctx)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
