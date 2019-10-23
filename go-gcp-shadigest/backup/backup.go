package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	sha256Digest, err := validateImage(dockerImage)
	if err == nil {
		fmt.Println(sha256Digest)
	} else {
		fmt.Println("Docker Image with Tag not found")
	}
}

func validateImage(dockerImage string) (string, error) {

	url := createURL(dockerImage)
	imageTag := stringSplit(dockerImage, ":", 1)

	token, err := tokenFor("https://www.googleapis.com/auth/cloud-platform")
	if err == nil {
		req, err := http.NewRequest("GET", url, nil)
		if err == nil {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			resp, _ := http.DefaultClient.Do(req)

			if resp.StatusCode == 200 {
				sha256Digest, err := parseResponse(resp, imageTag)
				return sha256Digest, err
			} else {
				fmt.Println(fmt.Scanf("Error Status Code - %d", resp.StatusCode))
			}
		}
	}
	return "", err
}

func parseResponse(response *http.Response, imageTag string) (string, error) {
	var data dockerManifest

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(body, &data)
	if err == nil {
		for key, value := range data.Manifest {
			if len(value.Tag) > 0 {
				for _, tag := range value.Tag {
					if imageTag == tag {
						tagMatch = true
						fmt.Println("Match Found")
						fmt.Printf("key[%s]\n", key)
						fmt.Println(tag)
						return key, nil
					}
				}
			}
		}
	}

	if !tagMatch {
		fmt.Println("No Match")
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
