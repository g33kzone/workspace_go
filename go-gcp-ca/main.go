package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gcp-ca/constants"
	"strings"
	"time"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1"
	"cloud.google.com/go/pubsub"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"
)

// OccurenceNotification is the payload of a Pub/Sub event.
type OccurenceNotification struct {
	Name             string `json:"name"`
	Kind             string `json:"kind"`
	NotificationTime string `json:"notificationTime"`
}

func main() {
	fmt.Println("Hello World")
	fmt.Println(time.Now().Format(time.RFC850))
	err := FetchDiscoveryOccurence()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(time.Now().Format(time.RFC850))
}

var jsonPath = "devops-deliver-creds.json"

func cntxt() context.Context {
	return context.Background()
}

// SplitString splits the string into a array
func SplitString(str string, sep string) []string {
	return strings.Split(str, sep)
}

//validates our PubSub Service at GCP
func authenticatePubSub() (*pubsub.Client, error) {

	cs, err := pubsub.NewClient(cntxt(), constants.DeliverProjectID, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Info("Error while GCP authentication for Pub/Sub")
		log.Error(err.Error())
	}
	return cs, err
}

//validates our Container Analysis Service at GCP
func authenticateContainerAnalysis() (*containeranalysis.Client, error) {

	cs, err := containeranalysis.NewClient(cntxt(), option.WithCredentialsFile(jsonPath))
	if err != nil {

		log.Info("Error while GCP authentication for Container Analysis")
		log.Error(err.Error())
	}
	return cs, err
}

// FetchDiscoveryOccurence will asynchronously pull Discovery Occurences for Vulnerability Scanning
func FetchDiscoveryOccurence() error {

	psClient, err := authenticatePubSub()
	if err == nil {
		ch := make(chan bool)
		go pullMsgs(psClient, ch)

		select {
		case msg := <-ch:
			fmt.Println(fmt.Sprintf("Read from channel - %v", msg))
		case <-time.After(time.Second * 20):
			fmt.Println("Timed out")
			// cancel()
		}
	}

	return err
}

func pullMsgs(client *pubsub.Client, ch chan bool) error {
	var isVulnerable bool
	// var mu sync.Mutex
	// receivedCount := 0
	sub := client.Subscription(constants.SubscriptionOccurrence)
	vulnerabilitySeverity := []string{grafeaspb.Severity_CRITICAL.String(), grafeaspb.Severity_HIGH.String(), grafeaspb.Severity_MEDIUM.String()}

	// cctx, cancel := context.WithCancel(cntxt())
	cctx, cancel := context.WithTimeout(cntxt(), time.Second*5)
	// defer cancel()
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		// fmt.Println("Pub/Sub Asynch Pull started...")
		msg.Ack()

		// mu.Lock()
		// defer mu.Unlock()

		cctx.Done()

		s := string(msg.Data)
		data := OccurenceNotification{}
		json.Unmarshal([]byte(s), &data)

		// Enum for DISCOVERY
		if data.Kind == grafeaspb.NoteKind(6).String() {
			// receivedCount++
			// if receivedCount == 10 {
			// 	cancel()
			// }

			fmt.Println("Discovery Occurence found...")
			fmt.Println(data.Name)
			resourceURI, isexist := GetDiscoveryOccurrence(data.Name)
			fmt.Println("Resource URI")
			fmt.Println(resourceURI)
			if isexist {
				nameArray := SplitString(data.Name, "/")
				occurrenceRequest := &grafeaspb.ListOccurrencesRequest{
					Parent:   fmt.Sprintf("projects/%s", nameArray[1]),
					Filter:   fmt.Sprintf("resourceUrl = %q kind = %q", resourceURI, grafeaspb.NoteKind(1).String()),
					PageSize: 1000,
				}

				caClient, err := authenticateContainerAnalysis()
				if err == nil {
					occurrenceResponse := caClient.GetGrafeasClient().ListOccurrences(cntxt(), occurrenceRequest)

					for {
						vulnerabilityOccurence, err := occurrenceResponse.Next()
						if err != nil {
							fmt.Println("Read all Vulnerability Occurences....exit from loop")
							break
						}
						fmt.Println(vulnerabilityOccurence.GetVulnerability().GetSeverity().String())
						if contains(vulnerabilitySeverity, vulnerabilityOccurence.GetVulnerability().GetSeverity().String()) {
							isVulnerable = true
							fmt.Println("Vulnerability Found!!")
							break
						}
					}
					fmt.Println("Is Vulnerable")
					fmt.Println(isVulnerable)
					ch <- isVulnerable
					cancel()
					return
				}
			}
		}
	})
	if err != nil {
		fmt.Println("Error encountered with Pub/Sub Asynch pull")
		return err
	}
	fmt.Println("Testing timeout....")
	return nil
}

// GetDiscoveryOccurrence fetches details for Discovery Occurrence
func GetDiscoveryOccurrence(discoveryOccurrence string) (string, bool) {
	fmt.Println(fmt.Sprintf("Occurence : %s", discoveryOccurrence))

	reqOcc := &grafeaspb.GetOccurrenceRequest{
		Name: discoveryOccurrence,
	}

	fmt.Println(fmt.Sprintf("Discovery Occ - %s", discoveryOccurrence))

	caClient, err := authenticateContainerAnalysis()
	if err == nil {
		response, err := caClient.GetGrafeasClient().GetOccurrence(cntxt(), reqOcc)
		fmt.Println(fmt.Sprintf("Discovery Occurrence : %s", discoveryOccurrence))
		fmt.Println(fmt.Sprintf("Discovery Response : %s", response.GetDiscovery().GetAnalysisStatus().String()))
		sha256Digest := StringSplit(response.GetResourceUri(), "@", 1)
		fmt.Println(sha256Digest)
		if err == nil {
			fmt.Println(fmt.Sprintf("Analysis : %s", response.GetDiscovery().GetAnalysisStatus().String()))
			if response.GetDiscovery().GetAnalysisStatus().String() == "FINISHED_SUCCESS" {
				return response.GetResourceUri(), true
			}
		}
	}
	return "", false
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// StringSplit -
func StringSplit(str string, sep string, pos int) string {
	strArray := strings.Split(str, sep)
	return strArray[pos]
}
