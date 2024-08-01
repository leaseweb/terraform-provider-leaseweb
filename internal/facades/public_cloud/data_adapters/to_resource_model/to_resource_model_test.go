package to_resource_model

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func Test_adaptImage(t *testing.T) {
	image := domain.NewImage(
		"UBUNTU_20_04_64BIT",
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		[]string{"one"},
		[]string{"storageType"},
	)

	got, err := adaptImage(context.TODO(), image)

	assert.NoError(t, err)

	assert.Equal(
		t,
		"UBUNTU_20_04_64BIT",
		got.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"name",
		got.Name.ValueString(),
		"name should be set",
	)
	assert.Equal(
		t,
		"version",
		got.Version.ValueString(),
		"version should be set",
	)
	assert.Equal(
		t,
		"family",
		got.Family.ValueString(),
		"family should be set",
	)
	assert.Equal(
		t,
		"flavour",
		got.Flavour.ValueString(),
		"flavour should be set",
	)
	assert.Equal(
		t,
		"architecture",
		got.Architecture.ValueString(),
		"architecture should be set",
	)

	var marketApps []string
	got.MarketApps.ElementsAs(context.TODO(), &marketApps, false)
	assert.Len(t, marketApps, 1)
	assert.Equal(
		t,
		"one",
		marketApps[0],
		"marketApps should be set",
	)

	var storageTypes []string
	got.StorageTypes.ElementsAs(context.TODO(), &storageTypes, false)
	assert.Len(t, storageTypes, 1)
	assert.Equal(
		t,
		"storageType",
		storageTypes[0],
		"storageTypes should be set",
	)
}

func Test_AdaptContract(t *testing.T) {
	endsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2023-12-14 17:09:47",
	)
	renewalsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2022-12-14 17:09:47",
	)
	createdAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2021-12-14 17:09:47",
	)

	contract, _ := domain.NewContract(
		enum.ContractBillingFrequencySix,
		enum.ContractTermThree,
		enum.ContractTypeMonthly,
		renewalsAt,
		createdAt,
		enum.ContractStateActive,
		&endsAt,
	)
	got, err := adaptContract(context.TODO(), *contract)

	assert.NoError(t, err)

	assert.Equal(
		t,
		int64(6),
		got.BillingFrequency.ValueInt64(),
		"billingFrequency should be set",
	)
	assert.Equal(
		t,
		int64(3),
		got.Term.ValueInt64(),
		"term should be set",
	)
	assert.Equal(
		t,
		"MONTHLY",
		got.Type.ValueString(),
		"type should be set",
	)
	assert.Equal(
		t,
		"2023-12-14 17:09:47 +0000 UTC",
		got.EndsAt.ValueString(),
		"endsAt should be set",
	)
	assert.Equal(
		t,
		"2022-12-14 17:09:47 +0000 UTC",
		got.RenewalsAt.ValueString(),
		"renewalsAt should be set",
	)
	assert.Equal(
		t,
		"2021-12-14 17:09:47 +0000 UTC",
		got.CreatedAt.ValueString(),
		"createdAt should be set",
	)
	assert.Equal(
		t,
		"ACTIVE",
		got.State.ValueString(),
		"state should be set",
	)
}

func Test_adaptPrivateNetwork(t *testing.T) {
	privateNetwork := domain.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)

	got, err := adaptPrivateNetwork(
		context.TODO(),
		privateNetwork,
	)

	assert.NoError(t, err)

	assert.Equal(
		t,
		"id",
		got.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"status",
		got.Status.ValueString(),
		"status should be set",
	)
	assert.Equal(
		t,
		"subnet",
		got.Subnet.ValueString(),
		"subnet should be set",
	)
}

func Test_adaptCpu(t *testing.T) {
	entityCpu := domain.NewCpu(1, "unit")
	got, err := adaptCpu(context.TODO(), entityCpu)

	assert.NoError(t, err)
	assert.Equal(
		t,
		int64(1),
		got.Value.ValueInt64(),
		"value should be set",
	)
	assert.Equal(
		t,
		"unit",
		got.Unit.ValueString(),
		"unit should be set",
	)
}

func Test_adaptMemory(t *testing.T) {
	memory := domain.NewMemory(1, "unit")

	got, err := adaptMemory(context.TODO(), memory)

	assert.NoError(t, err)
	assert.Equal(
		t,
		float64(1),
		got.Value.ValueFloat64(),
		"value should be set",
	)
	assert.Equal(
		t,
		"unit",
		got.Unit.ValueString(),
		"unit should be set",
	)
}

