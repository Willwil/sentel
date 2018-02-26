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
	"strconv"
	"time"

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

type device struct {
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

func BulkApplyDevices(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}

func BulkApplyGetStatus(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}

func BulkApplyGetDevices(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}

func GetDeviceList(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}

// Device bulk req.
type DeviceBulkRequest struct {
	DeviceName string `json:"deviceName"`
	ProductId  string `json:"productId"`
	Number     string `json:"number"`
}

func BulkGetDeviceStatus(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}
func GetDeviceByName(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}
func SaveDevicePropsByName(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}
func GetDevicePropsByName(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}
func RemoveDevicePropsByName(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}

func BulkRegisterDevices(ctx echo.Context) error {
	req := DeviceBulkRequest{}
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
