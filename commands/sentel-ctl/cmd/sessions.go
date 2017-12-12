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

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "List all MQTT session of the broker",
	Long:  `List all MQTT session of the broker, or specified session`,
	Run: func(cmd *cobra.Command, args []string) {
		req := &pb.SessionsRequest{Category: "list"}
		if reply, err := brokerApi.Sessions(req); err != nil {
			fmt.Printf("Broker Api call failed:%s", err.Error())
			return
		} else if len(reply.Sessions) == 0 {
			fmt.Println("No sessions found in broker")
			return
		} else {
			for _, session := range reply.Sessions {
				fmt.Printf("clientid=%s, created_at=%s", session.ClientId, session.CreatedAt)
			}
		}
	},
}

var sessionsShowCmd = &cobra.Command{
	Use:     "show",
	Short:   "List specified MQTT session of the broker",
	Example: "sentel-ctrl session show session_id",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage error, please see help")
			return
		}
		req := &pb.SessionsRequest{Category: args[0], ClientId: args[0]}
		if reply, err := brokerApi.Sessions(req); err != nil {
			fmt.Printf("Broker Api call failed:%s", err.Error())
			return
		} else if len(reply.Sessions) != 1 {
			fmt.Println("Error:sentel server return multiple sessions")
			return
		} else {
			session := reply.Sessions[0]
			fmt.Printf(`(clientid=%s, created_at=%s, clean_session=%t, 
				max_inflight=%d,infliaht=%d,inqueue=%d,droped=%d,awaiting_rel=%d,
				awaiting_comp=%d,awaiting_ack=%d)`,
				session.ClientId, session.CreatedAt, session.CleanSession,
				session.MessageMaxInflight,
				session.MessageInflight,
				session.MessageInQueue,
				session.MessageDropped,
				session.AwaitingRel,
				session.AwaitingComp,
				session.AwaitingAck)
		}
	},
}
