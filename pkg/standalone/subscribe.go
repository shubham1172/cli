/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package standalone

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/phayes/freeport"
)

const (
	subscribeHttpRoute = "/dapr/subscribe"
	eventsHttpRoute    = "/events"
)

type subscription struct {
	PubsubName string                 `json:"pubsubname"`
	Topic      string                 `json:"topic"`
	Route      string                 `json:"route"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// startSubscribeServer starts a web server on a random port
// and registers the required HTTP handlers.
func startSubscribeServer(s subscription) (int, error) {
	port, err := freeport.GetFreePort()
	if err != nil {
		return -1, err
	}

	http.HandleFunc(subscribeHttpRoute, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]subscription{s})
	})

	http.HandleFunc(eventsHttpRoute, func(w http.ResponseWriter, r *http.Request) {
		var event interface{}
		json.NewDecoder(r.Body).Decode(&event)
		fmt.Printf("Received event: %v", event)
	})

	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return port, nil
}

// Subscribe to a topic in a pubsub and receive messages.
func (s *Standalone) Subscribe(appID, pubsubName, topic, socket string, metadata map[string]interface{}) error {
	sub := subscription{
		PubsubName: pubsubName,
		Topic:      topic,
		Metadata:   metadata,
		Route:      eventsHttpRoute,
	}

	_, err := startSubscribeServer(sub)
	if err != nil {
		return err
	}
	return nil
}
