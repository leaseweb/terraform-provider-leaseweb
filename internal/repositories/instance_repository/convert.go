package instance_repository

import (
	"fmt"
	"time"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

type errLoadBalancerNotFound struct {
	msg string
}

func (e errLoadBalancerNotFound) Error() string {
	return e.msg
}

func convertInstance(
	sdkInstance publicCloud.InstanceDetails,
	autoScalingGroup *domain.AutoScalingGroup,
) (*domain.Instance, error) {
	instanceId, err := value_object.NewUuid(sdkInstance.GetId())
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstance: %w",
			err,
		)
	}

	image, err := convertImage(sdkInstance.GetImage())
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstance: %w",
			err,
		)
	}

	state, err := enum.NewState(string(sdkInstance.GetState()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstance: %w",
			err,
		)
	}

	rootDiskSize, err := value_object.NewRootDiskSize(int(sdkInstance.GetRootDiskSize()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstance: %w",
			err,
		)
	}

	rootDiskStorageType, err := enum.NewRootDiskStorageType(
		string(sdkInstance.GetRootDiskStorageType()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstance: %w",
			err,
		)
	}

	ips, err := convertIps(sdkInstance.GetIps())
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstance:  %w",
			err,
		)
	}

	contract, err := convertContract(sdkInstance.GetContract())
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstance:  %w",
			err,
		)
	}

	optionalValues := domain.OptionalInstanceValues{
		Reference:        convertNullableStringToValue(sdkInstance.Reference),
		MarketAppId:      convertNullableStringToValue(sdkInstance.MarketAppId),
		StartedAt:        convertNullableTimeToValue(sdkInstance.StartedAt),
		AutoScalingGroup: autoScalingGroup,
	}
	if sdkInstance.Iso.Get() != nil {
		iso := convertIso(*sdkInstance.Iso.Get())
		optionalValues.Iso = &iso
	}
	if sdkInstance.PrivateNetwork.Get() != nil {
		privateNetwork := convertPrivateNetwork(*sdkInstance.PrivateNetwork.Get())
		optionalValues.PrivateNetwork = &privateNetwork
	}

	instance := domain.NewInstance(
		*instanceId,
		sdkInstance.GetRegion(),
		convertResources(sdkInstance.GetResources()),
		*image,
		state,
		sdkInstance.GetProductType(),
		sdkInstance.GetHasPublicIpV4(),
		sdkInstance.GetIncludesPrivateNetwork(),
		*rootDiskSize,
		string(sdkInstance.GetType()),
		rootDiskStorageType,
		ips,
		*contract,
		optionalValues,
	)

	return &instance, nil
}

func convertResources(sdkResources publicCloud.Resources) domain.Resources {
	resources := domain.NewResources(
		convertCpu(sdkResources.GetCpu()),
		convertMemory(sdkResources.GetMemory()),
		convertNetworkSpeed(sdkResources.GetPublicNetworkSpeed()),
		convertNetworkSpeed(sdkResources.GetPrivateNetworkSpeed()),
	)

	return resources
}

func convertCpu(sdkCpu publicCloud.Cpu) domain.Cpu {
	return domain.NewCpu(int(sdkCpu.GetValue()), sdkCpu.GetUnit())
}

func convertMemory(sdkMemory publicCloud.Memory) domain.Memory {
	return domain.NewMemory(float64(sdkMemory.GetValue()), sdkMemory.GetUnit())
}

func convertNetworkSpeed(sdkNetworkSpeed publicCloud.NetworkSpeed) domain.NetworkSpeed {
	return domain.NewNetworkSpeed(
		int(sdkNetworkSpeed.GetValue()),
		sdkNetworkSpeed.GetUnit(),
	)
}

func convertImage(sdkImage publicCloud.ImageDetails) (*domain.Image, error) {
	imageId, err := enum.NewImageId(string(sdkImage.GetId()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertImage: %w",
			err,
		)
	}

	image := domain.NewImage(
		imageId,
		sdkImage.GetName(),
		sdkImage.GetVersion(),
		sdkImage.GetFamily(),
		sdkImage.GetFlavour(),
		sdkImage.GetArchitecture(),
		sdkImage.GetMarketApps(),
		sdkImage.GetStorageTypes(),
	)

	return &image, nil
}

func convertIps(sdkIps []publicCloud.IpDetails) (domain.Ips, error) {
	var ips domain.Ips
	for _, sdkIp := range sdkIps {
		ip, err := convertIp(sdkIp)
		if err != nil {
			return nil, fmt.Errorf("convertIps: %w", err)
		}
		ips = append(ips, *ip)
	}

	return ips, nil
}

