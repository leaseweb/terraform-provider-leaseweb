package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	providerConfig = `
provider "leaseweb" {
  host     = "localhost:8080"
  scheme = "http"
  token = "tralala"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (
		tfprotov6.ProviderServer,
		error,
	){
		"leaseweb": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func TestLeasewebProvider_Metadata(t *testing.T) {
	leasewebProvider := New("dev")
	metadataResponse := provider.MetadataResponse{}
	leasewebProvider().Metadata(
		context.TODO(),
		provider.MetadataRequest{},
		&metadataResponse,
	)

	want := "dev"
	got := metadataResponse.Version

	assert.Equal(
		t,
		want,
		got,
		"version should be passed to provider",
	)
}

func TestLeasewebProvider_Schema(t *testing.T) {
	leasewebProvider := New("dev")
	schemaResponse := provider.SchemaResponse{}
	leasewebProvider().Schema(
		context.TODO(),
		provider.SchemaRequest{},
		&schemaResponse,
	)

	assert.True(
		t,
		schemaResponse.Schema.Attributes["host"].IsOptional(),
		"host is optional",
	)
	assert.True(
		t,
		schemaResponse.Schema.Attributes["scheme"].IsOptional(),
		"scheme is optional",
	)
	assert.True(
		t,
		schemaResponse.Schema.Attributes["token"].IsSensitive(),
		"token is sensitive",
	)
}

func TestAccPublicCloudInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "leaseweb_public_cloud_instances" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.#",
						"4",
					),

					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.id",
						"ace712e9-a166-47f1-9065-4af0f7e7fce1",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.reference",
						"my webserver",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.region",
						"eu-west-3",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.root_disk_size",
						"5",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.root_disk_storage_type",
						"CENTRAL",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.state",
						"RUNNING",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.type",
						"lsw.m3.large",
					),

					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.contract.billing_frequency",
						"1",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.contract.state",
						"ACTIVE",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.contract.term",
						"0",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.contract.type",
						"HOURLY",
					),

					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.image.custom",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.image.flavour",
						"ubuntu",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.image.id",
						"UBUNTU_20_04_64BIT",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.image.name",
						"Ubuntu 20.04 LTS (x86_64)",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.image.region",
						"eu-west-3",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.image.state",
						"READY",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.image.storage_types.#",
						"2",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.image.storage_types.1",
						"CENTRAL",
					),

					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.ips.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.ips.0.ip",
						"10.32.60.12",
					),

					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.iso.id",
						"ACRONIS_BOOT_MEDIA",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.iso.name",
						"Acronis Boot Media",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.has_private_network",
						"false",
					),
				),
			},
		},
	})
}

func TestAccPublicCloudInstanceResource(t *testing.T) {
	t.Run("creates and updates an instance", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 1
					    term = 0
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_storage_type = "CENTRAL"
					  root_disk_size         = 5
					}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"image.custom",
							"false",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"image.flavour",
							"ubuntu",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"image.market_apps.#",
							"0",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"image.name",
							"Ubuntu 20.04 LTS (x86_64)",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"ips.#",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"ips.0.ip",
							"10.32.60.12",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"reference",
							"my webserver",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"root_disk_size",
							"5",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"state",
							"RUNNING",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"contract.state",
							"ACTIVE",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"iso.id",
							"ACRONIS_BOOT_MEDIA",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"iso.name",
							"Acronis Boot Media",
						),
					),
				},
				// ImportState testing
				{
					ResourceName:      "leaseweb_public_cloud_instance.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				// Update and Read testing
				{
					Config: providerConfig + `
				resource "leaseweb_public_cloud_instance" "test" {
				  region = "eu-west-3"
				  type = "lsw.m3.large"
				  contract = {
				    billing_frequency = 1
				    term = 0
				    type = "HOURLY"
				  }
				  image = {
				    id = "UBUNTU_20_04_64BIT"
				  }
				  root_disk_storage_type = "CENTRAL"
				  root_disk_size         = 5
				}
				`,
				},
			},

			// Delete testing automatically occurs in TestCase
		})
	})

	t.Run("an invalid region throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "tralala"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 1
					    term = 0
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_storage_type = "CENTRAL"
					  root_disk_size         = 50
					}
					`,
					ExpectError: regexp.MustCompile(
						`Attribute region value must be one of:`,
					),
				},
			},
		})
	})

	t.Run("updating image.id triggers replacement", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 1
					    term = 0
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_storage_type = "CENTRAL"
					  root_disk_size         = 5
					}
					`,
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(
								"leaseweb_public_cloud_instance.test",
								plancheck.ResourceActionDestroyBeforeCreate,
							),
						},
					},
					// Ignore the inconsistent result as prism returns the old result.
					ExpectError: regexp.MustCompile(
						"Provider produced inconsistent result after apply",
					),
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 1
					    term = 0
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_24_04_64BIT"
					  }
					  root_disk_storage_type = "CENTRAL"
					  root_disk_size         = 5
					}
					`,
				},
			},
		})
	})

	t.Run("an invalid type throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "tralala"
					  contract = {
					    billing_frequency = 1
					    term = 0
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_storage_type = "CENTRAL"
					  root_disk_size         = 5
					}
					`,
					ExpectError: regexp.MustCompile(
						`Attribute type value must be one of:`,
					),
				},
			},
		})
	})

	t.Run("missing root_disk_size throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 1
					    term = 0
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_storage_type = "CENTRAL"
					}
					`,
					ExpectError: regexp.MustCompile(
						"The argument \"root_disk_size\" is required, but no definition was found.",
					),
				},
			},
		})
	})

	t.Run("an invalid root_disk_size throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 1
					    term = 0
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_size = 5000000
					  root_disk_storage_type = "CENTRAL"
					}
					`,
					ExpectError: regexp.MustCompile(
						"Attribute root_disk_size value must be between 5 and 1000",
					),
				},
			},
		})
	})

	t.Run("an invalid root_disk_storage_type throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 1
					    term = 0
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_storage_type = "tralala"
					  root_disk_size         = 5
					}
					`,
					ExpectError: regexp.MustCompile(
						"Attribute root_disk_storage_type value must be one of:",
					),
				},
			},
		})
	})

	t.Run("an invalid contract.billing_frequency throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 55
					    term = 0
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_storage_type = "CENTRAL"
					  root_disk_size         = 5
					}
					`,
					ExpectError: regexp.MustCompile(
						"Attribute contract.billing_frequency value must be one of:",
					),
				},
			},
		})
	})

	t.Run("an invalid contract.term throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 1
					    term = 55
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_storage_type = "CENTRAL"
					  root_disk_size         = 5
					}
					`,
					ExpectError: regexp.MustCompile(
						"Attribute contract.term value must be one of:",
					),
				},
			},
		})
	})

	t.Run("an invalid contract.type throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 1
					    term = 0
					    type = "tralala"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_storage_type = "CENTRAL"
					  root_disk_size         = 5
					}
					`,
					ExpectError: regexp.MustCompile(
						"Attribute contract.type value must be one of:",
					),
				},
			},
		})
	})

	t.Run("updating market_app_id triggers replacement", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 1
					    term = 0
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_storage_type = "CENTRAL"
					  root_disk_size         = 5
					}
					`,
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(
								"leaseweb_public_cloud_instance.test",
								plancheck.ResourceActionDestroyBeforeCreate,
							),
						},
					},
					// Ignore the inconsistent result as prism returns the old result.
					ExpectError: regexp.MustCompile(
						"Provider produced inconsistent result after apply",
					),
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  contract = {
					    billing_frequency = 1
					    term = 0
					    type = "HOURLY"
					  }
					  image = {
					    id = "UBUNTU_20_04_64BIT"
					  }
					  root_disk_storage_type = "CENTRAL"
					  root_disk_size         = 5
					  market_app_id = "test"
					}
					`,
				},
			},
		})
	})
}

func TestAccPublicCloudCredentialResource(t *testing.T) {
	t.Run("creates and updates a credential", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_credential" "test" {
						instance_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
					   	username = "root"
					   	type = "OPERATING_SYSTEM"
					   	password = "12341234"
					}`,
				},
				// Update and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_credential" "test" {
						instance_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
					   	username = "root"
					   	type = "OPERATING_SYSTEM"
					   	password = "12341234"
					}`,
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run(
		"username should not be empty",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_public_cloud_credential" "test" {
							instance_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
						   	username = ""
						   	type = "OPERATING_SYSTEM"
						   	password = "blah"
						}`,
						ExpectError: regexp.MustCompile(
							`Attribute username string length must be at least 1, got: 0`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"password should not be empty",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_public_cloud_credential" "test" {
							instance_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
						   	username = "root"
						   	type = "OPERATING_SYSTEM"
						   	password = ""
						}`,
						ExpectError: regexp.MustCompile(
							`Attribute password string length must be at least 1, got: 0`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"type must be a valid one",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_public_cloud_credential" "test" {
							instance_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
						   	username = "root"
						   	type = "invalid"
						   	password = "12341234"
						}`,
						ExpectError: regexp.MustCompile(
							`Attribute type value must be one of:`,
						),
					},
				},
			})
		},
	)
}

