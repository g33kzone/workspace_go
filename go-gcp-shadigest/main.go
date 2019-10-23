package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var jsonPath = "devops-deliver-creds.json"
var dockerImage = "gcr.io/devops-deliver/github_g33kzone_nginx-container-analysis:30f5d0d"
var sha256 map[string]shaData
var tagMatch bool

type dockerManifest struct {
	Manifest map[string]shaData `json:"manifest"`
	Name     string             `json:"name"`
	Tags     []string           `json:"tags"`
	Child    []interface{}      `json:"child"`
}

type shaData struct {
	ImageSizeBytes string   `json:"imageSizeBytes"`
	LayerID        string   `json:"layerId"`
	MediaType      string   `json:"mediaType"`
	Tag            []string `json:"tag"`
	TimeCreatedMs  string   `json:"timeCreatedMs"`
	TimeUploadedMs string   `json:"timeUploadedMs"`
}

type manifest struct {
	Data map[string]shaData `json:"manifest"`
}

func cntxt() context.Context {
	return context.Background()
}

func main() {
	imageTag := stringSplit(dockerImage, ":", 1)
	response, err := fetchDockerManifests(dockerImage)
	if err == nil {
		if response.StatusCode == 200 {
			sha256Digest, err := fetchSHA256Digest(response, imageTag)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}
			fmt.Println("SHA 256 Digest")
			fmt.Println(sha256Digest)
		}
	}
}

func fetchDockerManifests(dockerImage string) (*http.Response, error) {
	var resp *http.Response
	url := createURL(dockerImage)

	token, err := tokenFor("https://www.googleapis.com/auth/cloud-platform")
	if err == nil {
		req, err := http.NewRequest("GET", url, nil)
		if err == nil {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			resp, err = http.DefaultClient.Do(req)
			if err == nil {
				return resp, nil
			}
		}
	}
	return resp, err
}

func fetchSHA256Digest(response *http.Response, imageTag string) (string, error) {
	var data dockerManifest
	var isTagExist bool

	responseBody, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(responseBody, &data)
		if err == nil {
			for key, value := range data.Manifest {
				if len(value.Tag) > 0 {
					for _, tag := range value.Tag {
						if imageTag == tag {
							isTagExist = true
							fmt.Println("Match Found!!")
							return key, nil
						}
					}
				}
			}
		}
	}

	if !isTagExist {
		err = errors.New("Image Tag not found")
	}

	return "", err
}

func createURL(dockerImage string) string {
	imageString := stringSplit(dockerImage, ":", 0)

	host := stringSplit(imageString, "/", 0)
	projectID := stringSplit(imageString, "/", 1)
	image := stringSplit(imageString, "/", 2)

	return fmt.Sprintf("https://%s/v2/%s/%s/tags/list", host, projectID, image)
}

func tokenFor(scopes ...string) (string, error) {
	keyBytes, err := ioutil.ReadFile(jsonPath)
	if err == nil {
		var creds *google.Credentials

		creds, err = google.CredentialsFromJSON(cntxt(), keyBytes, scopes...)
		// creds, err := google.FindDefaultCredentials(cntxt(), scopes...)
		if err == nil {
			var token *oauth2.Token
			token, err = creds.TokenSource.Token()
			if err == nil {
				return token.AccessToken, err
			}
		}
	}
	return "", err
}

func stringSplit(str string, sep string, pos int) string {
	strArray := strings.Split(str, sep)
	return strArray[pos]
}
