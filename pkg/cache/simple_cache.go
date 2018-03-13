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

package cache

import (
	"errors"
	"time"

	"github.com/muesli/cache2go"
)

const (
	lifeSpan = time.Duration(time.Second * 5)
)

type simpleCache struct {
	cacheTable *cache2go.CacheTable
}

func newSimpleCache(name string) Cache {
	return &simpleCache{
		cacheTable: cache2go.Cache(name),
	}
}

func (s *simpleCache) Add(key interface{}, data interface{}) {
	s.cacheTable.Add(key, lifeSpan, data)
}

func (s *simpleCache) Delete(key interface{}) {
	s.cacheTable.Delete(key)
}
func (s *simpleCache) Exists(key interface{}) bool {
	return s.cacheTable.Exists(key)
}

func (s *simpleCache) Value(key interface{}) (interface{}, error) {
	if cacheItem, err := s.cacheTable.Value(key); err == nil {
		return cacheItem.Data(), nil
	}
	return nil, errors.New("no valid cached value")
}

func (s *simpleCache) KeepAlive(key interface{}) {
	if cacheItem, err := s.cacheTable.Value(key); err == nil {
		cacheItem.KeepAlive()
	}
}

func (s *simpleCache) Flush() {
	s.cacheTable.Flush()
}
