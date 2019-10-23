package main

import (
	"context"
	"fmt"
	"log"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

var jsonPath = "devops-deliver-creds.json"

func main() {

	//var resp *cloudresourcemanager.ListProjectsResponse

	ctx := context.Background()

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx, option.WithCredentialsFile(jsonPath))

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	} else {
		fmt.Println("Authentication Successful...")
	}

	req := cloudresourcemanagerService.Projects.List().Filter("labels.org=odjui AND labels.product=deliver AND labels.criticality=work")
	if err := req.Pages(ctx, func(page *cloudresourcemanager.ListProjectsResponse) error {
		for _, project := range page.Projects {
			fmt.Printf("%#v\n", project.ProjectId)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
