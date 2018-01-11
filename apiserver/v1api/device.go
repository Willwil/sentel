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

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/keystone/ram"
	"github.com/cloustone/sentel/pkg/registry"

	"github.com/labstack/echo"
)

// RegisterDevice register a new device in IoT hub
func CreateDevice(ctx echo.Context) error {
	accessId := getAccessId(ctx)

	device := registry.Device{}
	if err := ctx.Bind(&device); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if device.ProductId == "" || device.DeviceName == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	// Authorization
	objectName := device.ProductId + "/devices"
	if err := base.Authorize(accessId, objectName, "w"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	device.DeviceId = ram.NewObjectId()
	device.TimeCreated = time.Now()
	device.TimeUpdated = time.Now()
	device.DeviceSecret = ram.NewObjectId()
	if err := r.RegisterDevice(&device); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{Result: &device})
}

// Delete the identify of a device from the identity registry of an IoT Hub
func RemoveDevice(ctx echo.Context) error {
	accessId := getAccessId(ctx)

	device := registry.Device{}
	if err := ctx.Bind(&device); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if device.ProductId == "" || device.DeviceId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	// Authorization
	objectName := device.ProductId + "/devices"
	if err := base.Authorize(accessId, objectName, "w"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	// Get device into registry, the created product
	if err := r.DeleteDevice(device.ProductId, device.DeviceId); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	base.DestroyResource(device.DeviceId, accessId)
	return ctx.JSON(OK, apiResponse{})
}

type device struct {
	DeviceId string `json:"id"`
	Status   string `json:"status"`
}

// Retrieve a device from the identify registry of an IoT hub
func GetOneDevice(ctx echo.Context) error {
	accessId := getAccessId(ctx)

	deviceId := ctx.Param("deviceId")
	productId := ctx.QueryParam("productId")
	if productId == "" || deviceId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	// Authorization
	objectName := productId + "/devices"
	if err := base.Authorize(accessId, objectName, "r"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
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
	accessId := getAccessId(ctx)

	device := registry.Device{}
	if err := ctx.Bind(&device); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if device.ProductId == "" || device.DeviceId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	// Authorization
	objectName := device.ProductId + "/devices"
	if err := base.Authorize(accessId, objectName, "r"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	device.TimeUpdated = time.Now()
	if err := r.UpdateDevice(&device); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{Result: &device})
}

func GetDevicesStatus(ctx echo.Context) error {
	return nil
}

func BulkApplyDevices(ctx echo.Context) error {
	return nil
}

func BulkApplyGetStatus(ctx echo.Context) error {
	return nil
}

func BulkApplyGetDevices(ctx echo.Context) error {
	return nil
}

func GetDeviceList(ctx echo.Context) error {
	return nil
}

// Device bulk req.
type DeviceBulkRequest struct {
	Number     string
	ProductId  string
	DeviceName string
}

func BulkGetDeviceStatus(ctx echo.Context) error {
	return nil
}
func GetDeviceByName(ctx echo.Context) error {
	return nil
}
func SaveDevicePropsByName(ctx echo.Context) error {
	return nil
}
func GetDevicePropsByName(ctx echo.Context) error {
	return nil
}
func DeleteDevicePropsByName(ctx echo.Context) error {
	return nil
}

func BulkRegisterDevices(ctx echo.Context) error {
	accessId := getAccessId(ctx)

	req := DeviceBulkRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.Number == "" || req.ProductId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	// Authorization
	objectName := req.ProductId + "/devices"
	if err := base.Authorize(accessId, objectName, "r"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}
	n, err := strconv.Atoi(req.Number)
	if err != nil || n < 1 {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	rdevices := []registry.Device{}
	for i := 0; i < n; i++ {
		d := registry.Device{
			ProductId:    req.ProductId,
			DeviceId:     ram.NewObjectId(),
			DeviceName:   req.DeviceName,
			TimeCreated:  time.Now(),
			DeviceSecret: ram.NewObjectId(),
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
	return nil
}

func UpdateShadowDevice(ctx echo.Context) error {
	return nil
}
