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

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
)

type Resource struct {
	Name string `json:"name" bson:"Name"`
}

type ResourceManager interface {
	GetAllResources() []Resource
	GetResource(name string) (Resource, error)
	AddResource(Resource)
	RemoveResource(Resource)
}

func NewResourceManager(c config.Config) (ResourceManager, error) {
	return &simpleResourceManager{
		config: c,
	}, nil
}

type simpleResourceManager struct {
	config config.Config
}

func (r *simpleResourceManager) GetAllResources() []Resource {
	return []Resource{}
}
func (r *simpleResourceManager) GetResource(name string) (Resource, error) {
	return Resource{}, errors.New("not implemented")
}
func (r *simpleResourceManager) AddResource(Resource)    {}
func (r *simpleResourceManager) RemoveResource(Resource) {}
