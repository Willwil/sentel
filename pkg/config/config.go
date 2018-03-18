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

package config

import (
	"fmt"

	"github.com/Unknwon/goconfig"
	"github.com/golang/glog"
)

type M map[string]map[string]interface{}

// Config interface
type Config interface {
	// using primary section as default
	Bool(key string) (bool, error)
	Int(key string) (int, error)
	String(key string) (string, error)
	Value(key string) (interface{}, error)
	MustBool(key string) bool
	MustInt(key string) int
	MustString(key string) string
	MustValue(key string) interface{}
	SetValue(key string, val string)
	AddConfigItem(key string, value interface{})

	// WithSection
	BoolWithSection(section string, key string) (bool, error)
	IntWithSection(section string, key string) (int, error)
	StringWithSection(section string, key string) (string, error)
	MustBoolWithSection(section string, key string) bool
	MustIntWithSection(section string, key string) int
	MustStringWithSection(section string, key string) string
	SetValueWithSection(section string, key string, val string)
	AddConfigItemWithSection(section string, key string, value interface{})

	AddConfig(options map[string]map[string]interface{}) Config
	AddConfigSection(setion string, options map[string]string) Config
	AddConfigFile(fileName string) (Config, error)
}

type config struct {
	primary  string
	sections map[string]map[string]interface{}
}

func New(primary string) Config {
	c := &config{
		primary:  primary,
		sections: make(map[string]map[string]interface{}),
	}
	c.sections[primary] = make(map[string]interface{})
	return c
}

// NewConfigWithFile load configurations from files
func NewWithFile(primary string, fileName string) (Config, error) {
	c := New(primary).(*config)
	if cfg, err := goconfig.LoadConfigFile(fileName); err == nil {
		sections := cfg.GetSectionList()
		for _, name := range sections {
			// create section if it doesn't exist
			if _, ok := c.sections[name]; !ok {
				c.sections[name] = make(map[string]interface{})
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
		return fmt.Errorf("section '%s' no exist'", section)
	}
	if _, found := c.sections[section][key]; !found {
		return fmt.Errorf("no item '%s' in section '%s'", key, section)
	}
	return nil
}

func (c *config) AddConfigItemWithSection(section string, key string, value interface{}) {
	if _, found := c.sections[section]; !found {
		c.sections[section] = make(map[string]interface{})
	}
	if key != "" && value != nil {
		c.sections[section][key] = value
	}
}

func (c *config) AddConfigItem(key string, value interface{}) {
	c.AddConfigItemWithSection(c.primary, key, value)
}

// Bool return bool value for key
func (c *config) BoolWithSection(section string, key string) (bool, error) {
	if err := c.checkItemExist(section, key); err != nil {
		return false, err
	}
	if val, ok := c.sections[section][key].(bool); ok {
		return val, nil
	}
	return false, fmt.Errorf("invalid seciont '%s' with key '%s'", section, key)
}

// Int return int value for key
func (c *config) IntWithSection(section string, key string) (int, error) {
	if err := c.checkItemExist(section, key); err != nil {
		return -1, err
	}
	val := c.sections[section][key]
	if _, ok := val.(int); !ok {
		return -1, fmt.Errorf("invalid value for key '%s'", key)
	}
	return val.(int), nil
}

// String return string valu for key
func (c *config) StringWithSection(section string, key string) (string, error) {
	if err := c.checkItemExist(section, key); err == nil {
		if val, ok := c.sections[section][key].(string); ok {
			return val, nil
		}
	}
	return "", fmt.Errorf("invalid value for key '%s'", key)
}

// Value return element value for key
func (c *config) ValueWithSection(section string, key string) (interface{}, error) {
	if err := c.checkItemExist(section, key); err == nil {
		return c.sections[section][key], nil
	}
	return nil, fmt.Errorf("invalid value for key '%s'", key)
}

func (c *config) MustBoolWithSection(section string, key string) bool {
	val, err := c.BoolWithSection(section, key)
	if err != nil {
		glog.Fatal(err)
	}
	return val
}
func (c *config) MustIntWithSection(section string, key string) int {
	val, err := c.IntWithSection(section, key)
	if err != nil {
		glog.Fatal(err)
	}
	return val
}

func (c *config) MustStringWithSection(section string, key string) string {
	val, err := c.StringWithSection(section, key)
	if err != nil {
		glog.Fatal(err)
	}
	return val
}

func (c *config) MustValueWithSection(section string, key string) interface{} {
	val, err := c.ValueWithSection(section, key)
	if err != nil {
		glog.Fatal(err)
	}
	return val
}

func (c *config) SetValueWithSection(section string, key string, valu string) {
}

func (c *config) Bool(key string) (bool, error)         { return c.BoolWithSection(c.primary, key) }
func (c *config) Int(key string) (int, error)           { return c.IntWithSection(c.primary, key) }
func (c *config) String(key string) (string, error)     { return c.StringWithSection(c.primary, key) }
func (c *config) Value(key string) (interface{}, error) { return c.ValueWithSection(c.primary, key) }
func (c *config) MustBool(key string) bool              { return c.MustBoolWithSection(c.primary, key) }
func (c *config) MustInt(key string) int                { return c.MustIntWithSection(c.primary, key) }
func (c *config) MustString(key string) string          { return c.MustStringWithSection(c.primary, key) }
func (c *config) MustValue(key string) interface{}      { return c.MustValueWithSection(c.primary, key) }
func (c *config) SetValue(key string, val string)       { c.SetValueWithSection(c.primary, key, val) }

func (c *config) AddConfig(options map[string]map[string]interface{}) Config {
	for section, values := range options {
		if _, found := c.sections[section]; !found {
			c.sections[section] = make(map[string]interface{})
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
				glog.Infof("config item '%s' overwriting occured", key, val)
			}
			c.sections[sectionName][key] = val
		}
	} else {
		c.sections[sectionName] = make(map[string]interface{})
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
				c.sections[name] = make(map[string]interface{})
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
