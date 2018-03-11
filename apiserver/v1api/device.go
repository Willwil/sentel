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

package v1api

import (
	"github.com/golang/glog"
	"strconv"
	"time"
	"encoding/json"

	"github.com/cloustone/sentel/apiserver/util"
	"github.com/cloustone/sentel/pkg/registry"

	"github.com/labstack/echo"
)

type deviceRequest struct {
	DeviceName string `json:"DeviceName"`
	ProductId  string `json:"ProductId"`
	DeviceId   string `json:"DeviceId"`
}

// RegisterDevice register a new device in IoT hub
func CreateDevice(ctx echo.Context) error {
	req := deviceRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductId == "" || req.DeviceName == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	r := getRegistry(ctx)
	device := registry.Device{}
	device.DeviceName = req.DeviceName
	device.ProductId = req.ProductId
	device.DeviceId = util.NewObjectId()
	device.TimeCreated = time.Now()
	device.TimeUpdated = time.Now()
	device.DeviceSecret = util.NewObjectId()
	if err := r.RegisterDevice(&device); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{Result: &device})
}

// Delete the identify of a device from the identity registry of an IoT Hub
func RemoveDevice(ctx echo.Context) error {
	deviceId := ctx.Param("deviceId")
	productId := ctx.Param("productId")
	if productId == "" || deviceId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	r := getRegistry(ctx)
	// Get device into registry, the created product
	if err := r.DeleteDevice(productId, deviceId); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{})
}

type deviceStatus struct {
	DeviceId string `json:"id"`
	Status   string `json:"status"`
}

// Retrieve a device from the identify registry of an IoT hub
func GetOneDevice(ctx echo.Context) error {
	deviceId := ctx.Param("deviceId")
	productId := ctx.Param("productId")
	if productId == "" || deviceId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	// Get device into registry, the created product
	r := getRegistry(ctx)
	dev, err := r.GetDevice(productId, deviceId)
	if err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{Result: dev})
}

// updateDevice update the identity of a device in the identity registry of an IoT Hub
func UpdateDevice(ctx echo.Context) error {
	req := deviceRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductId == "" || req.DeviceId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	device := registry.Device{}
	device.ProductId = req.ProductId
	device.DeviceId = req.DeviceId
	device.DeviceName = req.DeviceName
	r := getRegistry(ctx)
	device.TimeUpdated = time.Now()
	if err := r.UpdateDevice(&device); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{Result: &device})
}

// Device bulk req.
type BulkDeviceRegisterRequest struct {
	DeviceName string `json:"deviceName"`
	ProductId  string `json:"productId"`
	Number     string `json:"number"`
}

type BulkDeviceQueryRequest struct {
	ProductId  string `json:"productId"`
	DevicesName []string `json:"deviceName"`
}

func BulkApplyDevices(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}

func BulkApplyGetStatus(ctx echo.Context) error {
	req := BulkDeviceQueryRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	rdevices := []registry.Device{}
	for _, dev := range req.DevicesName {
		d := registry.Device{
			ProductId:    req.ProductId,
			DeviceId:     util.NewObjectId(),
			DeviceName:   dev,
		}
		rdevices = append(rdevices, d)
	}
	r := getRegistry(ctx)
	odevs, err := r.BulkGetDevices(rdevices)
	if err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	rdevicestatus := []deviceStatus{}
	for _, dev := range odevs {
		rdevicestatus = append(rdevicestatus, deviceStatus{DeviceId: dev.DeviceId, Status: dev.DeviceStatus})
	}
	return ctx.JSON(OK, apiResponse{Result: rdevicestatus})
}

func BulkApplyGetDevices(ctx echo.Context) error {
	req := BulkDeviceQueryRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	rdevices := []registry.Device{}
	for _, dev := range req.DevicesName {
		d := registry.Device{
			ProductId:    req.ProductId,
			DeviceId:     util.NewObjectId(),
			DeviceName:   dev,
			TimeCreated:  time.Now(),
			DeviceSecret: util.NewObjectId(),
		}
		rdevices = append(rdevices, d)
	}
	r := getRegistry(ctx)
	odevs, err := r.BulkGetDevices(rdevices)
	if err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}

	return ctx.JSON(NotImplemented, apiResponse{Result: odevs})
}

func GetDeviceList(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}



func BulkGetDeviceStatus(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}
func GetDeviceByName(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}

