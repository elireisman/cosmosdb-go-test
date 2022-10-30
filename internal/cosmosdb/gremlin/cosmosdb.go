package gremlin

import (
	"context"
	"fmt"
	"log"
)

// emulator credentials - unlike the cert we need to install locally from the emulator,
// this won't change every time you recreate the emulator Docker container
const (
	emulatorKey = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="
	emulatorURI = "https://localhost:8081"
)

type Client struct{
  cred azcore.TokenCredential
  dbs map[string]*armcosmos.GremlinResourcesClient
}


func NewClient(ctx context.context) (*Client, error) {
  // https://github.com/Azure/azure-sdk-for-go/blob/df8dbee5478fe0f4330592c5ad36467c176dc210/sdk/azidentity/username_password_credential.go#L36-L47
  // https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos#example-GremlinResourcesClient.BeginCreateUpdateGremlinDatabase
	databaseAccountsClient, err := armcosmos.NewDatabaseAccountsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	pollerResp, err := databaseAccountsClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		accountName,
		armcosmos.DatabaseAccountCreateUpdateParameters{
			Location: to.Ptr("useast"),
			Kind:     to.Ptr(armcosmos.DatabaseAccountKindGlobalDocumentDB),
			Properties: &armcosmos.DatabaseAccountCreateUpdateProperties{
				DatabaseAccountOfferType: to.Ptr("Standard"),
				Locations: []*armcosmos.Location{
					{
						FailoverPriority: to.Ptr[int32](0),
						LocationName:     to.Ptr(location),
					},
				},
				Capabilities: []*armcosmos.Capability{
					{
						Name: to.Ptr("EnableGremlin"),
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	resp.DatabaseAccountGetResults, nil
  return &Client{}, nil
}
