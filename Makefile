SHELL := /bin/bash
.DEFAULT_GOAL := help

include ./.env
export

env?=development

.PHONY: run
run: generate bindata ## Run service
	@go run src/main.go -a

.PHONY: build
build: generate bindata ## Build service
	@go build src/main.go 

.PHONY: bindata
bindata: ## Link sql files to golang code
	@go-bindata -pkg db -o ./src/db/schemas.go ./src/db/schedule/schema.sql ./src/db/subscription/schema.sql

.PHONY: generate
generate: ## Generate golang code for database
	@sqlc generate

.PHONY: tidy
tidy: ## Clean and organize golang dependencies
	@cd go mod tidy

.PHONY: help
help: ## This help dialog
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/:/'`); \
	printf "%-30s %s\n" "target" "help" ; \
	printf "%-30s %s\n" "------" "----" ; \
	for help_line in $${help_lines[@]}; do \
		IFS=$$':' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf '\033[36m'; \
		printf "%-30s %s" $$help_command ; \
		printf '\033[0m'; \
		printf "%s\n" $$help_info; \
	done