func TestAccPublicCloudCredentialDataSource(t *testing.T) {
	t.Run("read data for public cloud credential", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        data "leaseweb_public_cloud_credential" "test" {
						          instance_id         = "695ddd91-051f-4dd6-9120-938a927a47d0"
						          type                = "OPERATING_SYSTEM"
						          username            = "root"
						        }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_credential.test",
							"password",
							"12341234",
						),
					),
				},
			},
		})
	})

	t.Run("instance_id is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        data "leaseweb_public_cloud_credential" "test" {
						          type                = "OPERATING_SYSTEM"
						          username            = "root"
						        }`,
					ExpectError: regexp.MustCompile(
						"The argument \"instance_id\" is required, but no definition was found",
					),
				},
			},
		})
	})

	t.Run("type is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        data "leaseweb_public_cloud_credential" "test" {
						          instance_id         = "695ddd91-051f-4dd6-9120-938a927a47d0"
						          username            = "root"
						        }`,
					ExpectError: regexp.MustCompile(
						"The argument \"type\" is required, but no definition was found",
					),
				},
			},
		})
	})

	t.Run("username is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        data "leaseweb_public_cloud_credential" "test" {
						          instance_id         = "695ddd91-051f-4dd6-9120-938a927a47d0"
						          type                = "OPERATING_SYSTEM"
						        }`,
					ExpectError: regexp.MustCompile(
						"The argument \"username\" is required, but no definition was found",
					),
				},
			},
		})
	})

	t.Run("invalid type is not accepted", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        data "leaseweb_public_cloud_credential" "test" {
						          instance_id         = "695ddd91-051f-4dd6-9120-938a927a47d0"
						          type                = "A_WRONG_TYPE"
						          username            = "root"
						        }`,
					ExpectError: regexp.MustCompile(
						"Attribute type value must be one of:",
					),
				},
			},
		})
	})
}

func TestAccPublicCloudImagesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "leaseweb_public_cloud_images" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.#",
						"6",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.0.custom",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.0.flavour",
						"ubuntu",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.0.id",
						"UBUNTU_24_04_64BIT",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.0.market_apps.#",
						"0",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.0.name",
						"Ubuntu 24.04 LTS (x86_64)",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.0.region",
						"eu-west-3",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.0.state",
						"READY",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.0.storage_types.#",
						"2",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.0.storage_types.0",
						"LOCAL",
					),
				),
			},
		},
	})
}

func TestAccPublicCloudImageResource(t *testing.T) {
	t.Run("creates & updates an image", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_image" "test" {
					  instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
					  name = "Custom image - 03"
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_image.test",
							"custom",
							"true",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_image.test",
							"flavour",
							"ubuntu",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_image.test",
							"id",
							"ace712e9-a166-47f1-9065-4af0f7e7fce1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_image.test",
							"market_apps.#",
							"0",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_image.test",
							"region",
							"eu-west-3",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_image.test",
							"state",
							"READY",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_image.test",
							"storage_types.#",
							"0",
						),
					),
				},
				// Update and Read testing
				{
					Config: providerConfig + `
					  resource "leaseweb_public_cloud_image" "test" {
					    instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
					    name = "Custom image - 03"
					  }`,
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})
}

func TestPublicCloudAccLoadBalancersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "leaseweb_public_cloud_load_balancers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.id",
						"32082a93-d1e2-4bc0-8f5e-8fe4312b0844",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.reference",
						"my-loadbalancer1",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.region",
						"eu-west-3",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.state",
						"RUNNING",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.type",
						"lsw.m3.large",
					),

					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.contract.billing_frequency",
						"1",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.contract.state",
						"ACTIVE",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.contract.term",
						"0",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.contract.type",
						"HOURLY",
					),

					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.ips.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.ips.0.ip",
						"85.99.99.99",
					),
				),
			},
		},
	})
}

func TestAccDedicatedServerNotificationSettingBandwidthResource(t *testing.T) {
	t.Run("creates a notification setting bandwidth", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create testing
				{
					Config: providerConfig + `
					resource "leaseweb_dedicated_server_notification_setting_bandwidth" "test" {
					    dedicated_server_id = "12345678"
					    frequency = "WEEKLY"
					    threshold = "1"
					    unit = "Gbps"
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_bandwidth.test",
							"id",
							"12345",
						),
					),
				},
			},
		})
	})

	t.Run(
		"server id should be there in the request",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_dedicated_server_notification_setting_bandwidth" "test" {
						    frequency = "WEEKLY"
						    threshold = "1"
						    unit = "Gbps"
						}`,
						ExpectError: regexp.MustCompile(
							"The argument \"dedicated_server_id\" is required, but no definition was found",
						),
					},
				},
			})
		},
	)

	t.Run(
		"frequency should be one of these values 'DAILY', 'MONTHLY', 'WEEKLY'",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_dedicated_server_notification_setting_bandwidth" "test" {
						   dedicated_server_id = "12345678"
						   frequency = "WRONG"
						   threshold = "1"
						   unit = "Gbps"
						}`,
						ExpectError: regexp.MustCompile(
							`Attribute frequency value must be one of: \["DAILY" "WEEKLY" "MONTHLY"]`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"threshold should be greater than 0",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_dedicated_server_notification_setting_bandwidth" "test" {
						    dedicated_server_id = "12345678"
						    frequency = "DAILY"
						    threshold = "0"
						    unit = "Gbps"
						}`,
						ExpectError: regexp.MustCompile(
							"The value must be greater than 0, but got 0",
						),
					},
				},
			})
		},
	)

	t.Run(
		"unit should be one of these values 'Mbps', 'Gbps'",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_dedicated_server_notification_setting_bandwidth" "test" {
						   dedicated_server_id = "12345678"
						   frequency = "DAILY"
						   threshold = "0"
						   unit = "Kbps"
						}`,
						ExpectError: regexp.MustCompile(
							`Attribute unit value must be one of: \["Mbps" "Gbps"], got: "Kbps"`,
						),
					},
				},
			})
		},
	)
}

func TestAccDedicatedServerAccControlPanelsDataSource(t *testing.T) {
	t.Run("get all control panels", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
							data "leaseweb_dedicated_server_control_panels" "dtest" {
							}
						`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_dedicated_server_control_panels.dtest",
							"control_panels.#",
							"8",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_dedicated_server_control_panels.dtest",
							"control_panels.0.id",
							"CPANEL_PREMIER_100",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_dedicated_server_control_panels.dtest",
							"control_panels.0.name",
							"cPanel Premier 100",
						),
					),
				},
			},
		})
	})

	t.Run(
		"filter control panels by operating_system_id",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
							data "leaseweb_dedicated_server_control_panels" "dtest" {
							    operating_system_id = "ALMALINUX_8_64BIT"
							}
						`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_control_panels.dtest",
								"control_panels.#",
								"8",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_control_panels.dtest",
								"control_panels.0.id",
								"CPANEL_PREMIER_100",
							),
						),
					},
				},
			})
		},
	)
}

func TestAccDedicatedServerCredentialDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				        data "leaseweb_dedicated_server_credential" "test" {
				          dedicated_server_id = "12345"
				          type                = "OPERATING_SYSTEM"
				          username            = "root"
				        }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.leaseweb_dedicated_server_credential.test",
						"password",
						"mys3cr3tp@ssw0rd",
					),
				),
			},
		},
	})
}

func TestAccDedicatedServerCredentialResource(t *testing.T) {
	t.Run("creates and updates a credential", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_dedicated_server_credential" "test" {
						dedicated_server_id = "12345"
					   	username = "root"
					   	type = "OPERATING_SYSTEM"
					   	password = "mys3cr3tp@ssw0rd"
					}`,
				},
				// Update and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_dedicated_server_credential" "test" {
						dedicated_server_id = "12345"
					   	username = "root"
					   	type = "OPERATING_SYSTEM"
					   	password = "mys3cr3tp@ssw0rd"
					}`,
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run("type must be valid", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						resource "leaseweb_dedicated_server_credential" "test" {
							dedicated_server_id = "12345"
						   	username = "root"
						   	type = "invalid"
						   	password = "mys3cr3tp@ssw0rd"
						}`,
					ExpectError: regexp.MustCompile(
						`Attribute type value must be one of: \["OPERATING_SYSTEM" "CONTROL_PANEL"(\s*)"REMOTE_MANAGEMENT" "RESCUE_MODE" "SWITCH" "PDU" "FIREWALL" "LOAD_BALANCER"],(\s*)got: "invalid"`,
					),
				},
			},
		})
	})
}