type devicePropsRequest struct {
	DeviceName string `json:"DeviceName"`
	ProductId  string `json:"ProductId"`
	DeviceId   string `json:"DeviceId"`
	Props      []map[string]string `json:"Props"`
}

func SaveDeviceProps(ctx echo.Context) error {
	req := devicePropsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductId == "" || req.DeviceId =="" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	r := getRegistry(ctx)
	device, err := r.GetDevice(req.ProductId, req.DeviceId)
	if err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	props := make(map[string]string)
	if device.Props != "" {
		byteprops := []byte(device.Props)
		var mapprops map[string]string
		err = json.Unmarshal(byteprops, &mapprops)
		if err != nil {
			return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
		}
		for key, value := range mapprops{
			props[key] = value 
		}
	
	}
	for _,p := range req.Props{
		for key, value:= range p{
			props[key] = value //replace original.
		}
	}
	b, err:=json.Marshal(props)
	if err != nil{
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	device.Props = string(b)
	if err := r.UpdateDevice(device); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}

	return ctx.JSON(OK, apiResponse{Result: &device})
}

func SaveDevicePropsByName(ctx echo.Context) error {
	req := devicePropsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	r := getRegistry(ctx)
	devices, err := r.GetDevicesByName(req.ProductId, req.DeviceName)
	if err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	dstprops := make(map[string]string)
	for _,device :=  range devices{
		if device.Props != "" {
			byteprops := []byte(device.Props)
			var mapprops map[string]string
			err = json.Unmarshal(byteprops, &mapprops)
			if err != nil {
				return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
			}
			for key, value := range mapprops{
				dstprops[key] = value 
			}
		
		}
		for _,p := range req.Props{
			for key, value:= range p{
				dstprops[key] = value 
			}
		}
		b, err:=json.Marshal(dstprops)
		glog.Infof("sp:%s\n", b)
		if err != nil{
			return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
		}
		device.Props = string(b)
	}

	if err := r.BulkUpdateDevice(devices); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}

	return ctx.JSON(OK, apiResponse{Result: devices})
}

func GetDeviceProps(ctx echo.Context) error {
	req := devicePropsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductId == "" || req.DeviceId =="" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	r := getRegistry(ctx)
	device, err := r.GetDevice(req.ProductId, req.DeviceId)
	if err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	props := make(map[string]string)
	if device.Props != "" {
		byteprops := []byte(device.Props)
		var mapprops map[string]string
		err = json.Unmarshal(byteprops, &mapprops)
		if err != nil {
			return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
		}
		for key, value := range mapprops{
			props[key] = value 
		}
	
	}

	return ctx.JSON(OK, apiResponse{Result: &device})
}

func GetDevicePropsByName(ctx echo.Context) error {
	req := devicePropsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	r := getRegistry(ctx)
	devices, err := r.GetDevicesByName(req.ProductId, req.DeviceName)
	if err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	dstprops := make(map[string]string)
	for _,device :=  range devices{
		if device.Props != "" {
			byteprops := []byte(device.Props)
			var mapprops map[string]string
			err = json.Unmarshal(byteprops, &mapprops)
			if err != nil {
				return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
			}
			for key, value := range mapprops{
				dstprops[key] = value 
			}
		
		}
	}

	return ctx.JSON(OK, apiResponse{Result: devices})
}
func RemoveDevicePropsByName(ctx echo.Context) error {
	req := devicePropsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	r := getRegistry(ctx)
	devices, err := r.GetDevicesByName(req.ProductId, req.DeviceName)
	if err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	for _,device :=  range devices{
		device.Props = ""
	}
	if err := r.BulkUpdateDevice(devices); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{Result: devices})
}

func BulkRegisterDevices(ctx echo.Context) error {
	req := BulkDeviceRegisterRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.Number == "" || req.ProductId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	n, err := strconv.Atoi(req.Number)
	if err != nil || n < 1 {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	rdevices := []registry.Device{}
	for i := 0; i < n; i++ {
		d := registry.Device{
			ProductId:    req.ProductId,
			DeviceId:     util.NewObjectId(),
			DeviceName:   req.DeviceName,
			TimeCreated:  time.Now(),
			DeviceSecret: util.NewObjectId(),
		}
		rdevices = append(rdevices, d)
	}
	r := getRegistry(ctx)
	if err = r.BulkRegisterDevices(rdevices); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{Result: rdevices})
}

func GetShadowDevice(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}

func UpdateShadowDevice(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}
