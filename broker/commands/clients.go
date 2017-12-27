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

var clientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "Inquery and control connected client",
	Long:  `Inquery client information and controll client`,
	Run: func(cmd *cobra.Command, args []string) {
		if reply, err := brokerApi.Clients(&pb.ClientsRequest{Category: "list"}); err != nil {
			fmt.Printf("Broker Api call failed:%s", err.Error())
			return
		} else if len(reply.Clients) == 0 {
			fmt.Println("No clients found in the broker")
			return
		} else {
			for _, info := range reply.Clients {
				fmt.Printf("clientId:%s, cleanSession:%T, peername:%s, connectTime:%s",
					info.UserName, info.CleanSession, info.PeerName, info.ConnectTime)
			}
		}
	},
}

var clientsShowCmd = &cobra.Command{
	Use:     "show",
	Short:   "show specified client's information",
	Example: "sentel-ctrl clients show client_id",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage error, please see help")
			return
		}
		if reply, err := brokerApi.Clients(&pb.ClientsRequest{Category: args[0], ClientId: args[0]}); err != nil {
			fmt.Printf("Broker Api call failed:%s", err.Error())
			return
		} else if len(reply.Clients) == 0 {
			fmt.Printf("No client '%s' found in the broker", args[0])
			return
		} else {
			info := reply.Clients[0]
			fmt.Printf("username:%s, cleanSession:%T, peername:%s, connectTime:%s",
				info.UserName, info.CleanSession, info.PeerName, info.ConnectTime)
		}
	},
}

var clientsKickoffCmd = &cobra.Command{
	Use:     "kickoff",
	Short:   "show specified client's information",
	Example: "sentel-ctrl clients kickoff client_id",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage error, please see help")
			return
		}
		_, err := brokerApi.Clients(&pb.ClientsRequest{Category: "kickoff", ClientId: args[0]})
		if err != nil {
			fmt.Printf("Broker Api call failed:%s", err.Error())
			return
		}
	},
}
