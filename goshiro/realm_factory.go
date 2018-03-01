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
	"strings"

	"github.com/cloustone/sentel/goshiro/shiro"
	"github.com/golang/glog"
)

type RealmFactory struct {
	env    shiro.Environment
	realms []shiro.Realm
}

func NewRealmFactory(env shiro.Environment) *RealmFactory {
	realms := []shiro.Realm{}
	realmString, err := env.GetValue("realms")
	if err == nil {
		realmNames := strings.Split(realmString.(string), ",")
		for _, realmName := range realmNames {
			realm, err := shiro.NewRealm(env, realmName)
			if err != nil {
				glog.Errorf("'%s' realm laod failed, %s", realmName, err.Error())
				continue
			}
			realms = append(realms, realm)
		}
	}
	return &RealmFactory{env: env, realms: realms}
}

func (r *RealmFactory) GetRealms() []shiro.Realm { return r.realms }

func (r *RealmFactory) AddRealm(realm shiro.Realm) {
	r.realms = append(r.realms, realm)
}

func (r *RealmFactory) GetRealm(realmName string) shiro.Realm {
	for _, realm := range r.realms {
		if realm.GetName() == realmName {
			return realm
		}
	}
	return nil
}