func TestAccDataTrafficNotificationSettingResource(t *testing.T) {
	t.Run("creates and updates a data traffic notification setting", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_dedicated_server_notification_setting_datatraffic" "test" {
					  dedicated_server_id = "145406"
					  frequency = "WEEKLY"
					  threshold = "1"
					  unit = "GB"
					}`,
				},
				// Update and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_dedicated_server_notification_setting_datatraffic" "test" {
					  dedicated_server_id = "145406"
					  frequency = "WEEKLY"
					  threshold = "1"
					  unit = "GB"
					}`,
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run(
		"threshold must be greater than 0",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_dedicated_server_notification_setting_datatraffic" "test" {
						  dedicated_server_id = "145406"
						  frequency = "WEEKLY"
						  threshold = "-1"
						  unit = "GB"
						}`,
						ExpectError: regexp.MustCompile(
							"The value must be greater than 0, but got -1",
						),
					},
				},
			})
		},
	)

	t.Run(
		"unit must be one of GB,MB,TB",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_dedicated_server_notification_setting_datatraffic" "test" {
						  dedicated_server_id = "145406"
						  frequency = "WEEKLY"
						  threshold = "1"
						  unit = "blah"
						}`,
						ExpectError: regexp.MustCompile(
							`Attribute unit value must be one of: \["MB" "GB" "TB"], got: "blah"`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"frequency must be one of DAILY,WEEKLY,MONTHLY",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_dedicated_server_notification_setting_datatraffic" "test" {
						  dedicated_server_id = "145406"
						  frequency = "blah"
						  threshold = "1"
						  unit = "GB"
						}`,
						ExpectError: regexp.MustCompile(
							`Attribute frequency value must be one of: \["DAILY" "WEEKLY" "MONTHLY"], got:`,
						),
					},
				},
			})
		},
	)
}

func TestAccDedicatedServerInstallationResource(t *testing.T) {
	t.Run("install os on a dedicated server",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Create testing
					{
						Config: providerConfig + `
						resource "leaseweb_dedicated_server_installation" "test" {
							dedicated_server_id = "12345"
						    operating_system_id = "UBUNTU_22_04_64BIT"
							callback_url = "https://example.com/callBack"
							control_panel_id = "123456"
							device = "SATA2TB"
							hostname = "server.leaseweb.com"
							password = "password"
							post_install_script = <<-EOS
								#!/bin/sh
								apt install nginx -y -qq
							EOS
							power_cycle = true
							ssh_keys = ["tralala"]
							timezone = "UTC"
							partitions  = [
								{
									filesystem = "ext2"
									mountpoint = "/boot"
									size = 1024
								},
								{
									filesystem = "swap"
									size = 4096
								},
								{
									filesystem = "ext4"
									mountpoint = "/tmp"
									size = 4096
								},
								{
									filesystem = "ext4"
									mountpoint = "/"
									size = "*"
								}
							]
							raid = {
								level = 0
								number_of_disks = 2
								type = "SW"
							}
						}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"leaseweb_dedicated_server_installation.test",
								"id",
								"bcf2bedf-8450-4b22-86a8-f30aeb3a38f9",
							),
							resource.TestCheckResourceAttr(
								"leaseweb_dedicated_server_installation.test",
								"post_install_script",
								"#!/bin/sh\napt install nginx -y -qq\n",
							),
						),
					},
					{
						ResourceName:                         "leaseweb_dedicated_server_installation.test",
						ImportState:                          true,
						ImportStateVerify:                    true,
						ImportStateId:                        "12345",
						ImportStateVerifyIdentifierAttribute: "dedicated_server_id",
						ImportStateVerifyIgnore:              []string{"callback_url", "control_panel_id", "hostname", "password", "ssh_keys", "post_install_script", "raid"},
					},
				},
			})
		})

	t.Run(
		"server id should be in the request",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						    resource "leaseweb_dedicated_server_installation" "test" {
						      operating_system_id = "UBUNTU_22_04_64BIT"
						    }`,
						ExpectError: regexp.MustCompile(
							"The argument \"dedicated_server_id\" is required, but no definition was found",
						),
					},
				},
			})
		},
	)

	t.Run(
		"operating system id should be in the request",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						    resource "leaseweb_dedicated_server_installation" "test" {
						      dedicated_server_id = "12345"
						    }`,
						ExpectError: regexp.MustCompile(
							"The argument \"operating_system_id\" is required, but no definition was found",
						),
					},
				},
			})
		},
	)

	t.Run(
		"raid.level should be one of these values '0', '1', '5', '10'",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_dedicated_server_installation" "test" {
							dedicated_server_id = "12345"
						    operating_system_id = "UBUNTU_22_04_64BIT"
						    raid = {
						    	level = 11
						    }
						}`,
						ExpectError: regexp.MustCompile(
							`Attribute raid.level value must be one of: \["0" "1" "5" "10"]`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"raid.type should be one of these values 'HW', 'SW', 'NONE'",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_dedicated_server_installation" "test" {
							dedicated_server_id = "12345"
						    operating_system_id = "UBUNTU_22_04_64BIT"
						    raid = {
						    	type = "TEST"
						    }
						}`,
						ExpectError: regexp.MustCompile(
							`Attribute raid.type value must be one of: \["HW" "SW" "NONE"]`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"ssh_keys should be set of string",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_dedicated_server_installation" "test" {
							dedicated_server_id = "12345"
						    operating_system_id = "UBUNTU_22_04_64BIT"
						    ssh_keys = "test keys"
						}`,
						ExpectError: regexp.MustCompile(
							`Inappropriate value for attribute "ssh_keys": set of string required`,
						),
					},
				},
			})
		},
	)
}

func TestAccOperatingSystemsDataSource(t *testing.T) {
	t.Run("get all operating systems", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
						data "leaseweb_dedicated_server_operating_systems" "dtest" {
							control_panel_id = "1234"
						}
						`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_dedicated_server_operating_systems.dtest",
							"operating_systems.#",
							"24",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_dedicated_server_operating_systems.dtest",
							"operating_systems.0.id",
							"ALMALINUX_8_64BIT",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_dedicated_server_operating_systems.dtest",
							"operating_systems.0.name",
							"AlmaLinux 8 (x86_64)",
						),
					),
				},
			},
		})
	})

	t.Run(
		"filter operating systems by control_panel_id",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
							data "leaseweb_dedicated_server_operating_systems" "dtest" {
								control_panel_id = "CPANEL_PREMIER_100"
							}
						`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_operating_systems.dtest",
								"operating_systems.#",
								"24",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_operating_systems.dtest",
								"operating_systems.0.id",
								"ALMALINUX_8_64BIT",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_operating_systems.dtest",
								"operating_systems.0.name",
								"AlmaLinux 8 (x86_64)",
							),
						),
					},
				},
			})
		},
	)
}

func TestAccDedicatedServerDataSource(t *testing.T) {
	t.Run("get dedicated server detail by id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
						data "leaseweb_dedicated_server" "test" {
							id = "123456"
						}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_dedicated_server.test",
							"asset_id",
							"1234561882",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_dedicated_server.test",
							"contract_id",
							"12123412312",
						),
					),
				},
			},
		})
	})

	t.Run(
		"id is required to get the dedicated server detail",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						data "leaseweb_dedicated_server" "test" {
						}`,
						ExpectError: regexp.MustCompile(
							`The argument "id" is required, but no definition was found`,
						),
					},
				},
			})
		},
	)
}

func TestAccDedicatedServersDataSource(t *testing.T) {
	t.Run("get dedicated servers by reference", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
						data "leaseweb_dedicated_servers" "test" {
							reference = "test-reference"
						}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_dedicated_servers.test",
							"ids.#",
							"2",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_dedicated_servers.test",
							"ids.0",
							"12345",
						),
					),
				},
			},
		})
	})

	t.Run("get dedicated servers", func(t *testing.T) {
		resource.Test(
			t,
			resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
						data "leaseweb_dedicated_servers" "test" {
						}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.test",
								"ids.#",
								"2",
							),
						),
					},
				},
			},
		)
	})

	t.Run("get dedicated servers with all filters", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
						data "leaseweb_dedicated_servers" "filter" {
							reference = "test-reference"
							ip = "127.0.0.4"
							mac_address = "aa:bb:cc:dd:ee:ff"
							site = "ams-01"
							private_rack_id = "r id"
							private_network_capable = "true"
							private_network_enabled = "true"
						}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_dedicated_servers.filter",
							"ids.#",
							"2",
						),
					),
				},
			},
		})
	})
}

