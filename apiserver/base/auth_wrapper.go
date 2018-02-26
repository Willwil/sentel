//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
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
	"errors"

	"github.com/cloustone/sentel/goshiro/auth"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/labstack/echo"
)

var (
	noAuth = true
)

var (
	ErrorInvalidArgument = errors.New("invalid argument")
	ErrorAuthDenied      = errors.New("authentication denied")
)

type ApiAuthParam struct {
	Format      string `json:"format"`
	Version     string `json:"version"`
	AccessId    string `json:"accessId"`
	Signature   string `json:"sign"`
	SignMethod  string `json:"signMethod"`
	SignVersion string `json:"signVersion"`
	SignNonce   string `json:"singNonce"`
	Timestamp   string `json:"timestamp"`
}

func InitializeAuthorization(c config.Config, decls []auth.Resource) error {
	/*
		if auth, err := c.String("auth"); err != nil || auth == "none" {
			noAuth = true
			return nil
		}
		return client.RegisterSubjectDeclarations(decls)
	*/
	return nil
}

func Authenticate(opts interface{}) error {
	/*
		if !noAuth {
			return client.Authenticate(opts)
		}
	*/
	return nil
}

func Authorize(ctx echo.Context) error {
	/*
		if !noAuth {
			return client.Authorize(ctx)
		}
	*/
	return nil
}