func Test_adaptNetworkSpeed(t *testing.T) {
	networkSpeed := domain.NewNetworkSpeed(1, "unit")

	got, err := adaptNetworkSpeed(context.TODO(), networkSpeed)

	assert.NoError(t, err)
	assert.Equal(
		t,
		int64(1),
		got.Value.ValueInt64(),
		"value should be set",
	)
	assert.Equal(
		t,
		"unit",
		got.Unit.ValueString(),
		"unit should be set",
	)
}

func Test_adaptResources(t *testing.T) {
	resources := domain.NewResources(
		domain.Cpu{Unit: "cpu"},
		domain.Memory{Unit: "memory"},
		domain.NetworkSpeed{Unit: "publicNetworkSpeed"},
		domain.NetworkSpeed{Unit: "privateNetworkSpeed"},
	)

	got, err := adaptResources(context.TODO(), resources)

	assert.NoError(t, err)

	cpu := model.Cpu{}
	got.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"cpu",
		cpu.Unit.ValueString(),
		"cpu should be set",
	)

	memory := model.Memory{}
	got.Memory.As(context.TODO(), &memory, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"memory",
		memory.Unit.ValueString(),
		"memory should be set",
	)

	publicNetworkSpeed := model.NetworkSpeed{}
	got.PublicNetworkSpeed.As(
		context.TODO(),
		&publicNetworkSpeed,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"publicNetworkSpeed",
		publicNetworkSpeed.Unit.ValueString(),
		"publicNetworkSpeed should be set",
	)

	privateNetworkSpeed := model.NetworkSpeed{}
	got.PrivateNetworkSpeed.As(
		context.TODO(),
		&privateNetworkSpeed,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"privateNetworkSpeed",
		privateNetworkSpeed.Unit.ValueString(),
		"privateNetworkSpeed should be set",
	)
}

func Test_adaptHealthCheck(t *testing.T) {
	host := "host"
	healthCheck := domain.NewHealthCheck(
		enum.MethodGet,
		"uri",
		22,
		domain.OptionalHealthCheckValues{Host: &host},
	)

	got, err := adaptHealthCheck(context.TODO(), healthCheck)

	assert.NoError(t, err)
	assert.Equal(t, "GET", got.Method.ValueString())
	assert.Equal(t, "uri", got.Uri.ValueString())
	assert.Equal(t, host, got.Host.ValueString())
	assert.Equal(t, int64(22), got.Port.ValueInt64())
}

func Test_adaptStickySession(t *testing.T) {
	stickySession := domain.NewStickySession(false, 1)

	got, err := adaptStickySession(context.TODO(), stickySession)

	assert.Nil(t, err)
	assert.False(t, got.Enabled.ValueBool())
	assert.Equal(t, int64(1), got.MaxLifeTime.ValueInt64())
}

func Test_adaptLoadBalancerConfiguration(t *testing.T) {

	loadBalancerConfiguration := domain.NewLoadBalancerConfiguration(
		enum.BalanceSource,
		false,
		5,
		6,
		domain.OptionalLoadBalancerConfigurationOptions{
			StickySession: &domain.StickySession{MaxLifeTime: 5},
			HealthCheck:   &domain.HealthCheck{Method: enum.MethodHead},
		},
	)

	got, err := adaptLoadBalancerConfiguration(
		context.TODO(),
		loadBalancerConfiguration,
	)

	assert.NoError(t, err)
	assert.Equal(t, "source", got.Balance.ValueString())
	assert.False(t, got.XForwardedFor.ValueBool())
	assert.Equal(t, int64(5), got.IdleTimeout.ValueInt64())
	assert.Equal(t, int64(6), got.TargetPort.ValueInt64())

	stickySession := model.StickySession{}
	got.StickySession.As(
		context.TODO(),
		&stickySession,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, int64(5), stickySession.MaxLifeTime.ValueInt64())

	healthCheck := model.HealthCheck{}
	got.HealthCheck.As(
		context.TODO(),
		&healthCheck,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, "HEAD", healthCheck.Method.ValueString())
}

