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
	dbNameDevices  = "devices"
	dbNameProducts = "products"
	dbNameTenants  = "tenants"
	dbNameRules    = "rules"
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
	hosts := c.MustString("registry", "hosts")
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
func New(owner string, c config.Config) (*Registry, error) {
	hosts := c.MustString(owner, "mongo")
	session, err := mgo.DialWithTimeout(hosts, 5*time.Second)
	if err != nil {
		glog.Infof("Failed to initialize registry:%s", err.Error())
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
func (r *Registry) Release() {
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
	return c.Insert(t, nil)
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

// BulkRegisterDevice add a lot of devices into registry
func (r *Registry) BulkRegisterDevices(devices []Device) error {
	for _, device := range devices {
		if err := r.RegisterDevice(&device); err != nil {
			return err
		}
	}
	return nil
}

// DeleteDevice delete a device from registry
func (r *Registry) DeleteDevice(productId string, deviceId string) error {
	c := r.db.C(dbNameDevices)
	if err := c.Find(bson.M{"ProductId": productId, "deviceId": deviceId}); err == nil {
		return c.Remove(bson.M{"ProductId": productId, "deviceId": deviceId})
	}
	return fmt.Errorf("invalid operataion")
}

// BulkDeleteDevice delete a lot of devices from registry
func (r *Registry) BulkDeleteDevice(devices []string) error {
	return nil
}

// UpdateDevice update device information in registry
func (r *Registry) UpdateDevice(dev *Device) error {
	c := r.db.C(dbNameDevices)
	return c.Update(bson.M{"ProductId": dev.ProductId, "DeviceId": dev.DeviceId}, dev)
}

// BulkUpdateDevice update a lot of devices in registry
func (r *Registry) BulkUpdateDevice(devices []Device) error {
	for _, device := range devices {
		if err := r.UpdateDevice(&device); err != nil {
			return err
		}
	}
	return nil
}

// Rule

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
	rule := Rule{}
	err := c.Find(bson.M{"RuleName": ruleName, "ProductId": productId}).One(&rule)
	return &rule, err
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProductRuleNames(productId string) ([]string, error) {
	c := r.db.C(dbNameRules)
	rules := []string{}
	err := c.Find(bson.M{"ProductId": productId}).Sort("RuleName").All(&rules)
	return rules, err
}

// DeleteRule delete a rule from registry
func (r *Registry) DeleteRule(productId string, ruleName string) error {
	c := r.db.C(dbNameRules)
	return c.Remove(bson.M{"RuleName": ruleName, "ProductId": productId})
}

// UpdateRule update rule information in registry
func (r *Registry) UpdateRule(rule *Rule) error {
	c := r.db.C(dbNameRules)
	return c.Update(bson.M{"RuleName": rule.RuleName, "ProductId": rule.ProductId}, rule)
}
