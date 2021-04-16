module github.com/vultr/cert-manager-webhook-vultr

go 1.16

require (
	github.com/jetstack/cert-manager v1.3.1
	github.com/vultr/govultr/v2 v2.4.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	k8s.io/apiextensions-apiserver v0.19.0
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v0.19.0
)
