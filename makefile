# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# CLASS NOTES
#
# Kind
# 	For full Kind v0.20 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.21.0
#
# RSA Keys
# 	To generate a private/public key PEM file.
# 	$ openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# 	$ openssl rsa -pubout -in private.pem -out public.pem


run:
	"/mnt/c/Program Files/Go/bin/go.exe" run app/services/sales-api/main.go | "/mnt/c/Program Files/Go/bin/go.exe" run app/tooling/logfmt/main.go

run-help:
		"/mnt/c/Program Files/Go/bin/go.exe" run app/services/sales-api/main.go --help | "/mnt/c/Program Files/Go/bin/go.exe" run app/tooling/logfmt/main.go

curl:
		curl -il http://localhost:3000/v1/hack

curl-auth:
	curl -il -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/hackauth

load:
	hey -m GET -c 100 -n 100000 "http://localhost:3000/v1/hack"

ready:
	curl -il http://localhost:3000/v1/readiness

live:
	curl -il http://localhost:3000/v1/liveness

admin:
	"/mnt/c/Program Files/Go/bin/go.exe" run app/tooling/sales-admin/main.go
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

			kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

# ------------------------------------------------------------------------------

dev-load:
	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/database | kubectl apply -f -
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database

	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --timeout=120s --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply

# ==============================================================================
# Metrics and Tracing

metrics-view-sc:
	expvarmon -ports="localhost:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"



# ==============================================================================
dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 --max-log-requests=6 | /mnt/c/Program\ Files/Go/bin/go.exe run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)


dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(APP)

dev-logs-db:
	kubectl logs --namespace=$(NAMESPACE) -l app=database --all-containers=true -f --tail=100

pgcli:
	pgcli postgresql://postgres:postgres@localhost

# ==============================================================================

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces


# ==============================================================================
# Modules support

tidy:
	/mnt/c/Program\ Files/Go/bin/go.exe mod tidy
	/mnt/c/Program\ Files/Go/bin/go.exe mod vendor