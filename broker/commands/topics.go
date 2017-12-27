//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package commands

import (
	"fmt"

	pb "github.com/cloustone/sentel/broker/rpc"

	"github.com/spf13/cobra"
)

var topicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "List all topics of the broker",
	Long:  `List All topics of the broker and inquery detail topic information`,
	Run: func(cmd *cobra.Command, args []string) {
		req := &pb.TopicsRequest{Category: "list"}
		if reply, err := brokerApi.Topics(req); err != nil {
			fmt.Printf("Broker Api call failed:%s", err.Error())
			return
		} else if len(reply.Topics) == 0 {
			fmt.Println("No topics found in broker")
			return
		} else {
			for _, topic := range reply.Topics {
				fmt.Printf("%s, %s", topic.Topic, topic.Attribute)
			}
		}
	},
}

var topicsShowCmd = &cobra.Command{
	Use:     "show",
	Short:   "show client's topics",
	Example: "sentel-ctrl show clientid",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage error, please see help")
			return
		}
		req := &pb.TopicsRequest{Category: "show", ClientId: args[0]}
		if reply, err := brokerApi.Topics(req); err != nil {
			fmt.Printf("Broker Api call failed:%s", err.Error())
			return
		} else if len(reply.Topics) == 0 {
			fmt.Printf("No topics found for client '%s'", args[0])
			return
		} else {
			for _, topic := range reply.Topics {
				fmt.Printf("%s, %s", topic.Topic, topic.Attribute)
			}
		}

	},
}
