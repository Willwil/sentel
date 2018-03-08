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

package web

import (
	"net/http"
	"strings"

	"github.com/cloustone/sentel/pkg/goshiro/shiro"
)

// WebRequest is wraper for authorization request based on http.Request
type WebRequest struct {
	resource string
	action   string
}

func NewRequest(mgr shiro.SecurityManager, req *http.Request, ctx shiro.RequestContext) (shiro.Request, error) {
	// get  resource requested from subject
	policy, err := mgr.GetPolicy(ctx.Path(), ctx)
	if err != nil {
		return nil, err
	}
	// generate the requested real resource, using parameter from context to replace identifier in resource field
	resource := "/"
	items := strings.Split(policy.Resource, "/")
	for _, item := range items {
		switch item[0] {
		case '$':
			val := ctx.Get(item[1:]).(string)
			resource += val
		case ':':
			val := ctx.Param(item[1:])
			resource += val
		}
	}
	action := req.Method
	return &WebRequest{
		resource: resource,
		action:   action,
	}, nil
}

func (w *WebRequest) GetAction() string   { return w.action }
func (w *WebRequest) GetResource() string { return w.resource }
