package cosmos

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"log"
	"net/http"
)

// quick and dirty wrapper for demo'ing this against the local-dev emulator:
// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos#pkg-index

// emulator credentials - unlike the cert we need to install locally from the emulator,
// this won't change every time you recreate the emulator Docker container
const (
	emulatorKey = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="
	emulatorURI = "https://localhost:8081"
)

// Check - identify errors and panic if found. ignores CosmosDB error reponses
// indicating use of non-idempotent APIs repeatedly (like recreating an existing DB)
// to simplify our demo environment.
func Check(err error) bool {
	if err != nil {
		azRespErr, found := err.(*azcore.ResponseError)
		if !found || azRespErr.StatusCode != http.StatusConflict {
			return true
		}
	}

	return false
}

func NewClient(lgr *log.Logger) (*Client, error) {
	// create client, using default emulator creds and URI
	lgr.Printf("creating credential and CosmosDB emulator client")
	cred, err := azcosmos.NewKeyCredential(emulatorKey)
	if Check(err) {
		return nil, err
	}

	client, err := azcosmos.NewClientWithKey(emulatorURI, cred, nil)
	if Check(err) {
		return nil, err
	}
	lgr.Printf("created CosmosDB emulator client with response: %+v", client)

	return &Client{
		lgr:    lgr,
		cred:   cred,
		client: client,
		dbs:    map[string]*azcosmos.DatabaseClient{},
	}, nil
}

type Client struct {
	lgr    *log.Logger
	cred   azcosmos.KeyCredential
	client *azcosmos.Client
	dbs    map[string]*azcosmos.DatabaseClient
}

func (c *Client) Database(ctx context.Context, dbName string) error {
	if c.dbs[dbName] != nil {
		return fmt.Errorf("Client: already associated with database %q", dbName)
	}

	// create remote database if not exists
	c.lgr.Printf("Client: creating CosmosDB emulator database %q", dbName)
	dbCfg := azcosmos.DatabaseProperties{
		ID: dbName,
	}
	_, err := c.client.CreateDatabase(ctx, dbCfg, nil)
	if Check(err) {
		return err
	}

	// grab an instance of the new DB to work with containers on
	db, err := c.client.NewDatabase(dbName)
	if Check(err) {
		return err
	}
	c.dbs[dbName] = db

	return nil
}

func (c *Client) Container(ctx context.Context, dbName, ctrName string, partitionKey []string) (*azcosmos.ContainerClient, error) {
	// create container and select partition key
	c.lgr.Printf("Client: creating CosmosDB emulator container %q on database %q", ctrName, dbName)

	thruProp := azcosmos.NewManualThroughputProperties(400)
	ctrProps := &azcosmos.CreateContainerOptions{
		ThroughputProperties: &thruProp,
	}
	ctrCfg := azcosmos.ContainerProperties{
		ID: ctrName,
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: partitionKey,
		},
	}

	db, found := c.dbs[dbName]
	if !found {
		return nil, fmt.Errorf("Client: db %q not yet associated with this client", dbName)
	}

	ctrResp, err := db.CreateContainer(ctx, ctrCfg, ctrProps)
	if Check(err) {
		return nil, err
	}
	c.lgr.Printf("Client: created CosmosDB emulator container with response: %+v", ctrResp)

	ctr, err := c.client.NewContainer(dbName, ctrName)
	if Check(err) {
		return nil, err
	}

	return ctr, nil
}
