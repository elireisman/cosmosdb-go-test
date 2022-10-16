.DEFAULT: all

IPADDR=$(shell ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}' | head -n 1)

.PHONY: pull
	@docker pull mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator
	@docker run \
		--publish 8081:8081 \
		--publish 10251-10254:10251-10254 \
		--memory 3g --cpus=2.0 \
		--name=test-linux-emulator \
		--env AZURE_COSMOS_EMULATOR_PARTITION_COUNT=10 \
		--env AZURE_COSMOS_EMULATOR_ENABLE_DATA_PERSISTENCE=true \
		--env AZURE_COSMOS_EMULATOR_IP_ADDRESS_OVERRIDE=$(IPADDR) \
		--interactive \
		--tty \
		mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator

.PHONY: build
build:
	@mkdir -p bin
	@rm -f bin/*
	@go build -o bin/demo ./...

.PHONY: test
test:
	@go test ./...

.PHONY: run
run:
	bin/demo
.PHONY:
all: test build run