func TestAccPublicCloudLoadBalancerResource(t *testing.T) {
	t.Run("creates and updates a loadBalancer", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_load_balancer" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  reference = "my-loadbalancer1"
					  contract = {
					    billing_frequency = 1
					    term              = 0
					    type              = "HOURLY"
					  }
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"id",
							"32082a93-d1e2-4bc0-8f5e-8fe4312b0844",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"contract.state",
							"ACTIVE",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"ips.#",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"ips.0.ip",
							"85.99.99.99",
						),
					),
				},
				// ImportState testing
				{
					ResourceName:      "leaseweb_public_cloud_load_balancer.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				// Update and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_load_balancer" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  reference = "my-loadbalancer1"
					  contract = {
					    billing_frequency = 1
					    term              = 0
					    type              = "HOURLY"
					  }
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"id",
							"32082a93-d1e2-4bc0-8f5e-8fe4312b0844",
						),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run("invalid type", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_load_balancer" "test" {
					  region = "eu-west-3"
					  type = "tralala"
					  reference = "my-loadbalancer1"
					  contract = {
					    billing_frequency = 1
					    term              = 0
					    type              = "HOURLY"
					  }
					}`,
					ExpectError: regexp.MustCompile(
						"Attribute type value must be one of:",
					),
				},
			},
		})
	})

	t.Run("invalid region", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_load_balancer" "test" {
					  region = "tralala"
					  type = "lsw.m4.2xlarge"
					  reference = "my-loadbalancer1"
					  contract = {
					    billing_frequency = 1
					    term              = 0
					    type              = "HOURLY"
					  }
					}`,
					ExpectError: regexp.MustCompile("Attribute region value must be one of"),
				},
			},
		})
	})

	t.Run("invalid contract.billingFrequency", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_load_balancer" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.2xlarge"
					  reference = "my-loadbalancer1"
					  contract = {
					    billing_frequency = 55
					    term              = 0
					    type              = "HOURLY"
					  }
					}`,
					ExpectError: regexp.MustCompile(
						"Attribute contract.billing_frequency value must be one of",
					),
				},
			},
		})
	})

	t.Run("invalid contract.term", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_load_balancer" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.2xlarge"
					  reference = "my-loadbalancer1"
					  contract = {
					    billing_frequency = 1
					    term              = 55
					    type              = "MONTHLY"
					  }
					}`,
					ExpectError: regexp.MustCompile(
						"Attribute contract.term value must be one of",
					),
				},
			},
		})
	})

	type errorTestCases struct {
		requiredField string
		expectedError string
	}

	for _, scenario := range []errorTestCases{
		{
			requiredField: "region",
			expectedError: fmt.Sprintf(
				"The argument %q is required, but no definition was found.",
				"region",
			),
		},
		{
			requiredField: "type",
			expectedError: fmt.Sprintf(
				"The argument %q is required, but no definition was found.",
				"type",
			),
		},
		{
			requiredField: "contract.type|contract.term|contract.billing_frequency",
			expectedError: "Inappropriate value for attribute \"contract\": attributes \"billing_frequency\",\n\"term\", and \"type\" are required.",
		},
	} {
		t.Run(scenario.requiredField+" should be set", func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_public_cloud_load_balancer" "test" {
						  contract = {}
						}`,
						ExpectError: regexp.MustCompile(scenario.expectedError),
					},
				},
			})
		})
	}

	t.Run("changing the region triggers replacement", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_load_balancer" "test" {
					  region = "eu-west-3"
					  type = "lsw.m3.large"
					  reference = "my-loadbalancer1"
					  contract = {
					    billing_frequency = 1
					    term              = 0
					    type              = "HOURLY"
					  }
					}`,
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(
								"leaseweb_public_cloud_load_balancer.test",
								plancheck.ResourceActionDestroyBeforeCreate,
							),
						},
					},
					// Ignore the inconsistent result as prism returns the old result.
					ExpectError: regexp.MustCompile(
						"Provider produced inconsistent result after apply",
					),
					Config: providerConfig + `
					resource "leaseweb_public_cloud_load_balancer" "test" {
					  region = "eu-west-2"
					  type = "lsw.m3.large"
					  reference = "my-loadbalancer1"
					  contract = {
					    billing_frequency = 1
					    term              = 0
					    type              = "HOURLY"
					  }
					}`,
				},
			},
		})
	})
}

func TestAccPublicCloudLoadBalancerListenersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
				data "leaseweb_public_cloud_load_balancer_listeners" "test" {
				    load_balancer_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancer_listeners.test",
						"listeners.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancer_listeners.test",
						"listeners.0.id",
						"fac06878-6655-4956-8ea7-124a97f133ab",
					),
				),
			},
		},
	})
}

