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

package shiro

type Request interface{}

type RequestContext interface {
	// Bind binds the request body into provided type `i`. The default binder
	// does it based on Content-Type header.
	Bind(i interface{}) error
	// Param returns path parameter by name.
	Param(name string) string
	// ParamNames returns path parameter names.
	ParamNames() []string
	// SetParamNames sets path parameter names.
	SetParamNames(names ...string)
	// ParamValues returns path parameter values.
	ParamValues() []string
	// SetParamValues sets path parameter values.
	SetParamValues(values ...string)
	// QueryParam returns the query param for the provided name.
	QueryParam(name string) string
}
