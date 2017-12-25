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
package base

import (
	com "github.com/cloustone/sentel/common"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo"
)

func SyncProduceMessage(ctx echo.Context, topic string, value sarama.Encoder) error {
	c := ctx.(*ApiContext).Config
	hosts := c.MustString("apiserver", "kafka")
	key := "iot"
	return com.SyncProduceMessage(hosts, key, topic, value)
}

func AsyncProduceMessage(ctx echo.Context, topic string, value sarama.Encoder) error {
	c := ctx.(*ApiContext).Config
	hosts := c.MustString("apiserver", "kafka")
	key := "iot"
	return com.AsyncProduceMessage(hosts, key, topic, value)
}
