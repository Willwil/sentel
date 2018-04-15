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

package registry

import (
	"errors"
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

const (
	dbNameDevices      = "devices"
	dbNameProducts     = "products"
	dbNameTenants      = "tenants"
	dbNameRules        = "rules"
	dbNameRunlogs      = "runlogs"
	dbNameTopicFlavors = "topicflavors"
)

const PageSize = 3

var (
	ErrorNotFound = errors.New("not found")
)

// Registry is wraper of mongo database about for iot object
type Registry struct {
	config  config.Config
	session *mgo.Session
	db      *mgo.Database
}

// InitializeRegistry try to connect with background database
// to confirm wether it is normal
func Initialize(c config.Config) error {
	hosts := c.MustString("mongo")
	glog.Infof("Initializing registry:%s...", hosts)
	session, err := mgo.Dial(hosts)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)

	ensureDevicesIndex(session)
	ensureProductsIndex(session)
	ensureTenantsIndex(session)
	ensureRulesIndex(session)

	session.Close()
	return nil
}

// NewRegistry create registry instance
func New(c config.Config) (*Registry, error) {
	hosts := c.MustString("mongo")
	session, err := mgo.DialWithTimeout(hosts, 5*time.Second)
	if err != nil {
		glog.Infof("registry '%s':%s", hosts, err.Error())
		return nil, err
	}
	return &Registry{session: session, db: session.DB("registry"), config: c}, nil
}

func ensureDevicesIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("registry").C(dbNameDevices)
	index := mgo.Index{
		Key:        []string{"DeviceId"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		glog.Errorf("Create dev INDEX error! %s\n", err)
	}
}

func ensureProductsIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("registry").C(dbNameProducts)
	index := mgo.Index{
		Key:        []string{"ProductId"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		glog.Errorf("Create dev INDEX error! %s\n", err)
	}
}

func ensureTenantsIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("registry").C(dbNameTenants)
	index := mgo.Index{
		Key:        []string{"TenantId"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		glog.Errorf("Create dev INDEX error! %s\n", err)
	}
}

func ensureRulesIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("registry").C(dbNameRules)
	index := mgo.Index{
		Key:        []string{"ProductId", "RuleName"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		glog.Errorf("Create dev INDEX error! %s\n", err)
	}
}

// Release release registry rources and disconnect with background database
func (r *Registry) Close() {
	r.session.Close()
}

// Tenant

// CheckTenantNamveAvailable return true if name is available
func (r *Registry) CheckTenantNameAvailable(id string) error {
	c := r.db.C(dbNameTenants)
	return c.Find(bson.M{"TenantId": id}).One(nil)
}

// RegisterTenant insert a tenant into registry
func (r *Registry) RegisterTenant(t *Tenant) error {
	c := r.db.C(dbNameTenants)
	if err := c.Find(bson.M{"TenantId": t.TenantId}).One(nil); err == nil {
		return fmt.Errorf("Tenant %s already exist", t.TenantId)
	}
	return c.Insert(t)
}

func (r *Registry) DeleteTenant(id string) error {
	c := r.db.C(dbNameTenants)
	return c.Remove(bson.M{"TenantId": id})
}

