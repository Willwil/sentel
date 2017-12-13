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

	"github.com/cloustone/sentel/core"
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
	config  core.Config
	session *mgo.Session
	db      *mgo.Database
}

// InitializeRegistry try to connect with background database
// to confirm wether it is normal
func InitializeRegistry(c core.Config) error {
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
func NewRegistry(c core.Config) (*Registry, error) {
	hosts := c.MustString("registry", "hosts")
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
	err := c.Find(bson.M{"Id": bson.M{"$regex": id, "$options": "$i"}}).One(nil)
	return err
}

// AddTenant insert a tenant into registry
func (r *Registry) RegisterTenant(t *Tenant) error {
	c := r.db.C(dbNameTenants)
	if err := c.Find(bson.M{"Name": bson.M{"$regex": t.Name, "$options": "$i"}}).One(nil); err == nil {
		return fmt.Errorf("Tenant %s already exist", t.Name)
	}
	return c.Insert(t, nil)
}

func (r *Registry) DeleteTenant(id string) error {
	c := r.db.C(dbNameTenants)
	return c.Remove(bson.M{"Name": bson.M{"$regex": id, "$options": "$i"}})
}

func (r *Registry) GetTenant(id string) (*Tenant, error) {
	c := r.db.C(dbNameTenants)
	t := Tenant{}
	err := c.Find(bson.M{"Name": bson.M{"$regex": id, "$options": "$i"}}).One(&t)
	if err != nil {
		return nil, ErrorNotFound
	}
	return &t, nil
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
			glog.Infof("\nSearch end!\n\n")
			break
		} else {
			lastId = tn.Name
			glog.Infof("\nName:%s\n", lastId)
			tns = append(tns, tn)
		}
		for iter.Next(&tn) {
			lastId = tn.Name
			glog.Infof("\nName:%s\n", lastId)
			tns = append(tns, tn)
		}
		if iter.Err() != nil {
			glog.Infof("\nSearch end!\n\n")
			break
		}
		iter = c.Find(bson.M{"Name": bson.M{"$gt": lastId}}).Sort("Name").Limit(PageSize).Iter()
	}
	iter.Close()
	return tns, nil
}

func (r *Registry) UpdateTenant(id string, t *Tenant) error {
	c := r.db.C(dbNameTenants)
	return c.Update(bson.M{"Name": bson.M{"$regex": id, "$options": "$i"}}, t)
}

// Product
// CheckProductNameAvailable check wethere product name is available
func (r *Registry) CheckProductNameAvailable(p *Product) bool {
	c := r.db.C(dbNameProducts)
	err := c.Find(bson.M{"Id": bson.M{"$regex": p.Id, "$options": "$i"}}).One(nil)
	return err != nil
}

// RegisterProduct register a product into registry
func (r *Registry) RegisterProduct(p *Product) error {
	c := r.db.C(dbNameProducts)
	glog.Infof("use products")
	product := &Product{}
	if err := c.Find(bson.M{"Name": bson.M{"$regex": p.Name, "$options": "$i"}}).One(product); err == nil {
		return fmt.Errorf("product %s already exist", p.Name)
	}
	return c.Insert(p)

}

