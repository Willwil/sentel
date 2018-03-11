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

type Resource struct {
	Name  string `json:"name" bson:"Name"`
	Roles string `json:"roles" bson:"Roles"`
}

// ResourceManager manage what resource we have and resource's role
type ResourceManager interface {
	GetAllResources() []Resource
	GetResource(name string) (Resource, error)
	AddResource(Resource)
	RemoveResource(Resource)
	GetResourceRoleNames(resourceName string) []string
}

func NewResourceManager(adaptor Adaptor) ResourceManager {
	return &defaultResourceManager{
		adaptor: adaptor,
	}
}
