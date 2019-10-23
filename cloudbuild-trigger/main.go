package main

import (
	"context"
	"encoding/json"
	"fmt"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/option"
)

var jsonPath string = "cloudsloutions-fd422e010b14.json"

var projectID string = "cloudsloutions"

func createbuild(projectID string, build *cloudbuild.Build, jsonPath_arg string) Response {
	ctx := context.Background()

	cb, err := cloudbuild.NewService(ctx, option.WithCredentialsFile(jsonPath_arg))
	if err != nil {
		fmt.Printf("error1: " + err.Error())
	}

	operation, err := cb.Projects.Builds.Create(projectID, build).Do()
	if err != nil {
		fmt.Printf("error2: " + err.Error())
	}

	responseBytes := operation.Metadata

	var response Response
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		fmt.Printf("error3: " + err.Error())
	}

	fmt.Println(fmt.Sprintf("Build ID: %s", response.Build.Id))
	fmt.Println(fmt.Sprintf("Build Status: %s", response.Build.Status))
	fmt.Println(fmt.Sprintf("Project Id: %s", response.Build.ProjectId))
	fmt.Println(fmt.Sprintf("Logs Bucket: %s", response.Build.LogsBucket))
	fmt.Println(fmt.Sprintf("Log Url: %s", response.Build.LogUrl))
	return response
}

// Get the status of CI using BuildId
func getbuildstatus(projectID string, response Response, jsonPath_arg string) {
	ctx := context.Background()

	cb, err := cloudbuild.NewService(ctx, option.WithCredentialsFile(jsonPath_arg))
	if err != nil {
		fmt.Printf("error1: " + err.Error())
	}

	build, err := cb.Projects.Builds.Get(projectID, response.Build.Id).Do()
	if err != nil {
		fmt.Printf("error2: " + err.Error())
	}

	fmt.Println("Build Status :" + build.Status)

}

type Response struct {
	Build cloudbuild.Build `json:"build"`
}

func main() {

	//explicitAuth(jsonPath, projectID)

	var buildStep1 cloudbuild.BuildStep
	buildStep1.Args = []string{"clone", "https://github.com/g33kzone/sb-hello-world", "/workspace/sb-hello-world"}
	buildStep1.Name = "gcr.io/cloud-builders/git"
	buildStep1.Id = "CLONE"

	var buildStep2 cloudbuild.BuildStep
	buildStep2.Args = []string{"install"}
	buildStep2.Name = "gcr.io/cloud-builders/mvn"
	buildStep2.Id = "BUILD"
	buildStep2.Dir = "/workspace/sb-hello-world"

	var buildStep3 cloudbuild.BuildStep
	buildStep3.Name = "gcr.io/cloud-builders/docker"
	buildStep3.Args = []string{"build", "--tag=gcr.io/sb-hello-world-234105/image", "."}

	var buildStep4 cloudbuild.BuildStep
	buildStep4.Name = "gcr.io/cloud-builders/kubectl"
	buildStep4.Args = []string{"delete", "deployment", "sb-hello-world", "--ignore-not-found"}
	buildStep4.Env = []string{"CLOUDSDK_COMPUTE_ZONE=us-central1-a", "CLOUDSDK_CONTAINER_CLUSTER=sb-hello-world"}

	var buildSteps []*cloudbuild.BuildStep
	buildSteps = append(buildSteps, &buildStep1, &buildStep2, &buildStep3, &buildStep4)

	var build cloudbuild.Build
	build.Steps = buildSteps

	response := createbuild(projectID, &build, jsonPath)
	getbuildstatus(projectID, response, jsonPath)

	//var response1 Response
	//response1.Build.Id = "6b088c66-4788-454a-9e0d-aa1152cf2d51"
	//getbuildstatus(projectID, response1, jsonPath)

}
