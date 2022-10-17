package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"strconv"
)

func main() {
	// create client, using default emulator creds and URI
	emulatorKey := ""
	cred, err := azcosmos.NewKeyCredential(emulatorKey)
	if err != nil {
		panic(err.Error())
	}
	emulatorURI := "http://localhost:8081"
	client, err := azcosmos.NewClientWithKey(emulatorURI, cred, nil)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("[CREATE DATABASE]\n")

	// create database
	dbName := "eli_demo"
	dbCfg := azcosmos.DatabaseProperties{ID: dbName}
	dbResp, err := client.CreateDatabase(context.Background(), dbCfg, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("[dbResp] %+v\n", dbResp)

	db, err := client.NewDatabase(dbName)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("[CREATE CONTAINER]\n")

	// create container and select partition key
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
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("[ctrResp] %+v\n", ctrResp)

	ctr, err := client.NewContainer(dbName, ctrName)
	if err != nil {
		panic(err.Error())
	}

	// create a record in the container
	fileID := 1234  // record ID
	ownerID := 5678 // partition key
	filePartitionKey := azcosmos.NewPartitionKeyString(strconv.Itoa(ownerID))
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
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("[createResp] %+v\n", createResp)

	// fetch back the created record
	readResp, err := ctr.ReadItem(context.Background(), filePartitionKey, strconv.Itoa(fileID), nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("[readResp] %+v\n", readResp)
}