// DeleteProduct delete a product from registry
func (r *Registry) DeleteProduct(id string) error {
	c := r.db.C(dbNameProducts)
	return c.Remove(bson.M{"Id": bson.M{"$regex": id, "$options": "$i"}})
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProduct(id string) (*Product, error) {
	c := r.db.C(dbNameProducts)
	product := &Product{}
	err := c.Find(bson.M{"Id": bson.M{"$regex": id, "$options": "$i"}}).One(product)
	return product, err
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetAllProducts() ([]Product, error) {
	c := r.db.C(dbNameProducts)
	pro := Product{}
	products := []Product{}
	var lastId string

	iter := c.Find(nil).Sort("Id").Limit(PageSize).Iter()
	for {
		if iter.Next(&pro) == false {
			glog.Infof("\nSearch end!\n\n")
			break
		} else {
			lastId = pro.Id
			glog.Infof("\nId:%s\n", lastId)
			products = append(products, pro)
		}
		for iter.Next(&pro) {
			lastId = pro.Id
			glog.Infof("\nId:%s\n", lastId)
			products = append(products, pro)
		}
		if iter.Err() != nil {
			glog.Infof("\nSearch end!\n\n")
			break
		}
		iter = c.Find(bson.M{"Id": bson.M{"$gt": lastId}}).Sort("Id").Limit(PageSize).Iter()
	}
	iter.Close()
	return products, nil
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProductByName(name string) (*Product, error) {
	c := r.db.C(dbNameProducts)
	product := &Product{}
	err := c.Find(bson.M{"Name": bson.M{"$regex": name, "$options": "$i"}}).One(product)
	return product, err
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
func (r *Registry) GetProductDevices(id string) ([]Device, error) {
	c := r.db.C(dbNameDevices)

	devices := []Device{}
	var device Device
	var lastId string

	iter := c.Find(bson.M{"ProductId": bson.M{"$regex": id, "$options": "$i"}}).Sort("Id").Limit(PageSize).Iter()
	for {
		if iter.Next(&device) == false {
			glog.Infof("\nSearch end!\n\n")
			break
		} else {
			lastId = device.Id
			glog.Infof("\nId:%s\ndev:%s\n\n", lastId, device)
			devices = append(devices, device)
		}
		for iter.Next(&device) {
			lastId = device.Id
			glog.Infof("\nId:%s\ndev:%s\n\n", lastId, device)
			devices = append(devices, device)
		}
		if iter.Err() != nil {
			glog.Infof("\nError=%s,Search end!\n\n", iter.Err())
			break
		}
		iter = c.Find(bson.M{"ProductId": bson.M{"$regex": id, "$options": "$i"}, "Id": bson.M{"$gt": lastId}}).Sort("Id").Limit(PageSize).Iter()
	}
	iter.Close()
	return devices, nil
}

// GetProductDevices get product's device list
func (r *Registry) GetProductDevicesPage(id string, indexId string) ([]Device, string, error) {
	c := r.db.C(dbNameDevices)

	devices := []Device{}
	var device Device
	var lastId string

	iter := c.Find(bson.M{"ProductId": bson.M{"$regex": id, "$options": "$i"}, "Id": bson.M{"$gt": indexId}}).Sort("Id").Limit(PageSize).Iter()
	for iter.Next(&device) {
		lastId = device.Id
		glog.Infof("\nId:%s\ndev:%s\n\n", lastId, device)
		devices = append(devices, device)
	}
	iter.Close()
	return devices, lastId, nil
}

// GetProductDevices get product's device list
// curl -XGET "http://localhost:4145/api/vi/products/mytest?api-version=v1
func (r *Registry) GetProductDevicesByName(name string) ([]Device, error) {
	glog.Infof("\n Product name:%s\n\n", name)
	pro, err := r.GetProductByName(name)
	if err != nil {
		glog.Infof("\nNot found! Product name:%s\n\n", name)
		return nil, ErrorNotFound
	}

	devices := []Device{}
	var device Device
	var lastId string

	c := r.db.C(dbNameDevices)

	glog.Infof("\nSearch all device by productKey:%s ...\n\n", pro.ProductKey)
	iter := c.Find(bson.M{"ProductKey": bson.M{"$regex": pro.ProductKey, "$options": "$i"}}).Sort("Id").Limit(PageSize).Iter()
	for {
		if iter.Next(&device) == false {
			glog.Infof("\nSearch end!\n\n")
			break
		} else {
			lastId = device.Id
			glog.Infof("\nId:%s\ndev:%s\n\n", lastId, device)
			devices = append(devices, device)
		}
		for iter.Next(&device) {
			lastId = device.Id
			glog.Infof("\nId:%s\ndev:%s\n\n", lastId, device)
			devices = append(devices, device)
		}
		if iter.Err() != nil {
			glog.Infof("\nErr, Search end!\n\n")
			break
		}
		iter = c.Find(bson.M{"ProductKey": bson.M{"$regex": pro.ProductKey, "$options": "$i"}, "Id": bson.M{"$gt": lastId}}).Sort("Id").Limit(PageSize).Iter()
	}
	iter.Close()
	return devices, nil
}

// GetProductDevices get product's device list
func (r *Registry) GetProductDevicesPageByName(name string, indexId string) ([]Device, string, error) {
	glog.Infof("\n Product name:%s\n\n", name)
	pro, err := r.GetProductByName(name)
	if err != nil {
		glog.Infof("\nNot found! Product name:%s\n\n", name)
		return nil, "0", ErrorNotFound
	}
	c := r.db.C(dbNameDevices)

	devices := []Device{}
	var device Device
	var lastId string

	iter := c.Find(bson.M{"ProductKey": bson.M{"$regex": pro.ProductKey, "$options": "$i"}, "Id": bson.M{"$gt": indexId}}).Sort("Id").Limit(PageSize).Iter()
	for iter.Next(&device) {
		lastId = device.Id
		glog.Infof("\nId:%s\ndev:%s\n\n", lastId, device)
		devices = append(devices, device)
	}
	iter.Close()
	return devices, lastId, nil
}

// UpdateProduct update product detail information in registry
func (r *Registry) UpdateProduct(p *Product) error {
	c := r.db.C(dbNameProducts)
	return c.Update(bson.M{"Id": p.Id}, p)
}

// Device

// RegisterDevice add a new device into registry
func (r *Registry) RegisterDevice(dev *Device) error {
	glog.Infof("RegisterDevice:[%s].%s", dev.Id, dev)
	c := r.db.C(dbNameDevices)
	device := Device{}
	if err := c.Find(bson.M{"Id": bson.M{"$regex": dev.Id, "$options": "$i"}}).One(&device); err == nil { // found existed device
		return fmt.Errorf("device %s already exist", dev.ProductId)
	}
	return c.Insert(dev)
}

// GetDevice retrieve a device information from registry/
func (r *Registry) GetMultipleDevices(dev *Device) []Device {
	c := r.db.C(dbNameDevices)
	devices := []Device{}
	var device Device

	iter := c.Find(bson.M{"Name": bson.M{"$regex": dev.Name, "$options": "$i"}}).Iter()
	for iter.Next(&device) {
		fmt.Println(device)
		devices = append(devices, device)
	}
	fmt.Println("find end!")
	return devices
}

// GetDevice retrieve a device information from registry/
func (r *Registry) GetDevice(id string) (*Device, error) {
	c := r.db.C(dbNameDevices)
	device := &Device{}
	err := c.Find(bson.M{"Name": bson.M{"$regex": id, "$options": "$i"}}).One(device)
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
			lastId = dev.Id
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
	return c.Update(bson.M{"Id": bson.M{"$regex": dev.Id, "$options": "$i"}}, dev)
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
	if err := c.Find(bson.M{"id": bson.M{"$regex": rule.Id, "$options": "$i"}}); err == nil { // found existed device
		return fmt.Errorf("rule %s already exist", rule.Id)
	}
	return c.Insert(rule)
}

// GetRule retrieve a rule information from registry/
func (r *Registry) GetRule(id string) (*Rule, error) {
	c := r.db.C(dbNameRules)
	rule := &Rule{}
	err := c.Find(bson.M{"Id": bson.M{"$regex": id, "$options": "$i"}}).One(rule)
	return rule, err
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetAllRules() ([]Rule, error) {
	c := r.db.C(dbNameRules)

	rule := Rule{}
	rules := []Rule{}
	var lastId string

	iter := c.Find(nil).Sort("Id").Limit(PageSize).Iter()
	for {
		if iter.Next(&rule) == false {
			glog.Infof("\nSearch end!\n\n")
			break
		}
		for iter.Next(&r) {
			lastId = rule.Id
			glog.Infof("\nId:%s\n", lastId)
			rules = append(rules, rule)
		}
		if iter.Err() != nil {
			glog.Infof("\nSearch end!\n\n")
			break
		}
		iter = c.Find(bson.M{"Id": bson.M{"$gt": lastId}}).Sort("Id").Limit(PageSize).Iter()
	}
	iter.Close()
	return rules, nil
}

// DeleteRule delete a rule from registry
func (r *Registry) DeleteRule(id string) error {
	c := r.db.C(dbNameRules)
	return c.Remove(bson.M{"Id": bson.M{"$regex": id, "$options": "$i"}})
}

// UpdateRule update rule information in registry
func (r *Registry) UpdateRule(rule *Rule) error {
	c := r.db.C(dbNameRules)
	return c.Update(bson.M{"Id": bson.M{"$regex": rule.Id, "$options": "$i"}}, rule)
}
