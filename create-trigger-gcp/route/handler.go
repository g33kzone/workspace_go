package route

import (
	"fmt"
	"net/http"
	"odj-deliver-cloudbuild/create_trigger"
	"odj-deliver-cloudbuild/db"
	"odj-deliver-cloudbuild/model"
	"odj-deliver-cloudbuild/webhook"

	"github.com/gin-gonic/gin"
)

//setCors handles cross origin requests, api requests via a web browser to an api on a different domain
func setCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Authorization, Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	}
}

//globalRecover returns json if there is an unhandled exception
func globalRecover(c *gin.Context) {
	defer func(c *gin.Context) {
		if rec := recover(); rec != nil {
			c.JSON(http.StatusInternalServerError, model.APIMessage{Code: 500, Text: fmt.Sprintf("Internal Server Error, %v", rec)})
		}
	}(c)
	c.Next()
}

// Create cloudbuild trigger
// @Summary Create cloudbuild trigger in GCP
// @Description create cloudbuild trigger and store triggerId in component table
// @ID odj-cloudbuild
// @Accept json
// @Param newComponent body model.ComponentResponse true "Information of Component to Create Cloud Build Trigger"
// @Success 200 {object} model.APIMessage
// @Failure 400 {object} model.APIMessage
// @Failure 401 {object} model.APIMessage
// @Failure 500 {object} model.APIMessage
// @Router /odj-deliver/create-trigger [post]
func handleClouddBuild(c *gin.Context, db *db.Db) string {
	response := create_trigger.HandleClouddBuild(c, db)
	if response != "" {
		return "Success"
	} else {
		return "Fail"
	}

}

// saveWebhookPayloadHandler godoc
// @Summary capture webook payload and store the build details into build table
// @Description capture a webook payload of Github and store build details into build table
// @ID odj-webhook-github
// @Accept json
// @Param githubResponse body model.GitResponse true "Information of Commit and Build to insert"
// @Success 200 {object} model.APIMessage
// @Failure 400 {object} model.APIMessage
// @Failure 401 {object} model.APIMessage
// @Failure 500 {object} model.APIMessage
// @Router /odj-deliver/webhook/github [post]
func handleWebHooksGithub(c *gin.Context, db *db.Db) {
	build := webhook.HandleWebHooksGithub(c, db)
	if build.ProjectID != "" {
		response := getBuildDetails(build, db)
		fmt.Println(response)
	}

}

// saveWebhookPayloadHandler godoc
// @Summary capture webook payload and store the build details into build table
// @Description capture a webook payload of Bitbucket and store build details into build table
// @ID odj-webhook-bitbucket
// @Accept json
// @Param bitbucketResponse body model.Response true "Information of Commit and Build to insert"
// @Success 200 {object} model.APIMessage
// @Failure 400 {object} model.APIMessage
// @Failure 401 {object} model.APIMessage
// @Failure 500 {object} model.APIMessage
// @Router /odj-deliver/webhook/bitbucket [post]
func handleWebHooksBitbucket(c *gin.Context, db *db.Db) {
	build := webhook.HandleWebHooksBitbucket(c, db)
	if build.ProjectID != "" {
		response := getBuildDetails(build, db)
		fmt.Println(response)

	}
}

func getBuildDetails(build model.Build, db *db.Db) string {
	response := create_trigger.GetBuildDetails(build, db)
	return response
}