func convertIp(sdkIp publicCloud.IpDetails) (*domain.Ip, error) {
	networkType, err := enum.NewNetworkType(string(sdkIp.GetNetworkType()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertIp: %w",
			err,
		)
	}

	optionalIpValues := domain.OptionalIpValues{
		ReverseLookup: convertNullableStringToValue(sdkIp.ReverseLookup),
	}

	sdkDdos, _ := sdkIp.GetDdosOk()
	if sdkDdos != nil {
		ddos := convertDdos(*sdkDdos)
		optionalIpValues.Ddos = &ddos
	}

	ip := domain.NewIp(
		sdkIp.GetIp(),
		sdkIp.GetPrefixLength(),
		int(sdkIp.GetVersion()),
		sdkIp.GetNullRouted(),
		sdkIp.GetMainIp(),
		networkType,
		optionalIpValues,
	)

	return &ip, nil
}

func convertDdos(sdkDdos publicCloud.Ddos) domain.Ddos {
	return domain.NewDdos(
		sdkDdos.GetDetectionProfile(),
		sdkDdos.GetProtectionType(),
	)
}

func convertContract(sdkContract publicCloud.Contract) (*domain.Contract, error) {
	billingFrequency, err := enum.NewContractBillingFrequency(
		int(sdkContract.GetBillingFrequency()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"convertContract: %w",
			err,
		)
	}

	contractTerm, err := enum.NewContractTerm(int(sdkContract.GetTerm()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertContract: %w",
			err,
		)
	}

	contractType, err := enum.NewContractType(string(sdkContract.GetType()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertContract: %w",
			err,
		)
	}

	contractState, err := enum.NewContractState(string(sdkContract.GetState()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertContract: %w",
			err,
		)
	}

	contract, err := domain.NewContract(
		billingFrequency,
		contractTerm,
		contractType,
		sdkContract.GetRenewalsAt(),
		sdkContract.GetCreatedAt(),
		contractState,
		convertNullableTimeToValue(sdkContract.EndsAt),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"convertContract: %w",
			err,
		)
	}

	return contract, nil
}

func convertIso(sdkIso publicCloud.Iso) domain.Iso {
	return domain.NewIso(sdkIso.GetId(), sdkIso.GetName())
}

func convertPrivateNetwork(sdkPrivateNetwork publicCloud.PrivateNetwork) domain.PrivateNetwork {
	return domain.PrivateNetwork{
		Id:     sdkPrivateNetwork.GetPrivateNetworkId(),
		Status: sdkPrivateNetwork.GetStatus(),
		Subnet: sdkPrivateNetwork.GetSubnet(),
	}
}

func convertAutoScalingGroup(
	sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
	loadBalancer *domain.LoadBalancer,
) (
	*domain.AutoScalingGroup,
	error,
) {
	if sdkAutoScalingGroup.LoadBalancer.Get() != nil && loadBalancer == nil {
		return nil, errLoadBalancerNotFound{msg: fmt.Sprintf(
			"required loadBalacner %q linked to autoScalingGroup %q has not been passed",
			sdkAutoScalingGroup.LoadBalancer.Get().GetId(),
			sdkAutoScalingGroup.GetId(),
		)}
	}

	autoScalingGroupId, err := value_object.NewUuid(sdkAutoScalingGroup.GetId())
	if err != nil {
		return nil, fmt.Errorf(
			"convertAutoScalingGroup: %w",
			err,
		)
	}

	autoScalingGroupType, err := enum.NewAutoScalingGroupType(
		string(sdkAutoScalingGroup.GetType()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"convertAutoScalingGroup: %w",
			err,
		)
	}

	state, err := enum.NewAutoScalingGroupState(string(sdkAutoScalingGroup.GetState()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertAutoScalingGroup: %w",
			err,
		)
	}

	reference, err := value_object.NewAutoScalingGroupReference(sdkAutoScalingGroup.GetReference())

	if err != nil {
		return nil, fmt.Errorf(
			"convertAutoScalingGroup: %w",
			err,
		)
	}

	options := domain.AutoScalingGroupOptions{
		DesiredAmount: convertNullableInt32ToValue(sdkAutoScalingGroup.DesiredAmount),
		MinimumAmount: convertNullableInt32ToValue(sdkAutoScalingGroup.MinimumAmount),
		MaximumAmount: convertNullableInt32ToValue(sdkAutoScalingGroup.MaximumAmount),
		CpuThreshold:  convertNullableInt32ToValue(sdkAutoScalingGroup.CpuThreshold),
		CoolDownTime:  convertNullableInt32ToValue(sdkAutoScalingGroup.CooldownTime),
		StartsAt:      convertNullableTimeToValue(sdkAutoScalingGroup.StartsAt),
		EndsAt:        convertNullableTimeToValue(sdkAutoScalingGroup.EndsAt),
		WarmupTime:    convertNullableInt32ToValue(sdkAutoScalingGroup.WarmupTime),
		LoadBalancer:  loadBalancer,
	}

	autoScalingGroup := domain.NewAutoScalingGroup(
		*autoScalingGroupId,
		autoScalingGroupType,
		state,
		sdkAutoScalingGroup.GetRegion(),
		*reference,
		sdkAutoScalingGroup.GetCreatedAt(),
		sdkAutoScalingGroup.GetUpdatedAt(),
		options,
	)

	return &autoScalingGroup, nil
}

