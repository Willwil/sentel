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

package com

import (
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
	AddConfig(options map[string]map[string]string) Config
	AddConfigSection(setion string, options map[string]string) Config
	AddConfigFile(fileName string) (Config, error)
}

type config struct {
	sections map[string]map[string]string
}

func NewConfig() Config {
	return &config{
		sections: make(map[string]map[string]string),
	}
}

// NewConfigWithFile load configurations from files
func NewConfigWithFile(fileName string, moreFiles ...string) (Config, error) {
	c := &config{
		sections: make(map[string]map[string]string),
	}
	if cfg, err := goconfig.LoadConfigFile(fileName, moreFiles...); err == nil {
		sections := cfg.GetSectionList()
		for _, name := range sections {
			// create section if it doesn't exist
			if _, ok := c.sections[name]; !ok {
				c.sections[name] = make(map[string]string)
			}
			items, err := cfg.GetSection(name)
			if err == nil {
				for key, val := range items {
					c.sections[name][key] = val
				}
			}
		}
	} else {
		return nil, err
	}
	return c, nil
}

func (c *config) checkItemExist(section string, key string) error {
	if _, found := c.sections[section]; !found {
		return fmt.Errorf("Invalid configuration, section '%s' doesn't exist'", section)
	}
	if _, found := c.sections[section][key]; !found {
		return fmt.Errorf("Invalid configuration, there are no item '%s' in section '%s'", key, section)
	}
	return nil
}

// Bool return bool value for key
func (c *config) Bool(section string, key string) (bool, error) {
	if err := c.checkItemExist(section, key); err != nil {
		return false, err
	}
	val := c.sections[section][key]
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
	if err := c.checkItemExist(section, key); err != nil {
		return -1, err
	}
	val := c.sections[section][key]
	return strconv.Atoi(val)
}

// String return string valu for key
func (c *config) String(section string, key string) (string, error) {
	if err := c.checkItemExist(section, key); err != nil {
		return "", err
	}
	return c.sections[section][key], nil
}

func (c *config) MustBool(section string, key string) bool {
	if err := c.checkItemExist(section, key); err != nil {
		glog.Fatal(err)
	}
	val := c.sections[section][key]
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
	if err := c.checkItemExist(section, key); err != nil {
		glog.Fatal(err)
	}
	val := c.sections[section][key]
	n, err := strconv.Atoi(val)
	if err != nil {
		glog.Fatalf("Invalid configuration for service '%s':'%s'", section, key)
	}
	return n
}

func (c *config) MustString(section string, key string) string {
	if err := c.checkItemExist(section, key); err != nil {
		glog.Fatal(err)
	}
	return c.sections[section][key]
}

func (c *config) SetValue(section string, key string, valu string) {
}

func (c *config) AddConfig(options map[string]map[string]string) Config {
	for section, values := range options {
		if _, found := c.sections[section]; !found {
			c.sections[section] = make(map[string]string)
		}
		for key, val := range values {
			c.sections[section][key] = val
		}
	}
	return c
}

// Config global functions
func (c *config) AddConfigSection(sectionName string, items map[string]string) Config {
	if _, found := c.sections[sectionName]; found {
		for key, val := range c.sections[sectionName] {
			if c.sections[sectionName][key] != "" {
				glog.Infof("Config item(%s) will overide existed item:%s", key, val)
			}
			c.sections[sectionName][key] = val
		}
	} else {
		c.sections[sectionName] = make(map[string]string)
		for key, val := range items {
			c.sections[sectionName][key] = val
		}
	}
	return c
}

func (c *config) AddConfigFile(fileName string) (Config, error) {
	if cfg, err := goconfig.LoadConfigFile(fileName); err == nil {
		sections := cfg.GetSectionList()
		for _, name := range sections {
			// create section if it doesn't exist
			if _, ok := c.sections[name]; !ok {
				c.sections[name] = make(map[string]string)
			}
			items, err := cfg.GetSection(name)
			if err == nil {
				for key, val := range items {
					c.sections[name][key] = val
				}
			}
		}
		return c, nil
	} else {
		return nil, err
	}
}
