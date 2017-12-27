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

var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "List all subscriptions of the broker",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		req := &pb.SubscriptionsRequest{Category: "list"}
		if reply, err := brokerApi.Subscriptions(req); err != nil {
			fmt.Printf("Broker Api call failed:%s", err.Error())
			return
		} else if len(reply.Subscriptions) == 0 {
			fmt.Println("No subscriptions found in broker")
			return
		} else {
			for _, sub := range reply.Subscriptions {
				fmt.Printf("clientid:%s, topic:%s, attribute:%s",
					sub.ClientId, sub.Topic, sub.Attribute)
			}
		}
	},
}

var subscriptionsShowCmd = &cobra.Command{
	Use:     "show",
	Short:   "show client's subscriptions",
	Example: "sentel-ctrl show clientid",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage error, please see help")
			return
		}

		req := &pb.SubscriptionsRequest{Category: "show", ClientId: args[0]}
		if reply, err := brokerApi.Subscriptions(req); err != nil {
			fmt.Printf("Broker Api call failed:%s", err.Error())
			return
		} else if len(reply.Subscriptions) == 0 {
			fmt.Println("No subscription found in the broker")
			return
		} else {
			for _, sub := range reply.Subscriptions {
				fmt.Printf("clientid:%s, topic:%s, attribute:%s",
					sub.ClientId, sub.Topic, sub.Attribute)
			}
		}
	},
}
