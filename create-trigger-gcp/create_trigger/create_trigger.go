package create_trigger

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"odj-deliver-cloudbuild/config"
	"odj-deliver-cloudbuild/db"
	"odj-deliver-cloudbuild/model"
	"strings"

	"github.com/gin-gonic/gin"
	cloudbuild "google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/option"
	sourcerepo "google.golang.org/api/sourcerepo/v1"
)

//Method to handle create trigger
func HandleClouddBuild(c *gin.Context, db *db.Db) string {

	jsonPath := "cloudsloutions-fd422e010b14.json"

	var componentObject model.ComponentResponse
	fmt.Println("Build Struct Recieved ")
	componentData, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("headers: %v\n", c.Request.Header)
	fmt.Printf("body: %v\n", c.Request.Body)
	json.Unmarshal(componentData, &componentObject)
	fmt.Println("Repository Name: ", componentObject.ProjectID)
	fmt.Println("Repository Name: ", componentObject.ComponentName)

	//authenticate Source Repository and return Repository.Service object
	repoService := authenticateSourceRepo(jsonPath)

	//Return the BuildTrigger object
	build := createBuildObject(componentObject)
	//Check the existance of repository in Source Repository
	flag := checkRepoList(build, repoService)
	// response := create_trigger.TriggerCloudBuild(build.ProjectId, &cbuild, cs, build, db)
	//var trigger model.Trigger
	if flag != true {
		fmt.Println("Repository not found in Source Repo")
		return "Not Found"
	}
	cloudService := authenticateCloudBuild(jsonPath)
	//Check the existance of Trigger in CloudBuild
	flag1 := checkTriggerList(build, cloudService)
	if flag1 == false {
		var buildTrigger cloudbuild.BuildTrigger
		var triggerTemplate cloudbuild.RepoSource
		//triggerTemplate.BranchName = buildTrigger.BranchName
		triggerTemplate.RepoName = build.RepositoryName
		triggerTemplate.BranchName = build.BranchName
		buildTrigger.Filename = build.Filename
		buildTrigger.Description = "Push to branch"
		triggerTemplate.BranchName = build.BranchName

		buildTrigger.TriggerTemplate = &triggerTemplate
		triggerId := triggerCloudBuild(build.ProjectID, &buildTrigger, cloudService)
		if triggerId != "Failed" {
			componentObject.TriggerID = triggerId
			updated, err := db.UpdateComponentData(componentObject)
			if !updated {
				fmt.Println("error", err)
				return "Failed"
			}
			return "Trigger created Successfully"
		} else {
			return "Trigger creation Failed"
		}
	} else {
		return "Trigger already exists"
	}
}

// Create BuildTrigger Object
func createBuildObject(component model.ComponentResponse) model.BuildTrigger {
	conf := &config.Config{Title: "ODJ"}
	conf.InitConfig()
	var buildTrigger model.BuildTrigger
	if strings.Contains(component.RepositoryURL, "bitbucket") {
		buildTrigger.Repo = "bitbucket"
	} else {
		buildTrigger.Repo = "github"
	}
	//fmt.Println(conf.RepoDetails.FileName + " -- " + conf.RepoDetails.UserName + " -- " + conf.RepoDetails.Branch + " -- " + component.ProjectID)
	buildTrigger.ProjectID = component.ProjectID
	buildTrigger.Filename = conf.RepoDetails.FileName
	buildTrigger.UserName = conf.RepoDetails.UserName
	buildTrigger.BranchName = conf.RepoDetails.Branch
	buildTrigger.RepositoryName = buildTrigger.Repo + "_" + buildTrigger.UserName + "_" + component.RepositoryName
	return buildTrigger
}

//Create Build Trigger
func triggerCloudBuild(projectID string, build *cloudbuild.BuildTrigger, cs *cloudbuild.Service) string {
	respBuildTrigger, err := cs.Projects.Triggers.Create(projectID, build).Do()
	if err != nil {
		fmt.Printf("error: " + err.Error())
		return "Failed"
	}
	fmt.Println("Trigger Id: " + respBuildTrigger.Id)
	return respBuildTrigger.Id
}

