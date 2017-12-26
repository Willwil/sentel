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
	"time"

	"github.com/cloustone/sentel/common/db"
	"github.com/cloustone/sentel/keystone/orm"

	"github.com/labstack/echo"
)

// RegisterDevice register a new device in IoT hub
func RegisterDevice(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	device := db.Device{}
	if err := ctx.Bind(&device); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}
	if device.ProductId == "" || device.DeviceName == "" {
		return reply(ctx, BadRequest, apiResponse{Message: "invalid parameter"})
	}
	// Authorization
	objectName := device.ProductId + "/$devices"
	if err := orm.AccessObject(objectName, accessId, orm.AccessRightWrite); err != nil {
		return reply(ctx, Unauthorized, apiResponse{Message: "access denied"})
	}
	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()

	device.DeviceId = orm.NewObjectId()
	device.TimeCreated = time.Now()
	device.TimeUpdated = time.Now()
	device.DeviceSecret = orm.NewObjectId()
	if err = r.RegisterDevice(&device); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	// Declare keystone object
	orm.CreateObject(orm.Object{
		ObjectName:     device.DeviceName,
		ObjectId:       device.DeviceId,
		ParentObjectId: device.ProductId,
		Category:       "device",
		CreatedTime:    time.Now(),
		Owner:          accessId,
	})

	return reply(ctx, OK, apiResponse{Result: &device})
}

// Delete the identify of a device from the identity registry of an IoT Hub
func DeleteDevice(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	device := db.Device{}
	if err := ctx.Bind(&device); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}
	if device.ProductId == "" || device.DeviceId == "" {
		return reply(ctx, BadRequest, apiResponse{Message: "invalid parameter"})
	}
	// Authorization
	objectName := device.ProductId + "/$devices"
	if err := orm.AccessObject(objectName, accessId, orm.AccessRightWrite); err != nil {
		return reply(ctx, Unauthorized, apiResponse{Message: "access denied"})
	}

	registry, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, &apiResponse{Message: err.Error()})
	}
	defer registry.Release()

	// Get device into registry, the created product
	if err := registry.DeleteDevice(device.ProductId, device.DeviceId); err != nil {
		return reply(ctx, ServerError, &apiResponse{Message: err.Error()})
	}
	orm.DestroyObject(device.DeviceId, accessId)
	return reply(ctx, OK, apiResponse{})
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
		return reply(ctx, BadRequest, &apiResponse{Message: "invalid parameter"})
	}
	// Authorization
	objectName := productId + "/$devices"
	if err := orm.AccessObject(objectName, accessId, orm.AccessRightReadOnly); err != nil {
		return reply(ctx, Unauthorized, apiResponse{Message: "access denied"})
	}

	// Connect with registry
	registry, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, &apiResponse{Message: err.Error()})
	}
	defer registry.Release()

	// Get device into registry, the created product
	dev, err := registry.GetDevice(productId, deviceId)
	if err != nil {
		return reply(ctx, ServerError, &apiResponse{Message: err.Error()})
	}
	return reply(ctx, OK, &apiResponse{Result: dev})
}

// updateDevice update the identity of a device in the identity registry of an IoT Hub
func UpdateDevice(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	device := db.Device{}
	if err := ctx.Bind(&device); err != nil {
		return reply(ctx, BadRequest, &apiResponse{Message: err.Error()})
	}
	if device.ProductId == "" || device.DeviceId == "" {
		return reply(ctx, BadRequest, &apiResponse{Message: "invalid parameter"})
	}
	// Authorization
	objectName := device.ProductId + "/$devices"
	if err := orm.AccessObject(objectName, accessId, orm.AccessRightWrite); err != nil {
		return reply(ctx, Unauthorized, apiResponse{Message: "access denied"})
	}

	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, &apiResponse{Message: err.Error()})
	}
	defer r.Release()

	device.TimeUpdated = time.Now()
	if err = r.UpdateDevice(&device); err != nil {
		return reply(ctx, ServerError, &apiResponse{Message: err.Error()})
	}
	return reply(ctx, OK, &apiResponse{Result: &device})
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
			return ctx.JSON(http.StatusBadRequest, &apiResponse{Success: false, Message: err.Error()})
		}
		config := ctx.(*ApiContext).Config
		// Connect with registry
		r, err := db.NewRegistry("apiserver", config)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &apiResponse{Success: false, Message: err.Error()})
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
			return ctx.JSON(http.StatusInternalServerError, &apiResponse{Success: false, Message: err.Error()})
		}
		return ctx.JSON(http.StatusOK,
			&apiResponse{
				Success: true,
				Result:  rdevices,
			})
	*/
	return nil
}
func GetShadowDevice(ctx echo.Context) error {
	return nil
}

func UpdateShadowDevice(ctx echo.Context) error {
	return nil
}
