package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	cloudbuild "google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var jsonPath string = "/Users/soumyabrata.sen/Downloads/cloudsloutions-fd422e010b14.json"

var projectID string = "cloudsloutions"

func explicitAuth(jsonPath_arg, projectID_arg string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath_arg))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Buckets Testing:")
	it := client.Buckets(ctx, projectID_arg)
	testServiceAccount, error := client.ServiceAccount(ctx, projectID_arg)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(battrs.Name)
	}

	if error != nil {
		log.Fatal(error)
	}
	fmt.Println("Service Account Name :- " + testServiceAccount)
}

func createbuildtrigger(projectID string, build *cloudbuild.Build, jsonPath_arg string) {
	ctx := context.Background()

	cb, err := cloudbuild.NewService(ctx)

	operation, err := cb.Projects.Builds.Create(projectID, build).Do()
	if err != nil {
		fmt.Printf("here" + err.Error())
	}

	fmt.Printf("%t", operation.Name)

	fmt.Printf("%v")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cb)
}
func main() {

	//explicitAuth(jsonPath, projectID)

	var buildStep1 cloudbuild.BuildStep
	buildStep1.Args = []string{"clone", "https://github.com/g33kzone/sb-hello-world", "/workspace/sb-hello-world"}
	buildStep1.Name = "gcr.io/cloud-builders/git"
	buildStep1.Id = "CLONE"

	var buildSteps []*cloudbuild.BuildStep
	buildSteps = append(buildSteps, &buildStep1)
	//buildSteps[0] = &buildStep1

	var build cloudbuild.Build
	build.Steps = buildSteps

	createbuildtrigger(projectID, &build, jsonPath)
}
