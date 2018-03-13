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

package goshiro

import "github.com/cloustone/sentel/pkg/message"

const (
	NameOfPrincipalTopic  = "goshiro-security-manager-principal"
	NameOfPermissionTopic = "goshiro-security-manager-permission"
	NameOfRoleTopic       = "goshiro-security-manager-role"
)

// Principal
type PrincipalTopic struct {
	MgrId  string `json:"mgrid"`
	Action string `json:"action"`
}

func (p *PrincipalTopic) Topic() string        { return NameOfPrincipalTopic }
func (p *PrincipalTopic) SetTopic(name string) {}
func (p *PrincipalTopic) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(p, opt)
}
func (p *PrincipalTopic) Deserialize(buf []byte, opt message.SerializeOption) error {
	return message.Deserialize(buf, opt, p)
}

// Permission
type PermissionTopic struct {
	MgrId  string `json:"mgrid"`
	Action string `json:"action"`
}

func (p *PermissionTopic) Topic() string        { return NameOfPermissionTopic }
func (p *PermissionTopic) SetTopic(name string) {}
func (p *PermissionTopic) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(p, opt)
}
func (p *PermissionTopic) Deserialize(buf []byte, opt message.SerializeOption) error {
	return message.Deserialize(buf, opt, p)
}

// Role
type RoleTopic struct {
	MgrId  string `json:"mgrid"`
	Action string `json:"action"`
}

func (p *RoleTopic) Topic() string        { return NameOfRoleTopic }
func (p *RoleTopic) SetTopic(name string) {}
func (p *RoleTopic) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(p, opt)
}
func (p *RoleTopic) Deserialize(buf []byte, opt message.SerializeOption) error {
	return message.Deserialize(buf, opt, p)
}
