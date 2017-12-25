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
	"net/http"
	"time"

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/common/db"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// RegisterDevice register a new device in IoT hub
func RegisterDevice(ctx echo.Context) error {
	device := db.Device{}
	if err := ctx.Bind(&device); err != nil {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	}
	if device.TenantId == "" || device.ProductKey == "" || device.DeviceId == "" {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: "invalid parameter"})
	}
	config := ctx.(*base.ApiContext).Config
	r, err := db.NewRegistry("apiserver", config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()

	device.TimeCreated = time.Now()
	device.TimeUpdated = time.Now()
	device.DeviceSecret = uuid.NewV4().String()
	if err = r.RegisterDevice(&device); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &base.ApiResponse{Success: true, Result: &device})
}

// Delete the identify of a device from the identity registry of an IoT Hub
func DeleteDevice(ctx echo.Context) error {
	device := db.Device{}
	if err := ctx.Bind(&device); err != nil {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	}
	if device.TenantId == "" || device.ProductKey == "" || device.DeviceId == "" {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: "invalid parameter"})
	}
	config := ctx.(*base.ApiContext).Config
	registry, err := db.NewRegistry("apiserver", config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer registry.Release()

	// Get device into registry, the created product
	if err := registry.DeleteDevice(device.TenantId, device.ProductKey, device.DeviceId); err != nil {
		return ctx.JSON(http.StatusOK, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, base.ApiResponse{Success: true})
}

type device struct {
	DeviceId string `json:"id"`
	Status   string `json:"status"`
}

// Retrieve a device from the identify registry of an IoT hub
func GetOneDevice(ctx echo.Context) error {
	deviceId := ctx.Param("deviceId")
	productKey := ctx.QueryParam("productKey")
	tenantId := ctx.QueryParam("tenantId")

	if tenantId == "" || productKey == "" || deviceId == "" {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: "invalid parameter"})
	}
	// Connect with registry
	config := ctx.(*base.ApiContext).Config
	registry, err := db.NewRegistry("apiserver", config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer registry.Release()

	// Get device into registry, the created product
	dev, err := registry.GetDevice(tenantId, productKey, deviceId)
	if err != nil {
		return ctx.JSON(http.StatusOK,
			&base.ApiResponse{RequestId: uuid.NewV4().String(), Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &base.ApiResponse{Success: true, Result: dev})
}

// updateDevice update the identity of a device in the identity registry of an IoT Hub
func UpdateDevice(ctx echo.Context) error {
	device := db.Device{}
	if err := ctx.Bind(&device); err != nil {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	}
	if device.TenantId == "" || device.ProductKey == "" || device.DeviceId == "" {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: "invalid parameter"})
	}
	config := ctx.(*base.ApiContext).Config
	r, err := db.NewRegistry("apiserver", config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()

	device.TimeUpdated = time.Now()
	if err = r.UpdateDevice(&device); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &base.ApiResponse{Success: true, Result: &device})
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
	/*
		// Get product
		glog.Infof("RegisterDevice ctx:[%s] ", ctx.Param("number"))
		req := new(registerDeviceRequest)
		if err := ctx.Bind(req); err != nil {
			return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Success: false, Message: err.Error()})
		}
		config := ctx.(*base.ApiContext).Config
		// Connect with registry
		r, err := db.NewRegistry("apiserver", config)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
		}
		defer r.Release()

		//
		// Insert device into registry, the created product
		// will be modified to retrieve specific information sucha as
		// product.id and creation time
		max, err := strconv.Atoi(ctx.Param("number"))
		rdevices := []db.Device{}
		for i := 0; i < max; i++ {
			dp := db.Device{
				DeviceId:     uuid.NewV4().String(),
				ProductKey:   req.ProductKey,
				TimeCreated:  time.Now(),
				TimeModified: time.Now(),
			}
			glog.Infof("+++%d,%s\n", i, dp.DeviceId)
			rdevices = append(rdevices, dp)
		}
		err = r.BulkRegisterDevices(rdevices)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
		}
		return ctx.JSON(http.StatusOK,
			&base.ApiResponse{
				Success: true,
				Result:  rdevices,
			})
	*/
	return nil
}
