#!make

.DEFAULT_GOAL := build
.PHONY: all fmt lint check test build image clean deploy

VERSION:=$(shell ./scripts/version)
IMAGE_PREFIX=jacops/
NAMESPACE=kubers
CLUSTER_NAME=kubers
GOOS=linux
GOARCH=amd64

ifneq ("$(wildcard .env)","")
	include .env
	export $(shell sed 's/=.*//' .env)
endif

show-version:
	@echo $(VERSION)

minikube:
	@minikube start -p $(CLUSTER_NAME) --insecure-registry "10.0.0.0/24"
	@kubectl create ns $(NAMESPACE) --dry-run=client --output yaml | kubectl apply -f -
	@$(MAKE) minikube-profile

minikube-profile:
	@minikube profile $(CLUSTER_NAME)
	@kubectl config use-context $(CLUSTER_NAME)
	@kubectl config set-context --current --namespace $(NAMESPACE)

build:
	goreleaser --snapshot --skip-publish --rm-dist

release:
	GITHUB_TOKEN=$(GITHUB_TOKEN) goreleaser --rm-dist

test: unit-test

test-coverage:
	./scripts/go.test.sh

unit-test:
	go test -race ./...

clean:
	@rm -rf dist

docker-build-%:
	docker build -f cmd/$*/dev.Dockerfile -t $(IMAGE_PREFIX)$*:$(VERSION) .
	docker tag $(IMAGE_PREFIX)$*:$(VERSION) $(IMAGE_PREFIX)$*:latest

docker-build: docker-build-kubersd docker-build-kubers-agent

generate-deployment-manifests:
	helm template kubers \
		--set injector.image.tag=latest \
		--set agent.image.tag=latest \
    --namespace kubers ./chart > ./deploy/kubers.yaml

k8s-deploy-kubers:
	helm upgrade --devel -i kubers \
		--set injector.image.tag=$(VERSION) \
		--set injector.log.level=debug \
		--set agent.image.tag=$(VERSION) \
		--set agent.log.level=debug \
    --wait --namespace kubers ./chart

k8s-deploy-kubers-release:
	helm upgrade -i kubers jacops/kubers \
	   --wait --namespace kubers

k8s-deploy-example-azure:
	examples/nginx-basic-auth/azure/keyvault-with-sp/deploy/kustomize.sh | kubectl apply -f -

k8s-deploy-example-aws:
	examples/nginx-basic-auth/aws/secretsmanager-with-user/deploy/kustomize.sh | kubectl apply -f -

k8s-restart-example:
	kubectl rollout restart deployment nginx

k8s-restart-injector:
	kubectl rollout restart deployment kubers

k8s-restart-all: k8s-restart-injector sleep3 k8s-restart-example

k8s-all-azure: build docker-build k8s-deploy-kubers k8s-deploy-example-azure

k8s-all-aws: build docker-build k8s-deploy-kubers k8s-deploy-example-aws

k8s-logs:
	kubectl logs -f $(shell kubectl get pod -l app.kubernetes.io/name=kubers -o jsonpath="{.items[-1:].metadata.name}" --sort-by=.status.startTime)

k8s-logs-example:
	kubectl logs -f $(shell kubectl get pod -l name=nginx -o jsonpath="{.items[-1:].metadata.name}" --sort-by=.status.startTime)

k8s-clean-kubers:
	helm delete kubers || true

k8s-clean-aws:
	examples/nginx-basic-auth/aws/secretsmanager-with-user/deploy/kustomize.sh | kubectl delete -f -

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
