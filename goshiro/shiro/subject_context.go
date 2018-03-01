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

type SubjectContext interface {
	GetSecurityManager() SecurityManager
	SetSecurityManager(SecurityManager)
	SetSessionId(sessionId Serializable)
	GetSessionId() Serializable
	GetSubject() Subject
	SetSubject(subject Subject)
	GetPrincipals() PrincipalCollection
	SetPrincipals(principals PrincipalCollection)
	GetSession() Session
	SetSession(Session)
	IsSessionCreationEnabled() bool
	IsAuthenticated() bool
	SetAuthenticated(bool)
	GetAuthenticaitonToken() AuthenticationToken
	SetAuthenticationToken(AuthenticationToken)
	SetAuthenticationInfo(AuthenticationInfo)
	GetAuthenticationInfo() AuthenticationInfo
	GetHost() string
	SetHost(string)
}

func NewSubjectContext() SubjectContext {
	return &defaultSubjectContext{
		values: make(map[string]interface{}),
	}
}

type defaultSubjectContext struct {
	values map[string]interface{}
}

const (
	SECURITY_MANAGER         = "securityManager"
	SESSION_ID               = "sessionId"
	SUBJECT                  = "subject"
	PRINCIPALS               = "principals"
	SESSION                  = "session"
	AUTHENTICATED            = "authenticated"
	AUTHENTICATION_TOKEN     = "authenticationToken"
	AUTHENTICATION_INFO      = "authenticationInfo"
	AUTHORIZATION_TOKEN      = "auorizationToken"
	HOST                     = "host"
	SESSION_CREATION_ENABLED = "sessionCreationEnabled"
)

func (p *defaultSubjectContext) GetSecurityManager() SecurityManager {
	if mgr, found := p.values[SECURITY_MANAGER]; found {
		return mgr.(SecurityManager)
	}
	return nil
}

func (p *defaultSubjectContext) SetSecurityManager(mgr SecurityManager) {
	p.values[SECURITY_MANAGER] = mgr
}

func (p *defaultSubjectContext) SetSessionId(sessionId Serializable) {
	p.values[SESSION_ID] = sessionId
}

func (p *defaultSubjectContext) GetSessionId() Serializable {
	if item, found := p.values[SESSION_ID]; found {
		return item.(Serializable)
	}
	return nil
}

func (p *defaultSubjectContext) GetSubject() Subject {
	if item, found := p.values[SUBJECT]; found {
		return item.(Subject)
	}
	return nil
}

func (p *defaultSubjectContext) SetSubject(subject Subject) {
	p.values[SUBJECT] = subject
}

func (p *defaultSubjectContext) GetPrincipals() PrincipalCollection {
	if item, found := p.values[PRINCIPALS]; found {
		return item.(PrincipalCollection)
	}
	return nil
}

func (p *defaultSubjectContext) SetPrincipals(principals PrincipalCollection) {
	p.values[PRINCIPALS] = principals
}

func (p *defaultSubjectContext) GetSession() Session {
	if item, found := p.values[SESSION]; found {
		return item.(Session)
	}
	return nil
}

func (p *defaultSubjectContext) SetSession(session Session) {
	p.values[SESSION] = session
}

func (p *defaultSubjectContext) IsSessionCreationEnabled() bool {
	if item, found := p.values[SESSION_CREATION_ENABLED]; found {
		return item.(bool)
	}
	return false
}

func (p *defaultSubjectContext) IsAuthenticated() bool {
	if item, found := p.values[AUTHENTICATED]; found {
		return item.(bool)
	}
	return false
}

func (p *defaultSubjectContext) SetAuthenticated(created bool) {
	p.values[AUTHENTICATED] = created

}

func (p *defaultSubjectContext) GetAuthenticaitonToken() AuthenticationToken {
	if item, found := p.values[SECURITY_MANAGER]; found {
		return item.(AuthenticationToken)
	}
	return nil
}

func (p *defaultSubjectContext) SetAuthenticationToken(token AuthenticationToken) {
	p.values[AUTHENTICATION_TOKEN] = token
}

func (p *defaultSubjectContext) SetAuthenticationInfo(info AuthenticationInfo) {
	p.values[AUTHENTICATION_INFO] = info
}

func (p *defaultSubjectContext) GetAuthenticationInfo() AuthenticationInfo {
	if item, found := p.values[AUTHENTICATION_INFO]; found {
		return item.(AuthenticationInfo)
	}
	return nil
}

func (p *defaultSubjectContext) GetHost() string {
	if item, found := p.values[HOST]; found {
		return item.(string)
	}
	return ""
}

func (p *defaultSubjectContext) SetHost(host string) {
	p.values[HOST] = host
}
