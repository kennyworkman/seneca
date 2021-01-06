.DEFAULT_GOAL := help
APP?=seneca
COMMIT_SHA?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
REGISTRY_REPO?=
LOG_FILE_NAME?=./logs/$(shell date +%Y%m%d_%H%M%S)_$(APP).log

OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

clean:
		@echo "$(OK_COLOR)==> Cleaning working directory... $(NO_COLOR)"
		@go clean
		rm -f $(APP)

build: clean
		@echo "$(OK_COLOR)==> Building binary... $(NO_COLOR)"
		go build -o $(APP) cmd/$(APP)/main.go
		chmod +x ./seneca
		cp ./seneca /usr/local/bin/

run:
		go run cmd/$(APP)/main.go

test: build
		go run cmd/$(APP)/main.go https://www.nature.com/articles/s41467-020-18008-4
