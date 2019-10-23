// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	// containeranalysis "google.golang.org/api/containeranalysis/v1beta1"
)

// OccurenceNotification is the payload of a Pub/Sub event.
type OccurenceNotification struct {
	Name             string `json:"name"`
	Kind             string `json:"kind"`
	NotificationTime string `json:"notificationTime"`
}

// PubSubMessage is the payload of a Pub/Sub event.
// type PubSubMessage struct {
// 	Data []byte `json:"data"`
// }

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(ctx context.Context, msg *pubsub.Message) error {

	s := string(msg.Data)
	// log.Println(decodedString)

	data := OccurenceNotification{}

	json.Unmarshal([]byte(s), &data)

	fmt.Printf("Name: %s", data.Name)
	fmt.Printf("Kind: %s", data.Kind)
	fmt.Printf("Time: %s", data.NotificationTime)

	return nil
}
