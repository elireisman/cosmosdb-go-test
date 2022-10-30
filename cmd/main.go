package main

import (
	"context"
	"encoding/json"
	"log"

	cosmos "github.com/elireisman/cosmosdb-go-test/internal/cosmosdb/sql"
	"github.com/elireisman/cosmosdb-go-test/internal/utils"
)

func main() {
	lgr := log.Default()

	// create emulator client
	client, err := cosmos.NewClient(lgr)
	if err != nil {
		lgr.Fatal(err)
	}

	// create db instance if not exists and associate it with the client
	dbName := "eli_demo"
	err = client.Database(context.Background(), dbName)
	if err != nil {
		lgr.Fatal(err)
	}

	// create a container where we can store and query data
	containerName := "manifests"
	partitionKey := []string{"/owner_id"}
	ctr, err := client.Container(context.Background(), dbName, containerName, partitionKey)
	if err != nil {
		lgr.Fatal(err)
	}

	// create a record in the container
	fileID := "1234"  // record ID
	ownerID := "5678" // partition key
	filePartitionKey := cosmos.PartitionKey(ownerID)
	file := map[string]interface{}{
		"id":              fileID,
		"owner_id":        ownerID,
		"path":            "utils",
		"filename":        "package-lock.json",
		"project_name":    "foo",
		"project_version": "1.2.3",
	}
	marshalled, err := json.Marshal(file)
	if err != nil {
		lgr.Fatal(err)
	}
	createResp, err := ctr.CreateItem(context.Background(), filePartitionKey, marshalled, nil)
	if cosmos.Check(err) {
		lgr.Fatal(err)
	}

	lgr.Printf("CosmosDB emulator create item succeeded with response: %+v", createResp)

	// fetch back the created record
	readResp, err := ctr.ReadItem(context.Background(), filePartitionKey, fileID, nil)
	if cosmos.Check(err) {
		lgr.Fatal(err)
	}

	lgr.Printf("Read back written object(id=%s):", fileID)
	if err := utils.PrettyJSON(readResp.Value); err != nil {
		lgr.Fatal(err)
	}
}
