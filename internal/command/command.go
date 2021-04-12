// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/interfaces"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/responses"
	edgexErr "github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/cache"
	sdkCommon "github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/container"
	context2 "github.com/tuya/tuya-edge-driver-sdk-go/internal/context"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/transformer"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
	dsModels "github.com/tuya/tuya-edge-driver-sdk-go/pkg/models"
)

type CommandProcessor struct {
	device         *models.Device
	deviceResource *models.DeviceResource
	correlationID  string
	cmd            string
	params         string
	dic            *di.Container
}

func NewCommandProcessor(device *models.Device, dr *models.DeviceResource, correlationID string, cmd string, params string, dic *di.Container) *CommandProcessor {
	return &CommandProcessor{
		device:         device,
		deviceResource: dr,
		correlationID:  correlationID,
		cmd:            cmd,
		params:         params,
		dic:            dic,
	}
}

func CommandHandler(isRead bool, sendEvent bool, correlationID string, vars map[string]string, body string, dic *di.Container) (res responses.EventResponse, err edgexErr.EdgeX) {
	var device models.Device
	deviceKey := vars[sdkCommon.NameVar]
	// the device service will perform some operations(e.g. update LastConnected timestamp,
	// push returning event to core-data) after a device is successfully interacted with if
	// it has been configured to do so, and those operation apply to every protocol and
	// need to be finished in the end of command layer before returning to protocol layer.
	defer func() {
		if err != nil {
			return
		}
		go sdkCommon.UpdateLastConnected(
			device,
			container.ConfigurationFrom(dic.Get),
			bootstrapContainer.LoggingClientFrom(dic.Get),
			container.MetadataDeviceClientFrom(dic.Get))

		if sendEvent {
			ec := container.CoredataEventClientFrom(dic.Get)
			lc := bootstrapContainer.LoggingClientFrom(dic.Get)
			go SendEvent(res, correlationID, lc, ec)
		}
	}()

	// check device service's AdminState
	ds := container.DeviceServiceFrom(dic.Get)
	if ds.AdminState == models.Locked {
		res = responses.NewEventResponse(correlationID, "device service locked", http.StatusInternalServerError, dtos.Event{})
		return res, edgexErr.NewCommonEdgeX(edgexErr.KindServiceLocked, "service locked", nil)
	}

	// check provided device exists
	device, exist := cache.Devices().ForName(deviceKey)
	if !exist {
		res = responses.NewEventResponse(correlationID, "device not found in local cache", http.StatusBadRequest, dtos.Event{})
		return res, edgexErr.NewCommonEdgeX(edgexErr.KindEntityDoesNotExist, fmt.Sprintf("device %s not found", deviceKey), nil)
	}

	// check device's AdminState
	if device.AdminState == models.Locked {
		res = responses.NewEventResponse(correlationID, "device locked", http.StatusInternalServerError, dtos.Event{})
		return res, edgexErr.NewCommonEdgeX(edgexErr.KindServiceLocked, fmt.Sprintf("device %s locked", device.Name), nil)
	}

	var method string
	if isRead {
		method = sdkCommon.GetCmdMethod
	} else {
		method = sdkCommon.SetCmdMethod
	}
	cmd := vars[sdkCommon.CommandVar]
	cmdExists, e := cache.Profiles().CommandExists(device.ProfileName, cmd, method)
	if e != nil {
		errMsg := fmt.Sprintf("failed to identify command %s in cache", cmd)
		res = responses.NewEventResponse(correlationID, errMsg, http.StatusBadRequest, dtos.Event{})
		return res, edgexErr.NewCommonEdgeX(edgexErr.KindServerError, errMsg, e)
	}
	helper := NewCommandProcessor(&device, nil, correlationID, cmd, body, dic)
	if cmdExists {
		if isRead {
			return helper.ReadCommand()
		} else {
			if err = helper.WriteCommand(); err != nil {
				res = responses.NewEventResponse(correlationID, err.Message(), http.StatusBadRequest, dtos.Event{})
			} else {
				res = responses.NewEventResponse(correlationID, "success", http.StatusOK, dtos.Event{})
			}
			return res, err
		}
	} else {
		dr, drExists := cache.Profiles().DeviceResource(device.ProfileName, cmd)
		if !drExists {
			return res, edgexErr.NewCommonEdgeX(edgexErr.KindEntityDoesNotExist, "command not found", nil)
		}

		helper = NewCommandProcessor(&device, &dr, correlationID, cmd, body, dic)
		if isRead {
			return helper.ReadDeviceResource()
		} else {
			if err = helper.WriteDeviceResource(); err != nil {
				res = responses.NewEventResponse(correlationID, err.Message(), http.StatusBadRequest, dtos.Event{})
			} else {
				res = responses.NewEventResponse(correlationID, "success", http.StatusOK, dtos.Event{})
			}
			return res, err
		}
	}
}

