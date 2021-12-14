module github.com/vultr/cert-manager-webhook-vultr

go 1.16

require (
	github.com/jetstack/cert-manager v1.6.1
	github.com/vultr/govultr/v2 v2.12.0
	golang.org/x/oauth2 v0.0.0-20210810183815-faf39c7919d5
	k8s.io/apiextensions-apiserver v0.22.4
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
)
