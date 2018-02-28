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
	"strings"

	"github.com/golang/glog"
)

type Realm interface {
	GetName() string
	Supports(token AuthenticationToken) bool
	GetAuthenticationInfo(token AuthenticationToken) (AuthenticationInfo, error)
	GetAuthorizationInfo(principals PrincipalCollection) (AuthorizationInfo, error)
	SetCacheEnable(bool)
}

type RealmFactory interface {
	AddRealm(r Realm)
	GetRealms() []Realm
	GetRealm(realmName string) Realm
}

func NewRealmFactory(env Environment) RealmFactory {
	realms := []Realm{}
	realmString, err := env.GetValue("realms")
	if err == nil {
		realmNames := strings.Split(realmString.(string), ",")
		for _, realmName := range realmNames {
			realm, err := newRealm(env, realmName)
			if err != nil {
				glog.Errorf("'%s' realm laod failed, %s", realmName, err.Error())
				continue
			}
			realms = append(realms, realm)
		}
	}
	return &realmFactory{env: env, realms: realms}
}

type realmFactory struct {
	env    Environment
	realms []Realm
}

func (r *realmFactory) GetRealms() []Realm { return r.realms }

func (r *realmFactory) AddRealm(realm Realm) {
	r.realms = append(r.realms, realm)
}

func (r *realmFactory) GetRealm(realmName string) Realm {
	for _, realm := range r.realms {
		if realm.GetName() == realmName {
			return realm
		}
	}
	return nil
}

func newRealm(env Environment, realmName string) (Realm, error) {
	return nil, errors.New("not implemented")
}