func (c *CommandProcessor) ReadDeviceResource() (res responses.EventResponse, e edgexErr.EdgeX) {
	lc := bootstrapContainer.LoggingClientFrom(c.dic.Get)
	lc.Debug(fmt.Sprintf("Application - readDeviceResource: reading deviceResource: %s", c.deviceResource.Name), sdkCommon.CorrelationHeader, c.correlationID)

	// check provided deviceResource is not write-only
	if c.deviceResource.Properties.ReadWrite == sdkCommon.DeviceResourceWriteOnly {
		errMsg := fmt.Sprintf("deviceResource %s is marked as write-only", c.deviceResource.Name)
		return res, edgexErr.NewCommonEdgeX(edgexErr.KindNotAllowed, errMsg, nil)
	}

	var req dsModels.CommandRequest
	var reqs []dsModels.CommandRequest

	// prepare CommandRequest
	req.DeviceResourceName = c.deviceResource.Name
	req.Attributes = c.deviceResource.Attributes
	if c.params != "" {
		if len(req.Attributes) <= 0 {
			req.Attributes = make(map[string]string)
		}
		req.Attributes[sdkCommon.URLRawQuery] = c.params
	}
	req.Type = c.deviceResource.Properties.Type
	reqs = append(reqs, req)

	// execute protocol-specific read operation
	driver := container.ProtocolDriverFrom(c.dic.Get)
	results, err := driver.HandleReadCommands(c.device.Name, c.device.Protocols, reqs)
	if err != nil {
		errMsg := fmt.Sprintf("error reading DeviceResourece %s for %s: %v", c.deviceResource.Name, c.device.Name, err)
		return res, edgexErr.NewCommonEdgeX(edgexErr.KindServerError, errMsg, err)
	}

	// convert CommandValue to Event
	event, err := c.commandValuesToEvent(results, c.deviceResource.Name)
	if err != nil {
		return res, edgexErr.NewCommonEdgeX(edgexErr.KindServerError, "failed to convert CommandValue to Event", err)
	}

	res = responses.NewEventResponse(c.correlationID, "", http.StatusOK, event)
	return
}