func Test_adaptDdos(t *testing.T) {
	ddos := domain.NewDdos(
		"detectionProfile",
		"protectionType",
	)

	got, err := adaptDdos(context.TODO(), ddos)

	assert.NoError(t, err)

	assert.Equal(
		t,
		"detectionProfile",
		got.DetectionProfile.ValueString(),
		"detectionProfile should be set",
	)
	assert.Equal(
		t,
		"protectionType",
		got.ProtectionType.ValueString(),
		"protectionType should be set",
	)
}

func Test_adaptIp(t *testing.T) {
	reverseLookup := "reverse-lookup"

	ip := domain.NewIp(
		"1.2.3.4",
		"prefix-length",
		46,
		true,
		false,
		"tralala",
		domain.OptionalIpValues{
			Ddos:          &domain.Ddos{ProtectionType: "protection-type"},
			ReverseLookup: &reverseLookup,
		},
	)
	got, err := adaptIp(context.TODO(), ip)

	assert.NoError(t, err)

	assert.Equal(
		t,
		"1.2.3.4",
		got.Ip.ValueString(),
		"ip should be set",
	)
	assert.Equal(
		t,
		"prefix-length",
		got.PrefixLength.ValueString(),
		"prefix-length should be set",
	)
	assert.Equal(
		t,
		int64(46),
		got.Version.ValueInt64(),
		"version should be set",
	)
	assert.Equal(
		t,
		true,
		got.NullRouted.ValueBool(),
		"nullRouted should be set",
	)
	assert.Equal(
		t,
		false,
		got.MainIp.ValueBool(),
		"mainIp should be set",
	)
	assert.Equal(
		t,
		"tralala",
		got.NetworkType.ValueString(),
		"networkType should be set",
	)
	assert.Equal(
		t,
		"reverse-lookup",
		got.ReverseLookup.ValueString(),
		"reverseLookup should be set",
	)

	ddos := model.Ddos{}
	got.Ddos.As(context.TODO(), &ddos, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"protection-type",
		ddos.ProtectionType.ValueString(),
		"ddos should be set",
	)
}

func Test_adaptLoadBalancer(t *testing.T) {
	t.Run("loadBalancer Conversion works", func(t *testing.T) {
		reference := "reference"
		startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
		id := value_object.NewGeneratedUuid()

		loadBalancer := domain.NewLoadBalancer(
			id,
			value_object.InstanceType{Type: "type"},
			domain.Resources{Cpu: domain.Cpu{Unit: "cpu"}},
			"region",
			enum.StateCreating,
			domain.Contract{BillingFrequency: enum.ContractBillingFrequencySix},
			domain.Ips{{Ip: "1.2.3.4"}},
			domain.OptionalLoadBalancerValues{
				Reference:      &reference,
				StartedAt:      &startedAt,
				PrivateNetwork: &domain.PrivateNetwork{Id: "privateNetworkId"},
				Configuration: &domain.LoadBalancerConfiguration{
					Balance: enum.BalanceSource,
				},
			},
		)

		got, err := adaptLoadBalancer(
			context.TODO(),
			loadBalancer,
		)

		assert.NoError(t, err)

		assert.Equal(t, id.String(), got.Id.ValueString())
		assert.Equal(t, "type", got.Type.ValueString())
		assert.Equal(
			t,
			"{\"unit\":\"cpu\",\"value\":0}",
			got.Resources.Attributes()["cpu"].String(),
		)
		assert.Equal(t, "region", got.Region.ValueString())
		assert.Equal(t, "reference", got.Reference.ValueString())
		assert.Equal(t, "CREATING", got.State.ValueString())

		assert.Equal(
			t,
			"6",
			got.Contract.Attributes()["billing_frequency"].String(),
		)

		assert.Equal(
			t,
			"2019-09-08 00:00:00 +0000 UTC",
			got.StartedAt.ValueString(),
		)

		var ips []model.Ip
		got.Ips.ElementsAs(
			context.TODO(),
			&ips,
			false,
		)
		assert.Equal(t, "1.2.3.4", ips[0].Ip.ValueString())

		loadBalancerConfiguration := model.LoadBalancerConfiguration{}
		got.LoadBalancerConfiguration.As(
			context.TODO(),
			&loadBalancerConfiguration,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			"source",
			loadBalancerConfiguration.Balance.ValueString(),
		)

		privateNetwork := model.PrivateNetwork{}
		got.PrivateNetwork.As(
			context.TODO(),
			&privateNetwork,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			"privateNetworkId",
			privateNetwork.Id.ValueString(),
		)
	})
}

