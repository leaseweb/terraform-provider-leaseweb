package public_cloud_service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/ports"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

var (
	_ ports.PublicCloudRepository = &repositorySpy{}
)

type repositorySpy struct {
	instances              entity.Instances
	instance               *entity.Instance
	autoScalingGroup       *entity.AutoScalingGroup
	loadBalancer           *entity.LoadBalancer
	availableInstanceTypes entity.InstanceTypes
	regions                entity.Regions

	getAutoScalingGroupError                error
	getLoadBalancerError                    error
	getAllInstancesError                    error
	getInstanceError                        error
	createInstanceError                     error
	updateInstanceError                     error
	deleteInstanceError                     error
	getAvailableInstanceTypesForUpdateError error
	getRegionsError                         error
}

func (r repositorySpy) GetRegions(ctx context.Context) (entity.Regions, error) {
	return r.regions, r.getRegionsError
}

func (r repositorySpy) GetAvailableInstanceTypesForUpdate(
	id value_object.Uuid,
	ctx context.Context,
) (entity.InstanceTypes, error) {
	return r.availableInstanceTypes, r.getAvailableInstanceTypesForUpdateError
}

func (r repositorySpy) GetAutoScalingGroup(
	id value_object.Uuid,
	ctx context.Context,
) (*entity.AutoScalingGroup, error) {
	return r.autoScalingGroup, r.getAutoScalingGroupError
}

func (r repositorySpy) GetLoadBalancer(
	id value_object.Uuid,
	ctx context.Context,
) (*entity.LoadBalancer, error) {
	return r.loadBalancer, r.getLoadBalancerError
}

func (r repositorySpy) GetAllInstances(ctx context.Context) (
	entity.Instances,
	error,
) {
	return r.instances, r.getAllInstancesError
}

func (r repositorySpy) GetInstance(
	id value_object.Uuid,
	ctx context.Context,
) (*entity.Instance, error) {
	return r.instance, r.getInstanceError
}

func (r repositorySpy) CreateInstance(
	instance entity.Instance,
	ctx context.Context,
) (*entity.Instance, error) {
	return r.instance, r.createInstanceError
}

func (r repositorySpy) UpdateInstance(
	instance entity.Instance,
	ctx context.Context,
) (*entity.Instance, error) {
	return r.instance, r.updateInstanceError
}

func (r repositorySpy) DeleteInstance(
	id value_object.Uuid,
	ctx context.Context,
) error {
	return r.deleteInstanceError
}

func TestService_GetAllInstances(t *testing.T) {
	t.Run(
		"service passes back instances from repository",
		func(t *testing.T) {
			want := entity.Instances{{Id: value_object.NewGeneratedUuid()}}

			service := New(repositorySpy{instances: want})

			got, err := service.GetAllInstances(context.TODO())

			assert.NoError(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		service := New(
			repositorySpy{getAllInstancesError: errors.New("some error")},
		)

		_, err := service.GetAllInstances(context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_GetInstance(t *testing.T) {
	t.Run(
		"service passes back instance from repository",
		func(t *testing.T) {
			want := entity.Instance{}

			service := New(repositorySpy{instance: &want})

			got, err := service.GetInstance(
				value_object.NewGeneratedUuid(),
				context.TODO(),
			)

			assert.NoError(t, err)
			assert.Same(t, &want, got)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		service := New(
			repositorySpy{getInstanceError: errors.New("some error")},
		)

		_, err := service.GetInstance(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_CreateInstance(t *testing.T) {
	t.Run("service passes back instance from repository", func(t *testing.T) {
		want := entity.Instance{}

		service := New(repositorySpy{instance: &want})

		got, err := service.CreateInstance(entity.Instance{}, context.TODO())

		assert.NoError(t, err)
		assert.Same(t, &want, got)
	})

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(
			repositorySpy{createInstanceError: errors.New("some error")},
		)

		_, err := instanceService.CreateInstance(entity.Instance{}, context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_UpdateInstance(t *testing.T) {
	t.Run("service passes back instance from repository", func(t *testing.T) {
		want := entity.Instance{}

		service := New(repositorySpy{instance: &want})

		got, err := service.UpdateInstance(entity.Instance{}, context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, &want, got)
	})

	t.Run("service passes back error from repository", func(t *testing.T) {
		service := New(
			repositorySpy{updateInstanceError: errors.New("some error")},
		)

		_, err := service.UpdateInstance(
			entity.Instance{},
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_DeleteInstance(t *testing.T) {
	t.Run("service passes back nil from repository", func(t *testing.T) {
		service := New(repositorySpy{})

		err := service.DeleteInstance(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

		assert.NoError(t, err)
	})

	t.Run("service passes back error from repository", func(t *testing.T) {
		service := New(
			repositorySpy{deleteInstanceError: errors.New("some error")},
		)

		err := service.DeleteInstance(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_GetAvailableInstanceTypesForUpdate(t *testing.T) {
	t.Run(
		"expected instance types returned from repository",
		func(t *testing.T) {
			want := entity.InstanceTypes{{Name: "tralala"}}
			spy := repositorySpy{availableInstanceTypes: want}

			service := New(spy)
			got, err := service.GetAvailableInstanceTypesForUpdate(
				value_object.NewGeneratedUuid(),
				context.TODO(),
			)

			assert.NoError(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		spy := repositorySpy{
			getAvailableInstanceTypesForUpdateError: errors.New("some error"),
		}

		service := New(spy)
		_, err := service.GetAvailableInstanceTypesForUpdate(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_GetRegions(t *testing.T) {
	t.Run(
		"expected regions returned from repository",
		func(t *testing.T) {
			want := entity.Regions{{Name: "tralala"}}
			spy := repositorySpy{regions: want}

			service := New(spy)
			got, err := service.GetRegions(context.TODO())

			assert.NoError(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		spy := repositorySpy{getRegionsError: errors.New("some error")}

		service := New(spy)
		_, err := service.GetRegions(context.TODO())

		assert.ErrorContains(t, err, "some error")
	})
}
