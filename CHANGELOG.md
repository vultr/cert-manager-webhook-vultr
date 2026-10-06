# Change Log
## [v0.4.3](https://github.com/vultr/cert-manager-webhook-vultr) (2026-10-06)
* Match the GoReleaser Docker publisher to the explicit binary build ID.

## [v0.4.2](https://github.com/vultr/cert-manager-webhook-vultr) (2026-10-06)
* Correct the DHI pull-through registry hostname.
* Honor the release version from the release commit when creating the tag.

## [v0.4.1](https://github.com/vultr/cert-manager-webhook-vultr) (2026-10-06)
* Fix the GoReleaser v2 release configuration and Docker image publishing.

## [v0.4.0](https://github.com/vultr/cert-manager-webhook-vultr) (2026-10-06)
* Update Go to 1.27 and upgrade dependencies.
* Use DHI pull-through images for Docker builds.
* Isolate Vultr API clients per challenge and improve DNS zone handling.
* Document supported cert-manager DNS-01 features and add unit coverage.

## [v0.3.1](https://github.com/vultr/cert-manager-webhook-vultr) (2022-03-25)
* Bump k8s.io/client-go from 0.23.1 to 0.23.5  [PR 34](https://github.com/vultr/cert-manager-webhook-vultr/pull/34) 
* Bump k8s.io/apiextensions-apiserver from 0.23.1 to 0.23.5 [PR 36](https://github.com/vultr/cert-manager-webhook-vultr/pull/36) 
* Bump github.com/vultr/govultr/v2 from 2.14.1 to 2.14.2 [PR 37](https://github.com/vultr/cert-manager-webhook-vultr/pull/37) 
* Bump github.com/jetstack/cert-manager from 1.7.1 to 1.7.2 [PR 38](https://github.com/vultr/cert-manager-webhook-vultr/pull/38) 

## [v0.3.0](https://github.com/vultr/cert-manager-webhook-vultr) (2022-02-14)
* Bump github.com/jetstack/cert-manager from 1.6.1 to 1.7.1 [PR 29](https://github.com/vultr/cert-manager-webhook-vultr/pull/29) 
* Bump github.com/vultr/govultr/v2 from 2.12.0 to 2.14.1 [PR 28](https://github.com/vultr/cert-manager-webhook-vultr/pull/28) 


## [v0.2.0](https://github.com/vultr/cert-manager-webhook-vultr) (2021-12-16)
* Bumped Kubernetes to to 1.22.4, Cert-manager to 1.6.1, Go to 1.17 [PR 6](https://github.com/vultr/cert-manager-webhook-vultr/pull/6) 
* Bump github.com/vultr/govultr/v2 from 2.4.0 to 2.12.0 [PR 3](https://github.com/vultr/cert-manager-webhook-vultr/pull/3) 


## [v0.1.0](https://github.com/vultr/cert-manager-webhook-vultr) (2021-04-20)
Initial Release