func (c *CommandProcessor) ReadCommand() (res responses.EventResponse, e edgexErr.EdgeX) {
	lc := bootstrapContainer.LoggingClientFrom(c.dic.Get)
	lc.Debug(fmt.Sprintf("Application - readCmd: reading cmd: %s", c.cmd), sdkCommon.CorrelationHeader, c.correlationID)

	// check GET ResourceOperation(s) exist for provided command
	ros, err := cache.Profiles().ResourceOperations(c.device.ProfileName, c.cmd, sdkCommon.GetCmdMethod)
	if err != nil {
		errMsg := fmt.Sprintf("GET ResourceOperation(s) for %s command not found", c.cmd)
		return res, edgexErr.NewCommonEdgeX(edgexErr.KindNotAllowed, errMsg, err)
	}
	// check ResourceOperation count does not exceed MaxCmdOps defined in configuration
	configuration := container.ConfigurationFrom(c.dic.Get)
	if len(ros) > configuration.Device.MaxCmdOps {
		errMsg := fmt.Sprintf("GET command %s exceed device %s MaxCmdOps (%d)", c.cmd, c.device.Name, configuration.Device.MaxCmdOps)
		return res, edgexErr.NewCommonEdgeX(edgexErr.KindServerError, errMsg, nil)
	}

	// prepare CommandRequests
	reqs := make([]dsModels.CommandRequest, len(ros))
	for i, op := range ros {
		drName := op.DeviceResource
		// check the deviceResource in ResourceOperation actually exist
		dr, ok := cache.Profiles().DeviceResource(c.device.ProfileName, drName)
		if !ok {
			errMsg := fmt.Sprintf("deviceResource %s in GET commnd %s for %s not defined", drName, c.cmd, c.device.Name)
			return res, edgexErr.NewCommonEdgeX(edgexErr.KindServerError, errMsg, nil)
		}

		// check the deviceResource isn't write-only
		if dr.Properties.ReadWrite == sdkCommon.DeviceResourceWriteOnly {
			errMsg := fmt.Sprintf("deviceResource %s in GET command %s is marked as write-only", drName, c.cmd)
			return res, edgexErr.NewCommonEdgeX(edgexErr.KindNotAllowed, errMsg, nil)
		}
		reqs[i].DeviceResourceName = dr.Name
		reqs[i].Attributes = dr.Attributes
		if c.params != "" {
			if len(reqs[i].Attributes) <= 0 {
				reqs[i].Attributes = make(map[string]string)
			}
			reqs[i].Attributes[sdkCommon.URLRawQuery] = c.params
		}
		reqs[i].Type = dr.Properties.Type
	}

	// execute protocol-specific read operation
	driver := container.ProtocolDriverFrom(c.dic.Get)
	results, eerr := driver.HandleReadCommands(c.device.Name, c.device.Protocols, reqs)
	if eerr != nil {
		errMsg := fmt.Sprintf("error reading DeviceCommand %s for %s: %v", c.cmd, c.device.Name, err)
		return res, edgexErr.NewCommonEdgeX(edgexErr.KindServerError, errMsg, eerr)
	}

	// convert CommandValue to Event
	event, err := c.commandValuesToEvent(results, c.cmd)
	if err != nil {
		return res, edgexErr.NewCommonEdgeX(edgexErr.KindServerError, "failed to transform CommandValue to Event", err)
	}

	res = responses.NewEventResponse(c.correlationID, "", http.StatusOK, event)
	return
}