//Authenticate Source Repository
func authenticateSourceRepo(jsonPath_arg string) *sourcerepo.Service {
	fmt.Println("Source Repo Authentication")
	ctx := context.Background()
	cs, err := sourcerepo.NewService(ctx, option.WithCredentialsFile(jsonPath_arg))
	if err != nil {
		fmt.Printf("error: " + err.Error())
	}
	return cs
}

//Authenticate CloudBuild
func authenticateCloudBuild(jsonPath_arg string) *cloudbuild.Service {
	fmt.Println("Cloud Build Authentication")
	ctx := context.Background()
	cs, err := cloudbuild.NewService(ctx, option.WithCredentialsFile(jsonPath_arg))
	if err != nil {
		fmt.Printf("error: " + err.Error())
	}
	return cs
}

//Check whether Source Repository exists Repository or not
func checkRepoList(buildTrigger model.BuildTrigger, cs *sourcerepo.Service) bool {
	name := "projects/" + buildTrigger.ProjectID
	fmt.Println(name)
	respRepoList, err := cs.Projects.Repos.List(name).Do()
	if err != nil {
		fmt.Printf("error: " + err.Error())
		return false
	}
	tempRepoName := "projects/" + buildTrigger.ProjectID + "/repos/" + buildTrigger.RepositoryName

	for i := 0; i < len(respRepoList.Repos); i++ {
		fmt.Println(tempRepoName, respRepoList.Repos[i].Name)
		if tempRepoName == respRepoList.Repos[i].Name {
			fmt.Println("Repository Found in Source Repo", respRepoList.Repos[i].Name)
			return true
		}
	}
	return false

}

//Check whether BuildTrigger exists in CloudBuild or not
func checkTriggerList(trigger model.BuildTrigger, cloudService *cloudbuild.Service) bool {
	projectId := trigger.ProjectID
	respTriggerList, err := cloudService.Projects.Triggers.List(projectId).Do()
	if err != nil {
		fmt.Printf("error: " + err.Error())
		return false
	}
	for i := 0; i < len(respTriggerList.Triggers); i++ {

		//repoName := respTriggerList.Triggers[i].TriggerTemplate.RepoName
		//fmt.Println(trigger.RepositoryName, "----", respTriggerList.Triggers[i].TriggerTemplate.RepoName)
		if trigger.RepositoryName == respTriggerList.Triggers[i].TriggerTemplate.RepoName {
			fmt.Println("Build Trigger Already Exists" + respTriggerList.Triggers[i].TriggerTemplate.RepoName)
			return true
		}
	}
	return false

}

//Get the details of Build
func GetBuildDetails(build model.Build, db *db.Db) string {
	jsonPath := "cloudsloutions-fd422e010b14.json"
	cloudService := authenticateCloudBuild(jsonPath)
	build = getBuildList(build, cloudService)
	if build.BuildID == "0" {
		return "Failed"
	}
	build.BuildSeq = db.ReadBuildSeq(build)
	updated, err := db.UpdateCloudBuild(build)
	if !updated {
		fmt.Println("error", err.Error())
		return "Failed"
	}
	return "Success"
}

//Get List of Builds in a project
func getBuildList(build model.Build, cs *cloudbuild.Service) model.Build {

	projectId := build.ProjectID
	respBuildList, err := cs.Projects.Builds.List(projectId).Do()
	if err != nil {
		fmt.Printf("error: " + err.Error())
	}
	for i := 0; i < len(respBuildList.Builds); i++ {
		//fmt.Println(respBuildList.Builds[i].SourceProvenance.ResolvedRepoSource.CommitSha)
		if respBuildList.Builds[i].SourceProvenance.ResolvedRepoSource.CommitSha == build.CommitID {
			build.BuildID = respBuildList.Builds[i].Id
			build.BuildStatus = respBuildList.Builds[i].Status
			build.ImageID = respBuildList.Builds[i].Artifacts.Images[0]
			return build
		}
	}
	fmt.Println("Build Details Not Found")
	build.BuildID = "0"
	return build
}
