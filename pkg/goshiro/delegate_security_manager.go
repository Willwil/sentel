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

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	"github.com/cloustone/sentel/pkg/message"
	uuid "github.com/satori/go.uuid"
)

// delegateSecurityManager is delegation of security manager to support kafka
type delegateSecurityManager struct {
	consumer message.Consumer      // kafka consumer
	delegate shiro.SecurityManager // backend security manager
	mgrid    string                // Identifier of current security manager
}

// newDelegateSecurityManager create security manager supporting distributed security managment
func newDelegateSecurityManager(c config.Config, realm ...shiro.Realm) (shiro.SecurityManager, error) {
	delegate, err := shiro.NewSecurityManager(c, realm...)
	if err != nil {
		return nil, err
	}
	consumer, err := message.NewConsumer(c, "")
	if err != nil {
		return nil, err
	}
	securityMgr := &delegateSecurityManager{
		delegate: delegate,
		consumer: consumer,
		mgrid:    uuid.NewV4().String(),
	}
	consumer.SetMessageFactory(securityMgr)
	err1 := consumer.Subscribe(NameOfPrincipalTopic, messageHandlerFunc, securityMgr)
	err2 := consumer.Subscribe(NameOfPermissionTopic, messageHandlerFunc, securityMgr)
	err3 := consumer.Subscribe(NameOfRoleTopic, messageHandlerFunc, securityMgr)
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, errors.New("consumer subscribe failed")
	}
	return securityMgr, nil

}

// messageHandlerFunc is callback to message module when message is received
func messageHandlerFunc(msg message.Message, ctx interface{}) {
	securityMgr := ctx.(*delegateSecurityManager)
	switch msg.Topic() {
	case NameOfPrincipalTopic:
		securityMgr.handlePrincipalTopic(msg.(*PrincipalTopic))
	case NameOfPermissionTopic:
		securityMgr.handlePermissionTopic(msg.(*PermissionTopic))
	case NameOfRoleTopic:
		securityMgr.handleRoleTopic(msg.(*RoleTopic))
	}
}

// CreateMessage is MessageFacotory interface method to create customized topic
func (p *delegateSecurityManager) CreateMessage(topicName string) message.Message {
	switch topicName {
	case NameOfPrincipalTopic:
		return &PrincipalTopic{}
	case NameOfPermissionTopic:
		return &PermissionTopic{}
	case NameOfRoleTopic:
		return &RoleTopic{}
	default:
		return nil
	}
}

func (p *delegateSecurityManager) handlePrincipalTopic(topic *PrincipalTopic) error {
	return errors.New("not implemented")
}

func (p *delegateSecurityManager) handlePermissionTopic(topic *PermissionTopic) error {
	return errors.New("not implemented")
}

func (p *delegateSecurityManager) handleRoleTopic(topic *RoleTopic) error {
	return errors.New("not implemented")
}

func (p *delegateSecurityManager) Load() {
	p.delegate.Load()
}

func (p *delegateSecurityManager) AddRealm(realm ...shiro.Realm) {
	p.delegate.AddRealm(realm...)
}

func (p *delegateSecurityManager) GetRealm(realmName string) shiro.Realm {
	return p.delegate.GetRealm(realmName)
}

func (p *delegateSecurityManager) Login(token shiro.AuthenticationToken) (shiro.Principal, error) {
	return p.delegate.Login(token)
}

func (p *delegateSecurityManager) RemovePrincipal(principal shiro.Principal) {
	p.delegate.RemovePrincipal(principal)
}

func (p *delegateSecurityManager) GetPrincipalsPermissions(principals shiro.PrincipalCollection) []shiro.Permission {
	return p.delegate.GetPrincipalsPermissions(principals)
}

func (p *delegateSecurityManager) AddPrincipalPermissions(principal shiro.Principal, permissions []shiro.Permission) {
	p.delegate.AddPrincipalPermissions(principal, permissions)
}

func (p *delegateSecurityManager) RemovePrincipalPermissions(principal shiro.Principal, permissions []shiro.Permission) {
	p.delegate.RemovePrincipalPermissions(principal, permissions)
}

func (p *delegateSecurityManager) Authorize(principal shiro.Principal, resource, action string) error {
	return p.delegate.Authorize(principal, resource, action)
}

func (p *delegateSecurityManager) AuthorizeWithToken(token shiro.AuthenticationToken, resource, action string) error {
	return p.delegate.AuthorizeWithToken(token, resource, action)
}

func (p *delegateSecurityManager) AddRole(r shiro.Role) {
	p.delegate.AddRole(r)
}

func (p *delegateSecurityManager) RemoveRole(roleName string) {
	p.delegate.RemoveRole(roleName)
}

func (p *delegateSecurityManager) GetRole(roleName string) (shiro.Role, error) {
	return p.delegate.GetRole(roleName)
}

func (p *delegateSecurityManager) AddRolePermissions(roleName string, permissions []shiro.Permission) {
	p.delegate.AddRolePermissions(roleName, permissions)
}

func (p *delegateSecurityManager) RemoveRolePermissions(roleName string, permissions []shiro.Permission) {
	p.delegate.RemoveRolePermissions(roleName, permissions)
}

func (p *delegateSecurityManager) RemoveRoleAllPermissions(roleName string) {
	p.delegate.RemoveRoleAllPermissions(roleName)
}

// GetRolePermission return specfied role's all permissions
func (p *delegateSecurityManager) GetRolePermissions(roleName string) []shiro.Permission {
	return p.delegate.GetRolePermissions(roleName)
}
