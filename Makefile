OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)

IMAGE_NAME := "webhook"
IMAGE_TAG := "latest"

OUT := $(shell pwd)/_out
TEST_ASSET_ETCD := $(OUT)/kubebuilder/bin/etcd
TEST_ASSET_KUBE_APISERVER := $(OUT)/kubebuilder/bin/kube-apiserver
TEST_ASSET_KUBECTL := $(OUT)/kubebuilder/bin/kubectl

$(shell mkdir -p "$(OUT)")

test: _test/kubebuilder
	TEST_ASSET_ETCD="$(TEST_ASSET_ETCD)" TEST_ASSET_KUBE_APISERVER="$(TEST_ASSET_KUBE_APISERVER)" TEST_ASSET_KUBECTL="$(TEST_ASSET_KUBECTL)" \
	go test -v .

_test/kubebuilder:
	sh ./scripts/fetch-test-binaries.sh

clean: clean-kubebuilder

clean-kubebuilder:
	rm -Rf _out/kubebuilder



.PHONY: deploy
deploy: build-linux docker-build docker-push

.PHONY: build-linux
build-linux:
	@echo "building vultr csi for linux"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags '-X main.version=$(VERSION)' -o cert-manager-webhook-vultr .

.PHONY: docker-build
docker-build:
	@echo "building docker image to dockerhub $(REGISTRY) with version $(VERSION)"
	docker build . -t $(REGISTRY)/cert-manager-webhook-vultr:$(VERSION)

.PHONY: docker-push
docker-push:
	docker push $(REGISTRY)/cert-manager-webhook-vultr:$(VERSION)

.PHONY: rendered-manifest.yaml
rendered-manifest.yaml:
	helm template \
	    --name cert-manager-webhook-vultr \
        --set image.repository=$(IMAGE_NAME) \
        --set image.tag=$(IMAGE_TAG) \
        deploy/cert-manager-webhook-vultr > "$(OUT)/rendered-manifest.yaml"
