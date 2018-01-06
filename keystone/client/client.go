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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloustone/sentel/keystone/auth"
	"github.com/cloustone/sentel/keystone/ram"
	"github.com/cloustone/sentel/pkg/config"
)

type Client struct {
	hosts string
}

type apiResponse struct {
	Message string `json:"message"`
}

func New(c config.Config) (*Client, error) {
	hosts, err := c.String("keystone", "hosts")
	if err != nil {
		return nil, err
	}
	return &Client{hosts: hosts}, nil
}

func (p *Client) Authenticate(opts interface{}) error {
	buf, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	url := ""
	format := "application/json;charset=utf-8"
	req := bytes.NewBuffer([]byte(buf))

	switch opts.(type) {
	case auth.ApiAuthParam:
		url = "http://" + p.hosts + "/keystone/api/v1/auth/api"
	default:
		return errors.New("invalid type")
	}
	resp, err := http.Post(url, format, req)
	if err == nil && resp.StatusCode == http.StatusOK {
		return nil
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		result := apiResponse{}
		if err := json.Unmarshal(body, &result); err == nil && result.Message != "" {
			return fmt.Errorf("%s", result.Message)
		}
	}

	return err
}

func (p *Client) Authorize(accessId string, resource string, action string) error {
	url := fmt.Sprintf("http://%s/keystone/api/v1/ram/resource?resource=%s&accessId=%s&accessRight=%s", p.hosts, resource, accessId, action)
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		return nil
	} else if resp != nil {
		return handleResponse(resp)
	}
	return err
}

func handleResponse(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	result := apiResponse{}
	err := json.Unmarshal(body, &result)
	if err == nil && result.Message != "" {
		return fmt.Errorf("%s", result.Message)
	}
	return err
}

func (p *Client) CreateResource(accessId string, res ram.ResourceCreateOption) error {
	url := fmt.Sprintf("http://%s/keystone/api/v1/ram/resource?accessId=%s", p.hosts, accessId)
	format := "application/json;charset=utf-8"

	if buf, err := json.Marshal(res); err == nil {
		req := bytes.NewBuffer([]byte(buf))
		resp, err := http.Post(url, format, req)
		if err != nil {
			return err
		} else {
			return handleResponse(resp)
		}
	}
	return errors.New("object creation failed")
}

func (p *Client) AccessResource(res string, accessId string, action ram.Action) error {
	url := fmt.Sprintf("http://%s/keystone/api/v1/ram/resource?resource=%s&accessId=%s&action=%d", p.hosts, res, accessId, action)
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		return nil
	} else if resp != nil {
		return handleResponse(resp)
	}
	return err
}

func (p *Client) DestroyResource(resourceId string, accessId string) error {
	opt := ram.ResourceDestroyOption{ObjectId: resourceId, AccessId: accessId}
	if buf, err := json.Marshal(opt); err == nil {
		url := fmt.Sprintf("http://%s/keystone/api/v1/ram/resource", p.hosts)
		body := bytes.NewBuffer([]byte(buf))

		client := &http.Client{}
		req, err := http.NewRequest("DELETE", url, body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			return nil
		} else if resp != nil {
			return handleResponse(resp)
		}
		return err
	}
	return errors.New("object destroy failed")
}

func (p *Client) AddResourceGrantee(res string, accessId string, right ram.Right) error {
	url := fmt.Sprintf("http://%s/keystone/api/v1/ram/resource?resource=%s&accessId=%s&right=%s",
		p.hosts, res, accessId, string(right))
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		return nil
	} else if resp != nil {
		return handleResponse(resp)
	}
	return err
}

func (p *Client) CreateAccount(account string) error {
	url := fmt.Sprintf("http://%s/keystone/api/v1/ram/account/%s", p.hosts, account)
	format := "application/json;charset=utf-8"
	resp, err := http.Post(url, format, nil)
	if err == nil && resp.StatusCode == http.StatusOK {
		return nil
	} else if resp != nil {
		return handleResponse(resp)
	}
	return err
}

func (p *Client) DestroyAccount(account string) error {
	url := fmt.Sprintf("http://%s/keystone/api/v1/ram/account/%s", p.hosts, account)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		return nil
	} else if resp != nil {
		return handleResponse(resp)
	}
	return err
}