func convertNullableStringToValue(nullableString publicCloud.NullableString) *string {
	return nullableString.Get()
}

func convertNullableTimeToValue(nullableTime publicCloud.NullableTime) *time.Time {
	return nullableTime.Get()
}

func convertNullableInt32ToValue(nullableInt publicCloud.NullableInt32) *int {
	if nullableInt.Get() == nil {
		return nil
	}

	value := int(*nullableInt.Get())
	return &value
}

func convertLoadBalancer(sdkLoadBalancer publicCloud.LoadBalancerDetails) (
	*domain.LoadBalancer,
	error,
) {
	loadBalancerId, err := value_object.NewUuid(sdkLoadBalancer.Id)
	if err != nil {
		return nil, fmt.Errorf(
			"convertLoadBalancer: %w",
			err,
		)
	}

	state, err := enum.NewState(string(sdkLoadBalancer.GetState()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertLoadBalancer: %w",
			err,
		)
	}

	contract, err := convertContract(sdkLoadBalancer.GetContract())
	if err != nil {
		return nil, fmt.Errorf(
			"convertLoadBalancer: %w",
			err,
		)
	}

	ips, err := convertIps(sdkLoadBalancer.GetIps())
	if err != nil {
		return nil, fmt.Errorf(
			"convertLoadBalancer:  %w",
			err,
		)
	}

	options := domain.OptionalLoadBalancerValues{
		Reference: convertNullableStringToValue(sdkLoadBalancer.Reference),
		StartedAt: convertNullableTimeToValue(sdkLoadBalancer.StartedAt),
	}

	if sdkLoadBalancer.Configuration.Get() != nil {
		configuration, err := convertLoadBalancerConfiguration(sdkLoadBalancer.GetConfiguration())
		if err != nil {
			return nil, fmt.Errorf(
				"convertLoadBalancer:  %w",
				err,
			)
		}
		options.Configuration = configuration
	}

	if sdkLoadBalancer.PrivateNetwork.Get() != nil {
		privateNetwork := convertPrivateNetwork(*sdkLoadBalancer.PrivateNetwork.Get())
		options.PrivateNetwork = &privateNetwork
	}

	loadBalancer := domain.NewLoadBalancer(
		*loadBalancerId,
		sdkLoadBalancer.GetType(),
		convertResources(sdkLoadBalancer.GetResources()),
		sdkLoadBalancer.GetRegion(),
		state,
		*contract,
		ips,
		options,
	)

	return &loadBalancer, nil
}

func convertLoadBalancerConfiguration(sdkLoadBalancerConfiguration publicCloud.LoadBalancerConfiguration) (
	*domain.LoadBalancerConfiguration,
	error,
) {
	balance, err := enum.NewBalance(string(sdkLoadBalancerConfiguration.GetBalance()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertLoadBalancerConfiguration: %w",
			err,
		)
	}

	options := domain.OptionalLoadBalancerConfigurationOptions{
		HealthCheck: nil,
	}
	if sdkLoadBalancerConfiguration.StickySession.Get() != nil {
		stickySession := convertStickySession(*sdkLoadBalancerConfiguration.StickySession.Get())
		options.StickySession = &stickySession
	}
	if sdkLoadBalancerConfiguration.HealthCheck.Get() != nil {
		healthCheck, err := convertHealthCheck(*sdkLoadBalancerConfiguration.HealthCheck.Get())
		if err != nil {
			return nil, fmt.Errorf(
				"convertLoadBalancerConfiguration: %w",
				err,
			)
		}

		options.HealthCheck = healthCheck
	}

	configuration := domain.NewLoadBalancerConfiguration(
		balance,
		sdkLoadBalancerConfiguration.GetXForwardedFor(),
		int(sdkLoadBalancerConfiguration.GetIdleTimeOut()),
		int(sdkLoadBalancerConfiguration.GetTargetPort()),
		options,
	)

	return &configuration, nil
}

func convertStickySession(sdkStickySession publicCloud.StickySession) domain.StickySession {
	return domain.NewStickySession(
		sdkStickySession.GetEnabled(),
		int(sdkStickySession.GetMaxLifeTime()),
	)
}

