// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/responses"
	eErr "github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
)

const (
	InvalidDeviceId = "1ef435eb-5060-49b0-8d55-8d4e43239800"
)

var (
	ValidDeviceRandomBoolGenerator            = dtos.Device{}
	ValidDeviceRandomIntegerGenerator         = dtos.Device{}
	ValidDeviceRandomUnsignedIntegerGenerator = dtos.Device{}
	ValidDeviceRandomFloatGenerator           = dtos.Device{}
	DuplicateDeviceRandomFloatGenerator       = dtos.Device{}
	NewValidDevice                            = dtos.Device{}
	OperatingStateDisabled                    = dtos.Device{}
)

type DeviceClientMock struct{}

func (dc *DeviceClientMock) Add(ctx context.Context, reqs []requests.AddDeviceRequest) ([]common.BaseWithIdResponse, eErr.EdgeX) {
	panic("implement me")
}

func (dc *DeviceClientMock) Update(ctx context.Context, reqs []requests.UpdateDeviceRequest) ([]common.BaseResponse, eErr.EdgeX) {
	panic("implement me")
}

func (dc *DeviceClientMock) AllDevices(ctx context.Context, labels []string, offset int, limit int) (responses.MultiDevicesResponse, eErr.EdgeX) {
	panic("implement me")
}

func (dc *DeviceClientMock) DeviceNameExists(ctx context.Context, name string) (common.BaseResponse, eErr.EdgeX) {
	panic("implement me")
}

func (dc *DeviceClientMock) DeviceByName(ctx context.Context, name string) (responses.DeviceResponse, eErr.EdgeX) {
	var device = dtos.Device{Id: "5b977c62f37ba10e36673802", Name: name}
	var err eErr.EdgeX = nil
	if name == "" {
		err = eErr.NewCommonEdgeXWrapper(fmt.Errorf("item not found"))
	}

	return responses.DeviceResponse{
		BaseResponse: common.NewBaseResponse("", "", 0),
		Device:       device,
	}, err
}

func (dc *DeviceClientMock) DeleteDeviceByName(ctx context.Context, name string) (common.BaseResponse, eErr.EdgeX) {
	panic("implement me")
}

func (dc *DeviceClientMock) DevicesByProfileName(ctx context.Context, name string, offset int, limit int) (responses.MultiDevicesResponse, eErr.EdgeX) {
	panic("implement me")
}

func (dc *DeviceClientMock) DevicesByServiceName(ctx context.Context, name string, offset int, limit int) (responses.MultiDevicesResponse, eErr.EdgeX) {
	err := populateDeviceMock()
	if err != nil {
		return responses.MultiDevicesResponse{}, eErr.NewCommonEdgeXWrapper(err)
	}
	return responses.MultiDevicesResponse{
		BaseResponse: common.NewBaseResponse("", "", 0),
		Devices: []dtos.Device{
			ValidDeviceRandomBoolGenerator,
			ValidDeviceRandomIntegerGenerator,
			ValidDeviceRandomUnsignedIntegerGenerator,
			ValidDeviceRandomFloatGenerator,
			OperatingStateDisabled,
		}}, nil
}

func (dc *DeviceClientMock) DeviceById(ctx context.Context, id string) (responses.DeviceResponse, eErr.EdgeX) {
	if id == InvalidDeviceId {
		return responses.DeviceResponse{}, eErr.NewCommonEdgeXWrapper(fmt.Errorf("invalid id"))
	}
	return responses.DeviceResponse{
		BaseResponse: common.NewBaseResponse("", "", 0),
		Device:       dtos.Device{},
	}, nil
}
func (dc *DeviceClientMock) DeleteDeviceById(ctx context.Context, id string) (common.BaseResponse, eErr.EdgeX) {
	panic("implement me")
}

func populateDeviceMock() error {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	devices, err := loadData(basepath + "/data/device")
	if err != nil {
		return err
	}
	profiles, err := loadData(basepath + "/data/deviceprofile")
	if err != nil {
		return err
	}
	_ = json.Unmarshal(devices[DeviceBool], &ValidDeviceRandomBoolGenerator)
	_ = json.Unmarshal(profiles[DeviceBool], &ValidDeviceRandomBoolGenerator.ProfileName)
	_ = json.Unmarshal(devices[DeviceInt], &ValidDeviceRandomIntegerGenerator)
	_ = json.Unmarshal(profiles[DeviceInt], &ValidDeviceRandomIntegerGenerator.ProfileName)
	_ = json.Unmarshal(devices[DeviceUint], &ValidDeviceRandomUnsignedIntegerGenerator)
	_ = json.Unmarshal(profiles[DeviceUint], &ValidDeviceRandomUnsignedIntegerGenerator.ProfileName)
	_ = json.Unmarshal(devices[DeviceFloat], &ValidDeviceRandomFloatGenerator)
	_ = json.Unmarshal(profiles[DeviceFloat], &ValidDeviceRandomFloatGenerator.ProfileName)
	_ = json.Unmarshal(devices[DeviceFloat], &DuplicateDeviceRandomFloatGenerator)
	_ = json.Unmarshal(profiles[DeviceFloat], &DuplicateDeviceRandomFloatGenerator.ProfileName)
	_ = json.Unmarshal(devices[DeviceNew], &NewValidDevice)
	_ = json.Unmarshal(profiles[DeviceNew], &NewValidDevice.ProfileName)
	_ = json.Unmarshal(devices[DeviceNew02], &OperatingStateDisabled)
	_ = json.Unmarshal(profiles[DeviceNew], &OperatingStateDisabled.ProfileName)

	return nil
}
