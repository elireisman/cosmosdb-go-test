# cosmosdb-go-test
Playground repo for getting to know the CosmosDB SDK for Go and the Docker-based emulator.

### Instructions
1. Ensure Docker is running locally
1. `make emu` - wait for the download to complete, and emulator logs to go from `Starting` to `Started`
1. _If this is your very first time running the project_: `make cert` to obtain and install the emulator certificate (required for local access)
1. `make` to build, test, and run the Go SDK playground app

_Note: These instructions apply to macOS machines *with Intel silicon*. For Codespaces we'll want to use the "pure Linux" bootstrap instructions. The Docker-based emulator is not compatible with other macOS machines :(
