package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/option"
)

var jsonPath = "sb-test-cb-go.json"

func createService() *cloudbuild.Service {
	ctx := context.Background()

	cb, err := cloudbuild.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		fmt.Print("Authentication error...")
	}
	return cb
}

// datastoreDB persists books to Cloud Datastore.
// https://cloud.google.com/datastore/docs/concepts/overview
type datastoreDB struct {
	client *datastore.Client
}

var (
	DB BookDatabase
)

// Ensure datastoreDB conforms to the BookDatabase interface.
var _ BookDatabase = &datastoreDB{}

func configureDatastoreDB(projectID string) (BookDatabase, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return newDatastoreDB(client)
}

// newDatastoreDB creates a new BookDatabase backed by Cloud Datastore.
// See the datastore and google packages for details on creating a suitable Client:
// https://godoc.org/cloud.google.com/go/datastore
func newDatastoreDB(client *datastore.Client) (BookDatabase, error) {
	ctx := context.Background()
	// Verify that we can communicate and authenticate with the datastore service.
	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	return &datastoreDB{
		client: client,
	}, nil
}

func main() {

	cb, err := configureDatastoreDB("sb-test-cb-go")

	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Println(cb)
}

//validates our service at gcp using json file
func authenticate() (_ error) {
	ctx := context.Background()
	cs, err := cloudbuild.NewService(ctx, option.WithCredentialsFile(jsonPath))

	if err != nil {
		logger.Msg("Error while GCP authentication")
		logger.Error(err)
		cs = nil
	}

	return cs, err
}
