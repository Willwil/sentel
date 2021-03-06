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

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List all running services ",
	Run: func(cmd *cobra.Command, args []string) {
		if reply, err := brokerApi.Services(&pb.ServicesRequest{Category: "list"}); err != nil {
			fmt.Printf("Broker Api call failed:%s", err.Error())
			return
		} else {
			for _, service := range reply.Services {
				fmt.Printf("service '%s'", service.ServiceName) // TODO: add service's status
			}
		}
	},
}
