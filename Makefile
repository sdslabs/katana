PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
PROJECTROOT := $(shell pwd)
GOBIN := $(PROJECTROOT)/bin

# Shell script related variables.
UTILDIR := $(PROJECTROOT)/scripts/utils
SPINNER := $(UTILDIR)/spinner.sh
BUILDIR := $(PROJECTROOT)/scripts/build
MANIFEST:= $(PROJECTROOT)/kubernetes/manifests

KEY_NAME := team

NO_OF_TEAMS:= 10
OPENVPN_NAMESPACE := openvpn

POD_COMMAND =$(shell kubectl get pods --namespace $(OPENVPN_NAMESPACE) -l "app=openvpn,release=openvpn" -o jsonpath='{ .items[0].metadata.name }') 
SERVICE_NAME_COMMAND =$(shell kubectl get svc --namespace $(OPENVPN_NAMESPACE) -l "app=openvpn,release=openvpn" -o jsonpath='{ .items[0].metadata.name }') 
SERVICE_IP_COMMAND=$(shell kubectl get svc --namespace $(OPENVPN_NAMESPACE) -l "app=openvpn,release=openvpn" -o jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}')
# CHALLENGE_DEPLOYER_IP :=  $(shell minikube service nginx-ingress-controller --url -n kube-system)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: default
default: help

## Build the cli binary
build-cli:
	@printf "🔨 Building binary $(GOBIN)/$(PROJECTNAME)\n" 
	@./scripts/build/build-cli.sh
	# @cp ./cli/$(PROJECTNAME) $(GOBIN)/
	@printf "👍 Done\n"

## Lint the code
install-golint:
	@printf "🔨 Installing golint\n" 
	@./scripts/install-golint.sh
	@printf "👍 Done\n"

## Format the code
fmt:
	@printf "🔨 Formatting\n" 
	@gofmt -l -s .
	@printf "👍 Done\n"

## Check codebase for style mistakes
lint: install-golint
	@printf "🔨 Linting\n"
	@golangci-lint run
	@printf "👍 Done\n"

## Clean build files
clean:
	@printf "🔨 Cleaning build cache\n" 
	@go clean .
	@printf "👍 Done\n"
	@-rm $(GOBIN)/* 2>/dev/null

## Prepare code for PR
prepare-for-pr: fmt lint
	@git diff-index --quiet HEAD -- ||\
	(echo "-----------------" &&\
	echo "NOTICE: There are some files that have not been committed." &&\
	echo "-----------------\n" &&\
	git status &&\
	echo "\n-----------------" &&\
	echo "NOTICE: There are some files that have not been committed." &&\
	echo "-----------------\n"  &&\
	exit 0)

gen-certificates:
	$(eval POD_NAME := $(POD_COMMAND))
	$(eval SERVICE_NAME := $(SERVICE_NAME_COMMAND))
	$(eval SERVICE_IP := $(SERVICE_IP_COMMAND))
	for n in $$(seq 1 $(NO_OF_TEAMS)); do \
	kubectl --namespace $(OPENVPN_NAMESPACE) exec -it $(POD_NAME) /etc/openvpn/setup/newClientCert.sh $(KEY_NAME)-$$n $(SERVICE_IP) && \
	kubectl --namespace $(OPENVPN_NAMESPACE) exec -it $(POD_NAME) cat "/etc/openvpn/certs/pki/$(KEY_NAME)-$$n.ovpn" > $(KEY_NAME)-$$n.ovpn; \
	done

set-env: build
	minikube start --driver=docker && \
	minikube addons enable ingress  && \
	kubectl apply -f $(MANIFEST) && \
	./bin/katana run

set-env-prod: build
	kubectl apply -f $(MANIFEST) && \
	sudo ./bin/katana run

set-clean:
	sudo rm -rf peer_configs teamcreds.txt teams && \
	kind delete clusters $$(kind get clusters)


set-kind: build
	kind create cluster --config metallb-kind-config.yaml && \
	kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.14.9/config/manifests/metallb-native.yaml && \
	kubectl wait --for=condition=Available=True --timeout=90s deployment/controller -n metallb-system && \
	kubectl apply -f metallb-config.yaml && \
	kubectl apply -f $(MANIFEST) && \
	sudo ./bin/katana run

build:
	cd cmd && go build -o ../bin/katana

run : build
	sudo ./bin/katana run

setup-docs:
	git submodule update --init --recursive
	cp ./docs/config.sample.toml ./docs/config.toml
	npm install --prefix ./docs/themes/hugo-geekdoc
	npm run build --prefix ./docs/themes/hugo-geekdoc

# Prints help message
help:
	@echo "KATANA"
	@echo "build-cli		- Build katana"
	@echo "fmt  	   		- Format code using golangci-lint"
	@echo "help    	   		- Prints help message"
	@echo "install-golint 	- Install golint"
	@echo "clean 			- Clean the build cache"
	@echo "prepare-for-pr 	- Prepare the code for PR after fmt, lint and checking uncommitted files"
	@echo "lint    			- Lint code using golangci-lint"
	@echo "set-env" 		- Setup Katana environment  
	@echo "build"         	- Build katana binary