func (r *Registry) GetTenant(id string) (*Tenant, error) {
	c := r.db.C(dbNameTenants)
	t := Tenant{}
	err := c.Find(bson.M{"TenantId": id}).One(&t)
	return &t, err
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetAllTenants() ([]Tenant, error) {
	c := r.db.C(dbNameTenants)
	tenants := []Tenant{}
	err := c.Find(nil).Sort("TenantId").All(&tenants)
	return tenants, err
}

func (r *Registry) UpdateTenant(id string, t *Tenant) error {
	c := r.db.C(dbNameTenants)
	return c.Update(bson.M{"TenantId": id}, t)
}

// Product
// CheckProductNameAvailable check wethere product name is available
func (r *Registry) CheckProductNameAvailable(p *Product) bool {
	c := r.db.C(dbNameProducts)
	err := c.Find(bson.M{"TenantId": p.TenantId, "ProductId": p.ProductId}).One(nil)
	return err != nil
}

// RegisterProduct register a product into registry
func (r *Registry) RegisterProduct(p *Product) error {
	c := r.db.C(dbNameProducts)
	product := &Product{}
	if err := c.Find(bson.M{"TenantId": p.TenantId, "ProductId": p.ProductId}).One(product); err == nil {
		return fmt.Errorf("product %s already exist", p.ProductId)
	}
	return c.Insert(p)
}

// DeleteProduct delete a product from registry
func (r *Registry) DeleteProduct(productId string) error {
	c := r.db.C(dbNameProducts)
	return c.Remove(bson.M{"ProductId": productId})
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProduct(productId string) (*Product, error) {
	c := r.db.C(dbNameProducts)
	product := &Product{}
	err := c.Find(bson.M{"ProductId": productId}).One(product)
	return product, err
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProducts(tenantId string) ([]Product, error) {
	c := r.db.C(dbNameProducts)
	products := []Product{}

	err := c.Find(bson.M{"TenantId": tenantId}).Sort("ProductId").All(&products)
	return products, err
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProductsWithCondition(tenantId string, conditions map[string]string) ([]Product, error) {
	c := r.db.C(dbNameProducts)
	products := []Product{}
	query := make(bson.M)
	for k, v := range conditions {
		query[k] = v
	}
	query["TenantId"] = tenantId

	err := c.Find(query).Sort("ProductId").All(&products)
	return products, err
}

// GetProductDevices get product's device list
func (r *Registry) GetProductDevices(productId string) ([]Device, error) {
	c := r.db.C(dbNameDevices)

	devices := []Device{}
	err := c.Find(bson.M{"ProductId": productId}).Sort("DeviceId").All(&devices)
	return devices, err
}

// UpdateProduct update product detail information in registry
func (r *Registry) UpdateProduct(p *Product) error {
	c := r.db.C(dbNameProducts)
	return c.Update(bson.M{"ProductId": p.ProductId}, p)
}

// Device

// RegisterDevice add a new device into registry
func (r *Registry) RegisterDevice(dev *Device) error {
	c := r.db.C(dbNameDevices)
	device := Device{}
	if err := c.Find(bson.M{"DeviceId": dev.DeviceId}).One(&device); err == nil { // found existed device
		return fmt.Errorf("device %s already exist", dev.DeviceId)
	}
	return c.Insert(dev)
}

// GetDevice retrieve a device information from registry/
func (r *Registry) GetDevice(productId string, deviceId string) (*Device, error) {
	c := r.db.C(dbNameDevices)
	device := &Device{}
	err := c.Find(bson.M{"ProductId": productId, "DeviceId": deviceId}).One(device)
	return device, err
}

// GetDeviceList retrieve a device list information from registry/
func (r *Registry) GetDeviceList(productId string) ([]Device, error) {
	c := r.db.C(dbNameDevices)
	devices := []Device{}
	err := c.Find(bson.M{"ProductId": productId}).All(devices)
	return devices, err
}

// GetDeviceByName retrieve a device information from registry/
func (r *Registry) GetDeviceByName(productId string, deviceName string) (*Device, error) {
	c := r.db.C(dbNameDevices)
	device := &Device{}
	err := c.Find(bson.M{"ProductId": productId, "DeviceName": deviceName}).One(device)
	return device, err
}

// GetDeviceByName retrieve a device information from registry/
func (r *Registry) GetDevicesByName(productId string, deviceName string) ([]Device, error) {
	c := r.db.C(dbNameDevices)
	devices := []Device{}
	err := c.Find(bson.M{"ProductId": productId, "DeviceName": deviceName}).All(&devices)
	return devices, err
}

// BulkRegisterDevice add a lot of devices into registry
func (r *Registry) BulkRegisterDevices(devices []Device) error {
	for _, device := range devices {
		if err := r.RegisterDevice(&device); err != nil {
			return err
		}
	}
	return nil
}

// BulkRegisterDevice add a lot of devices into registry
func (r *Registry) BulkGetDevices(devices []Device) ([]Device, error) {
	var err error
	rdevs := []Device{}
	devs := []Device{}
	for _, dev := range devices {
		devs, err = r.GetDevicesByName(dev.ProductId, dev.DeviceName)
		if err != nil {
			for _, dev := range devs {
				rdevs = append(rdevs, dev)
			}
		}
	}
	return rdevs, err
}

// DeleteDevice delete a device from registry
func (r *Registry) DeleteDevice(productId string, deviceId string) error {
	c := r.db.C(dbNameDevices)
	return c.Remove(bson.M{"ProductId": productId, "DeviceId": deviceId})
}

// BulkDeleteDevice delete a lot of devices from registry
func (r *Registry) BulkDeleteDevice(devices []Device) error {
	c := r.db.C(dbNameDevices)
	for _, dev := range devices {
		c.Remove(bson.M{"ProductId": dev.ProductId, "DeviceId": dev.DeviceId})
	}
	return nil
}

// UpdateDevice update device information in registry
func (r *Registry) UpdateDevice(dev *Device) error {
	c := r.db.C(dbNameDevices)
	return c.Update(bson.M{"ProductId": dev.ProductId, "DeviceId": dev.DeviceId}, dev)
}

// BulkUpdateDevice update a lot of devices in registry
func (r *Registry) BulkUpdateDevice(devices []Device) error {
	c := r.db.C(dbNameDevices)
	for _, dev := range devices {
		c.Update(bson.M{"ProductId": dev.ProductId, "DeviceId": dev.DeviceId}, dev)
	}
	return nil
}

//save device runlog.
func (r *Registry) SaveDeviceRunlog(productId string, deviceId string, deviceStatus string) error {
	log := Runlog{}
	log.ProductId = productId
	log.DeviceId = deviceId
	log.DeviceStatus = deviceStatus
	log.TimeCreated = time.Now()
	log.TimeUpdated = time.Now()
	log.IsShow = "0"
	l := r.db.C(dbNameRunlogs)
	return l.Insert(log)
}

// Duration 3 second device status, then show it.
func (r *Registry) GetShadowDevice(productId string, deviceId string) (*Runlog, error) {
	c := r.db.C(dbNameRunlogs)
	showlog := Runlog{}
	unshowlogs := []Runlog{}
	log := Runlog{}
	//find showing device status.
	err := c.Find(bson.M{"ProductId": productId, "DeviceId": deviceId, "IsShow": "1"}).One(&showlog)
	if err == nil {
		basetime := time.Now()
		//then try to find all fited from showlog to now .
		err = c.Find(bson.M{"TimeCreated": bson.M{"&gte": showlog.TimeCreated}}).Select(bson.M{"ProductId": productId, "DeviceId": deviceId, "IsShow": "0"}).All(unshowlogs)
		if err == nil {
			for _, log = range unshowlogs {
				duration := basetime.Sub(log.TimeCreated)
				s, _ := time.ParseDuration("3s")
				if duration >= s {
					log.IsShow = "1"
					log.TimeUpdated = time.Now()
					c.Update(bson.M{"ProductId": productId, "DeviceId": deviceId, "IsShow": "0", "TimeCreated": log.TimeCreated}, log)
					return &log, err
				} else {
					basetime = log.TimeCreated
				}
			}
			return &showlog, err
		} else {
			return &showlog, err
		}
	} else {
		//no found,find the first device status,update status.
		err = c.Find(bson.M{"ProductId": productId, "DeviceId": deviceId}).One(log)
		if err == nil {
			log.IsShow = "1"
			c.Update(bson.M{"ProductId": productId, "DeviceId": deviceId}, log)
		}
	}
	return &log, err
}

// Rule
// GetRulesWithStatus return all rules in registry
func (r *Registry) GetRulesWithStatus(status string) []Rule {
	rules := []Rule{}
	c := r.db.C(dbNameRules)
	c.Find(bson.M{"Status": status}).All(&rules)
	return rules
}

// RegisterRule add a new rule into registry
func (r *Registry) RegisterRule(rule *Rule) error {
	c := r.db.C(dbNameRules)
	if err := c.Find(bson.M{"RuleName": rule.RuleName, "ProductId": rule.ProductId}); err == nil {
		return fmt.Errorf("rule %s already exist", rule.RuleName)
	}
	return c.Insert(rule)
}

// GetRule retrieve a rule information from registry/
func (r *Registry) GetRule(productId string, ruleName string) (*Rule, error) {
	c := r.db.C(dbNameRules)
	rule := &Rule{}
	err := c.Find(bson.M{"RuleName": ruleName, "ProductId": productId}).One(rule)
	return rule, err
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProductRuleNames(productId string) ([]string, error) {
	c := r.db.C(dbNameRules)
	rules := []string{}
	err := c.Find(bson.M{"ProductId": productId}).Sort("RuleName").All(&rules)
	return rules, err
}

// RemoveRule delete a rule from registry
func (r *Registry) RemoveRule(productId string, ruleName string) error {
	c := r.db.C(dbNameRules)
	return c.Remove(bson.M{"RuleName": ruleName, "ProductId": productId})
}

// UpdateRule update rule information in registry
func (r *Registry) UpdateRule(rule *Rule) error {
	c := r.db.C(dbNameRules)
	return c.Update(bson.M{"RuleName": rule.RuleName, "ProductId": rule.ProductId}, rule)
}

// CreateTopicFlavor create topic flavor
func (r *Registry) CreateTopicFlavor(t *TopicFlavor) error {
	c := r.db.C(dbNameTopicFlavors)
	flavor := &TopicFlavor{}
	if err := c.Find(bson.M{"flavorName": t.FlavorName, "builtin": t.Builtin, "tenantId": t.TenantId}).One(&flavor); err != nil {
		return c.Insert(t)
	}
	return fmt.Errorf("flavro '%s' already exist", t.FlavorName)
}

// GetBuiltinTopicFlavor retrieve system builtin topic flavor
func (r *Registry) GetBuiltinTopicFlavors() []TopicFlavor {
	c := r.db.C(dbNameTopicFlavors)
	flavors := []TopicFlavor{}
	c.Find(bson.M{"Builtin": true}).All(&flavors)
	return flavors
}

// GetTenantTopicFlavor retrive tenant's topic flavor
func (r *Registry) GetTenantTopicFlavors(tenantId string) []TopicFlavor {
	c := r.db.C(dbNameTopicFlavors)
	flavors := []TopicFlavor{}
	c.Find(bson.M{"TenantId": tenantId}).All(&flavors)
	return flavors
}

// GetTenantTopicFlavor retrive tenant's topic flavor
func (r *Registry) GetTenantTopicFlavor(tenantId string, flavorName string) (TopicFlavor, error) {
	c := r.db.C(dbNameTopicFlavors)
	flavor := TopicFlavor{}
	if err := c.Find(bson.M{"TenantId": tenantId, "FlavorName": flavorName}).One(&flavor); err != nil {
		return flavor, err
	}
	return flavor, nil
}

// RemoveTopicFlavor remove a specified topic flavor from registry
func (r *Registry) RemoveTopicFlavor(t *TopicFlavor) error {
	c := r.db.C(dbNameTopicFlavors)
	return c.Remove(bson.M{"FlavorName": t.FlavorName, "Builtin": t.Builtin, "TenantId": t.TenantId})
}

// GetProductTopicFlavor retrieve a product's topic flavor
func (r *Registry) GetProductTopicFlavors(productId string) []TopicFlavor {
	flavors := []TopicFlavor{}
	c := r.db.C(dbNameProducts)
	p := Product{}
	if err := c.Find(bson.M{"ProductId": productId}).One(&p); err != nil {
		return flavors
	}
	c = r.db.C(dbNameTopicFlavors)
	flavor := TopicFlavor{}
	if err := c.Find(bson.M{"TenantId": p.ProductId, "FlavorName": p.TopicFlavor}).One(&flavor); err == nil {
		flavors = append(flavors, flavor)
	}
	return flavors
}

// UpdateTopicFlavor update a topic flavor
func (r *Registry) UpdateTopicFlavor(t *TopicFlavor) error {
	c := r.db.C(dbNameTopicFlavors)
	return c.Update(bson.M{"FlavorName": t.FlavorName, "Builtin": t.Builtin, "TenantId": t.TenantId}, t)
}
