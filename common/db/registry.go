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

package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	com "github.com/cloustone/sentel/common"
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
	config  com.Config
	session *mgo.Session
	db      *mgo.Database
}

// InitializeRegistry try to connect with background database
// to confirm wether it is normal
func InitializeRegistry(c com.Config) error {
	hosts := c.MustString("registry", "hosts")
	glog.Infof("Initializing registry:%s...", hosts)
	session, err := mgo.Dial(hosts)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)

	EnsureDevicesIndex(session)
	EnsureProductsIndex(session)
	EnsureTenantsIndex(session)
	EnsureRulesIndex(session)

	session.Close()
	return nil
}

// NewRegistry create registry instance
func NewRegistry(owner string, c com.Config) (*Registry, error) {
	hosts := c.MustString(owner, "mongo")
	session, err := mgo.DialWithTimeout(hosts, 5*time.Second)
	if err != nil {
		glog.Infof("Failed to initialize registry:%s", err.Error())
		return nil, err
	}
	return &Registry{session: session, db: session.DB("registry"), config: c}, nil
}

func EnsureDevicesIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("registry").C(dbNameDevices)

	index := mgo.Index{
		Key:        []string{"Id", "Name", "ProductId", "ProductKey"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		glog.Infof("Create dev INDEX error! %s\n", err)
		panic(err)
	}
}

func EnsureProductsIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("registry").C(dbNameProducts)

	index := mgo.Index{
		Key:        []string{"Id", "Name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		glog.Infof("Create INDEX error! %s\n", err)
		panic(err)
	}
}

func EnsureTenantsIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("registry").C(dbNameTenants)

	index := mgo.Index{
		Key:        []string{"Name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		glog.Infof("Create INDEX error! %s\n", err)
		panic(err)
	}
}

func EnsureRulesIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("registry").C(dbNameRules)

	index := mgo.Index{
		Key:        []string{"Id", "Name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		glog.Infof("Create INDEX error! %s\n", err)
		panic(err)
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
	tn := Tenant{}
	tns := []Tenant{}
	var lastId string

	iter := c.Find(nil).Sort("Name").Limit(PageSize).Iter()
	for {
		if iter.Next(&tn) == false {
			break
		} else {
			lastId = tn.TenantId
			tns = append(tns, tn)
		}
		for iter.Next(&tn) {
			lastId = tn.TenantId
			tns = append(tns, tn)
		}
		if iter.Err() != nil {
			break
		}
		iter = c.Find(bson.M{"Name": bson.M{"$gt": lastId}}).Sort("Name").Limit(PageSize).Iter()
	}
	iter.Close()
	return tns, nil
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
func (r *Registry) DeleteProduct(productKey string) error {
	c := r.db.C(dbNameDevices)
	c.Remove(bson.M{"ProductKey": productKey})
	c = r.db.C(dbNameProducts)
	return c.Remove(bson.M{"ProductKey": productKey})
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProduct(productKey string) (*Product, error) {
	c := r.db.C(dbNameProducts)
	product := &Product{}
	err := c.Find(bson.M{"ProductKey": productKey}).One(product)
	return product, err
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetTenantProducts(tenantId string) ([]Product, error) {
	c := r.db.C(dbNameProducts)
	pro := Product{}
	products := []Product{}
	var lastId string

	iter := c.Find(bson.M{"TenantId": tenantId}).Sort("ProductId").Iter()
	for {
		if iter.Next(&pro) == false {
			break
		} else {
			lastId = pro.ProductId
			products = append(products, pro)
		}
		for iter.Next(&pro) {
			lastId = pro.ProductId
			products = append(products, pro)
		}
		if iter.Err() != nil {
			break
		}
		iter = c.Find(bson.M{"Id": bson.M{"$gt": lastId}}).Sort("Id").Limit(PageSize).Iter()
	}
	iter.Close()
	return products, nil
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProductsByCategory(tenantId string, cat string) ([]Product, error) {
	c := r.db.C(dbNameProducts)
	pro := Product{}
	products := []Product{}
	var lastId string

	iter := c.Find(bson.M{"TenantId": tenantId, "CategoryId": cat}).Sort("ProductId").Iter()
	for {
		if iter.Next(&pro) == false {
			break
		} else {
			lastId = pro.ProductId
			products = append(products, pro)
		}
		for iter.Next(&pro) {
			lastId = pro.ProductId
			products = append(products, pro)
		}
		if iter.Err() != nil {
			break
		}
		iter = c.Find(bson.M{"CategoryId": bson.M{"$regex": cat, "$options": "$i"}, "Id": bson.M{"$gt": lastId}}).Sort("Id").Limit(PageSize).Iter()
	}
	iter.Close()
	return products, nil
}

// GetProduct retrieve product detail information from registry
// bson.M{"Name": bson.M{"$regex": name, "$options": "$i"}}
func (r *Registry) GetProductByDef(query bson.M) (*Product, error) {
	c := r.db.C(dbNameProducts)
	product := &Product{}
	err := c.Find(query).One(product)
	return product, err
}

// GetProductDevices get product's device list
// mongo golang:
// http://www.jianshu.com/p/b63e5cfa4ce5
func (r *Registry) GetProductDevices(productKey string) ([]Device, error) {
	c := r.db.C(dbNameDevices)

	devices := []Device{}
	var device Device

	iter := c.Find(bson.M{"ProductKey": productKey}).Sort("DeviceId").Limit(PageSize).Iter()
	for {
		if iter.Next(&device) == false {
			break
		} else {
			devices = append(devices, device)
		}
		for iter.Next(&device) {
			devices = append(devices, device)
		}
		if iter.Err() != nil {
			break
		}
	}
	iter.Close()
	return devices, nil
}

// UpdateProduct update product detail information in registry
func (r *Registry) UpdateProduct(p *Product) error {
	c := r.db.C(dbNameProducts)
	return c.Update(bson.M{"ProductKey": p.ProductKey}, p)
}

// Device

// RegisterDevice add a new device into registry
func (r *Registry) RegisterDevice(dev *Device) error {
	c := r.db.C(dbNameDevices)
	device := Device{}
	if err := c.Find(bson.M{"Id": bson.M{"$regex": dev.DeviceId, "$options": "$i"}}).One(&device); err == nil { // found existed device
		return fmt.Errorf("device %s already exist", dev.ProductId)
	}
	return c.Insert(dev)
}

// GetDevice retrieve a device information from registry/
func (r *Registry) GetMultipleDevices(name string) ([]Device, error) {
	c := r.db.C(dbNameDevices)
	devices := []Device{}
	var device Device

	iter := c.Find(bson.M{"Name": bson.M{"$regex": name, "$options": "$i"}}).Iter()
	for iter.Next(&device) {
		fmt.Println(device)
		devices = append(devices, device)
	}
	return devices, nil
}

// GetDevice retrieve a device information from registry/
func (r *Registry) GetDevice(id string) (*Device, error) {
	c := r.db.C(dbNameDevices)
	device := &Device{}
	err := c.Find(bson.M{"Id": bson.M{"$regex": id, "$options": "$i"}}).One(device)
	v, _ := json.Marshal(*device)
	glog.Infof("\ndevice %s\n", v)
	return device, err
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetAllDevices() ([]Device, error) {
	c := r.db.C(dbNameDevices)
	dev := Device{}
	devices := []Device{}
	//err := c.Find(bson.M{"Id": bson.M{"$regex": id, "$options": "$i"}}).One(product)
	var lastId string

	iter := c.Find(nil).Sort("Id").Limit(PageSize).Iter()
	for {
		if iter.Next(&dev) == false {
			glog.Infof("\nSearch end!\n\n")
			break
		}
		for iter.Next(&dev) {
			lastId = dev.DeviceId
			glog.Infof("\nId:%s\n", lastId)
			devices = append(devices, dev)
		}
		if iter.Err() != nil {
			glog.Infof("\nSearch end!\n\n")
			break
		}
		iter = c.Find(bson.M{"Id": bson.M{"$gt": lastId}}).Sort("Id").Limit(PageSize).Iter()
	}
	iter.Close()
	return devices, nil
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
func (r *Registry) DeleteDevice(id string) error {
	c := r.db.C(dbNameDevices)
	return c.Remove(bson.M{"Id": bson.M{"$regex": id, "$options": "$i"}})
}

// BulkDeleteDevice delete a lot of devices from registry
func (r *Registry) BulkDeleteDevice(devices []string) error {
	for _, id := range devices {
		if err := r.DeleteDevice(id); err != nil {
			return err
		}
	}
	return nil
}

// UpdateDevice update device information in registry
func (r *Registry) UpdateDevice(dev *Device) error {
	c := r.db.C(dbNameDevices)
	return c.Update(bson.M{"Id": bson.M{"$regex": dev.DeviceId, "$options": "$i"}}, dev)
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
	if err := c.Find(bson.M{"ruleName": rule.RuleName, "productKey": rule.ProductKey}); err == nil {
		return fmt.Errorf("rule %s already exist", rule.RuleName)
	}
	return c.Insert(rule)
}

// GetRule retrieve a rule information from registry/
func (r *Registry) GetRule(productId string, ruleName string) (*Rule, error) {
	c := r.db.C(dbNameRules)
	rule := Rule{}
	err := c.Find(bson.M{"ruleName": ruleName, "productId": productId}).One(&rule)
	return &rule, err
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProductRuleNames(productKey string) ([]string, error) {
	c := r.db.C(dbNameRules)

	rule := Rule{}
	rules := []string{}
	iter := c.Find(bson.M{"productKey": productKey}).Sort("ruleName").Iter()
	for iter.Next(&rule) {
		if iter.Err() == nil {
			break
		}
		rules = append(rules, rule.RuleName)
	}
	return rules, nil
}

// DeleteRule delete a rule from registry
func (r *Registry) DeleteRule(productId string, ruleName string) error {
	c := r.db.C(dbNameRules)
	return c.Remove(bson.M{"ruleName": ruleName, "productId": productId})
}

// UpdateRule update rule information in registry
func (r *Registry) UpdateRule(rule *Rule) error {
	c := r.db.C(dbNameRules)
	return c.Update(bson.M{"ruleName": rule.RuleName, "productKey": rule.ProductKey}, rule)
}