func convertHealthCheck(sdkHealthCheck publicCloud.HealthCheck) (
	*domain.HealthCheck,
	error,
) {
	method, err := enum.NewMethod(sdkHealthCheck.GetMethod())
	if err != nil {
		return nil, fmt.Errorf(
			"convertHealthCheck: %w",
			err,
		)
	}

	healthCheck := domain.NewHealthCheck(
		method,
		sdkHealthCheck.GetUri(),
		int(sdkHealthCheck.GetPort()),
		domain.OptionalHealthCheckValues{
			Host: convertNullableStringToValue(sdkHealthCheck.Host),
		},
	)

	return &healthCheck, nil
}

func convertEntityToLaunchInstanceOpts(instance domain.Instance) (
	*publicCloud.LaunchInstanceOpts,
	error,
) {
	instanceTypeName, err := publicCloud.NewInstanceTypeNameFromValue(
		instance.Type,
	)
	if err != nil {
		return nil, fmt.Errorf("convertEntityToLaunchInstanceOpts: %w", err)
	}

	rootDiskStorageType, err := publicCloud.NewRootDiskStorageTypeFromValue(
		instance.RootDiskStorageType.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("convertEntityToLaunchInstanceOpts: %w", err)
	}

	imageId, err := publicCloud.NewImageIdFromValue(instance.Image.Id.String())
	if err != nil {
		return nil, fmt.Errorf("convertEntityToLaunchInstanceOpts: %w", err)
	}

	contractType, err := publicCloud.NewContractTypeFromValue(
		instance.Contract.Type.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("convertEntityToLaunchInstanceOpts: %w", err)
	}

	contractTerm, err := publicCloud.NewContractTermFromValue(
		int32(instance.Contract.Term.Value()),
	)
	if err != nil {
		return nil, fmt.Errorf("convertEntityToLaunchInstanceOpts: %w", err)
	}

	billingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
		int32(instance.Contract.BillingFrequency.Value()),
	)
	if err != nil {
		return nil, fmt.Errorf("convertEntityToLaunchInstanceOpts: %w", err)
	}

	launchInstanceOpts := publicCloud.NewLaunchInstanceOpts(
		instance.Region,
		*instanceTypeName,
		*imageId,
		*contractType,
		*contractTerm,
		*billingFrequency,
		*rootDiskStorageType,
	)
	launchInstanceOpts.MarketAppId = instance.MarketAppId
	launchInstanceOpts.Reference = instance.Reference

	if instance.SshKey != nil {
		sshKey := instance.SshKey.String()
		launchInstanceOpts.SshKey = &sshKey
	}

	return launchInstanceOpts, nil
}

func convertEntityToUpdateInstanceOpts(instance domain.Instance) (
	*publicCloud.UpdateInstanceOpts,
	error,
) {
	updateInstanceOpts := publicCloud.NewUpdateInstanceOpts()
	updateInstanceOpts.Reference = instance.Reference

	if instance.RootDiskSize.Value != 0 {
		rootDiskSize := int32(instance.RootDiskSize.Value)
		updateInstanceOpts.RootDiskSize = &rootDiskSize
	}

	if instance.Type != "" {
		instanceTypeName, err := publicCloud.NewInstanceTypeNameFromValue(instance.Type)
		if err != nil {
			return nil, fmt.Errorf("convertEntityToUpdateInstanceOpts: %w", err)
		}
		updateInstanceOpts.Type = instanceTypeName
	}

	if instance.Contract.Type != "" {
		contractType, err := publicCloud.NewContractTypeFromValue(instance.Contract.Type.String())
		if err != nil {
			return nil, fmt.Errorf("convertEntityToUpdateInstanceOpts: %w", err)
		}
		updateInstanceOpts.ContractType = contractType
	}

	if instance.Contract.Term != 0 {
		contractTerm, err := publicCloud.NewContractTermFromValue(int32(instance.Contract.Term.Value()))
		if err != nil {
			return nil, fmt.Errorf("convertEntityToUpdateInstanceOpts: %w", err)
		}
		updateInstanceOpts.ContractTerm = contractTerm
	}

	if instance.Contract.BillingFrequency != 0 {
		billingFrequency, err := publicCloud.NewBillingFrequencyFromValue(int32(instance.Contract.BillingFrequency.Value()))
		if err != nil {
			return nil, fmt.Errorf("convertEntityToUpdateInstanceOpts: %w", err)
		}
		updateInstanceOpts.BillingFrequency = billingFrequency
	}

	return updateInstanceOpts, nil
}

func convertInstanceType(sdkInstanceType publicCloud.InstanceType) domain.InstanceType {
	return domain.NewInstanceType(sdkInstanceType.GetName())
}

func convertRegion(sdkRegion publicCloud.Region) domain.Region {
	return domain.NewRegion(sdkRegion.GetName(), sdkRegion.GetLocation())
}