func (c *CommandProcessor) WriteDeviceResource() edgexErr.EdgeX {
	lc := bootstrapContainer.LoggingClientFrom(c.dic.Get)
	lc.Debug(fmt.Sprintf("Application - writeDeviceResource: writting deviceResource: %s", c.deviceResource.Name), sdkCommon.CorrelationHeader, c.correlationID)

	// check provided deviceResource is not read-only
	if c.deviceResource.Properties.ReadWrite == sdkCommon.DeviceResourceReadOnly {
		errMsg := fmt.Sprintf("deviceResource %s is marked as read-only", c.deviceResource.Name)
		return edgexErr.NewCommonEdgeX(edgexErr.KindNotAllowed, errMsg, nil)
	}

	// parse request body string
	paramMap, err := parseParams(c.params)
	if err != nil {
		return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, "failed to parse PUT parameters", err)
	}

	// check request body contains provided deviceResource
	v, ok := paramMap[c.deviceResource.Name]
	if !ok {
		if c.deviceResource.Properties.DefaultValue != "" {
			v = c.deviceResource.Properties.DefaultValue
		} else {
			errMsg := fmt.Sprintf("deviceResource %s not found in request body and no default value defined", c.deviceResource.Name)
			return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, errMsg, nil)
		}
	}

	valueStr, e := json.Marshal(v)
	if e != nil {
		errMsg := fmt.Sprintf("device resource %s value marshal error: %s", c.deviceResource.Name, e)
		return edgexErr.NewCommonEdgeX(edgexErr.KindNotAllowed, errMsg, nil)
	}
	// create CommandValue
	cv, err := createCommandValueFromDeviceResource(c.deviceResource, string(valueStr))
	if err != nil {
		return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, "failed to create CommandValue", err)
	}

	// prepare CommandRequest
	reqs := make([]dsModels.CommandRequest, 1)
	reqs[0].DeviceResourceName = cv.DeviceResourceName
	reqs[0].Attributes = c.deviceResource.Attributes
	reqs[0].Type = cv.Type

	// transform write value
	configuration := container.ConfigurationFrom(c.dic.Get)
	if configuration.Device.DataTransform {
		err = transformer.TransformWriteParameter(cv, c.deviceResource.Properties, lc)
		if err != nil {
			return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, "failed to transform write value", nil)
		}
	}

	// execute protocol-specific write operation
	driver := container.ProtocolDriverFrom(c.dic.Get)
	err = driver.HandleWriteCommands(c.device.Name, c.device.Protocols, reqs, []*dsModels.CommandValue{cv})
	if err != nil {
		errMsg := fmt.Sprintf("error writing DeviceResourece %s for %s: %v", c.deviceResource.Name, c.device.Name, err)
		return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, errMsg, err)
	}

	return nil
}

func (c *CommandProcessor) WriteCommand() edgexErr.EdgeX {
	lc := bootstrapContainer.LoggingClientFrom(c.dic.Get)
	lc.Debug(fmt.Sprintf("Application - writeCmd: writting command: %s", c.cmd), sdkCommon.CorrelationHeader, c.correlationID)

	// check SET ResourceOperation(s) exist for provided command
	ros, eErr := cache.Profiles().ResourceOperations(c.device.ProfileName, c.cmd, sdkCommon.SetCmdMethod)
	if eErr != nil {
		errMsg := fmt.Sprintf("SET ResourceOperation(s) for %s command not found", c.cmd)
		return edgexErr.NewCommonEdgeX(edgexErr.KindNotAllowed, errMsg, eErr)
	}

	// check ResourceOperation count does not exceed MaxCmdOps defined in configuration
	configuration := container.ConfigurationFrom(c.dic.Get)
	if len(ros) > configuration.Device.MaxCmdOps {
		errMsg := fmt.Sprintf("PUT command %s exceed device %s MaxCmdOps (%d)", c.cmd, c.device.Name, configuration.Device.MaxCmdOps)
		return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, errMsg, nil)
	}

	// parse request body
	paramMap, err := parseParams(c.params)
	if err != nil {
		return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, "failed to parse PUT parameters", err)
	}

	// create CommandValues
	cvs := make([]*dsModels.CommandValue, 0, len(paramMap))
	for _, ro := range ros {
		drName := ro.DeviceResource
		// check the deviceResource in ResourceOperation actually exist
		dr, ok := cache.Profiles().DeviceResource(c.device.ProfileName, drName)
		if !ok {
			errMsg := fmt.Sprintf("deviceResource %s in PUT commnd %s for %s not defined", drName, c.cmd, c.device.Name)
			return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, errMsg, nil)
		}

		// check the deviceResource isn't read-only
		if dr.Properties.ReadWrite == sdkCommon.DeviceResourceReadOnly {
			errMsg := fmt.Sprintf("deviceResource %s in PUT command %s is marked as read-only", drName, c.cmd)
			return edgexErr.NewCommonEdgeX(edgexErr.KindNotAllowed, errMsg, nil)
		}

		// check request body contains the deviceResource
		value, ok := paramMap[ro.DeviceResource]
		if !ok {
			if ro.Parameter != "" {
				value = ro.Parameter
			} else if dr.Properties.DefaultValue != "" {
				value = dr.Properties.DefaultValue
			} else {
				errMsg := fmt.Sprintf("deviceResource %s not found in request body and no default value defined", dr.Name)
				return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, errMsg, nil)
			}
		}

		valueStr, err := json.Marshal(value)
		if err != nil {
			errMsg := fmt.Sprintf("dp %s value marshal error: %s", ro.DeviceResource, err)
			return edgexErr.NewCommonEdgeX(edgexErr.KindNotAllowed, errMsg, nil)
		}
		// write value mapping
		if len(ro.Mappings) > 0 {
			newValue, ok := ro.Mappings[string(valueStr)]
			if ok {
				value = newValue
			} else {
				lc.Warn(fmt.Sprintf("ResourceOperation %s mapping value (%s) failed with the mapping table: %v", ro.DeviceResource, value, ro.Mappings))
			}
		}

		// create CommandValue
		cv, err := createCommandValueFromDeviceResource(&dr, string(valueStr))
		if err == nil {
			cvs = append(cvs, cv)
		} else {
			return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, "failed to create CommandValue", err)
		}
	}

	// prepare CommandRequests
	reqs := make([]dsModels.CommandRequest, len(cvs))
	for i, cv := range cvs {
		dr, _ := cache.Profiles().DeviceResource(c.device.ProfileName, cv.DeviceResourceName)

		reqs[i].DeviceResourceName = cv.DeviceResourceName
		reqs[i].Attributes = dr.Attributes
		reqs[i].Type = cv.Type

		// transform write value
		if configuration.Device.DataTransform {
			err = transformer.TransformWriteParameter(cv, dr.Properties, lc)
			if err != nil {
				return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, "failed to transform write values", err)
			}
		}
	}

	// execute protocol-specific write operation
	driver := container.ProtocolDriverFrom(c.dic.Get)
	err = driver.HandleWriteCommands(c.device.Name, c.device.Protocols, reqs, cvs)
	if err != nil {
		errMsg := fmt.Sprintf("error writing DeviceResourece for %s: %v", c.device.Name, err)
		return edgexErr.NewCommonEdgeX(edgexErr.KindServerError, errMsg, err)
	}

	return nil
}

