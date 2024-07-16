package instance

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/modify_plan"
)

func (i *instanceResource) ModifyPlan(
	ctx context.Context,
	request resource.ModifyPlanRequest,
	response *resource.ModifyPlanResponse,
) {
	planInstance := model.Instance{}
	request.Plan.Get(ctx, &planInstance)

	stateInstance := model.Instance{}
	request.State.Get(ctx, &stateInstance)

	err := i.validateInstanceType(
		ctx,
		stateInstance.Id,
		stateInstance.Type,
		planInstance.Type,
		response,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = i.validateRegion(ctx, response, planInstance.Region.ValueString())
	if err != nil {
		log.Fatal(err)
	}
}

func (i *instanceResource) validateInstanceType(
	ctx context.Context,
	stateId types.String,
	stateType types.String,
	planType types.String,
	response *resource.ModifyPlanResponse,
) error {
	typeValidator := modify_plan.NewTypeValidator(
		stateId,
		stateType,
		planType,
	)

	hasTypeChanged := typeValidator.HashTypeChanged()

	if !hasTypeChanged {
		return nil
	}

	uuid, err := value_object.NewUuid(stateId.ValueString())
	if err != nil {
		return fmt.Errorf("validateInstanceType: %w", err)
	}

	allowedInstanceTypes, err := i.client.PublicCloud.GetAvailableInstanceTypesForUpdate(*uuid, ctx)

	if err != nil {
		return fmt.Errorf("validateInstanceType: %w", err)
	}

	if typeValidator.IsTypeValid(allowedInstanceTypes) {
		return nil
	}

	response.Diagnostics.AddAttributeError(
		path.Root("type"),
		"Invalid Instance Type",
		fmt.Sprintf(
			"Allowed types are %v",
			allowedInstanceTypes,
		),
	)

	return nil
}

func (i *instanceResource) validateRegion(
	ctx context.Context,
	response *resource.ModifyPlanResponse,
	region string,
) error {
	// Region has not changed here.
	if region == "" {
		return nil
	}

	regions, err := i.client.PublicCloud.GetRegions(ctx)

	if err != nil {
		return fmt.Errorf("validateRegion: %w", err)
	}

	if regions.Contains(region) {
		return nil
	}

	response.Diagnostics.AddAttributeError(
		path.Root("region"),
		"Invalid Region",
		fmt.Sprintf(
			"Allowed regions are %v",
			regions,
		),
	)

	return nil
}
