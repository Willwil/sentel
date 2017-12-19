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

package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Unknwon/goconfig"
	"github.com/golang/glog"
)

// Config interface
type Config interface {
	Bool(section string, key string) (bool, error)
	Int(section string, key string) (int, error)
	String(section string, key string) (string, error)
	MustBool(section string, key string) bool
	MustInt(section string, key string) int
	MustString(section string, key string) string
	SetValue(section string, key string, val string)
	AddConfigs(options map[string]map[string]string)
}

type configSection struct {
	items map[string]string
}

type config struct{}

var allConfigSections map[string]*configSection = make(map[string]*configSection)

var (
	ErrorInvalidConfiguration = errors.New("Invalid configuration")
)

// config implementations

func checkItemExist(section string, key string) bool {
	if _, found := allConfigSections[section]; !found {
		return false
	}
	if _, found := allConfigSections[section].items[key]; !found {
		return false
	}
	return true
}

// Bool return bool value for key
func (c *config) Bool(section string, key string) (bool, error) {
	if !checkItemExist(section, key) {
		return false, ErrorInvalidConfiguration
	}
	val := allConfigSections[section].items[key]
	switch val {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, fmt.Errorf("Invalid configuration item for service '%s' item '%s'", section, key)
}

// Int return int value for key
func (c *config) Int(section string, key string) (int, error) {
	if !checkItemExist(section, key) {
		return -1, ErrorInvalidConfiguration
	}
	val := allConfigSections[section].items[key]
	return strconv.Atoi(val)
}

// String return string valu for key
func (c *config) String(section string, key string) (string, error) {
	if !checkItemExist(section, key) {
		return "", ErrorInvalidConfiguration
	}

	return allConfigSections[section].items[key], nil
}

func (c *config) MustBool(section string, key string) bool {
	if !checkItemExist(section, key) {
		glog.Fatalf("Invalid configuration item for service '%s':'%s'", section, key)
	}
	val := allConfigSections[section].items[key]
	switch val {
	case "true":
		return true
	case "false":
		return false
	}
	glog.Fatalf("Invalid configuration item for service '%s':'%s'", section, key)
	return false
}
func (c *config) MustInt(section string, key string) int {
	if !checkItemExist(section, key) {
		glog.Fatalf("Invalid configuration for service '%s':'%s'", section, key)
	}
	val := allConfigSections[section].items[key]
	n, err := strconv.Atoi(val)
	if err != nil {
		glog.Fatalf("Invalid configuration for service '%s':'%s'", section, key)
	}
	return n
}

func (c *config) MustString(section string, key string) string {
	if !checkItemExist(section, key) {
		glog.Fatalf("Invalid configuration item for service '%s':'%s'", section, key)
	}
	return allConfigSections[section].items[key]
}

func (c *config) SetValue(section string, key string, valu string) {
}

func (c *config) AddConfigs(options map[string]map[string]string) {
	for section, values := range options {
		if _, found := allConfigSections[section]; !found {
			allConfigSections[section] = &configSection{items: make(map[string]string)}
		}
		items := allConfigSections[section].items
		for key, val := range values {
			items[key] = val
		}
	}
}

// NewWithConfigFile load configurations from files
func NewConfigWithFile(fileName string, moreFiles ...string) (Config, error) {
	// load all config sections in _allConfigSections, get section and item to overide
	cfg, err := goconfig.LoadConfigFile(fileName, moreFiles...)
	if err == nil {
		sections := cfg.GetSectionList()
		for _, name := range sections {
			// create section if it doesn't exist
			if _, ok := allConfigSections[name]; !ok {
				allConfigSections[name] = &configSection{items: make(map[string]string)}
			}
			items, err := cfg.GetSection(name)
			if err == nil {
				for key, val := range items {
					allConfigSections[name].items[key] = val
				}
			}
		}
	}
	return &config{}, nil
}

// Config global functions
func RegisterConfig(sectionName string, items map[string]string) {
	if allConfigSections[sectionName] != nil { // section already exist
		section := allConfigSections[sectionName]
		for key, val := range items {
			if section.items[key] != "" {
				glog.Infof("Config item(%s) will overide existed item:%s", key, section.items[key])
			}
			section.items[key] = val
		}
	} else {
		section := new(configSection)
		section.items = make(map[string]string)
		for key, val := range items {
			section.items[key] = val
		}
		allConfigSections[sectionName] = section
	}
}

// RegisterConfigGropu reigster all group's subconfigurations
func RegisterConfigGroup(configs map[string]map[string]string) {
	for group, values := range configs {
		RegisterConfig(group, values)
	}
}
