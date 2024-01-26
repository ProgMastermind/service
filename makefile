# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

run:
		go run app/services/sales-api/main.go | go run app\tooling\logfmt\main.go

# ==============================================================================

# Define dependencies

GOLANG := golang:1.21.6
ALPINE := alpine:3.19
KIND := kindest/node:v1.29.0@sha256:eaa1450915475849a73a9227b8f201df25e55e268e5d619312131292e324d570
POSTGRES := postgres:16.1
GRAFANA := grafana/grafana:10.2.0
PROMETHEUS := prom/prometheus:v2.48.0
TEMPO := grafana/tempo:2.3.0
LOKI := grafana/loki:2.9.0
PROMTAIL := grafana/promtail:2.9.0

KIND_CLUSTER := ardan-starter-cluster
NAMESPACE := sales-system
APP := sales
BASE_IMAGE_NAME := ardanlabs/service
SERVICE_NAME := sales-api
VERSION := 0.0.1
SERVICE_IMAGE := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME)-metrics:$(VERSION)

# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"

# ==============================================================================
# Building containers

all: service 

service:
	docker build \
		-f zarf/docker/dockerfile.service \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ========================================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
	--image $(KIND) \
	--name $(KIND_CLUSTER) \
	--config zarf/k8s/dev/kind-config.yaml

			kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

# ==============================================================================

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces


