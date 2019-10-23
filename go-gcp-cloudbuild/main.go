package main

import (
	"context"
	"fmt"

	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/option"
)

var jsonPath = "devops-deliver-creds.json"

func main() {
	fmt.Println("Hello World!!")
	cb, err := authenticateCloudBuild()
	if err != nil {
cb.Projects.Builds.()
	}
}

func cntxt() context.Context {
	return context.Background()
}

func authenticateCloudBuild() (*cloudbuild.Service, error) {
	cb, err := cloudbuild.NewService(cntxt(), option.WithCredentialsFile(jsonPath))
	if err != nil {
		fmt.Println("Error while GCP authentication for Cloud Build")
		fmt.Println(err.Error())
		return nil, err
	}
	return cb, err
}