func TestAccPublicCloudTargetGroupsDataSource(t *testing.T) {
	t.Run("can read all target groups", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `data "leaseweb_public_cloud_target_groups" "test" {}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.#",
							"1",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.0.id",
							"7e59b33d-05f3-4078-b251-c7831ae8fe14",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.0.name",
							"Target group name",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.0.port",
							"80",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.0.protocol",
							"HTTP",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.0.region",
							"eu-west-2",
						),
					),
				},
				{
					Config: providerConfig + `
					data "leaseweb_public_cloud_target_groups" "test" {
					  id = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
					}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.#",
							"1",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.0.id",
							"7e59b33d-05f3-4078-b251-c7831ae8fe14",
						),
					),
				},
			},
		})
	})

	t.Run("can filter target groups by id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
					data "leaseweb_public_cloud_target_groups" "test" {
					  id = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
					}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.#",
							"1",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.0.id",
							"7e59b33d-05f3-4078-b251-c7831ae8fe14",
						),
					),
				},
			},
		})
	})

	t.Run("can filter target groups by name", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
					data "leaseweb_public_cloud_target_groups" "test" {
					  name = "Foo bar"
					}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.#",
							"1",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.0.id",
							"7e59b33d-05f3-4078-b251-c7831ae8fe14",
						),
					),
				},
			},
		})
	})

	t.Run("can filter target groups by protocol", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
					data "leaseweb_public_cloud_target_groups" "test" {
					  protocol = "HTTP"
					}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.#",
							"1",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.0.id",
							"7e59b33d-05f3-4078-b251-c7831ae8fe14",
						),
					),
				},
			},
		})
	})

	t.Run("inputting an invalid protocol throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					data "leaseweb_public_cloud_target_groups" "test" {
					  protocol = "tralala"
					}
					`,
					ExpectError: regexp.MustCompile(
						`Attribute protocol value must be one of:`,
					),
				},
			},
		})
	})

	t.Run("can filter target groups by port", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
					data "leaseweb_public_cloud_target_groups" "test" {
					  port = 80
					}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.#",
							"1",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.0.id",
							"7e59b33d-05f3-4078-b251-c7831ae8fe14",
						),
					),
				},
			},
		})
	})

	t.Run("inputting an invalid port throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
					data "leaseweb_public_cloud_target_groups" "test" {
					  port = 800000
					}
					`,
					ExpectError: regexp.MustCompile(
						`Attribute port value must be between 1 and 65535`,
					),
				},
			},
		})
	})

	t.Run("can filter target groups by region", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
					data "leaseweb_public_cloud_target_groups" "test" {
					  region = "eu-west-3"
					}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.#",
							"1",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_target_groups.test",
							"target_groups.0.id",
							"7e59b33d-05f3-4078-b251-c7831ae8fe14",
						),
					),
				},
			},
		})
	})

	t.Run("inputting an invalid region throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
					data "leaseweb_public_cloud_target_groups" "test" {
					  region = "tralala"
					}
					`,
					ExpectError: regexp.MustCompile(
						`Attribute region value must be one of:`,
					),
				},
			},
		})
	})
}

func TestAccPublicCloudLoadBalancerListenerResource(t *testing.T) {
	t.Run(
		"can create/import/update/delete load balancer listeners",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Create and Read testing
					{
						Config: providerConfig + `
						resource "leaseweb_public_cloud_load_balancer_listener" "test" {
						  certificate = {
						    certificate = "-----BEGIN CERTIFICATE-----MIIBhDCB7gIBADBFMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCtWdKNbZxvkXKAADjJMJ7VTJz6uFoMD403C+gMIF8hwqIsHggzCao6iXrW9sZoyZtUBVBiiq5RumHbbpwvOdMmXrShEB4sTJkWRMDy7yD4D91WCU1fc10E/zBJMwssAvmHZo5kGW1Pj1N9ktb+O/TMsEc6yd5suvdQj6aaJbQlTQIDAQABoAAwDQYJKoZIhvcNAQELBQADgYEAWOQ2CJLRo8MQgJgvhdoSIkHITnrbjB5hS3f/dx0lIcnyI6Q9nOyuQHXkCgkdBaV8lz7l+IbqcGc3CaIRP2ZIVFvo2252n630tOOSsqoqJS1tYIoIKsohi3T3d8T1i/s0BWbTJi8Xgd186wyUn/jHwXROKx2rq6yYsAO6fISDKw8=-----END CERTIFICATE-----"
						    chain       = "-----BEGIN CERTIFICATE-----MIICNDCCAZ2gAwIBAgIUEby6nzM+o7vkKfzcMS/DGA8tgwQwDQYJKoZIhvcNAQELBQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yNDA0MjUwODE3MjZaFw0yNTA0MjUwODE3MjZaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMMiux2r1AFLVpIhdZ0bvgIvhiT9XCnfHJlGE7OarGKDKJDQ6XAquCfosLws2XAugGcMJWrsqVWtJEYSu6OMsDLYCJhh39AKqZIW0pktkr8LGlo4VLvzGPqwpHnzWthyCEFsE6p+JJQumDA/izJm2zjZL+xHDocOlNqDTB87AIdrAgMBAAGjITAfMB0GA1UdDgQWBBT3sXUrIR2vcwak0QCXoIsxHa4dDDANBgkqhkiG9w0BAQsFAAOBgQCh/l+5s713J02b8sWicUK2KjTPfyKmZFkoS+Mlo+//B/aM612ZJpGL2tAKGF3v0NDOrRYLZj0t/tlZI55pUNJI9cNj/RExvnfTSSNJIbV+8kQt5AHo50wGxj/apkuEtQre2Fpf4pyovcfIoF6HJvvp1jy96yL14UoEehZypR8FlA==-----END CERTIFICATE-----"
						    private_key = "-----BEGIN EC PRIVATE KEY-----MIHcAgEBBEIBVlC0IObonfQZIQ81l/WILKfWT5Fv96eNnYmQZ7uleu73igfiVESVuPfNlbW9oNEK1XcXli4YNZMxWMkKuzC3w8CgBwYFK4EEACOhgYkDgYYABAHvOqz9d2xeSpm1FNdw0NR5j/q6PMd6whZFsTPNYNj0/PsTpsHk78ZB4MYnJUXwHJjpj+gnKkLNc02f4w/vSF8VXADX4l40XU/w82TAOCftQwoxO5o0jZcwEUIYzl02Zd7uNxhjtKJQnYFi9x8WI8L8zQ6GZB/fJKYwoHaUr0I1h/5LzQ==-----END EC PRIVATE KEY-----"
						  }
						  default_rule = {
						    target_group_id = "b05917e1-96a4-442a-900c-c41f273d95c9"
						  }
							load_balancer_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
							port = 443
						  protocol = "HTTPS"
						}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"leaseweb_public_cloud_load_balancer_listener.test",
								"listener_id",
								"fac06878-6655-4956-8ea7-124a97f133ab",
							),
						),
					},
					// ImportState testing
					{
						ResourceName:                         "leaseweb_public_cloud_load_balancer_listener.test",
						ImportStateId:                        "695ddd91-051f-4dd6-9120-938a927a47d0,fac06878-6655-4956-8ea7-124a97f133ab",
						ImportState:                          true,
						ImportStateVerify:                    true,
						ImportStateVerifyIdentifierAttribute: "listener_id",
					},
					// Update and Read testing
					{
						Config: providerConfig + `
						resource "leaseweb_public_cloud_load_balancer_listener" "test" {
						  certificate = {
						    certificate = "-----BEGIN CERTIFICATE-----MIIBhDCB7gIBADBFMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCtWdKNbZxvkXKAADjJMJ7VTJz6uFoMD403C+gMIF8hwqIsHggzCao6iXrW9sZoyZtUBVBiiq5RumHbbpwvOdMmXrShEB4sTJkWRMDy7yD4D91WCU1fc10E/zBJMwssAvmHZo5kGW1Pj1N9ktb+O/TMsEc6yd5suvdQj6aaJbQlTQIDAQABoAAwDQYJKoZIhvcNAQELBQADgYEAWOQ2CJLRo8MQgJgvhdoSIkHITnrbjB5hS3f/dx0lIcnyI6Q9nOyuQHXkCgkdBaV8lz7l+IbqcGc3CaIRP2ZIVFvo2252n630tOOSsqoqJS1tYIoIKsohi3T3d8T1i/s0BWbTJi8Xgd186wyUn/jHwXROKx2rq6yYsAO6fISDKw8=-----END CERTIFICATE-----"
						    chain       = "-----BEGIN CERTIFICATE-----MIICNDCCAZ2gAwIBAgIUEby6nzM+o7vkKfzcMS/DGA8tgwQwDQYJKoZIhvcNAQELBQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yNDA0MjUwODE3MjZaFw0yNTA0MjUwODE3MjZaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMMiux2r1AFLVpIhdZ0bvgIvhiT9XCnfHJlGE7OarGKDKJDQ6XAquCfosLws2XAugGcMJWrsqVWtJEYSu6OMsDLYCJhh39AKqZIW0pktkr8LGlo4VLvzGPqwpHnzWthyCEFsE6p+JJQumDA/izJm2zjZL+xHDocOlNqDTB87AIdrAgMBAAGjITAfMB0GA1UdDgQWBBT3sXUrIR2vcwak0QCXoIsxHa4dDDANBgkqhkiG9w0BAQsFAAOBgQCh/l+5s713J02b8sWicUK2KjTPfyKmZFkoS+Mlo+//B/aM612ZJpGL2tAKGF3v0NDOrRYLZj0t/tlZI55pUNJI9cNj/RExvnfTSSNJIbV+8kQt5AHo50wGxj/apkuEtQre2Fpf4pyovcfIoF6HJvvp1jy96yL14UoEehZypR8FlA==-----END CERTIFICATE-----"
						    private_key = "-----BEGIN EC PRIVATE KEY-----MIHcAgEBBEIBVlC0IObonfQZIQ81l/WILKfWT5Fv96eNnYmQZ7uleu73igfiVESVuPfNlbW9oNEK1XcXli4YNZMxWMkKuzC3w8CgBwYFK4EEACOhgYkDgYYABAHvOqz9d2xeSpm1FNdw0NR5j/q6PMd6whZFsTPNYNj0/PsTpsHk78ZB4MYnJUXwHJjpj+gnKkLNc02f4w/vSF8VXADX4l40XU/w82TAOCftQwoxO5o0jZcwEUIYzl02Zd7uNxhjtKJQnYFi9x8WI8L8zQ6GZB/fJKYwoHaUr0I1h/5LzQ==-----END EC PRIVATE KEY-----"
						  }
						  default_rule = {
						    target_group_id = "b05917e1-96a4-442a-900c-c41f273d95c9"
						  }
							load_balancer_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
							port = 443
						  protocol = "HTTPS"
						}`,
					},
					// Delete testing automatically occurs in TestCase
				},
			})
		},
	)

	t.Run("invalid protocol causes error to be thrown", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_load_balancer_listener" "test" {
					  default_rule = {
					    target_group_id = "b05917e1-96a4-442a-900c-c41f273d95c9"
					  }
						load_balancer_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
						port = 80
					  protocol = "tralala"
					}`,
					ExpectError: regexp.MustCompile(
						`Attribute protocol value must be one of:`,
					),
				},
			},
		})
	})

	t.Run("invalid port causes error to be thrown", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_load_balancer_listener" "test" {
					  default_rule = {
					    target_group_id = "b05917e1-96a4-442a-900c-c41f273d95c9"
					  }
						load_balancer_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
						port = -8
					  protocol = "HTTP"
					}`,
					ExpectError: regexp.MustCompile(
						`Attribute port value must be between`,
					),
				},
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_load_balancer_listener" "test" {
					  default_rule = {
					    target_group_id = "b05917e1-96a4-442a-900c-c41f273d95c9"
					  }
						load_balancer_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
						port = 400000
					  protocol = "HTTP"
					}`,
					ExpectError: regexp.MustCompile(
						`Attribute port value must be between`,
					),
				},
			},
		})
	})
}

func TestAccPublicCloudTargetGroupResource(t *testing.T) {
	t.Run("an invalid protocol throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					  resource "leaseweb_public_cloud_target_group" "test" {
					    name = "name"
					    port = 80
					    region = "eu-west-3"
					    protocol = "tralala"
					  }`,
					ExpectError: regexp.MustCompile(
						`Attribute protocol value must be one of:`,
					),
				},
			},
		})
	})

	t.Run("an invalid port throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					data "leaseweb_public_cloud_target_groups" "test" {
					  name = "name"
					  port = 800000
					  region = "eu-west-3"
					  protocol = "HTTP"
					}
					`,
					ExpectError: regexp.MustCompile(
						`Attribute port value must be between 1 and 65535`,
					),
				},
			},
		})
	})

	t.Run("an invalid region throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					  resource "leaseweb_public_cloud_target_group" "test" {
					    name = "name"
					    port = 80
					    region = "tralala"
					    protocol = "HTTP"
					  }`,
					ExpectError: regexp.MustCompile(
						`Attribute region value must be one of:`,
					),
				},
			},
		})
	})

	t.Run("an invalid health_check protocol throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					  resource "leaseweb_public_cloud_target_group" "test" {
					    name = "name"
					    port = 80
					    region = "eu-west-3"
					    protocol = "HTTP"
					    health_check = {
					      protocol = "tralala"
					      port = 80
					      uri = "/"
					    }
					  }`,
					ExpectError: regexp.MustCompile(
						`Attribute health_check.protocol value must be one of:`,
					),
				},
			},
		})
	})

	t.Run("an invalid health_check method throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					  resource "leaseweb_public_cloud_target_group" "test" {
					    name = "name"
					    port = 80
					    region = "eu-west-3"
					    protocol = "HTTP"
					    health_check = {
					      protocol = "HTTP"
					      method = "tralala"
					      port = 80
					      uri = "/"
					    }
					  }`,
					ExpectError: regexp.MustCompile(
						`Attribute health_check.method value must be one of:`,
					),
				},
			},
		})
	})

	t.Run("an invalid health_check port throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_target_group" "test" {
					  name = "name"
					  port = 80
					  region = "eu-west-3"
					  protocol = "HTTP"
					  health_check = {
					    protocol = "HTTP"
					    port = 80000
					    uri = "/"
					  }
					}
					`,
					ExpectError: regexp.MustCompile(
						`Attribute health_check.port value must be between 1 and 65535`,
					),
				},
			},
		})
	})

	t.Run("creates and updates a target group", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_target_group" "test" {
					  name = "Target group name"
					  port = 80
					  region = "eu-west-2"
					  protocol = "HTTP"
					  health_check = {
					    host = "example.com"
					    method = "GET"
					    protocol = "HTTP"
					    port = 80
					    uri = "/"
					  }
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_target_group.test",
							"id",
							"7e59b33d-05f3-4078-b251-c7831ae8fe14",
						),
					),
				},
				// ImportState testing
				{
					ResourceName:      "leaseweb_public_cloud_target_group.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				// Update and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_target_group" "test" {
					  name = "Target group name"
					  port = 80
					  region = "eu-west-2"
					  protocol = "HTTP"
					  health_check = {
					    host = "example.com"
					    method = "GET"
					    protocol = "HTTP"
					    port = 80
					    uri = "/"
					  }
					}`,
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run("changing the region triggers replacement", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_target_group" "test" {
					  name = "Target group name"
					  port = 80
					  region = "eu-west-2"
					  protocol = "HTTP"
					  health_check = {
					    host = "example.com"
					    method = "GET"
					    protocol = "HTTP"
					    port = 80
					    uri = "/"
					  }
					}`,
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(
								"leaseweb_public_cloud_target_group.test",
								plancheck.ResourceActionDestroyBeforeCreate,
							),
						},
					},
					// Ignore the inconsistent result as prism returns the old result.
					ExpectError: regexp.MustCompile(
						"Provider produced inconsistent result after apply",
					),
					Config: providerConfig + `
					resource "leaseweb_public_cloud_target_group" "test" {
					  name = "Target group name"
					  port = 80
					  region = "eu-west-3"
					  protocol = "HTTP"
					  health_check = {
					    host = "example.com"
					    method = "GET"
					    protocol = "HTTP"
					    port = 80
					    uri = "/"
					  }
					}`,
				},
			},
		})
	})

	t.Run("changing the protocol triggers replacement", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_target_group" "test" {
					  name = "Target group name"
					  port = 80
					  region = "eu-west-2"
					  protocol = "HTTP"
					  health_check = {
					    host = "example.com"
					    method = "GET"
					    protocol = "HTTP"
					    port = 80
					    uri = "/"
					  }
					}`,
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(
								"leaseweb_public_cloud_target_group.test",
								plancheck.ResourceActionDestroyBeforeCreate,
							),
						},
					},
					// Ignore the inconsistent result as prism returns the old result.
					ExpectError: regexp.MustCompile(
						"Provider produced inconsistent result after apply",
					),
					Config: providerConfig + `
					resource "leaseweb_public_cloud_target_group" "test" {
					  name = "Target group name"
					  port = 80
					  region = "eu-west-2"
					  protocol = "HTTPS"
					  health_check = {
					    host = "example.com"
					    method = "GET"
					    protocol = "HTTP"
					    port = 80
					    uri = "/"
					  }
					}`,
				},
			},
		})
	})

	t.Run("removing health_check triggers replacement", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_target_group" "test" {
					  name = "Target group name"
					  port = 80
					  region = "eu-west-2"
					  protocol = "HTTP"
					  health_check = {
					    host = "example.com"
					    method = "GET"
					    protocol = "HTTP"
					    port = 80
					    uri = "/"
					  }
					}`,
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(
								"leaseweb_public_cloud_target_group.test",
								plancheck.ResourceActionDestroyBeforeCreate,
							),
						},
					},
					// Ignore the inconsistent result as prism returns the old result.
					ExpectError: regexp.MustCompile(
						"Provider produced inconsistent result after apply",
					),
					Config: providerConfig + `
					resource "leaseweb_public_cloud_target_group" "test" {
					  name = "Target group name"
					  port = 80
					  region = "eu-west-2"
					  protocol = "HTTP"
					}`,
				},
			},
		})
	})
}

func TestAccDedicatedServerResource(t *testing.T) {
	t.Run("imports and updates a server", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// ImportState testing
				{
					Config: providerConfig + `
					  resource "leaseweb_dedicated_server" "test" {
					  }
					  `,
					ResourceName:       "leaseweb_dedicated_server.test",
					ImportState:        true,
					ImportStatePersist: true,
					ImportStateId:      "123456",
				},
				// Update and Read testing
				{
					Config: providerConfig + `
					  resource "leaseweb_dedicated_server" "test" {
					  }
					  `,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"id",
							"123456",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"reference",
							"this is ref",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"reverse_lookup",
							"domain.example.com",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"dhcp_lease",
							"http://mirror.leaseweb.com/ipxe-files/ubuntu-18.04.ipxe",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"powered_on",
							"false",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"public_network_interface_opened",
							"false",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"public_ip_null_routed",
							"false",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"public_ip",
							"127.0.0.1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"remote_management_ip",
							"127.0.0.1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"internal_mac",
							"macmacmac",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"location.rack",
							"LK18",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"location.site",
							"AMS-02",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"location.suite",
							"HALL1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server.test",
							"location.unit",
							"12",
						),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		})
	})

	t.Run("creating a new server causes an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_dedicated_server" "test" {
					  reference    = "ref"
					}
					`,
					ExpectError: regexp.MustCompile(
						"Resource can only be imported, not created.",
					),
				},
			},
		})
	})
}

