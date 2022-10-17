package main

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"log"
	"net/http"
)

const (
	emulatorKey = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="
	emulatorURI = "https://localhost:8081"
)

func check(err error) {
	if err != nil {
		azRespErr, found := err.(*azcore.ResponseError)
		if !found || azRespErr.StatusCode != http.StatusConflict {
			panic(err.Error())
		}
	}
}

func main() {
	lgr := log.Default()

	// create client, using default emulator creds and URI
	lgr.Printf("creating credential and CosmosDB emulator client")
	cred, err := azcosmos.NewKeyCredential(emulatorKey)
	check(err)

	client, err := azcosmos.NewClientWithKey(emulatorURI, cred, nil)
	check(err)
	lgr.Printf("created CosmosDB emulator client with response: %+v", client)

	// create database
	lgr.Printf("creating CosmosDB emulator database")
	dbName := "eli_demo"
	dbCfg := azcosmos.DatabaseProperties{ID: dbName}
	dbResp, err := client.CreateDatabase(context.Background(), dbCfg, nil)
	check(err)
	lgr.Printf("created CosmosDB emulator database with response: %+v", dbResp)

	db, err := client.NewDatabase(dbName)
	if err != nil {
		panic(err.Error())
	}

	// create container and select partition key
	lgr.Printf("creating CosmosDB emulator container")
	ctrName := "files"
	thruProp := azcosmos.NewManualThroughputProperties(400)
	ctrProps := &azcosmos.CreateContainerOptions{ThroughputProperties: &thruProp}
	ctrCfg := azcosmos.ContainerProperties{
		ID: ctrName,
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/owner_id"},
		},
	}
	ctrResp, err := db.CreateContainer(context.Background(), ctrCfg, ctrProps)
	check(err)
	lgr.Printf("created CosmosDB emulator container with response: %+v", ctrResp)

	ctr, err := client.NewContainer(dbName, ctrName)
	if err != nil {
		panic(err.Error())
	}

	// create a record in the container
	fileID := "1234"  // record ID
	ownerID := "5678" // partition key
	filePartitionKey := azcosmos.NewPartitionKeyString(ownerID)
	file := map[string]interface{}{
		"id":              fileID,
		"owner_id":        ownerID,
		"path":            "pkg/core",
		"filename":        "package-lock.json",
		"project_name":    "foo",
		"project_version": "1.2.3",
	}
	marshalled, err := json.Marshal(file)
	if err != nil {
		panic(err.Error())
	}
	createResp, err := ctr.CreateItem(context.Background(), filePartitionKey, marshalled, nil)
	check(err)
	lgr.Printf("CosmosDB emulator create item succeeded with response: %+v", createResp)

	// fetch back the created record
	readResp, err := ctr.ReadItem(context.Background(), filePartitionKey, fileID, nil)
	if err != nil {
		panic(err.Error())
	}

	out := map[string]interface{}{}
	if err := json.Unmarshal(readResp.Value, &out); err != nil {
		panic(err.Error())
	}
	lgr.Printf("Read back written object(id=%s): %+v", fileID, out)
}
