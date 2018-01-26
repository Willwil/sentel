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

package api

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloustone/sentel/meter/collector"
	"github.com/labstack/echo"
)

// getClusterClientInfo return clients info in cluster
func getClientInfo(ctx echo.Context) error {
	clientId := ctx.Param("clientId")
	if clientId == "" {
		return ctx.JSON(BadRequest,
			&response{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		return ctx.JSON(ServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("clients")
	defer session.Close()

	clients := []collector.Client{}
	if err := c.Find(bson.M{"ClientId": clientId}).Limit(100).Iter().All(&clients); err != nil {
		return ctx.JSON(NotFound,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(OK, &response{
		Success: true,
		Result:  clients,
	})
}
