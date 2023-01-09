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

package cmd

import (
	"encoding/json"
	"os"
	"runtime"

	"github.com/dapr/cli/pkg/print"
	"github.com/dapr/cli/pkg/standalone"
	"github.com/spf13/cobra"
)

var (
	subscribeAppID    string
	subscribePubSub   string
	subscribeTopic    string
	subscribeMetadata string
	subscribeSocket   string
)

var SubscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to a topic. Supported platforms: Self-hosted",
	Example: `
# Subscribe to sample topic in target pubsub via a subscribing app
dapr subscribe --subscribe-app-id myapp --pubsub target --topic sample

# Subscribe to sample topic in target pubsub via subscriber app without cloud event
dapr subscribe --subscribe-app-id myapp --pubsub target --topic sample --metadata '{"rawPayload":"true"}'

# Subscribe to sample topic in target pubsub via subscriber app with custom routes
dapr subscribe --subscribe-app-id myapp --pubsub target --topic sample --routes '{"rules": [{"match": "event.type == \'widget\'", "path": "/widgets"}], "default": "/products"}'
`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO(@daixiang0): add Windows support.
		if subscribeSocket != "" {
			if runtime.GOOS == "windows" {
				print.FailureStatusEvent(os.Stderr, "The unix-domain-socket option is not supported on Windows")
				os.Exit(1)
			}
		} else {
			print.WarningStatusEvent(os.Stderr, "Unix domain sockets are currently a preview feature")
		}

		metadata := make(map[string]interface{})
		if subscribeMetadata != "" {
			err := json.Unmarshal([]byte(subscribeMetadata), &metadata)
			if err != nil {
				print.FailureStatusEvent(os.Stderr, "Error parsing metadata as JSON. Error: %s", err)
				os.Exit(1)
			}
		}

		client := standalone.NewClient()
		err := client.Subscribe(subscribeAppID, subscribePubSub, subscribeTopic, subscribeSocket, metadata)
		if err != nil {
			print.FailureStatusEvent(os.Stderr, "Error subscribing to topic %s: %s", subscribeTopic, err)
			os.Exit(1)
		}

		print.SuccessStatusEvent(os.Stdout, "Subscription ended successfully")
	},
}

func init() {
	SubscribeCmd.Flags().StringVarP(&subscribeAppID, "subscribe-app-id", "a", "", "The ID of the subscribing app")
	SubscribeCmd.Flags().StringVarP(&subscribePubSub, "pubsub", "p", "", "The name of the pub/sub component")
	SubscribeCmd.Flags().StringVarP(&subscribeTopic, "topic", "t", "", "The topic to subscribe to")
	SubscribeCmd.Flags().StringVarP(&subscribeSocket, "unix-domain-socket", "u", "", "Path to a unix domain socket dir. If specified, Dapr API servers will use Unix Domain Sockets")
	SubscribeCmd.Flags().StringVarP(&subscribeMetadata, "metadata", "m", "", "The JSON serialized subscription metadata (optional)")
	SubscribeCmd.Flags().BoolP("help", "h", false, "Print this help message")
	SubscribeCmd.MarkFlagRequired("subscribe-app-id")
	SubscribeCmd.MarkFlagRequired("pubsub")
	SubscribeCmd.MarkFlagRequired("topic")
	RootCmd.AddCommand(SubscribeCmd)
}
