package main

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
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

	// create database
	dbName := "eli_demo"
	dbCfg := azcosmos.DatabaseProperties{ID: dbName}
	db, err := client.CreateDatabase(context.Background(), dbCfg, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("[db] %+v\n", db)

	// create container and select partition key
	thruProps := &CreateContainerOptions{ThroughputProperties: azcosmos.NewManualThroughputProperties(400)}
	ctrCfg := azcosmos.ContainerProperties{
		Id: "files",
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/owner_id"},
		},
	}
	ctr, err := database.CreateContainer(context, thruProp)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("[ctr] %+v\n", ctr)

	// create a record in the container
	fileId := 1234    // record ID
	fileOwner := 5678 // partition key
	file := map[string]interface{}{
		"id":              fileId,
		"path":            "pkg/core",
		"filename":        "package-lock.json",
		"project_name":    "foo",
		"project_version": "1.2.3",
	}
	marshalled, err := json.Marshal(item)
	if err != nil {
		panic(err.Error())
	}
	createResp, err := ctr.CreateItem(context.Background(), fileOwner, fileID, marshalled, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("[createResp] %+v\n", createResp)

	// fetch back the created record
	readResp, err := ctr.ReadItem(context.Background(), fileOwner, fileID)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("[readResp] %+v\n", readResp)
}