func TestAdaptInstance(t *testing.T) {
	var sshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

	t.Run("instance is adapted correctly", func(t *testing.T) {
		startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
		marketAppId := "marketAppId"
		reference := "reference"
		id := value_object.NewGeneratedUuid()
		rootDiskSize, _ := value_object.NewRootDiskSize(32)
		autoScalingGroupId := value_object.NewGeneratedUuid()
		sshKeyValueObject, _ := value_object.NewSshKey(sshKey)

		instance := generateDomainInstance()
		instance.Id = id
		instance.Type = value_object.NewUnvalidatedInstanceType(
			string(publicCloud.TYPENAME_M5A_4XLARGE),
		)
		instance.RootDiskSize = *rootDiskSize
		instance.StartedAt = &startedAt
		instance.MarketAppId = &marketAppId
		instance.Reference = &reference
		instance.SshKey = sshKeyValueObject
		instance.PrivateNetwork.Id = "privateNetworkId"
		instance.AutoScalingGroup.Id = autoScalingGroupId
		instance.Resources.Cpu.Unit = "cpu"

		got, err := AdaptInstance(instance, context.TODO())

		assert.NoError(t, err)
		assert.Equal(
			t,
			id.String(),
			got.Id.ValueString(),
			"id should be set",
		)
		assert.Equal(
			t,
			"region",
			got.Region.ValueString(),
			"region should be set",
		)
		assert.Equal(
			t,
			"CREATING",
			got.State.ValueString(),
			"state should be set",
		)
		assert.Equal(
			t,
			"productType",
			got.ProductType.ValueString(),
			"productType should be set",
		)
		assert.False(
			t,
			got.HasPublicIpv4.ValueBool(),
			"hasPublicIpv should be set",
		)
		assert.True(
			t,
			got.HasPrivateNetwork.ValueBool(),
			"hasPrivateNetwork should be set",
		)
		assert.Equal(
			t,
			"lsw.m5a.4xlarge",
			got.Type.ValueString(),
			"type should be set",
		)
		assert.Equal(
			t,
			int64(32),
			got.RootDiskSize.ValueInt64(),
			"rootDiskSize should be set",
		)
		assert.Equal(
			t,
			"CENTRAL",
			got.RootDiskStorageType.ValueString(),
			"rootDiskStorageType should be set",
		)
		assert.Equal(
			t,
			"2019-09-08 00:00:00 +0000 UTC",
			got.StartedAt.ValueString(),
			"startedAt should be set",
		)
		assert.Equal(
			t,
			"marketAppId",
			got.MarketAppId.ValueString(),
			"marketAppId should be set",
		)
		assert.Equal(
			t,
			"reference",
			got.Reference.ValueString(),
			"reference should be set",
		)

		image := model.Image{}
		got.Image.As(
			context.TODO(),
			&image,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			"UBUNTU_20_04_64BIT",
			image.Id.ValueString(),
			"image should be set",
		)

		contract := model.Contract{}
		got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
		assert.Equal(
			t,
			"MONTHLY",
			contract.Type.ValueString(),
			"contract should be set",
		)

		iso := model.Iso{}
		got.Iso.As(context.TODO(), &iso, basetypes.ObjectAsOptions{})
		assert.Equal(
			t,
			"isoId",
			iso.Id.ValueString(),
			"iso should be set",
		)

		privateNetwork := model.PrivateNetwork{}
		got.PrivateNetwork.As(
			context.TODO(),
			&privateNetwork,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			"privateNetworkId",
			privateNetwork.Id.ValueString(),
			"privateNetwork should be set",
		)

		autoScalingGroup := model.AutoScalingGroup{}
		got.AutoScalingGroup.As(
			context.TODO(),
			&autoScalingGroup,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			autoScalingGroupId.String(),
			autoScalingGroup.Id.ValueString(),
			"autoScalingGroup should be set",
		)

		var ips []model.Ip
		got.Ips.ElementsAs(context.TODO(), &ips, false)
		assert.Len(t, ips, 1)
		assert.Equal(
			t,
			"1.2.3.4",
			ips[0].Ip.ValueString(),
			"ip should be set",
		)

		resources := model.Resources{}
		cpu := model.Cpu{}
		got.Resources.As(context.TODO(), &resources, basetypes.ObjectAsOptions{})
		resources.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
		assert.Equal(
			t,
			"cpu",
			cpu.Unit.ValueString(),
			"privateNetwork should be set",
		)

		volume := model.Volume{}
		got.Volume.As(
			context.TODO(),
			&volume,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			"unit",
			volume.Unit.ValueString(),
			"volume should be set",
		)

		assert.Equal(t, sshKey, got.SshKey.ValueString())
	})
}

