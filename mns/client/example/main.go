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

package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/cloustone/sentel/mns/client"
	"github.com/gogap/logs"
)

type appConf struct {
	Url             string `json:"url"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

func main() {
	conf := appConf{}

	if bFile, e := ioutil.ReadFile("app.conf"); e != nil {
		panic(e)
	} else {
		if e := json.Unmarshal(bFile, &conf); e != nil {
			panic(e)
		}
	}

	mnsclient := client.NewMnsClient(conf.Url,
		conf.AccessKeyId,
		conf.AccessKeySecret)

	msg := client.MessageSendRequest{
		MessageBody:  []byte("hello gogap/client"),
		DelaySeconds: 0,
		Priority:     8}

	queue := client.NewMnsQueue("test", mnsclient)
	ret, err := queue.SendMessage(msg)

	if err != nil {
		logs.Error(err)
	} else {
		logs.Pretty("response:", ret)
	}

	respChan := make(chan client.MessageReceiveResponse)
	errChan := make(chan error)
	go func() {
		for {
			select {
			case resp := <-respChan:
				{
					logs.Pretty("response:", resp)
					logs.Debug("change the visibility: ", resp.ReceiptHandle)
					if ret, e := queue.ChangeMessageVisibility(resp.ReceiptHandle, 5); e != nil {
						logs.Error(e)
					} else {
						logs.Pretty("visibility changed", ret)
					}

					logs.Debug("delete it now: ", resp.ReceiptHandle)
					if e := queue.DeleteMessage(resp.ReceiptHandle); e != nil {
						logs.Error(e)
					}
				}
			case err := <-errChan:
				{
					logs.Error(err)
				}
			}
		}

	}()

	queue.ReceiveMessage(respChan, errChan)
	for {
		time.Sleep(time.Second * 1)
	}

}