func TestAccPublicCloudISOsDataSource(t *testing.T) {
	t.Run("can read all ISOs", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `data "leaseweb_public_cloud_isos" "test" {}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_isos.test",
							"isos.#",
							"2",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_isos.test",
							"isos.0.id",
							"GRML",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_public_cloud_isos.test",
							"isos.0.name",
							"GRML 2022.11",
						),
					),
				},
			},
		})
	})
}

func TestAccPublicCloudIpResource(t *testing.T) {
	t.Run("imports and updates an ip", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// ImportState testing
				{
					Config: providerConfig + `
					  resource "leaseweb_public_cloud_ip" "test" {
					    instance_id    = "695ddd91-051f-4dd6-9120-938a927a47d0"
					    ip             = "10.0.0.1"
					    reverse_lookup = "a-valid-domain.xpto"
					  }
					  `,
					ResourceName:                         "leaseweb_public_cloud_ip.test",
					ImportState:                          true,
					ImportStatePersist:                   true,
					ImportStateId:                        "695ddd91-051f-4dd6-9120-938a927a47d0,10.0.0.1",
					ImportStateVerifyIdentifierAttribute: "instance_id",
					ImportStateCheck: func(states []*terraform.InstanceState) error {
						for _, state := range states {
							if state.Attributes["ip"] != "10.0.0.1" || state.Attributes["instance_id"] != "695ddd91-051f-4dd6-9120-938a927a47d0" {
								return fmt.Errorf("%v", state.Attributes)
							}
						}

						return nil
					},
				},
				// Update and Read testing
				{
					Config: providerConfig + `
					  resource "leaseweb_public_cloud_ip" "test" {
					    instance_id    = "695ddd91-051f-4dd6-9120-938a927a47d0"
					    ip             = "10.0.0.1"
					    reverse_lookup = "a-valid-domain.xpto"
					  }
					  `,
				},
			},
			// Delete testing automatically occurs in TestCase
		})
	})

	t.Run("creating a new ip causes an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_ip" "test" {
					  instance_id    = "695ddd91-051f-4dd6-9120-938a927a47d0"
					  ip             = "10.0.0.1"
					  reverse_lookup = "example.com"
					}
					`,
					ExpectError: regexp.MustCompile(
						"Resource can only be imported, not created.",
					),
				},
			},
		})
	})
}

