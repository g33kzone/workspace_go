package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router = gin.Default()
	initializeRoutes()
	router.Run("0.0.0.0:8082")
}
func initializeRoutes() {
	router.POST("/webhook/github", handleWebHooksGithub)
	router.POST("/webhook/bitbucket", handleWebHooksBitbucket)
	router.GET("/test", handleWebHooksTest)
}

type Commit struct {
	CommitID string `json:"hash"`
}
type CommitBranch struct {
	Branch string `json:"name"`
}
type Change struct {
	Closed       bool         `json:"closed"`
	Commits      []Commit     `json:"commits"`
	CommitBranch CommitBranch `json:"new"`
}

type Push struct {
	Changes []Change `json:"changes"`
}
type Repository struct {
	Name string `json:"full_name"`
}
type Response struct {
	Push       Push       `json:"push"`
	Repository Repository `json:"repository"`
}

type GitCommits struct {
	CommitId string `json:"id"`
}
type GitRepository struct {
	FullName  string `json:"full_name"`
	GitBranch string `json:"master_branch"`
}
type GitResponse struct {
	GitCommits    []GitCommits  `json:"commits"`
	GitRepository GitRepository `json:"repository"`
}

func handleWebHooksBitbucket(c *gin.Context) {
	fmt.Println("Webhook data received for bitbucket...")

	var webhookObject Response

	webhookData, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("headers: %v\n", c.Request.Header)
	fmt.Printf("body: %v\n", c.Request.Body)
	json.Unmarshal(webhookData, &webhookObject)
	//fmt.Println(webhookObject)
	//fmt.Println(webhookData)
	fmt.Println("Repository Name: ", webhookObject.Repository.Name)
	fmt.Println("Branch Name: ", webhookObject.Push.Changes[0].CommitBranch.Branch)
	fmt.Println("Commit ID: ", webhookObject.Push.Changes[0].Commits[0].CommitID)

}

func handleWebHooksGithub(c *gin.Context) {
	fmt.Println("Webhook data received from github...")

	var webhookObject GitResponse

	webhookData, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("headers: %v\n", c.Request.Header)
	fmt.Printf("body: %v\n", c.Request.Body)
	json.Unmarshal(webhookData, &webhookObject)
	fmt.Println(webhookObject)
	fmt.Println("Repository Name: ", webhookObject.GitRepository.FullName)
	fmt.Println("Branch Name: ", webhookObject.GitRepository.GitBranch)
	fmt.Println("Commit ID: ", webhookObject.GitCommits[0].CommitId)
}
func handleWebHooksTest(c *gin.Context) {
	fmt.Println("Webhook Test")
}
