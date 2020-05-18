.DEFAULT_GOAL := build
.PHONY: all fmt lint check test build image clean deploy

IMAGE_TAG:=$(shell ./docker/image-tag)
IMAGE_PREFIX=jacops/
NAMESPACE=kubers
CLUSTER_NAME=kubers
GOOS=linux
GOARCH=amd64

minikube:
	@minikube start -p $(CLUSTER_NAME) --insecure-registry "10.0.0.0/24"
	@kubectl create ns $(NAMESPACE) --dry-run=client --output yaml | kubectl apply -f -
	@$(MAKE) minikube-profile

minikube-profile:
	@minikube profile $(CLUSTER_NAME)
	@kubectl config use-context $(CLUSTER_NAME)
	@kubectl config set-context --current --namespace $(NAMESPACE)

build-%:
	@echo "Building $*..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -i -v -ldflags="-X main.version=$(IMAGE_TAG)" -o dist/$* ./cmd/$*

build: build-kubersd build-kubers-agent

test: unit-test

unit-test:
	go test -race ./...

clean:
	@rm -rf dist

docker-build-%: build-%
	docker build -f docker/Dockerfile.$* -t $(IMAGE_PREFIX)$*:$(IMAGE_TAG) .

docker-build: docker-build-kubersd docker-build-kubers-agent

docker-push-%:
	docker push $(IMAGE_PREFIX)$*:$(IMAGE_TAG)

docker-push: docker-push-kubersd docker-push-kubers-agent

k8s-deploy-kubers:
	helm upgrade --devel -i kubers \
		--set injector.image.tag=$(IMAGE_TAG) \
		--set agent.image.tag=$(IMAGE_TAG) \
    --wait --namespace kubers ./chart

k8s-deploy-kubers-release:
	helm upgrade -i kubers jacops/kubers \
	   --wait --namespace kubers

k8s-deploy-example-azure:
	examples/nginx-basic-auth/azure/keyvault-with-sp/deploy/kustomize.sh | kubectl apply -f -

k8s-deploy-example-aws:
	examples/nginx-basic-auth/aws/secretmanager-with-user/deploy/kustomize.sh | kubectl apply -f -

k8s-restart-example:
	kubectl rollout restart deployment nginx

k8s-restart-injector:
	kubectl rollout restart deployment kubers

k8s-restart-all: k8s-restart-injector sleep3 k8s-restart-example

k8s-all-azure: docker-build k8s-deploy-kubers k8s-deploy-example-azure

k8s-all-aws: docker-build k8s-deploy-kubers k8s-deploy-example-aws

k8s-logs:
	kubectl logs -f $(shell kubectl get pod -l app.kubernetes.io/name=kubers -o jsonpath="{.items[-1:].metadata.name}" --sort-by=.status.startTime)

k8s-logs-example:
	kubectl logs -f $(shell kubectl get pod -l name=nginx -o jsonpath="{.items[-1:].metadata.name}" --sort-by=.status.startTime)

k8s-clean-kubers:
	helm delete kubers || true

k8s-clean-aws:
	examples/nginx-basic-auth/aws/secretmanager-with-user/deploy/kustomize.sh | kubectl delete -f -

k8s-clean-azure:
	examples/nginx-basic-auth/azure/keyvault-with-sp/deploy/kustomize.sh | kubectl delete -f -

sleep3:
	@sleep 3

helm-package:
	helm package chart --destination dist

k8s-all-restart: k8s-all
	@$(MAKE) k8s-restart-injector
	@sleep 5
	@$(MAKE) k8s-restart-example
