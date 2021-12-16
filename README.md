# Vultr Webhook for Cert Manager

This is a webhook solver for [Vultr](https://www.vultr.com) to be used with [Cert-Manager](https://cert-manager.io/docs/)
.
## Prerequisites

There are a few things required before you can start using `cert-manager-webhook-vultr`.

- Helm v3+ is required to install the `cert-manager-webhook-vultr` charts
- [Cert-Manager](https://cert-manager.io/docs/) needs to be running on your cluster prior.

## Installation

### Installing the webhook

First, you will need to deploy a secret with your api key. 
You can do this by using the sample yaml in `testdata/vultr/api-key.yaml.sample` or by running the following kubectl command: 

```shell
kubectl create secret generic "vultr-credentials" --from-literal=apiKey=<VULTR API KEY> --namespace=cert-manager
```

Second, you will need to deploy the `cert-manager-webhook-vultr`. We have a helm chart that makes installation of this fairly straightforward. 

```shell
helm install --namespace cert-manager cert-manager-webhook-vultr ./deploy/cert-manager-webhook-vultr
```

This will deploy all necessary services, deployments, rbac, and various other resources required for cert-manager.

To uninstall the webhook run the following:

```shell
helm uninstall cert-manager-webhook-vultr --namespace=cert-manager
```

### Deploying a ClusterIssuer

Below we will deploy a ClusterIssuer which will use LetsEncrypt staging environment 
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    # You must replace this email address with your own.
    # Let's Encrypt will use this to contact you about expiring
    # certificates, and issues related to your account.
    email: <enter your email address>
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      # Secret resource that will be used to store the account's private key.
      name: letsencrypt-staging
    solvers:
    - dns01:
        webhook:
          groupName: acme.vultr.com
          solverName: vultr
          config:
            apiKeySecretRef:
              key: apiKey
              name: vultr-credentials
```

We also need to grant permissions for the `service account` to be able to grab the secret .

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cert-manager-webhook-vultr:secret-reader
  namespace: cert-manager
rules:
- apiGroups: [""]
  resources: ["secrets"]
  resourceNames: ["vultr-credentials"]
  verbs: ["get", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: cert-manager-webhook-vultr:secret-reader
  namespace: cert-manager
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cert-manager-webhook-vultr:secret-reader
subjects:
  - apiGroup: ""
    kind: ServiceAccount
    name: cert-manager-webhook-vultr
```

### Request a certificate

The Certificate resource represents a human readable definition of a certificate request that is to be honored by an issuer which is to be kept up-to-date.

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: staging-cert-example-com
spec:
  commonName: example.com # REPLACE THIS WITH YOUR DOMAIN
  dnsNames:
  - example.com # REPLACE THIS WITH YOUR DOMAIN
  issuerRef:
    name: letsencrypt-staging
    kind: ClusterIssuer
  secretName: example-com-tls
```

To check on the certificate run the following:

```shell
kubectl describe certificate staging-cert-example-com
```

To delete a certificate run the following:

```shell
kubectl delete certificate staging-cert-example-com
```

## Troubleshooting
Cert-Manager has a great page that describes how to [troubleshoot](https://cert-manager.io/docs/faq/troubleshooting/).

### Running the test suite

You can run the test suite with:

```bash
$ TEST_ZONE_NAME=example.com. make test
```
