package main

import (
	"context"
	"fmt"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1"
	"google.golang.org/api/option"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"
)

var jsonPath = "devops-deliver-creds.json"

func cntxt() context.Context {
	return context.Background()
}

func main() {

	resourceURL := fmt.Sprintf("https://gcr.io/devops-deliver/github_g33kzone_nginx-container-analysis@sha256:c6d3ae83d0535c72dde8447419c5d2177493849e38a08823ae9440f528ef35a9")
	// resourceURL := fmt.Sprintf("https://gcr.io/devops-deliver/github_g33kzone_nginx-container-analysis@")
	projectID := "devops-deliver"

	ctx := context.Background()

	client, err := containeranalysis.NewClient(cntxt(), option.WithCredentialsFile(jsonPath))
	if err != nil {
		fmt.Printf("NewClient: %v", err)
	}
	defer client.Close()

	req := &grafeaspb.ListOccurrencesRequest{
		Parent:   fmt.Sprintf("projects/%s", projectID),
		Filter:   fmt.Sprintf("resourceUrl = %q kind = %q", resourceURL, "VULNERABILITY"),
		PageSize: 1000,
		// Filter: fmt.Sprintf("resourceUrl = %q kind = %q", resourceURL, "NORMAL"),
		// Filter: fmt.Sprintf("resourceUrl = %q ", resourceURL),
	}

	// req := &grafeaspb.ListNotesRequest{
	// 	Parent: fmt.Sprintf("projects/%s", projectID),
	// 	Filter: fmt.Sprintf("resourceUrl = %q kind = %q", resourceURL, "VULNERABILITY"),
	// }

	// it := client.GetGrafeasClient().ListNotes(ctx, req)
	it := client.GetGrafeasClient().ListOccurrences(ctx, req)

	var occ *grafeaspb.Occurrence
	// var occ *grafeaspb.Note
	occ, err = it.Next()

	// fmt.Println(occ.GetDiscovery())
	fmt.Println(occ.GetResourceUri())
	// fmt.Println(occ.GetVulnerability().GetSeverity())
	// fmt.Println(occ.GetVulnerability().GetCvssScore())

	reqOcc := &grafeaspb.GetOccurrenceRequest{
		Name: fmt.Sprintf("projects/%s/occurrences/%s", projectID, "d7f52807-90b1-4c0e-bfe6-be9230905072"),
	}

	occResponse, err := client.GetGrafeasClient().GetOccurrence(cntxt(), reqOcc)
	if err == nil {
		fmt.Println(occResponse.GetDiscovery().GetAnalysisStatus())
		fmt.Println(occResponse.GetResourceUri())
	}

}