func (c *CommandProcessor) commandValuesToEvent(cvs []*dsModels.CommandValue, cmd string) (dtos.Event, edgexErr.EdgeX) {
	var err error
	var transformsOK = true
	lc := bootstrapContainer.LoggingClientFrom(c.dic.Get)

	configuration := container.ConfigurationFrom(c.dic.Get)
	readings := make([]dtos.BaseReading, 0, configuration.Device.MaxCmdOps)

	for _, cv := range cvs {
		// double check the CommandValue return from ProtocolDriver match device command
		dr, ok := cache.Profiles().DeviceResource(c.device.ProfileName, cv.DeviceResourceName)
		if !ok {
			return dtos.Event{}, edgexErr.NewCommonEdgeXWrapper(fmt.Errorf("no deviceResource %s for %s in CommandValue (%s)", cv.DeviceResourceName, c.device.Name, cv.String()))
		}

		// perform data transformation
		if configuration.Device.DataTransform {
			err = transformer.TransformReadResult(cv, dr.Properties, lc)
			lc.Debug(fmt.Sprintf("command value: %+v", cv))
			if err != nil {
				lc.Error(fmt.Sprintf("failed to transform CommandValue (%s): %v", cv.String(), err), sdkCommon.CorrelationHeader, c.correlationID)

				if errors.As(err, &transformer.OverflowError{}) {
					cv = dsModels.NewStringValue(cv.DeviceResourceName, cv.Origin, transformer.Overflow)
				} else if errors.As(err, &transformer.NaNError{}) {
					cv = dsModels.NewStringValue(cv.DeviceResourceName, cv.Origin, transformer.NaN)
				} else {
					transformsOK = false
				}
			}
		}

		// assertion
		dc := container.MetadataDeviceClientFrom(c.dic.Get)
		err = transformer.CheckAssertion(cv, dr.Properties.Assertion, c.device, lc, dc)
		if err != nil {
			cv = dsModels.NewStringValue(cv.DeviceResourceName, cv.Origin, fmt.Sprintf("Assertion failed for device resource: %s, with value: %s", cv.DeviceResourceName, cv.String()))
		}
		// ResourceOperation mapping
		ro, exrr := cache.Profiles().ResourceOperation(c.device.ProfileName, cv.DeviceResourceName, sdkCommon.GetCmdMethod)
		if exrr != nil {
			// this allows SDK to directly read deviceResource without deviceCommands defined.
			lc.Debug(fmt.Sprintf("failed to read ResourceOperation: %v", exrr), sdkCommon.CorrelationHeader, c.correlationID)
		} else if len(ro.Mappings) > 0 {
			newCV, ok := transformer.MapCommandValue(cv, ro.Mappings)
			if ok {
				cv = newCV
			} else {
				lc.Warn(fmt.Sprintf("ResourceOperation (%s) mapping value (%s) failed with the mapping table: %v", ro.DeviceResource, cv.String(), ro.Mappings), sdkCommon.CorrelationHeader, c.correlationID)
			}
		}

		lc.Debug(fmt.Sprintf("command value: %+v", cv))

		reading := commandValueToReading(cv, c.device.Name, c.device.ProfileName, dr.Properties.MediaType, "")
		readings = append(readings, reading)

		if cv.Type == contracts.ValueTypeBinary {
			lc.Debug(fmt.Sprintf("device: %s DeviceResource: %v reading: binary value", c.device.Name, cv.DeviceResourceName), sdkCommon.CorrelationHeader, c.correlationID)
		} else {
			lc.Debug(fmt.Sprintf("device: %s DeviceResource: %v reading: %v", c.device.Name, cv.DeviceResourceName, reading), sdkCommon.CorrelationHeader, c.correlationID)
		}
	}

	if !transformsOK {
		return dtos.Event{}, edgexErr.NewCommonEdgeXWrapper(fmt.Errorf("GET command %s transform failed for %s", cmd, c.device.Name))
	}

	return dtos.Event{
		Versionable: common.Versionable{
			ApiVersion: contracts.ApiVersion,
		},
		Id:          uuid.NewString(),
		Created:     time.Now().UnixNano() / 1e6,
		Origin:      sdkCommon.GetUniqueOrigin(),
		DeviceName:  c.device.Name,
		ProfileName: c.device.ProfileName,
		Readings:    readings,
	}, nil
}

