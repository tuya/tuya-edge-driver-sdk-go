//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

import (
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
)

// DeviceService and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/DeviceService
type DeviceService struct {
	common.Versionable `json:",inline"`
	Id                 string                 `json:"id,omitempty" validate:"omitempty,uuid"`
	Name               string                 `json:"name"`
	Created            int64                  `json:"created,omitempty"`
	Modified           int64                  `json:"modified,omitempty"`
	Description        string                 `json:"description,omitempty"`
	LastConnected      int64                  `json:"lastConnected,omitempty"`
	LastReported       int64                  `json:"lastReported,omitempty"`
	Labels             []string               `json:"labels,omitempty"`
	BaseAddress        string                 `json:"baseAddress" validate:"required,uri"`
	AdminState         string                 `json:"adminState" validate:"oneof='LOCKED' 'UNLOCKED'"`
	ServiceName        string                 `json:"serviceName"`
	DeviceLibraryId    string                 `json:"deviceLibraryId" validate:"required"`
	Config             map[string]interface{} `json:"config" validate:"required"`
}

// UpdateDeviceService and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/UpdateDeviceService
type UpdateDeviceService struct {
	Id                *string  `json:"id" validate:"required_without=Name,edgex-dto-uuid"`
	Name              *string  `json:"name" validate:"required_without=Id,edgex-dto-none-empty-string,edgex-dto-rfc3986-unreserved-chars"`
	DeviceLibraryId   *string  `json:"deviceLibraryId"`
	Description       *string  `json:"description"`
	ServiceName       *string  `json:"serviceName"`
	BaseAddress       *string  `json:"baseAddress" validate:"omitempty,uri"`
	Labels            []string `json:"labels"`
	AdminState        *string  `json:"adminState" validate:"omitempty,oneof='LOCKED' 'UNLOCKED'"`
	DockerContainerId *string  `json:"dockerContainerId"`
}

// ToDeviceServiceModel transforms the DeviceService DTO to the DeviceService Model
func ToDeviceServiceModel(dto DeviceService) models.DeviceService {
	var ds models.DeviceService
	ds.Id = dto.Id
	ds.Name = dto.Name
	ds.Description = dto.Description
	ds.LastReported = dto.LastReported
	ds.LastConnected = dto.LastConnected
	ds.BaseAddress = dto.BaseAddress
	ds.Labels = dto.Labels
	ds.AdminState = models.AdminState(dto.AdminState)
	ds.ServiceName = dto.ServiceName
	ds.DeviceLibraryId = dto.DeviceLibraryId
	return ds
}

// FromDeviceServiceModelToDTO transforms the DeviceService Model to the DeviceService DTO
func FromDeviceServiceModelToDTO(ds models.DeviceService) DeviceService {
	var dto DeviceService
	dto.Versionable = common.NewVersionable()
	dto.Id = ds.Id
	dto.Name = ds.Name
	dto.Description = ds.Description
	dto.LastReported = ds.LastReported
	dto.LastConnected = ds.LastConnected
	dto.BaseAddress = ds.BaseAddress
	dto.Labels = ds.Labels
	dto.AdminState = string(ds.AdminState)
	dto.ServiceName = ds.ServiceName
	dto.DeviceLibraryId = ds.DeviceLibraryId
	return dto
}

// FromDeviceServiceModelToUpdateDTO transforms the DeviceService Model to the UpdateDeviceService DTO
func FromDeviceServiceModelToUpdateDTO(ds models.DeviceService) UpdateDeviceService {
	adminState := string(ds.AdminState)
	return UpdateDeviceService{
		Id:          &ds.Id,
		Name:        &ds.Name,
		BaseAddress: &ds.BaseAddress,
		Labels:      ds.Labels,
		AdminState:  &adminState,
	}
}
