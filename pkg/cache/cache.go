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

type Manager struct {
	caches map[string]*Cache
}

type Cache interface {
	Add(key interface{}, data interface{})
	Delete(key interface{})
	Exists(key interface{}) bool
	Value(key interface{}) interface{}
}

type CacheItem interface {
	Key() interface{}
	Data() interface{}
	KeepAlive()
}

func NewManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]*Cache),
	}
}

func (c *CacheManager) NewCache(cacheName string) *Cache {
	return nil
}
