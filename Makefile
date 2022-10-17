.DEFAULT: all

IPADDR := "$(shell ifconfig | grep 'inet ' | grep -v '127.0.0.1' | awk '{print $$2}' | head -n 1)"

.PHONY: all
all: build test run

.PHONY: cert
cert:
	@curl -k -q https://$(IPADDR):8081/_explorer/emulator.pem > emulatorcert.crt
	@sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain emulatorcert.crt
	@open https://localhost:8081/_explorer/index.html

.PHONY: emu
emu:
	@if ! docker info >/dev/null 2>&1; then echo "ERROR: Docker must be running locally"; exit 1; fi
	docker pull mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator
	docker run \
		--publish 8081:8081 \
		--publish 10251-10254:10251-10254 \
		--memory 3g --cpus=2.0 \
		--name=test-linux-emulator \
		--env AZURE_COSMOS_EMULATOR_PARTITION_COUNT=10 \
		--env AZURE_COSMOS_EMULATOR_ENABLE_DATA_PERSISTENCE=true \
		--env AZURE_COSMOS_EMULATOR_IP_ADDRESS_OVERRIDE=$(IPADDR) \
		--rm \
		--interactive \
		--tty \
		mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator

.PHONY: build
build:
	@mkdir -p bin
	@rm -f bin/*
	go build -o bin/demo cmd/main.go

.PHONY: test
test:
	@#go test ./...

.PHONY: run
run:
	bin/demo

