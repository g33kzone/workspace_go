package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

var jsonPath = "devops-deliver-creds.json"

// Task structure for Task Description
type Task struct {
	Description string
}

const (
	domain    = "Domain"
	product   = "Product"
	component = "Component"
)

func main() {

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, "devops-deliver", option.WithCredentialsFile(jsonPath))

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	} else {
		fmt.Printf("Data Authentication Successful...")
	}

	// Sets the kind for the new entity.
	// kind := "ProductHierarchy"
	// Sets the name/ID for the new entity.
	// name := "sampletask3"
	// Creates a Key instance.
	taskKey := keyForProduct("odj", "search")

	// Creates a Task instance.
	task := Task{
		Description: "Test Component",
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, taskKey, &task); err != nil {
		log.Fatalf("Failed to save task: %v", err)
	}

	fmt.Printf("Saved %v: %v\n", taskKey, task.Description)

	var taskResponse Task

	err = client.Get(ctx, taskKey, &taskResponse)

	if err == datastore.ErrNoSuchEntity {
		fmt.Printf("Could not find product " + "search" + " under " + "odj")
	} else if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("This Task was retreived %+v\n", taskResponse)

	_, err = client.Put(ctx, keyForComponent("odj", "prod01", "comp01"), &task)
	if err != nil {
		log.Fatalf("Failed to save task: %v", err)
	}
}

func keyForComponent(domainName string, productName string, componentName string) *datastore.Key {
	return datastore.NameKey(component, componentName, keyForProduct(domainName, productName))
}

func keyForProduct(domainName string, productName string) *datastore.Key {
	return datastore.NameKey(product, productName, keyForDomain(domainName))
}

func keyForDomain(domainName string) *datastore.Key {
	return datastore.NameKey(domain, domainName, nil)
}
