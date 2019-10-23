package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/option"

	"cloud.google.com/go/pubsub"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"
)

var (
	subscription *pubsub.Subscription
)

// const sub = "NotesSubscription"
const sub = "OccurrencesSubscription"

var jsonPath = "devops-deliver-creds.json"

// OccurenceNotification is the payload of a Pub/Sub event.
type OccurenceNotification struct {
	Name             string `json:"name"`
	Kind             string `json:"kind"`
	NotificationTime string `json:"notificationTime"`
}

func cntxt() context.Context {
	return context.Background()
}

func main() {
	fmt.Println("Hello World!")

	client, err := pubsub.NewClient(cntxt(), "devops-deliver", option.WithCredentialsFile(jsonPath))

	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}

	fmt.Println("Before pull logic....")
	// Pull messages via the subscription.
	if err := pullMsgs(client, sub); err != nil {
		log.Fatal(err)
	}

	fmt.Println("After pull logic....")

	go testGo("Manish")

	subscriptionResponse := client.Subscriptions(cntxt())

	for {
		subscription, err := subscriptionResponse.Next()
		if err != nil {
			break
		}
		subscriptionConfig, err := subscription.Config(cntxt())
		fmt.Println(subscriptionConfig.Labels["criticality"])
		fmt.Println(subscription.ID())
		fmt.Println(subscription.String())
	}

}

func testGo(name string) {
	fmt.Println(fmt.Sprintf("Welcome %s", name))
}

func pullMsgs(client *pubsub.Client, subName string) error {

	ctx := context.Background()
	sub := client.Subscription(subName)
	// cctx, cancel := context.WithCancel(ctx)
	cctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()

		s := string(msg.Data)
		data := OccurenceNotification{}
		json.Unmarshal([]byte(s), &data)

		if data.Kind == grafeaspb.NoteKind(6).String() {
			// parse Occurrence ID
			// fetch results for above fetched Occurrence ID

			fmt.Println(fmt.Sprintf("Name: %s", data.Name))
			fmt.Println(fmt.Sprintf("Kind: %s", data.Kind))
		}

		select {
		case <-time.After(2 * time.Second):
			fmt.Println("overslept")
		case <-ctx.Done():
			fmt.Println("Oops...got killed")
		}
	})
	if err != nil {
		return err
	}
	return nil
}