func parseParams(params string) (paramMap map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(params), &paramMap)
	if err != nil {
		return
	}

	if len(paramMap) == 0 {
		err = fmt.Errorf("no parameters specified")
		return
	}

	return
}

func createCommandValueFromDeviceResource(dr *models.DeviceResource, v string) (*dsModels.CommandValue, error) {
	var err error
	var result *dsModels.CommandValue

	origin := time.Now().UnixNano()
	switch strings.ToLower(dr.Properties.Type) {
	case strings.ToLower(contracts.ValueTypeString):
		result = dsModels.NewStringValue(dr.Name, origin, v)
	case strings.ToLower(contracts.ValueTypeBool):
		value, err := strconv.ParseBool(v)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewBoolValue(dr.Name, origin, value)
	case strings.ToLower(contracts.ValueTypeBoolArray):
		var arr []bool
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewBoolArrayValue(dr.Name, origin, arr)
	case strings.ToLower(contracts.ValueTypeUint8):
		n, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewUint8Value(dr.Name, origin, uint8(n))
	case strings.ToLower(contracts.ValueTypeUint8Array):
		var arr []uint8
		strArr := strings.Split(strings.Trim(v, "[]"), ",")
		for _, u := range strArr {
			n, err := strconv.ParseUint(strings.Trim(u, " "), 10, 8)
			if err != nil {
				return result, err
			}
			arr = append(arr, uint8(n))
		}
		result, err = dsModels.NewUint8ArrayValue(dr.Name, origin, arr)
	case strings.ToLower(contracts.ValueTypeUint16):
		n, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewUint16Value(dr.Name, origin, uint16(n))
	case strings.ToLower(contracts.ValueTypeUint16Array):
		var arr []uint16
		strArr := strings.Split(strings.Trim(v, "[]"), ",")
		for _, u := range strArr {
			n, err := strconv.ParseUint(strings.Trim(u, " "), 10, 16)
			if err != nil {
				return result, err
			}
			arr = append(arr, uint16(n))
		}
		result, err = dsModels.NewUint16ArrayValue(dr.Name, origin, arr)
	case strings.ToLower(contracts.ValueTypeUint32):
		n, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewUint32Value(dr.Name, origin, uint32(n))
	case strings.ToLower(contracts.ValueTypeUint32Array):
		var arr []uint32
		strArr := strings.Split(strings.Trim(v, "[]"), ",")
		for _, u := range strArr {
			n, err := strconv.ParseUint(strings.Trim(u, " "), 10, 32)
			if err != nil {
				return result, err
			}
			arr = append(arr, uint32(n))
		}
		result, err = dsModels.NewUint32ArrayValue(dr.Name, origin, arr)
	case strings.ToLower(contracts.ValueTypeUint64):
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewUint64Value(dr.Name, origin, n)
	case strings.ToLower(contracts.ValueTypeUint64Array):
		var arr []uint64
		strArr := strings.Split(strings.Trim(v, "[]"), ",")
		for _, u := range strArr {
			n, err := strconv.ParseUint(strings.Trim(u, " "), 10, 64)
			if err != nil {
				return result, err
			}
			arr = append(arr, n)
		}
		result, err = dsModels.NewUint64ArrayValue(dr.Name, origin, arr)
	case strings.ToLower(contracts.ValueTypeInt8):
		n, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewInt8Value(dr.Name, origin, int8(n))
	case strings.ToLower(contracts.ValueTypeInt8Array):
		var arr []int8
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewInt8ArrayValue(dr.Name, origin, arr)
	case strings.ToLower(contracts.ValueTypeInt16):
		n, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewInt16Value(dr.Name, origin, int16(n))
	case strings.ToLower(contracts.ValueTypeInt16Array):
		var arr []int16
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewInt16ArrayValue(dr.Name, origin, arr)
	case strings.ToLower(contracts.ValueTypeInt32):
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewInt32Value(dr.Name, origin, int32(n))
	case strings.ToLower(contracts.ValueTypeInt32Array):
		var arr []int32
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewInt32ArrayValue(dr.Name, origin, arr)
	case strings.ToLower(contracts.ValueTypeInt64):
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewInt64Value(dr.Name, origin, n)
	case strings.ToLower(contracts.ValueTypeInt64Array):
		var arr []int64
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewInt64ArrayValue(dr.Name, origin, arr)
	case strings.ToLower(contracts.ValueTypeFloat32):
		n, e := strconv.ParseFloat(v, 32)
		if e == nil {
			result, err = dsModels.NewFloat32Value(dr.Name, origin, float32(n))
			break
		}
		if numError, ok := e.(*strconv.NumError); ok {
			if numError.Err == strconv.ErrRange {
				err = e
				break
			}
		}
		var decodedToBytes []byte
		decodedToBytes, err = base64.StdEncoding.DecodeString(v)
		if err == nil {
			var val float32
			val, err = float32FromBytes(decodedToBytes)
			if err != nil {
				break
			} else if math.IsNaN(float64(val)) {
				err = fmt.Errorf("fail to parse %v to float32, unexpected result %v", v, val)
			} else {
				result, err = dsModels.NewFloat32Value(dr.Name, origin, val)
			}
		}
	case strings.ToLower(contracts.ValueTypeFloat32Array):
		var arr []float32
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewFloat32ArrayValue(dr.Name, origin, arr)
	case strings.ToLower(contracts.ValueTypeFloat64):
		var val float64
		val, err = strconv.ParseFloat(v, 64)
		if err == nil {
			result, err = dsModels.NewFloat64Value(dr.Name, origin, val)
			break
		}
		if numError, ok := err.(*strconv.NumError); ok {
			if numError.Err == strconv.ErrRange {
				break
			}
		}
		var decodedToBytes []byte
		decodedToBytes, err = base64.StdEncoding.DecodeString(v)
		if err == nil {
			val, err = float64FromBytes(decodedToBytes)
			if err != nil {
				break
			} else if math.IsNaN(val) {
				err = fmt.Errorf("fail to parse %v to float64, unexpected result %v", v, val)
			} else {
				result, err = dsModels.NewFloat64Value(dr.Name, origin, val)
			}
		}
	case strings.ToLower(contracts.ValueTypeFloat64Array):
		var arr []float64
		err = json.Unmarshal([]byte(v), &arr)
		if err != nil {
			return result, err
		}
		result, err = dsModels.NewFloat64ArrayValue(dr.Name, origin, arr)
	default:
		err = errors.New("unsupported deviceResource value type")
	}

	if err != nil {
		return result, err
	}

	return result, err
}

