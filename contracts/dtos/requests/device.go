//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package requests

import (
	"encoding/json"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
)

// AddDeviceRequest defines the Request Content for POST Device DTO.
// This object and its properties correspond to the AddDeviceRequest object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/AddDeviceRequest
type AddDeviceRequest struct {
	common.BaseRequest `json:",inline"`
	Device             dtos.Device `json:"device"`
}

// Validate satisfies the Validator interface
func (d AddDeviceRequest) Validate() error {
	err := contracts.Validate(d)
	return err
}

// UnmarshalJSON implements the Unmarshaler interface for the AddDeviceRequest type
func (d *AddDeviceRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		common.BaseRequest
		Device dtos.Device
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*d = AddDeviceRequest(alias)

	// validate AddDeviceRequest DTO
	if err := d.Validate(); err != nil {
		return err
	}
	return nil
}

// AddDeviceReqToDeviceModels transforms the AddDeviceRequest DTO array to the Device model array
func AddDeviceReqToDeviceModels(addRequests []AddDeviceRequest) (Devices []models.Device) {
	for _, req := range addRequests {
		d := dtos.ToDeviceModel(req.Device)
		Devices = append(Devices, d)
	}
	return Devices
}

// UpdateDeviceRequest defines the Request Content for PUT event as pushed DTO.
// This object and its properties correspond to the UpdateDeviceRequest object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/UpdateDeviceRequest
type UpdateDeviceRequest struct {
	common.BaseRequest `json:",inline"`
	Device             dtos.UpdateDevice `json:"device"`
}

// Validate satisfies the Validator interface
func (d UpdateDeviceRequest) Validate() error {
	err := contracts.Validate(d)
	return err
}

// UnmarshalJSON implements the Unmarshaler interface for the UpdateDeviceRequest type
func (d *UpdateDeviceRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		common.BaseRequest
		Device dtos.UpdateDevice
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*d = UpdateDeviceRequest(alias)

	// validate UpdateDeviceRequest DTO
	if err := d.Validate(); err != nil {
		return err
	}
	return nil
}

// ReplaceDeviceModelFieldsWithDTO replace existing Device's fields with DTO patch
func ReplaceDeviceModelFieldsWithDTO(device *models.Device, patch dtos.UpdateDevice) {
	if patch.Description != nil {
		device.Description = *patch.Description
	}
	if patch.AdminState != nil {
		device.AdminState = models.AdminState(*patch.AdminState)
	}
	if patch.OperatingState != nil {
		device.OperatingState = models.OperatingState(*patch.OperatingState)
	}
	if patch.LastConnected != nil {
		device.LastConnected = *patch.LastConnected
	}
	if patch.LastReported != nil {
		device.LastReported = *patch.LastReported
	}
	if patch.ServiceName != nil {
		device.ServiceName = *patch.ServiceName
	}
	if patch.ProfileName != nil {
		device.ProfileName = *patch.ProfileName
	}
	if patch.Labels != nil {
		device.Labels = patch.Labels
	}
	if patch.Location != nil {
		device.Location = patch.Location
	}
	if patch.AutoEvents != nil {
		device.AutoEvents = dtos.ToAutoEventModels(patch.AutoEvents)
	}
	if patch.Protocols != nil {
		device.Protocols = dtos.ToProtocolModels(patch.Protocols)
	}
	if patch.Notify != nil {
		device.Notify = *patch.Notify
	}
	if patch.DisplayName != nil {
		device.DisplayName = *patch.DisplayName
	}
	if patch.ActiveStatus != nil {
		device.ActiveStatus = *patch.ActiveStatus
	}
}

type ActiveDeviceRequest struct {
	common.BaseRequest `json:",inline"`
	ActiveDevice       dtos.ActiveDevice `json:"activeDevice"`
}

// Validate satisfies the Validator interface
func (d ActiveDeviceRequest) Validate() error {
	err := contracts.Validate(d)
	return err
}

// UnmarshalJSON implements the Unmarshaler interface for the AuthActiveDeviceRequest type
func (d *ActiveDeviceRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		common.BaseRequest
		ActiveDevice dtos.ActiveDevice
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*d = ActiveDeviceRequest(alias)

	if err := d.Validate(); err != nil {
		return err
	}
	return nil
}

type DeviceSearchQueryRequest struct {
	common.BaseSearchConditionQuery `schema:",inline"`
	ActiveStatusArray               string `schema:"active_status_array"`
	ServiceId                       string `schema:"service_id"`
	ProfileId                       string `schema:"profile_id"`
}
