# Changelog

## [1.26.1](https://github.com/leaseweb/terraform-provider-leaseweb/compare/v1.26.0...v1.26.1) (2025-04-30)


### Bug Fixes

* replace incorrect LoadBalancer struct with correct one ([#355](https://github.com/leaseweb/terraform-provider-leaseweb/issues/355)) ([701442b](https://github.com/leaseweb/terraform-provider-leaseweb/commit/701442bb2481d14eaa7ccf06d53f2afac67cf0f1))

## [1.26.0](https://github.com/leaseweb/terraform-provider-leaseweb/compare/v1.25.0...v1.26.0) (2025-04-04)


### Features

* added hostname and raid in dedicated server installation state ([#341](https://github.com/leaseweb/terraform-provider-leaseweb/issues/341)) ([b0523d7](https://github.com/leaseweb/terraform-provider-leaseweb/commit/b0523d78245885f5e9f449f348d98dcd3310cc0c))

## [1.24.0](https://github.com/leaseweb/terraform-provider-leaseweb/compare/v1.23.1...v1.24.0) (2025-02-10)


### Features

* wait until installation jobs finished ([#298](https://github.com/leaseweb/terraform-provider-leaseweb/issues/298)) ([44af5ea](https://github.com/leaseweb/terraform-provider-leaseweb/commit/44af5ea619b946d6637402e752b39ca9fb5d8601))

## [1.23.1](https://github.com/leaseweb/terraform-provider-leaseweb/compare/v1.23.0...v1.23.1) (2025-02-07)


### Bug Fixes

* **logging:** log errorMessage for unauthorized and relevant errors ([#291](https://github.com/leaseweb/terraform-provider-leaseweb/issues/291)) ([8753148](https://github.com/leaseweb/terraform-provider-leaseweb/commit/8753148597824a1745fcfaa76da52d8d2cbb6168))

## [1.23.0](https://github.com/leaseweb/terraform-provider-leaseweb/compare/v1.22.0...v1.23.0) (2025-01-27)


### Features

* revert post_install_script in dedicated server installation to accept bash script ([#280](https://github.com/leaseweb/terraform-provider-leaseweb/issues/280)) ([10ceace](https://github.com/leaseweb/terraform-provider-leaseweb/commit/10ceace5c7da8a28d7ef38e75c35143f636a8078))

## [1.22.0](https://github.com/leaseweb/terraform-provider-leaseweb/compare/v1.21.2...v1.22.0) (2025-01-22)


### Features

* add ips to load balancer resource ([#276](https://github.com/leaseweb/terraform-provider-leaseweb/issues/276)) ([bb4c6d7](https://github.com/leaseweb/terraform-provider-leaseweb/commit/bb4c6d7fbc23f034f38bfd373d22245fbdb550aa))

## [1.21.2](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.21.1...v1.21.2) (2025-01-15)


### Bug Fixes

* pin sdk version ([#264](https://github.com/Leaseweb/terraform-provider-leaseweb/issues/264)) ([34e3719](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/34e37197f8b7417b52616e51be8c16f373e40321))

## [1.21.0](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.20.1...v1.21.0) (2025-01-14)


### Features

* add null route IP resource ([#247](https://github.com/Leaseweb/terraform-provider-leaseweb/issues/247)) ([b616987](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/b61698729bf8c9bd50a45a851c136dee00652bdd))

## [1.20.1](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.20.0...v1.20.1) (2025-01-10)


### Bug Fixes

* import dedicated server ([#257](https://github.com/Leaseweb/terraform-provider-leaseweb/issues/257)) ([55cf551](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/55cf5516a2fd07249c0b294d703001bf79331268))

## [1.20.0](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.19.0...v1.20.0) (2024-12-31)


### Features

* implement null route history data source ([1c2df56](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/1c2df56ec8b1f3e1e9e1d0b47ce71c7cd287101d))

## [1.19.0](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.18.0...v1.19.0) (2024-12-30)


### Features

* implement ip resource ([7547b39](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/7547b391c31555077b18856f4b12f645f62e9f5d))

## [1.18.0](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.17.0...v1.18.0) (2024-12-30)


### Features

* implement resource record set resource ([8246aaa](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/8246aaa5f6fc2f4ced1e4c8159b935c919e38979))

## [1.17.0](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.16.0...v1.17.0) (2024-12-24)


### Features

* implement dns resource record sets data source ([9214c36](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/9214c361d8151940974b0480b32f59f5d2ac3663))

## [1.16.0](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.15.0...v1.16.0) (2024-12-12)


### Features

* implement public cloud ISO resource ([e3cf3f5](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/e3cf3f5981226eedf56f685d9e011355a0671c86))

## [1.15.0](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.14.2...v1.15.0) (2024-12-09)


### Features

* add ISO to instances data source ([122f662](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/122f6620359e61f6e85f0610cf5ba61946695de5))

## [1.14.2](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.14.1...v1.14.2) (2024-12-06)


### Bug Fixes

* first & second action can never be equal when generating unsupported action string ([815078b](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/815078bfb712aba7ff6b03046a937c3b85d5dc77))

## [1.14.1](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.14.0...v1.14.1) (2024-12-06)


### Bug Fixes

* fix typo ([39c371d](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/39c371d98cf565c91e8ba53cf26750c3cfa4f4b2))

## [1.14.0](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.13.0...v1.14.0) (2024-12-03)


### Features

* implement public cloud isos data source ([#188](https://github.com/Leaseweb/terraform-provider-leaseweb/issues/188)) ([ae82e33](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/ae82e339e724117ef3fa013017fe75ddb08ad86c))

## [1.13.0](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.12.1...v1.13.0) (2024-11-29)


### Features

* tell user if server takes too long to respond ([72c86e7](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/72c86e710259b1c19687ac6336a99d4ab7205808))

## [1.12.1](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.12.0...v1.12.1) (2024-11-26)


### Bug Fixes

* add missing image.instance_id to schema ([8408ec5](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/8408ec5cbc3f9671b514dae64165cb4e8bdf45fd))

## [1.12.0](https://github.com/Leaseweb/terraform-provider-leaseweb/compare/v1.11.0...v1.12.0) (2024-11-22)


### Features

* implement target group resource ([466a2a8](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/466a2a8448f1988aa1f49d60976f5b7b325ad02f))
* inject release version into user agent ([f9b4b8b](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/f9b4b8b8cc79cec8bc1b45be2b281d26297d2b3d))


## 1.11.0 (November 19, 2024)

* feat(docs): let end user know which actions aren't supported by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/161

## 1.10.0 (November 19, 2024)

* chore(deps-dev): bump @stoplight/prism-cli from 5.10.0 to 5.12.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/155
* chore(deps): bump stoplight/prism from 5.11.2 to 5.12.1 in /docker/prism by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/157
* chore(deps): bump caddy from `5d4e02f` to `e97e0e3` in /docker/caddy by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/156
* chore(test): remove useless acceptance tests by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/159
* chore: bump node to jod-lts by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/160
* chore(deps-dev): bump husky from 9.1.6 to 9.1.7 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/162
* feat: implement load balancer listener resource by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/158

## 1.9.0 (November 11, 2024)

* Implement instances & new terraform framework by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/14
* chore(deps): bump github.com/hashicorp/terraform-plugin-testing from 1.8.0 to 1.9.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/21
* chore(deps): bump actions/checkout from 3 to 4 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/19
* chore(deps): bump golangci/golangci-lint-action from 6.0.1 to 6.1.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/18
* chore(deps): bump goreleaser/goreleaser-action from 5.1.0 to 6.0.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/17
* chore(deps): bump pnpm/action-setup from 3 to 4 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/15
* chore(deps): bump actions/setup-go from 5.0.1 to 5.0.2 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/16
* chore(deps): bump github.com/hashicorp/terraform-plugin-docs from 0.19.2 to 0.19.4 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/20
* Feat/implement new image fields by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/22
* feat: implement new image fields by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/23
* Implement storageSize by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/24
* Fix module name by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/25
* feat: implement badge by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/26
* fix: replace disable with skip by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/27
* feat: implement extra type fields by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/28
* ci: enable automated changelog generation by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/30
* feat: implement caching for instanceType by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/29
* chore(deps): bump github.com/hashicorp/terraform-plugin-framework from 1.10.0 to 1.11.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/32
* feat: implement thread safe caching by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/31
* Fix/cannot create instance with ssh key by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/33
* chore(deps): bump github.com/hashicorp/terraform-plugin-testing from 1.9.0 to 1.10.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/34
* chore: bump sdk by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/35
* Refactor/move domain entities by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/36
* chore: remove uuid value_object by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/38
* Validate immutable fields by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/39
* feat: implement dedicated server data source by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/37
* fix: return instances from service in a predictable order by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/40
* chore(deps): bump hashicorp/setup-terraform from 3.1.1 to 3.1.2 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/42
* feat: add server attributes in leaseweb_dedicated_server by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/41
* Prevent instance termination in certain states by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/43
* Disable openssh key support temporarily by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/46
* feat: implement region.location field by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/48
* Docs/update docs by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/50
* docs: add CONTRIBUTING.md by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/51
* feat: implement dedicated_server control panels datasource by @arunvarghe in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/45
* feat: implement dedicated_server_operating_systems data source by @imtiazPabel in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/52
* docs: add package docs by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/53
* feat: implement dedicated_server_credential data source by @imtiazPabel in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/64
* feat: add dedicated server data traffic notification setting resource by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/58
* feat: update data source for dedicated server(s) by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/67
* feat: add dedicated server bandwidth notification setting resource by @arunvarghe in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/65
* feat: import dedicated server by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/68
* feat: add control panels & operating systems by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/69
* chore: remove old stuff by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/70
* docs: fix public cloud instance based on new SDK by @arunvarghe in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/72
* feat: add dedicated server credential resource by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/73
* chore: change license from MPL 2.0 to Apache 2.0 by @imtiazPabel in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/79
* feat: add dedicated server installation resource by @imtiazPabel in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/78
* feat: add dedicated server resource by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/81
* chore: api stability section in README.md by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/80
* chore(deps): bump github.com/hashicorp/terraform-plugin-framework from 1.11.0 to 1.12.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/75
* chore(deps): bump github.com/hashicorp/terraform-plugin-go from 0.23.0 to 0.24.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/76
* chore: add error message helper function to show error message details by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/82
* chore(deps): bump golangci/golangci-lint-action from 6.1.0 to 6.1.1 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/83
* Fix public cloud by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/86
* Minimize attributes by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/88
* Minimize attributes step 2 by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/89
* chore: enable automated node package updates by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/90
* chore(deps-dev): bump husky from 9.0.11 to 9.1.6 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/91
* chore(deps-dev): bump prettier from 3.2.5 to 3.3.3 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/92
* chore(deps-dev): bump @commitlint/config-conventional from 19.2.2 to 19.5.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/94
* chore(deps-dev): bump markdownlint-cli from 0.40.0 to 0.42.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/93
* chore(deps-dev): bump @stoplight/prism-cli from 5.8.1 to 5.10.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/95
* chore: enable automated docker dependency updates by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/96
* chore(deps-dev): bump lint-staged from 15.2.4 to 15.2.10 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/98
* chore(deps-dev): bump @commitlint/cli from 19.3.0 to 19.5.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/97
* fix: set directories where Dockerfiles are found by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/99
* chore(deps): bump node from `cffed8c` to `967bab2` in /docker/node by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/101
* chore(deps): bump golang from 1.22.3-bookworm to 1.23.2-bookworm in /docker/golang by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/100
* chore(deps): bump node from `967bab2` to `5bba9bf` in /docker/node by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/104
* chore(deps): bump node from `5bba9bf` to `ec35a66` in /docker/node by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/105
* chore(deps): bump github.com/hashicorp/terraform-plugin-framework-validators from 0.13.0 to 0.14.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/106
* Implement golang best practices for publiccloud package by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/103
* fix: update dedicatedserver to work with latest SDK by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/109
* code quality changes by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/111
* chore(deps): bump actions/setup-go from 5.0.2 to 5.1.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/116
* Code architecture by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/113
* feat: implement public_cloud_credential data source by @imtiazPabel in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/115
* Implement publiccloud image DataSource by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/108
* feat: add public cloud credential resource by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/107
* chore(deps): bump crazy-max/ghaction-import-gpg from 6.1.0 to 6.2.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/119
* Simplify code & add markdown helpers to instanceResource by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/118
* publiccloud image resource by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/110
* Public Cloud load balancer data source by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/112
* chore: set user agent in terraform plugin by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/122
* chore: update changelog by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/121
* Implement latest changes from the scaffolding by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/120
* chore(deps): bump github.com/stretchr/testify from 1.8.3 to 1.9.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/123
* Implement workaround so dependabot updates docker images by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/125
* chore(deps): bump caddy from `f8660fb` to `33937b9` in /docker/caddy by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/127
* chore(deps): bump golangci/golangci-lint from v1.58.2-alpine to v1.61.0-alpine in /docker/golangci by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/128
* chore(deps): bump stoplight/prism from 5.8.1 to 5.11.2 in /docker/prism by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/129
* fix: set the proper hostnames in the prism-proxy Caddyfile by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/131
* chore(deps): bump github.com/hashicorp/terraform-plugin-go from 0.24.0 to 0.25.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/132
* Map PublicCloudAPI to client by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/130
* refactor: simplify pagination by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/133
* Refactor/unified error message by @imtiazPabel in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/134
* chore(deps): bump github.com/hashicorp/terraform-plugin-framework-validators from 0.14.0 to 0.15.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/136
* Convert dedicatedserver to new structure by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/126
* Map API validation to resource attributes by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/124
* chore(deps): bump golang from 1.23.2-bookworm to 1.23.3-bookworm in /docker/golang by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/139
* chore(deps): bump github.com/hashicorp/terraform-plugin-docs from 0.19.4 to 0.20.0 in /tools by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/140
* fix: use returned request so new offset is properly set by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/138
* fix: implement instance_id attribute to create image resource by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/137
* Fix merge issues by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/142
* chore(deps): bump goreleaser/goreleaser-action from 6.0.0 to 6.1.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/144
* Fix/load balancer datasource by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/145
* Public Cloud loadBalancer resource by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/117
* chore(deps): bump golangci/golangci-lint from v1.61.0-alpine to v1.62.0-alpine in /docker/golangci by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/149
* Refactor/unified error method by @imtiazPabel in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/143
* feat: implement load balancer listeners data source by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/141
* feat: implement target groups data source by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/146


## 1.8.0-beta (October 15, 2024)

* Remove attributes not in instanceDetails by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/88
* Only show minimal attributes to user by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/89

## 1.7.1-beta (October 8, 2024)

* chore(deps): bump github.com/hashicorp/terraform-plugin-framework from 1.11.0 to 1.12.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/75
* chore(deps): bump github.com/hashicorp/terraform-plugin-go from 0.23.0 to 0.24.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/76
* chore: add error message helper function to show error message details by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/82
* chore(deps): bump golangci/golangci-lint-action from 6.1.0 to 6.1.1 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/83
* Fix public cloud by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/86

## 1.7.0-beta (October 1, 2024)

* feat: implement dedicated_server_credential data source by @imtiazPabel in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/64
* feat: add dedicated server data traffic notification setting resource by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/58
* feat: update data source for dedicated server(s) by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/67
* feat: add dedicated server bandwidth notification setting resource by @arunvarghe in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/65
* feat: import dedicated server by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/68
* feat: add control panels & operating systems by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/69
* docs: fix public cloud instance based on new SDK by @arunvarghe in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/72
* feat: add dedicated server credential resource by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/73
* chore: change license from MPL 2.0 to Apache 2.0 by @imtiazPabel in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/79
* feat: add dedicated server installation resource by @imtiazPabel in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/78
* feat: add dedicated server resource by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/81
* chore: api stability section in README.md by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/80

## 1.6.1-alpha (August 28, 2024)

* docs: add package docs by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/53

## 1.6.0-alpha (August 28, 2024)

* feat: implement dedicated_server_operating_systems data source by @imtiazPabel in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/52

## 1.5.0-alpha (August 28, 2024)

* docs: add CONTRIBUTING.md by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/51
* feat: implement dedicated_server control panels datasource by @arunvarghe in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/45

## 1.4.2-alpha (August 27, 2024)

* Docs/update docs by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/50

## 1.4.1-alpha (August 27, 2024)

* feat: implement region.location field by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/48

## 1.4.0-alpha (August 27, 2024)

* Prevent instance termination in certain states by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/43
* Disable openssh key support temporarily by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/46

## 1.3.2-alpha (August 23, 2024)

* feat: add server attributes in leaseweb_dedicated_server

## 1.3.1-alpha (August 14, 2024)

* feat: implement dedicated server data source by @majidkarimizadeh in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/37
* fix: return instances from service in a predictable order by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/40

## 1.3.0-alpha (August 13, 2024)

* chore: remove uuid value_object by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/38
* Validate immutable fields by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/39

## 1.2.0-alpha (August 5, 2024)

* feat: implement extra type fields by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/28
* ci: enable automated changelog generation by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/30
* feat: implement caching for instanceType by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/29
* chore(deps): bump github.com/hashicorp/terraform-plugin-framework from 1.10.0 to 1.11.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/32
* feat: implement thread safe caching by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/31
* Fix/cannot create instance with ssh key by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/33
* chore(deps): bump github.com/hashicorp/terraform-plugin-testing from 1.9.0 to 1.10.0 by @dependabot in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/34
* chore: bump sdk by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/35
* Refactor/move domain entities by @alrayyes in https://github.com/Leaseweb/terraform-provider-leaseweb/pull/36

## 1.1.1-alpha (August 5, 2024)

* [fix: replace disable with skip](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/58f4d93d40602bd41f5ffd0690f23af22147a7fc)

## 1.1.0-alpha (August 5, 2024)

* [feat: implement badge](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/6db262b6f7439356730cfd57bed6aa0a1699fab9)
* [feat: implement storageSize](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/f70f68040eacb87ae75297b487bf36b0c6847f6f)
* [test: add example for AdaptNullableBoolToBoolValue()](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/4ab3a39c57b2a41daf8235aab5942c5c8571e633)
* [feat: implement new image fields](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/466a7f106f117dc5a7b403a4c3676c5a942bce97)
* [feat: implement volume & architecture](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/b2018c66dcd46ff12c039548cf92a1e3da862172)
* [chore: update docs](https://github.com/Leaseweb/terraform-provider-leaseweb/commit/adf180a95af83c06da54860d5aaf7c722a461b84)
* update dependencies

## 1.0.0-alpha (August 1, 2024)

* remove old legacy code
* implement [Terraform Plugin Framework](https://github.com/hashicorp/terraform-provider-scaffolding-framework)
* implement public cloud instance management
* implement git pre push hooks
* implement [prettier](https://prettier.io/)
* implement [markdownlint](https://github.com/igorshubovych/markdownlint-cli)
* implement [ShellCheck](https://www.shellcheck.net/)
* implement [commitlint](https://commitlint.js.org/)
* implement [lint-staged](https://github.com/lint-staged/lint-staged)
* implement [prism](https://stoplight.io/open-source/prism)
* implement public cloud instances provider

## 0.4.1 (March 18, 2024)

Fix release

## 0.4.0 (March 18, 2024)

* Implement a `internal_mac` field for `dedicated_server` to be able to use the MAC address of the internal interface in `post_install_script` of `leaseweb_dedicated_server_installation`
* Implement data source to fetch credentials, such as those for REMOTE_MANAGEMENT to store them in 1Password via Terraform

## 0.3.3 (September 12, 2023)

Minor release

* ENHANCEMENT: Add all available filter options to data source `data_source_dedicated_servers`.
* CHANGE: Bump Leaseweb Go SDK to v0.0.3

## 0.3.2 (July 19, 2023)

Fix release

## 0.3.1 (July 14, 2023)

Fix release

## 0.3.0 (July 14, 2023)

* Bump Leaseweb Go SDK to v.0.2.0
* Add User-Agent string to API calls

## 0.2.0 (April 14, 2023)

Integrated Leaseweb SDK to handle API calls.

## 0.1.2 (November 18, 2022)

Add some logging and improve error messages

## 0.1.1 (September 19, 2022)

Fix release

## 0.1.0 (September 16, 2022)

FEATURES:

* **New Data Source:** `data_source_dedicated_server_control_panels`
* **New Data Source:** `data_source_dedicated_server_operating_systems`
* **New Data Source:** `data_source_dedicated_servers`
* **New Resource:** `resource_dedicated_server`
* **New Resource:** `resource_dedicated_server_credential`
* **New Resource:** `resource_dedicated_server_installation`
* **New Resource:** `resource_dedicated_server_notification_setting_bandwidth`
* **New Resource:** `resource_dedicated_server_notification_setting_datatraffic`