func Test_adaptAutoScalingGroup(t *testing.T) {
	desiredAmount := 1
	createdAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2020-09-08T00:00:00Z")
	startsAt, _ := time.Parse(time.RFC3339, "2010-09-08T00:00:00Z")
	endsAt, _ := time.Parse(time.RFC3339, "2011-09-08T00:00:00Z")
	minimumAmount := 2
	maximumAmount := 3
	cpuThreshold := 4
	warmupTime := 5
	cooldownTime := 6
	id := value_object.NewGeneratedUuid()
	reference, _ := value_object.NewAutoScalingGroupReference("reference")
	loadBalancerId := value_object.NewGeneratedUuid()

	autoScalingGroup := domain.NewAutoScalingGroup(
		id,
		"type",
		"state",
		"region",
		*reference,
		createdAt,
		updatedAt,
		domain.AutoScalingGroupOptions{
			DesiredAmount: &desiredAmount,
			StartsAt:      &startsAt,
			EndsAt:        &endsAt,
			MinimumAmount: &minimumAmount,
			MaximumAmount: &maximumAmount,
			CpuThreshold:  &cpuThreshold,
			WarmupTime:    &warmupTime,
			CoolDownTime:  &cooldownTime,
			LoadBalancer: &domain.LoadBalancer{
				Id:        loadBalancerId,
				StartedAt: &time.Time{},
			},
		},
	)

	got, err := adaptAutoScalingGroup(
		context.TODO(),
		autoScalingGroup,
	)

	assert.NoError(t, err)

	assert.Equal(t, id.String(), got.Id.ValueString())
	assert.Equal(t, "type", got.Type.ValueString())
	assert.Equal(t, "state", got.State.ValueString())
	assert.Equal(t, int64(1), got.DesiredAmount.ValueInt64())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(
		t,
		"2019-09-08 00:00:00 +0000 UTC",
		got.CreatedAt.ValueString(),
	)
	assert.Equal(
		t,
		"2020-09-08 00:00:00 +0000 UTC",
		got.UpdatedAt.ValueString(),
	)
	assert.Equal(
		t,
		"2010-09-08 00:00:00 +0000 UTC",
		got.StartsAt.ValueString(),
	)
	assert.Equal(
		t,
		"2011-09-08 00:00:00 +0000 UTC",
		got.EndsAt.ValueString(),
	)
	assert.Equal(t, int64(2), got.MinimumAmount.ValueInt64())
	assert.Equal(t, int64(3), got.MaximumAmount.ValueInt64())
	assert.Equal(t, int64(4), got.CpuThreshold.ValueInt64())
	assert.Equal(t, int64(5), got.WarmupTime.ValueInt64())
	assert.Equal(t, int64(6), got.CooldownTime.ValueInt64())

	loadBalancer := model.LoadBalancer{}
	got.LoadBalancer.As(
		context.TODO(),
		&loadBalancer,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, loadBalancerId.String(), loadBalancer.Id.ValueString())
}

