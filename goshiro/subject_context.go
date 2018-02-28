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

type SubjectContext interface {
	GetSecurityManager() securityManager
	SetSecurityManager(securityManager)
	SetSessionId(sessionId Serializable)
	GetSessionId() Serializable
	GetSubject() Subject
	SetSubject(subject Subject)
	GetPrincipals() PrincipalCollection
	SetPrincipals(principals PrincipalCollection)
	GetSession() Session
	SetSession(Session)
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
	return nil
}
