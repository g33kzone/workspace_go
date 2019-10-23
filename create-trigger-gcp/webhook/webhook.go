package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"odj-deliver-cloudbuild/db"
	"odj-deliver-cloudbuild/model"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

//Handle the bitbucket's webhook
func HandleWebHooksBitbucket(c *gin.Context, db *db.Db) model.Build {
	var webhookObject model.Response
	fmt.Println("Webhook data received for bitbucket...")
	webhookData, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("headers: %v\n", c.Request.Header)
	fmt.Printf("body: %v\n", c.Request.Body)
	json.Unmarshal(webhookData, &webhookObject)
	fmt.Println("Repository Name: ", webhookObject.Repository.Name)
	build1 := db.ReadProducComponentAndDockerRegistry(webhookObject.Repository.Name)
	projectId := db.ReadProjectId(build1.ProductName)
	build := model.Build{
		ProductName:    build1.ProductName,
		ComponentName:  build1.ComponentName,
		ProjectID:      projectId,
		DockerRegistry: build1.DockerRegistry,
		RepositoryName: webhookObject.Repository.Name,
		RepositoryUrl:  webhookObject.Repository.Link.Html.RepositoryUrl,
		BranchName:     webhookObject.Push.Changes[0].CommitBranch.Branch,
		CommitID:       webhookObject.Push.Changes[0].Commits[0].CommitID,
		CommitMsg:      webhookObject.Push.Changes[0].Commits[0].CommitMsg,
	}
	inserted, err := db.InsertCloudBuild(build)
	if !inserted {
		c.JSON(http.StatusInternalServerError, model.APIMessage{Code: 500,
			Text: fmt.Sprintf("Unexpected Error ! %s", err.Error())})
		var emptyBuild model.Build
		return emptyBuild
	}

	c.JSON(http.StatusOK, model.APIMessage{Code: 200, Text: "Build details inserted"})
	return build

}

//Handle the Github's webhook
func HandleWebHooksGithub(c *gin.Context, db *db.Db) model.Build {
	var webhookObject model.GitResponse
	fmt.Println("Webhook data received from github...")

	webhookData, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("headers: %v\n", c.Request.Header)
	fmt.Printf("body: %v\n", c.Request.Body)
	json.Unmarshal(webhookData, &webhookObject)
	fmt.Println("Repository Name: ", webhookObject.GitRepository.Name)
	fmt.Println("Branch Name: ", webhookObject.GitRepository.GitBranch)
	fmt.Println("Commit ID: ", webhookObject.GitCommits[0].CommitId)
	fmt.Println("Commit ID: ", webhookObject.GitCommits[0].CommitMsg)
	fmt.Println("Repository Url: ", webhookObject.GitRepository.RepoUrl)
	build1 := db.ReadProducComponentAndDockerRegistry(webhookObject.GitRepository.Name)
	projectId := db.ReadProjectId(build1.ProductName)
	build := model.Build{
		ProductName:    build1.ProductName,
		ComponentName:  build1.ComponentName,
		ProjectID:      projectId,
		DockerRegistry: build1.DockerRegistry,
		RepositoryName: webhookObject.GitRepository.Name,
		RepositoryUrl:  webhookObject.GitRepository.RepoUrl,
		BranchName:     webhookObject.GitRepository.GitBranch,
		CommitID:       webhookObject.GitCommits[0].CommitId,
		CommitMsg:      webhookObject.GitCommits[0].CommitMsg,
	}

	inserted, err := db.InsertCloudBuild(build)
	if !inserted {
		c.JSON(http.StatusInternalServerError, model.APIMessage{Code: 500,
			Text: fmt.Sprintf("Unexpected Error ! %s", err.Error())})
		var emptyBuild model.Build
		return emptyBuild
	}
	c.JSON(http.StatusOK, model.APIMessage{Code: 200, Text: "Build details inserted"})
	return build

}