func generateDomainInstance() domain.Instance {
	cpu := domain.NewCpu(1, "cpuUnit")
	memory := domain.NewMemory(2, "memoryUnit")
	publicNetworkSpeed := domain.NewNetworkSpeed(
		3,
		"publicNetworkSpeedUnit",
	)
	privateNetworkSpeed := domain.NewNetworkSpeed(
		4,
		"privateNetworkSpeedUnit",
	)

	resources := domain.NewResources(
		cpu,
		memory,
		publicNetworkSpeed,
		privateNetworkSpeed,
	)

	image := domain.NewImage(
		"UBUNTU_20_04_64BIT",
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		[]string{"one"},
		[]string{"storageType"},
	)

	rootDiskSize, _ := value_object.NewRootDiskSize(55)

	reverseLookup := "reverseLookup"
	ip := domain.NewIp(
		"1.2.3.4",
		"prefix-length",
		46,
		true,
		false,
		"tralala",
		domain.OptionalIpValues{
			Ddos:          &domain.Ddos{ProtectionType: "protection-type"},
			ReverseLookup: &reverseLookup,
		},
	)

	endsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2023-12-14 17:09:47",
	)
	renewalsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2022-12-14 17:09:47",
	)
	createdAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2021-12-14 17:09:47",
	)
	contract, _ := domain.NewContract(
		enum.ContractBillingFrequencySix,
		enum.ContractTermThree,
		enum.ContractTypeMonthly,
		renewalsAt,
		createdAt,
		enum.ContractStateActive,
		&endsAt,
	)

	reference := "reference"
	marketAppId := "marketAppId"
	sshKeyValueObject, _ := value_object.NewSshKey(defaultSshKey)
	startedAt := time.Now()

	privateNetwork := domain.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)

	stickySession := domain.NewStickySession(true, 5)

	host := "host"
	healthCheck := domain.NewHealthCheck(
		enum.MethodGet,
		"uri",
		22,
		domain.OptionalHealthCheckValues{Host: &host},
	)

	loadBalancerConfiguration := domain.NewLoadBalancerConfiguration(
		enum.BalanceSource,
		false,
		5,
		6,
		domain.OptionalLoadBalancerConfigurationOptions{
			StickySession: &stickySession,
			HealthCheck:   &healthCheck,
		},
	)

	loadBalancer := domain.NewLoadBalancer(
		value_object.NewGeneratedUuid(),
		value_object.NewUnvalidatedInstanceType("type"),
		resources,
		"region",
		enum.StateCreating,
		*contract,
		domain.Ips{ip},
		domain.OptionalLoadBalancerValues{
			Reference:      &reference,
			StartedAt:      &startedAt,
			PrivateNetwork: &privateNetwork,
			Configuration:  &loadBalancerConfiguration,
		},
	)

	autoScalingGroupReference, _ := value_object.NewAutoScalingGroupReference("reference")
	autoScalingGroupCreatedAt := time.Now()
	autoScalingGroupUpdatedAt := time.Now()
	autoScalingGroupDesiredAmount := 1
	autoScalingGroupStartsAt := time.Now()
	autoScalingGroupEndsAt := time.Now()
	autoScalingMinimumAmount := 2
	autoScalingMaximumAmount := 3
	autoScalingCpuThreshold := 4
	autoScalingWarmupTime := 5
	autoScalingCooldownTime := 6
	autoScalingGroup := domain.NewAutoScalingGroup(
		value_object.NewGeneratedUuid(),
		"type",
		"state",
		"region",
		*autoScalingGroupReference,
		autoScalingGroupCreatedAt,
		autoScalingGroupUpdatedAt,
		domain.AutoScalingGroupOptions{
			DesiredAmount: &autoScalingGroupDesiredAmount,
			StartsAt:      &autoScalingGroupStartsAt,
			EndsAt:        &autoScalingGroupEndsAt,
			MinimumAmount: &autoScalingMinimumAmount,
			MaximumAmount: &autoScalingMaximumAmount,
			CpuThreshold:  &autoScalingCpuThreshold,
			WarmupTime:    &autoScalingWarmupTime,
			CoolDownTime:  &autoScalingCooldownTime,
			LoadBalancer:  &loadBalancer,
		})

	volume := domain.NewVolume(1, "unit")

	return domain.NewInstance(
		value_object.NewGeneratedUuid(),
		"region",
		resources,
		image,
		enum.StateCreating,
		"productType",
		false,
		true,
		*rootDiskSize,
		value_object.NewUnvalidatedInstanceType(
			string(publicCloud.TYPENAME_C3_LARGE),
		),
		enum.RootDiskStorageTypeCentral,
		domain.Ips{ip},
		*contract,
		domain.OptionalInstanceValues{
			Reference:        &reference,
			Iso:              &domain.Iso{Id: "isoId"},
			MarketAppId:      &marketAppId,
			SshKey:           sshKeyValueObject,
			StartedAt:        &startedAt,
			PrivateNetwork:   &privateNetwork,
			AutoScalingGroup: &autoScalingGroup,
			Volume:           &volume,
		},
	)
}

func Test_adaptVolume(t *testing.T) {
	got, err := adaptVolume(
		context.TODO(),
		domain.Volume{
			Size: 2,
			Unit: "unit",
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, float64(2), got.Size.ValueFloat64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}
