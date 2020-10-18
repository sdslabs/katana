PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
PROJECTROOT := $(shell pwd)
GOBIN := $(PROJECTROOT)/bin

# Shell script related variables.
UTILDIR := $(PROJECTROOT)/scripts/utils
SPINNER := $(UTILDIR)/spinner.sh
BUILDIR := $(PROJECTROOT)/scripts/build

CREATEBIN := $(shell [ ! -d ./bin ] && mkdir bin)

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
	@./scripts/install_golint.sh
	@printf "👍 Done\n"

## Format the code
fmt:
	@printf "🔨 Formatting\n" 
	@gofmt -l -s .
	@printf "👍 Done\n"

## Check codebase for style mistakes
lint: install-golint
	@printf "🔨 Linting\n"
	@golint ./...
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
