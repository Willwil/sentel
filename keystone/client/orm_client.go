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
package client

import (
	"errors"

	"github.com/cloustone/sentel/keystone/orm"
)

func CreateObject(obj orm.Object) error {
	/*
		url := "http://" + hosts + "/keystone/api/v1/orm/object"
		if buf, err := json.Marshal(obj); err == nil {
			req := bytes.NewBuffer([]byte(buf))
			resp, err := http.Post(url, format, req)
			if err != nil {
				return err
			} else {
				if resp.StatusCode == http.StatusOK {
					return nil
				}
				body, _ := ioutil.ReadAll(resp.Body)
				result := authResponse{}
				if err := json.Unmarshal(body, &result); err == nil && result.Message != "" {
					return fmt.Errorf("%s", result.Message)
				}
			}
		}
	*/
	return errors.New("object creation failed")
}

func AccessObject(objid string, accessId string, right orm.AccessRight) error {
	return nil
}

func DestroyObject(objid string, accessId string) {}
func AssignObjectRight(objid string, accessId string, right orm.AccessRight) error {
	return nil
}
