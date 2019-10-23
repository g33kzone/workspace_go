package main

import (
	"context"
	"fmt"

	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/option"
)

var jsonPath = "devops-deliver-creds.json"

func main() {

	fmt.Println("Hello World!!")

	ctx := context.Background()

	service, err := cloudfunctions.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		fmt.Println("NewClient: %v", err)
	}

	service.Projects.Locations.Functions.Create()
}
