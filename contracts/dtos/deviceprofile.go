//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

import (
	"fmt"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
)

// DeviceProfile and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/DeviceProfile
type DeviceProfile struct {
	common.Versionable `json:",inline"`
	Id                 string            `json:"id,omitempty" validate:"omitempty,edgex-dto-none-empty-string,edgex-dto-rfc3986-unreserved-chars"`
	Name               string            `json:"name" yaml:"name"`
	DisplayName        string            `json:"displayName" yaml:"displayName" validate:"required"`
	Manufacturer       string            `json:"manufacturer,omitempty" yaml:"manufacturer,omitempty"`
	Description        string            `json:"description,omitempty" yaml:"description,omitempty"`
	Model              string            `json:"model,omitempty" yaml:"model,omitempty"`
	Labels             []string          `json:"labels,omitempty" yaml:"labels,flow,omitempty"`
	DeviceResources    []DeviceResource  `json:"deviceResources,omitempty" yaml:"deviceResources" validate:"gte=0,dive"`
	DeviceCommands     []ProfileResource `json:"deviceCommands,omitempty" yaml:"deviceCommands,omitempty" validate:"dive"`
	CoreCommands       []Command         `json:"coreCommands,omitempty" yaml:"coreCommands,omitempty" validate:"dive"`
	DeviceLibraryId    string            `json:"deviceLibraryId" yaml:"deviceLibraryId,omitempty" validate:"required"`
}

// ToDeviceProfileModel transforms the DeviceProfile DTO to the DeviceProfile model
func ToDeviceProfileModel(deviceProfileDTO DeviceProfile) models.DeviceProfile {
	return models.DeviceProfile{
		Id:              deviceProfileDTO.Id,
		Name:            deviceProfileDTO.Name,
		DisplayName:     deviceProfileDTO.DisplayName,
		Description:     deviceProfileDTO.Description,
		Manufacturer:    deviceProfileDTO.Manufacturer,
		Model:           deviceProfileDTO.Model,
		Labels:          deviceProfileDTO.Labels,
		DeviceResources: ToDeviceResourceModels(deviceProfileDTO.DeviceResources),
		DeviceCommands:  ToProfileResourceModels(deviceProfileDTO.DeviceCommands),
		CoreCommands:    ToCommandModels(deviceProfileDTO.CoreCommands),
		DeviceLibraryId: deviceProfileDTO.DeviceLibraryId,
	}
}

// FromDeviceProfileModelToDTO transforms the DeviceProfile Model to the DeviceProfile DTO
func FromDeviceProfileModelToDTO(deviceProfile models.DeviceProfile) DeviceProfile {
	return DeviceProfile{
		Versionable:     common.NewVersionable(),
		Id:              deviceProfile.Id,
		Name:            deviceProfile.Name,
		DisplayName:     deviceProfile.DisplayName,
		Description:     deviceProfile.Description,
		Manufacturer:    deviceProfile.Manufacturer,
		Model:           deviceProfile.Model,
		Labels:          deviceProfile.Labels,
		DeviceResources: FromDeviceResourceModelsToDTOs(deviceProfile.DeviceResources),
		DeviceCommands:  FromProfileResourceModelsToDTOs(deviceProfile.DeviceCommands),
		CoreCommands:    FromCommandModelsToDTOs(deviceProfile.CoreCommands),
		DeviceLibraryId: deviceProfile.DeviceLibraryId,
	}
}

func ValidateDeviceProfileDTO(profile DeviceProfile) error {
	// deviceResources validation
	dupCheck := make(map[string]bool)
	for _, resource := range profile.DeviceResources {
		// deviceResource name should not duplicated
		if dupCheck[resource.Name] {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("device resource %s is duplicated", resource.Name), nil)
		}
		dupCheck[resource.Name] = true
	}
	// deviceCommands validation
	dupCheck = make(map[string]bool)
	for _, command := range profile.DeviceCommands {
		// deviceCommand name should not duplicated
		if dupCheck[command.Name] {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("device command %s is duplicated", command.Name), nil)
		}
		dupCheck[command.Name] = true

		// deviceResources referenced in deviceCommands must exist
		getCommands := command.Get
		for _, getCommand := range getCommands {
			if !deviceResourcesContains(profile.DeviceResources, getCommand.DeviceResource) {
				return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("device command's Get resource %s doesn't match any deivce resource", getCommand.DeviceResource), nil)
			}
		}
		setCommands := command.Set
		for _, setCommand := range setCommands {
			if !deviceResourcesContains(profile.DeviceResources, setCommand.DeviceResource) {
				return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("device command's Set resource %s doesn't match any deivce resource", setCommand.DeviceResource), nil)
			}
		}
	}
	// coreCommands validation
	dupCheck = make(map[string]bool)
	for _, command := range profile.CoreCommands {
		// coreCommand name should not duplicated
		if dupCheck[command.Name] {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("core command %s is duplicated", command.Name), nil)
		}
		dupCheck[command.Name] = true

		// coreCommands name should match the one of deviceResources and deviceCommands
		if !deviceCommandsContains(profile.DeviceCommands, command.Name) &&
			!deviceResourcesContains(profile.DeviceResources, command.Name) {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("core command %s doesn't match any deivce command or resource", command.Name), nil)
		}
	}
	return nil
}

func deviceResourcesContains(resources []DeviceResource, name string) bool {
	contains := false
	for _, resource := range resources {
		if resource.Name == name {
			contains = true
			break
		}
	}
	return contains
}

func deviceCommandsContains(resources []ProfileResource, name string) bool {
	contains := false
	for _, resource := range resources {
		if resource.Name == name {
			contains = true
			break
		}
	}
	return contains
}
