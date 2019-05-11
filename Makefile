# Needed SHELL since I'm using zsh
SHELL := /bin/bash

.DEFAULT_GOAL := build

.PHONY: build
build: ## Build app
	@go build -v -o ./build/TrendingGithub

.PHONY: build
debug: build ## Run app in debug mode
	./build/TrendingGithub -- -debug

.PHONY: clean
clean: ## Clean up
	@rm -fR ./build/ ./cover*

.PHONY: configure
configure: ## Creates folders
	@mkdir -p ./build

.PHONY: cover
cover: ## Run tests and generates html coverage file
	@go test -v -race -coverprofile=./coverage.text -covermode=atomic $(shell go list ./...)
	@go tool cover -html=./coverage.text -o ./coverage.html

.PHONY: install
install: ## Install app
	go install

.PHONY: lint
lint: ## Run linters
	gometalinter \
		--disable-all \
		--exclude=vendor \
		--deadline=180s \
		--enable=gofmt \
		--linter='errch:errcheck {path}:PATH:LINE:MESSAGE' \
		--enable=errch \
		--enable=vet \
		--enable=gocyclo \
		--cyclo-over=15 \
		--enable=golint \
		--min-confidence=0.85 \
		--enable=ineffassign \
		--enable=misspell \
		./..

.PHONY: test
test: ## Run tests
	@go test -v -race $(shell go list ./...)

help: ## This help message
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"