func TestAccPublicCloudInstanceIsoResource(t *testing.T) {
	t.Run("can create/import/update/delete iso", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance_iso" "test" {
						desired_id = "ACRONIS_BOOT_MEDIA"
						instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance_iso.test",
							"name",
							"Acronis Boot Media",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance_iso.test",
							"id",
							"ACRONIS_BOOT_MEDIA",
						),
					),
				},
				// ImportState testing
				{
					ResourceName:                         "leaseweb_public_cloud_instance_iso.test",
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateId:                        "ace712e9-a166-47f1-9065-4af0f7e7fce1",
					ImportStateVerifyIdentifierAttribute: "instance_id",
				},
				// Update and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance_iso" "test" {
						desired_id = "ACRONIS_BOOT_MEDIA"
						instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance_iso.test",
							"name",
							"Acronis Boot Media",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance_iso.test",
							"id",
							"ACRONIS_BOOT_MEDIA",
						),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run("importing sets desired_id to id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					  resource "leaseweb_public_cloud_instance_iso" "test" {
						instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
					  }
					  `,
					ResourceName:                         "leaseweb_public_cloud_instance_iso.test",
					ImportState:                          true,
					ImportStateId:                        "ace712e9-a166-47f1-9065-4af0f7e7fce1",
					ImportStateVerifyIdentifierAttribute: "instance_id",
					ImportStatePersist:                   true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance_iso.test",
							"desired_id",
							"ACRONIS_BOOT_MEDIA",
						),
					),
				},
			},
		})
	})

	t.Run("detaching an ISO works as expected", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance_iso" "test" {
						desired_id = "ACRONIS_BOOT_MEDIA"
						instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
					}`,
				},
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance_iso" "test" {
						instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
					}`,
				},
			},
		})
	})

	t.Run("changing an attached ISO works as expected", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance_iso" "test" {
						desired_id = "ACRONIS_BOOT_MEDIA"
						instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
					}`,
				},
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance_iso" "test" {
						desired_id = "GRML"
						instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
					}`,
				},
			},
		})
	})

	t.Run("attaching a new ISO works as expected", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance_iso" "test" {
						instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
					}`,
				},
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance_iso" "test" {
						desired_id = "GRML"
						instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
					}`,
				},
			},
		})
	})

	t.Run(
		"updating to a non existent id returns an error",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_public_cloud_instance_iso" "test" {
							desired_id = "ACRONIS_BOOT_MEDIA"
							instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
						}`,
					},
					{
						Config: providerConfig + `
						resource "leaseweb_public_cloud_instance_iso" "test" {
							desired_id = "tralala"
							instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
						}`,
						ExpectError: regexp.MustCompile(
							`Attribute id value must be one of`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"creating a resource with a non existent id returns an error",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
						resource "leaseweb_public_cloud_instance_iso" "test" {
							desired_id = "tralala"
							instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
						}`,
						ExpectError: regexp.MustCompile(
							`Attribute id value must be one of`,
						),
					},
				},
			})
		},
	)

	t.Run("updating instance_id triggers replace", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance_iso" "test" {
						desired_id = "ACRONIS_BOOT_MEDIA"
						instance_id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"

					}`,
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(
								"leaseweb_public_cloud_instance_iso.test",
								plancheck.ResourceActionDestroyBeforeCreate,
							),
						},
					},
					Config: providerConfig + `
					resource "leaseweb_public_cloud_instance_iso" "test" {
						desired_id = "ACRONIS_BOOT_MEDIA"
						instance_id = "71d840af-f1c6-4c96-b82d-de4da216ee09"
					}`,
					ExpectError: regexp.MustCompile(
						"Provider produced inconsistent result after apply",
					),
				},
			},
		})
	})
}

func TestAccDnsResourceRecordSetsDataSource(t *testing.T) {
	t.Run("domain_name is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
        					data "leaseweb_dns_resource_record_sets" "test" {
        					}`,
					ExpectError: regexp.MustCompile(
						"The argument \"domain_name\" is required, but no definition was found",
					),
				},
			},
		})
	})

	t.Run("reading data succeeds", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
        					data "leaseweb_dns_resource_record_sets" "test" {
								domain_name = "example.com"
        					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_dns_resource_record_sets.test",
							"resource_record_sets.#",
							"13",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_dns_resource_record_sets.test",
							"resource_record_sets.0.name",
							"example.com.",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_dns_resource_record_sets.test",
							"resource_record_sets.0.type",
							"A",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_dns_resource_record_sets.test",
							"resource_record_sets.0.content.#",
							"3",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_dns_resource_record_sets.test",
							"resource_record_sets.0.content.0",
							"85.17.150.50",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_dns_resource_record_sets.test",
							"resource_record_sets.0.ttl",
							"300",
						),
					),
				},
			},
		})
	})
}

func TestAccDNSResourceRecordSetResource(t *testing.T) {
	t.Run("content is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
						        }`,
					ExpectError: regexp.MustCompile(
						"The argument \"content\" is required, but no definition was found",
					),
				},
			},
		})
	})
	t.Run("content cannot be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
									content = []
									domain_name = "example.com"
									name = "name"
									ttl = 3600
									type = "A"
						        }`,
					ExpectError: regexp.MustCompile(
						"Attribute content list must contain at least 1 elements, got: 0",
					),
				},
			},
		})
	})
	t.Run("content cannot contain null elements", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
									content = [null]
									domain_name = "example.com"
									name = "name"
									ttl = 3600
									type = "A"
						        }`,
					ExpectError: regexp.MustCompile(
						"This attribute contains a null value.",
					),
				},
			},
		})
	})

	t.Run("domain_name is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
						        }`,
					ExpectError: regexp.MustCompile(
						"The argument \"domain_name\" is required, but no definition was found",
					),
				},
			},
		})
	})
	t.Run("domain_name cannot be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
									content = ["1.2.3.4"]
									domain_name = ""
									name = "name"
									ttl = 3600
									type = "A"
						        }`,
					ExpectError: regexp.MustCompile(
						"Attribute domain_name string length must be at least 1, got: 0",
					),
				},
			},
		})
	})

	t.Run("name is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
						        }`,
					ExpectError: regexp.MustCompile(
						"The argument \"name\" is required, but no definition was found",
					),
				},
			},
		})
	})
	t.Run("name cannot be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
									content = ["1.2.3.4"]
									domain_name = "example.com"
									name = ""
									ttl = 3600
									type = "A"
						        }`,
					ExpectError: regexp.MustCompile(
						"Attribute name string length must be at least 1, got: 0",
					),
				},
			},
		})
	})
	t.Run("name must end with a dot", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
									content = ["1.2.3.4"]
									domain_name = "example.com"
									name = "name"
									ttl = 3600
									type = "A"
						        }`,
					ExpectError: regexp.MustCompile(
						"Attribute name must end in ., got: name",
					),
				},
			},
		})
	})

	t.Run("ttl is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
						        }`,
					ExpectError: regexp.MustCompile(
						"The argument \"ttl\" is required, but no definition was found",
					),
				},
			},
		})
	})
	t.Run("ttl must be valid", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
									content = ["1.2.3.4"]
									domain_name = "example.com"
									name = "example.com."
									ttl = 123
									type = "A"
						        }`,
					ExpectError: regexp.MustCompile(
						"Attribute ttl value must be one of:",
					),
				},
			},
		})
	})

	t.Run("type is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
						        }`,
					ExpectError: regexp.MustCompile(
						"The argument \"type\" is required, but no definition was found",
					),
				},
			},
		})
	})
	t.Run("type must be valid", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						        resource "leaseweb_dns_resource_record_set" "test" {
									content = ["1.2.3.4"]
									domain_name = "example.com"
									name = "example.com."
									ttl = 3600
									type = "tralala"
						        }`,
					ExpectError: regexp.MustCompile(
						"Attribute type value must be one of:",
					),
				},
			},
		})
	})

	t.Run("can create/import/update/delete dns record set", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_dns_resource_record_set" "test" {
						content = ["85.17.150.51", "85.17.150.52", "85.17.150.53"]
						domain_name = "example.com"
						name = "example.com."
						ttl = 3600
						type = "A"
					}`,
				},
				// ImportState testing
				{
					ResourceName:                         "leaseweb_dns_resource_record_set.test",
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateId:                        "example.com,example.com.,A",
					ImportStateVerifyIdentifierAttribute: "domain_name",
					ImportStateCheck: func(states []*terraform.InstanceState) error {
						for _, state := range states {
							if state.Attributes["domain_name"] != "example.com" || state.Attributes["name"] != "example.com." || state.Attributes["type"] != "A" {
								return fmt.Errorf("%v", state.Attributes)
							}
						}

						return nil
					},
				},
				// Update and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_dns_resource_record_set" "test" {
						content = ["85.17.150.51", "85.17.150.52", "85.17.150.53"]
						domain_name = "example.com"
						name = "example.com."
						ttl = 3600
						type = "A"
					}`,
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run("updating name triggers replacement", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_dns_resource_record_set" "test" {
						content = ["85.17.150.51", "85.17.150.52", "85.17.150.53"]
						domain_name = "example.com"
						name = "example.com."
						ttl = 3600
						type = "A"
					}
					`,
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(
								"leaseweb_dns_resource_record_set.test",
								plancheck.ResourceActionDestroyBeforeCreate,
							),
						},
					},
					// Ignore the inconsistent result as prism returns the old result.
					ExpectError: regexp.MustCompile(
						"Provider produced inconsistent result after apply",
					),
					Config: providerConfig + `
					resource "leaseweb_dns_resource_record_set" "test" {
						content = ["85.17.150.51", "85.17.150.52", "85.17.150.53"]
						domain_name = "example.com"
						name = "tralala."
						ttl = 3600
						type = "A"
					}
					`,
				},
			},
		})
	})
}

