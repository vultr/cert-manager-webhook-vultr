OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)

IMAGE_NAME := "webhook"
IMAGE_TAG := "latest"

OUT := $(shell pwd)/_out

KUBEBUILDER_VERSION=2.3.2

$(shell mkdir -p "$(OUT)")

test: _test/kubebuilder
	go test -v .

_test/kubebuilder:
	curl -fsSL https://github.com/kubernetes-sigs/kubebuilder/releases/download/v$(KUBEBUILDER_VERSION)/kubebuilder_$(KUBEBUILDER_VERSION)_$(OS)_$(ARCH).tar.gz -o kubebuilder-tools.tar.gz
	mkdir -p _test/kubebuilder
	tar -xvf kubebuilder-tools.tar.gz
	mv kubebuilder_$(KUBEBUILDER_VERSION)_$(OS)_$(ARCH)/bin _test/kubebuilder/
	rm kubebuilder-tools.tar.gz
	rm -R kubebuilder_$(KUBEBUILDER_VERSION)_$(OS)_$(ARCH)

clean: clean-kubebuilder

clean-kubebuilder:
	rm -Rf _test/kubebuilder


.PHONY: deploy
deploy: docker-build docker-push

#.PHONY: build-linux
#build-linux:
#	@echo "building vultr csi for linux"
#	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-X main.version=$(VERSION)' -o cert-manager-webhook-vultr ./cmd/csi-vultr-driver

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
