//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package responses

import (
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
)

// DeviceResponse defines the Response Content for GET Device DTOs.
// This object and its properties correspond to the DeviceResponse object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/DeviceResponse
type DeviceResponse struct {
	common.BaseResponse `json:",inline"`
	Device              dtos.Device `json:"device"`
}

func NewDeviceResponse(requestId string, message string, statusCode int, device dtos.Device) DeviceResponse {
	return DeviceResponse{
		BaseResponse: common.NewBaseResponse(requestId, message, statusCode),
		Device:       device,
	}
}

// MultiDevicesResponse defines the Response Content for GET multiple Device DTOs.
// This object and its properties correspond to the MultiDevicesResponse object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/MultiDevicesResponse
type MultiDevicesResponse struct {
	common.BaseResponse `json:",inline"`
	Total               uint32        `json:"total"`
	Devices             []dtos.Device `json:"devices"`
}

func NewMultiDevicesResponse(requestId string, message string, statusCode int, devices []dtos.Device, total uint32) MultiDevicesResponse {
	return MultiDevicesResponse{
		BaseResponse: common.NewBaseResponse(requestId, message, statusCode),
		Total:        total,
		Devices:      devices,
	}
}

type MultiActiveDevicesResponse struct {
	common.BaseResponse `json:",inline"`
	ActiveDeviceResult  []dtos.ActiveDeviceResult `json:"activeDeviceResult"`
	ProcessNum          int                       `json:"processNum"`
	SuccessNum          int                       `json:"successNum"`
	FailNum             int                       `json:"failNum"`
}

func NewMultiActiveDevicesResponse(requestId string, message string, statusCode int, ActiveDeviceResult []dtos.ActiveDeviceResult, processNum int, successNum int, failNum int) MultiActiveDevicesResponse {
	return MultiActiveDevicesResponse{
		BaseResponse:       common.NewBaseResponse(requestId, message, statusCode),
		ActiveDeviceResult: ActiveDeviceResult,
		ProcessNum:         processNum,
		SuccessNum:         successNum,
		FailNum:            failNum,
	}
}
