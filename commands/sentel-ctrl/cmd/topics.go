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

package cmd

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
		// Print topic list
		if reply, err := brokerApi.Topics(req); err != nil {
			fmt.Println("Broker Api call failed:%s", err.Error())
			return
		} else {
			if len(reply.Topics) == 0 {
				fmt.Println("No topics found in broker")
				return
			}
			for _, topic := range reply.Topics {
				fmt.Printf("%s, %s", topic.Topic, topic.Attribute)
			}
		}
	},
}

var topicsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show client's topics",
	Run: func(cmd *cobra.Command, args []string) {
		req := &pb.TopicsRequest{Category: "show"}
		if len(args) != 1 {
			fmt.Println("Usage error, please see help")
			return
		}
		req.ClientId = args[1]
		if reply, err := brokerApi.Topics(req); err != nil {
			fmt.Println("Error:%v", err)
			return
		} else if len(reply.Topics) != 1 {
			fmt.Println("Error:sentel server return multiple topics")
			return
		} else {
			topic := reply.Topics[0]
			fmt.Printf("%s, %s", topic.Topic, topic.Attribute)
		}

	},
}