func float32FromBytes(numericValue []byte) (res float32, err error) {
	reader := bytes.NewReader(numericValue)
	err = binary.Read(reader, binary.BigEndian, &res)
	return
}

func float64FromBytes(numericValue []byte) (res float64, err error) {
	reader := bytes.NewReader(numericValue)
	err = binary.Read(reader, binary.BigEndian, &res)
	return
}

func commandValueToReading(cv *dsModels.CommandValue, deviceName, profileName, mediaType, encoding string) dtos.BaseReading {
	if encoding == "" {
		encoding = dsModels.DefaultFloatEncoding
	}

	reading := dtos.BaseReading{
		Versionable: common.Versionable{
			ApiVersion: contracts.ApiVersion,
		},
		Created:      time.Now().UnixNano() / 1e6,
		ResourceName: cv.DeviceResourceName,
		DeviceName:   deviceName,
		ProfileName:  profileName,
		ValueType:    cv.Type,
	}
	if cv.Type == contracts.ValueTypeBinary {
		reading.BinaryValue = cv.BinValue
		reading.MediaType = mediaType
	} else {
		reading.Value = cv.ValueToString(encoding)
	}

	// if value has a non-zero Origin, use it
	if cv.Origin > 0 {
		reading.Origin = cv.Origin
	} else {
		reading.Origin = time.Now().UnixNano()
	}

	return reading
}

func SendEvent(event responses.EventResponse, correlationID string, lc logger.LoggingClient, ec interfaces.EventClient) {
	// TODO: comment out until core-contracts(EventClient) supports v2models
	// TODO: the usage of CBOR encoding for binary reading is under discussion
	ctx := context.WithValue(context.Background(), sdkCommon.CorrelationHeader, correlationID)
	ctx = context.WithValue(ctx, contracts.ContentType, contracts.ContentTypeJSON)
	aer := requests.AddEventRequest{
		BaseRequest: common.BaseRequest{
			Versionable: event.Versionable,
			RequestId:   event.RequestId,
		},
		Event: event.Event,
	}
	responseBody, err := ec.Add(ctx, aer)
	if err != nil {
		lc.Error("SendEvent: failed to push event to core data", "device", event.Event.DeviceName, "response", responseBody, "error", err)
	} else {
		lc.Info("SendEvent: pushed event to core data", contracts.ContentType, context2.FromContext(ctx, contracts.ContentType), contracts.CorrelationHeader, event.RequestId)
	}
}
