package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/option"
)

var jsonPath = "cloud-build-go-65b93e6fbda1.json"

func main() {

	ctx := context.Background()

	cb, err := cloudbuild.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		fmt.Print("Authentication error...")
	} else {
		fmt.Print("Authentication success...")
		log.Fatal("Authentication success...")
	}

}