func TestAccIPmgmtIpsDataSource(t *testing.T) {
	t.Run("data source works", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `data "leaseweb_ipmgmt_ips" "test" {}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.#",
							"2",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.assigned_contract.id",
							"5643634",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.equipment_id",
							"1234",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.ip",
							"192.0.2.1",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.null_routed",
							"false",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.prefix_length",
							"32",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.primary",
							"true",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.reverse_lookup",
							"mydomain1.example.com",
						),

						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.subnet.gateway",
							"192.0.2.254",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.subnet.id",
							"192.0.2.0_24",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.subnet.network_ip",
							"192.0.2.0",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.subnet.prefix_length",
							"24",
						),

						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.type",
							"NORMAL_IP",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.unnulling_allowed",
							"false",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_ips.test",
							"ips.0.version",
							"4",
						),
					),
				},
			},
		})
	})

	t.Run("can filter assigned_contract_ids", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							assigned_contract_ids = ["1", "2", "3"]
						}
					`,
				},
			},
		})
	})

	t.Run("can filter equipment_ids", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							equipment_ids = ["1", "2", "3"]
						}
					`,
				},
			},
		})
	})

	t.Run("can filter filtered_ips", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							filtered_ips = ["192.168.1.1", "192.168.1.2"]
						}
					`,
				},
			},
		})
	})

	t.Run("can filter from_ip", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							from_ip = "192.168.0.0"
						}
					`,
				},
			},
		})
	})

	t.Run("can filter null_routed", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							null_routed = true
						}
					`,
				},
			},
		})
	})

	t.Run("can filter primary", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							primary = true
						}
					`,
				},
			},
		})
	})

	t.Run("can filter reverse_lookup", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							reverse_lookup = "mydomain1.example.com"
						}
					`,
				},
			},
		})
	})

	t.Run("can filter subnet_id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							subnet_id = "1234"
						}
					`,
				},
			},
		})
	})

	t.Run("can filter to_ip", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							to_ip = "192.168.255.255"
						}
					`,
				},
			},
		})
	})

	t.Run("can sort", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							sort = ["ip"]
						}
					`,
				},
			},
		})
	})

	t.Run("can filter type", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							type = "NORMAL_IP"
						}
					`,
				},
			},
		})
	})
	t.Run("inputting incorrect type throws error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							type = "tralala"
						}
					`,
					ExpectError: regexp.MustCompile(
						`Attribute type value must be one of:`,
					),
				},
			},
		})
	})

	t.Run("can filter version", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							version = 6
						}
					`,
				},
			},
		})
	})
	t.Run("inputting incorrect IP version throws error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_ips" "test" {
							version = 65
						}
					`,
					ExpectError: regexp.MustCompile(
						`Attribute version value must be one of:`,
					),
				},
			},
		})
	})
}

func TestIPMgmtIPResourceResource(t *testing.T) {
	t.Run("creating a new IP throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_ipmgmt_ip" "test" {
					  ip = "1.2.3.4"
					}
					`,
					ExpectError: regexp.MustCompile(
						"Resource can only be imported, not created.",
					),
				},
			},
		})
	})

	t.Run("imports and updates an IP", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// ImportState testing
				{
					Config: providerConfig + `
					  resource "leaseweb_ipmgmt_ip" "test" {
					    ip    = "192.0.2.1"
					  }
					  `,
					ResourceName:                         "leaseweb_ipmgmt_ip.test",
					ImportState:                          true,
					ImportStatePersist:                   true,
					ImportStateId:                        "192.0.2.1",
					ImportStateVerifyIdentifierAttribute: "ip",
				},
				// Update and Read testing
				{
					Config: providerConfig + `
						resource "leaseweb_ipmgmt_ip" "test" {
						    ip    = "192.0.2.1"
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"assigned_contract.id",
							"5643634",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"equipment_id",
							"1234",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"ip",
							"192.0.2.1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"null_level",
							"0",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"null_routed",
							"false",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"prefix_length",
							"32",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"primary",
							"true",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"reverse_lookup",
							"mydomain1.example.com",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"subnet.gateway",
							"192.0.2.254",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"subnet.id",
							"192.0.2.0_24",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"subnet.network_ip",
							"192.0.2.0",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"subnet.prefix_length",
							"24",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"type",
							"NORMAL_IP",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"unnulling_allowed",
							"false",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_ip.test",
							"version",
							"4",
						),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		})
	})
}
func TestAccIPmgmtNullRouteHistoryDataSource(t *testing.T) {
	t.Run("data source works", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `data "leaseweb_ipmgmt_null_route_history" "test" {}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_null_route_history.test",
							"nullroutes.#",
							"2",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_null_route_history.test",
							"nullroutes.0.assigned_contract.id",
							"123456",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_null_route_history.test",
							"nullroutes.0.automatic_unnulling_at",
							"2015-06-28 13:00:00 +0000 UTC",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_null_route_history.test",
							"nullroutes.0.comment",
							"This IP is evil",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_null_route_history.test",
							"nullroutes.0.equipment_id",
							"456",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_null_route_history.test",
							"nullroutes.0.id",
							"4534536",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_null_route_history.test",
							"nullroutes.0.ip",
							"192.0.2.1",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_null_route_history.test",
							"nullroutes.0.nulled_at",
							"2015-06-28 12:00:00 +0000 UTC",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_null_route_history.test",
							"nullroutes.0.nulled_by",
							"john.doe@example.com",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_null_route_history.test",
							"nullroutes.0.null_level",
							"1",
						),
						resource.TestCheckResourceAttr(
							"data.leaseweb_ipmgmt_null_route_history.test",
							"nullroutes.0.ticket_id",
							"188612",
						),
					),
				},
			},
		})
	})

	t.Run("can filter contract_id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_null_route_history" "test" {
							contract_id = "123"
						}
					`,
				},
			},
		})
	})

	t.Run("can filter equipment_id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_null_route_history" "test" {
							equipment_id = "123"
						}
					`,
				},
			},
		})
	})

	t.Run("can filter from_date", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_null_route_history" "test" {
							from_date = "2022-12-02"
						}
					`,
				},
			},
		})
	})

	t.Run("can filter from_ip", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_null_route_history" "test" {
							from_ip = "123"
						}
					`,
				},
			},
		})
	})

	t.Run("can filter nulled_by", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_null_route_history" "test" {
							nulled_by = "john@example.com"
						}
					`,
				},
			},
		})
	})

	t.Run("can sort", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_null_route_history" "test" {
							sort = ["ip"]
						}
					`,
				},
			},
		})
	})

	t.Run("can filter ticket_id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_null_route_history" "test" {
							ticket_id = "1234"
						}
					`,
				},
			},
		})
	})

	t.Run("can filter to_date", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_null_route_history" "test" {
							to_date = "2022-12-02"
						}
					`,
				},
			},
		})
	})

	t.Run("can filter to_ip", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_null_route_history" "test" {
							to_ip = "123"
						}
					`,
				},
			},
		})
	})

	t.Run("can filter unnulled_by", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
						data "leaseweb_ipmgmt_null_route_history" "test" {
							unnulled_by = "john@example.com"
						}
					`,
				},
			},
		})
	})
}

func TestAccIpmgmtNullRouteResource(t *testing.T) {
	t.Run("creates and updates a null route", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
					  ip = "192.0.2.1"
					}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_null_route.test",
							"id",
							"4534536",
						),

						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_null_route.test",
							"assigned_contract.id",
							"123456",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_null_route.test",
							"automatic_unnulling_at",
							"2015-06-25 11:13:00 +0000 UTC",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_null_route.test",
							"comment",
							"This IP is evil",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_null_route.test",
							"equipment_id",
							"456",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_null_route.test",
							"nulled_at",
							"2015-06-28 12:00:00 +0000 UTC",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_null_route.test",
							"nulled_by",
							"john.doe@example.com",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_null_route.test",
							"null_level",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_ipmgmt_null_route.test",
							"ticket_id",
							"188612",
						),
					),
				},
				// ImportState testing
				{
					ResourceName:      "leaseweb_ipmgmt_null_route.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				// Update and Read testing
				{
					Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
						id = "4534536"
					}
					`,
				},
			},

			// Delete testing automatically occurs in TestCase
		})
	})

	t.Run("invalid automatic_nulling_at throws an error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
						ip = "192.0.2.1"
						automatic_unnulling_at = "tralala"
					}
					`,
					ExpectError: regexp.MustCompile(
						"Attribute automatic_unnulling_at must be specified using the RFC3339 format",
					),
				},
			},
		})
	})

	t.Run("ip must be set when creating a null route", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
					}
					`,
					ExpectError: regexp.MustCompile(
						"Attribute ip value must be set to create a null route",
					),
				},
			},
		})
	})

	t.Run(
		"automatic_unnulling_at can be set when creating a null route",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
						ip = "192.0.2.1"
						automatic_unnulling_at = "2015-06-25 11:13:00 +0000 UTC"
					}
					`,
					},
				},
			})
		},
	)

	t.Run("comment can be set when creating a null route", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
						ip = "192.0.2.1"
						comment = "This IP is evil"
					}
					`,
				},
			},
		})
	})

	t.Run("ticket_id can be set when creating a null route", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
						ip = "192.0.2.1"
						ticket_id = "188612"
					}
					`,
				},
			},
		})
	})

	t.Run(
		"automatic_unnulling_at can be set when updating a null route",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
						ip = "192.0.2.1"
					}
					`,
					},
					{
						Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
						ip = "192.0.2.1"
						automatic_unnulling_at = "2015-06-25 11:13:00 +0000 UTC"
					}
					`,
					},
				},
			})
		},
	)

	t.Run("comment can be set when updating a null route", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
						ip = "192.0.2.1"
					}
					`,
				},
				{
					Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
						ip = "192.0.2.1"
						comment = "This IP is evil"
					}
					`,
				},
			},
		})
	})

	t.Run("ticket_id can be set when updating a null route", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
						ip = "192.0.2.1"
					}
					`,
				},
				{
					Config: providerConfig + `
					resource "leaseweb_ipmgmt_null_route" "test" {
						ip = "192.0.2.1"
						ticket_id = "188612"
					}
					`,
				},
			},
		})
	})
}
