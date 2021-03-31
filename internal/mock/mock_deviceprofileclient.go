// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"context"
	"encoding/json"
	"path/filepath"
	"runtime"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/responses"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
)

var (
	DeviceProfileRandomBoolGenerator           = dtos.DeviceProfile{}
	DeviceProfileRandomIntegerGenerator        = dtos.DeviceProfile{}
	DeviceProfileRandomUnsignedGenerator       = dtos.DeviceProfile{}
	DeviceProfileRandomFloatGenerator          = dtos.DeviceProfile{}
	DuplicateDeviceProfileRandomFloatGenerator = dtos.DeviceProfile{}
	NewDeviceProfile                           = dtos.DeviceProfile{}
)

type DeviceProfileClientMock struct{}

func (DeviceProfileClientMock) Add(_ context.Context, _ []requests.DeviceProfileRequest) ([]common.BaseWithIdResponse, errors.EdgeX) {
	panic("implement me")
}

func (DeviceProfileClientMock) Update(_ context.Context, _ []requests.DeviceProfileRequest) ([]common.BaseResponse, errors.EdgeX) {
	panic("implement me")
}

func (DeviceProfileClientMock) AddByYaml(_ context.Context, _ string) (common.BaseWithIdResponse, errors.EdgeX) {
	panic("implement me")
}

func (DeviceProfileClientMock) UpdateByYaml(_ context.Context, _ string) (common.BaseResponse, errors.EdgeX) {
	panic("implement me")
}

func (DeviceProfileClientMock) DeleteByName(_ context.Context, _ string) (common.BaseResponse, errors.EdgeX) {
	panic("implement me")
}

func (DeviceProfileClientMock) DeviceProfileByName(_ context.Context, _ string) (responses.DeviceProfileResponse, errors.EdgeX) {
	panic("implement me")
}

func (DeviceProfileClientMock) AllDeviceProfiles(_ context.Context, _ []string, _ int, _ int) (responses.MultiDeviceProfilesResponse, errors.EdgeX) {
	err := populateDeviceProfileMock()
	if err != nil {
		return responses.MultiDeviceProfilesResponse{
			BaseResponse: common.NewBaseResponse("", "", 0),
		}, errors.NewCommonEdgeXWrapper(err)
	}
	return responses.MultiDeviceProfilesResponse{
		BaseResponse: common.NewBaseResponse("", "", 0),
		Profiles: []dtos.DeviceProfile{
			DeviceProfileRandomBoolGenerator,
			DeviceProfileRandomIntegerGenerator,
			DeviceProfileRandomUnsignedGenerator,
			DeviceProfileRandomFloatGenerator,
		},
	}, nil
}

func (DeviceProfileClientMock) DeviceProfilesByModel(_ context.Context, _ string, _ int, _ int) (responses.MultiDeviceProfilesResponse, errors.EdgeX) {
	panic("implement me")
}

// Query profiles by manufacturer
func (DeviceProfileClientMock) DeviceProfilesByManufacturer(_ context.Context, _ string, _ int, _ int) (responses.MultiDeviceProfilesResponse, errors.EdgeX) {
	panic("implement me")
}

// Query profiles by manufacturer and model
func (DeviceProfileClientMock) DeviceProfilesByManufacturerAndModel(_ context.Context, _ string, _ string, _ int, _ int) (responses.MultiDeviceProfilesResponse, errors.EdgeX) {
	panic("implement me")
}

func (DeviceProfileClientMock) DeviceProfileById(_ context.Context, _ string) (res responses.DeviceProfileResponse, edgexError errors.EdgeX) {
	panic("implement me")
}

func populateDeviceProfileMock() error {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	profiles, err := loadData(basepath + "/data/deviceprofile")
	if err != nil {
		return err
	}
	_ = json.Unmarshal(profiles[ProfileBool], &DeviceProfileRandomBoolGenerator)
	_ = json.Unmarshal(profiles[ProfileInt], &DeviceProfileRandomIntegerGenerator)
	_ = json.Unmarshal(profiles[ProfileUint], &DeviceProfileRandomUnsignedGenerator)
	_ = json.Unmarshal(profiles[ProfileFloat], &DeviceProfileRandomFloatGenerator)
	_ = json.Unmarshal(profiles[ProfileFloat], &DuplicateDeviceProfileRandomFloatGenerator)
	_ = json.Unmarshal(profiles[ProfileNew], &NewDeviceProfile)

	return nil
}
