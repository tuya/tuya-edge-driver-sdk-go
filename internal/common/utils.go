// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/interfaces"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	context2 "github.com/tuya/tuya-edge-driver-sdk-go/internal/context"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
	dsModels "github.com/tuya/tuya-edge-driver-sdk-go/pkg/models"
)

var (
	previousOrigin int64
	originMutex    sync.Mutex
)

func BuildAddr(host string, port string) string {
	var buffer bytes.Buffer

	buffer.WriteString(HttpScheme)
	buffer.WriteString(host)
	buffer.WriteString(Colon)
	buffer.WriteString(port)

	return buffer.String()
}

func CommandValueToReading(cv *dsModels.CommandValue, devName, profileName, mediaType, encoding string) models.Reading {
	baseReading := models.BaseReading{
		Id:           uuid.NewString(),
		DeviceName:   devName,
		ResourceName: cv.DeviceResourceName,
		ProfileName:  profileName,
		Created:      time.Now().UnixNano() / 1e6,
		ValueType:    cv.Type,
	}
	if cv.Origin > 0 {
		baseReading.Origin = cv.Origin
	} else {
		baseReading.Origin = time.Now().UnixNano() / 1e6
	}
	if cv.Type == contracts.ValueTypeBinary {

		return models.BinaryReading{
			BaseReading: baseReading,
			BinaryValue: cv.BinValue,
			MediaType:   mediaType,
		}
	} else {
		return models.SimpleReading{
			BaseReading: baseReading,
			Value:       cv.ValueToString(encoding),
		}
	}
	/*
		reading := models.Reading{ResourseName: cv.DeviceResourceName, Device: devName, ValueType: cv.Type}
		if cv.Type == contracts.ValueTypeBool {
			reading.BinaryValue = cv.BinValue
			reading.MediaType = mediaType
		} else if cv.Type == contracts.ValueTypeFloat32 || cv.Type == contracts.ValueTypeFloat64 {
			reading.Value = cv.ValueToString(encoding)
			reading.FloatEncoding = encoding
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
	*/
}

// models to dtos
func SendEvent(event dtos.Event, lc logger.LoggingClient, ec interfaces.EventClient) {
	correlation := uuid.New().String()
	ctx := context.WithValue(context.Background(), CorrelationHeader, correlation)
	/*
		if event.HasBinaryValue() {
			ctx = context.WithValue(ctx, contracts.ContentType, contracts.ContentTypeCBOR)
		} else {
			ctx = context.WithValue(ctx, contracts.ContentType, contracts.ContentTypeJSON)
		}
		// Call MarshalEvent to encode as byte array whether event contains binary or JSON readings
		var err error
		if len(event.EncodedEvent) <= 0 {
			event.EncodedEvent, err = ec.MarshalEvent(event.Event)
			if err != nil {
				lc.Error("SendEvent: Error encoding event", "device", event.Device, contracts.CorrelationHeader, correlation, "error", err)
			} else {
				lc.Debug("SendEvent: EventClient.MarshalEvent encoded event", contracts.CorrelationHeader, correlation)
			}
		} else {
			lc.Debug("SendEvent: EventClient.MarshalEvent passed through encoded event", contracts.CorrelationHeader, correlation)
		}
	*/
	req := requests.AddEventRequest{
		BaseRequest: common.NewBaseRequest(),
		Event:       event,
	}
	// Call Add to post event to core data
	responseBody, errPost := ec.Add(ctx, req)
	if errPost != nil {
		lc.Error("SendEvent Failed to push event", "device", event.DeviceName, "response", responseBody, "error", errPost)
	} else {
		lc.Debug("SendEvent: Pushed event to core data", contracts.ContentType, context2.FromContext(ctx, contracts.ContentType), contracts.CorrelationHeader, correlation)
		lc.Trace("SendEvent: Pushed this event to core data", contracts.ContentType, context2.FromContext(ctx, contracts.ContentType), contracts.CorrelationHeader, correlation, "event", event)
	}
}

//func CompareCoreCommands(a []models.Command, b []models.Command) bool {
//	if len(a) != len(b) {
//		return false
//	}
//
//	for i := range a {
//		if a[i].String() != b[i].String() {
//			return false
//		}
//	}
//
//	return true
//}

//func CompareDevices(a models.Device, b models.Device) bool {
//	labelsOk := CompareStrings(a.Labels, b.Labels)
//	profileOk := CompareDeviceProfiles(a.Profile, b.Profile)
//	serviceOk := CompareDeviceServices(a.Service, b.Service)
//
//	return reflect.DeepEqual(a.Protocols, b.Protocols) &&
//		a.AdminState == b.AdminState &&
//		a.Description == b.Description &&
//		a.Id == b.Id &&
//		a.Location == b.Location &&
//		a.Name == b.Name &&
//		a.OperatingState == b.OperatingState &&
//		labelsOk &&
//		profileOk &&
//		serviceOk
//}

//func CompareDeviceProfiles(a models.DeviceProfile, b models.DeviceProfile) bool {
//	labelsOk := CompareStrings(a.Labels, b.Labels)
//	cmdsOk := CompareCoreCommands(a.CoreCommands, b.CoreCommands)
//	devResourcesOk := CompareDeviceResources(a.DeviceResources, b.DeviceResources)
//	resourcesOk := CompareDeviceCommands(a.DeviceCommands, b.DeviceCommands)
//
//	// TODO: DeviceResource fields aren't compared as to dr properly
//	// requires introspection as DeviceResource is a slice of interface{}
//
//	return a.DescribedObject == b.DescribedObject &&
//		a.Id == b.Id &&
//		a.Name == b.Name &&
//		a.Manufacturer == b.Manufacturer &&
//		a.Model == b.Model &&
//		labelsOk &&
//		cmdsOk &&
//		devResourcesOk &&
//		resourcesOk
//}

//func CompareDeviceResources(a []models.DeviceResource, b []models.DeviceResource) bool {
//	if len(a) != len(b) {
//		return false
//	}
//
//	for i := range a {
//		// TODO: Attributes aren't compared, as to dr properly
//		// requires introspection as Attributes is an interface{}
//
//		if a[i].Description != b[i].Description ||
//			a[i].Name != b[i].Name ||
//			a[i].Tag != b[i].Tag ||
//			a[i].Properties != b[i].Properties {
//			return false
//		}
//	}
//
//	return true
//}

//func CompareDeviceServices(a models.DeviceService, b models.DeviceService) bool {
//	serviceOk := CompareServices(a, b)
//	return a.AdminState == b.AdminState && serviceOk
//}
//
//func CompareDeviceCommands(a []models.ProfileResource, b []models.ProfileResource) bool {
//	if len(a) != len(b) {
//		return false
//	}
//
//	for i := range a {
//		getOk := CompareResourceOperations(a[i].Get, b[i].Set)
//		setOk := CompareResourceOperations(a[i].Get, b[i].Set)
//
//		if a[i].Name != b[i].Name && !getOk && !setOk {
//			return false
//		}
//	}
//
//	return true
//}

//func CompareResourceOperations(a []models.ResourceOperation, b []models.ResourceOperation) bool {
//	if len(a) != len(b) {
//		return false
//	}
//
//	for i := range a {
//		secondaryOk := CompareStrings(a[i].Secondary, b[i].Secondary)
//		mappingsOk := CompareStrStrMap(a[i].Mappings, b[i].Mappings)
//
//		if a[i].Index != b[i].Index ||
//			a[i].Operation != b[i].Operation ||
//			a[i].DeviceResource != b[i].DeviceResource ||
//			a[i].Parameter != b[i].Parameter ||
//			a[i].DeviceCommand != b[i].DeviceCommand ||
//			!secondaryOk ||
//			!mappingsOk {
//			return false
//		}
//	}
//
//	return true
//}
//
//func CompareServices(a models.DeviceService, b models.DeviceService) bool {
//	labelsOk := CompareStrings(a.Labels, b.Labels)
//
//	return a.DescribedObject == b.DescribedObject &&
//		a.Id == b.Id &&
//		a.Name == b.Name &&
//		a.LastConnected == b.LastConnected &&
//		a.LastReported == b.LastReported &&
//		a.OperatingState == b.OperatingState &&
//		a.Addressable == b.Addressable &&
//		labelsOk
//}

func CompareStrings(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func CompareStrStrMap(a map[string]string, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, av := range a {
		if bv, ok := b[k]; !ok || av != bv {
			return false
		}
	}

	return true
}

func VerifyIdFormat(id string, drName string) errors.EdgeX {
	if len(id) == 0 {
		errMsg := fmt.Sprintf("The Id of %s is empty string", drName)
		return errors.NewCommonEdgeXWrapper(fmt.Errorf(errMsg))
	}
	return nil
}

func GetUniqueOrigin() int64 {
	originMutex.Lock()
	defer originMutex.Unlock()
	now := time.Now().UnixNano() / 1e6
	if now <= previousOrigin {
		now = previousOrigin + 1
	}
	previousOrigin = now
	return now
}

func FilterQueryParams(queryParams string, lc logger.LoggingClient) url.Values {
	m, err := url.ParseQuery(queryParams)
	if err != nil {
		lc.Error("Error parsing query parameters: %s\n", err)
	}
	// Filter out parameters with predefined prefix
	for k := range m {
		if strings.HasPrefix(k, SDKReservedPrefix) {
			delete(m, k)
		}
	}

	return m
}

func UpdateLastConnected(device models.Device, configuration *ConfigurationStruct, lc logger.LoggingClient, dc interfaces.DeviceClient) {
	if !configuration.Device.UpdateLastConnected {
		lc.Debug("Update of last connected times is disabled for: " + device.Name)
		return
	}

	t := time.Now().UnixNano() / 1e6
	req := make([]requests.UpdateDeviceRequest, 0, 1)
	req = append(req, requests.UpdateDeviceRequest{
		BaseRequest: common.NewBaseRequest(),
		Device: dtos.UpdateDevice{
			Name:          &device.Name,
			ServiceName:   &device.ServiceName,
			ProfileName:   &device.ProfileName,
			LastConnected: &t,
		},
	})
	_, err := dc.Update(context.Background(), req)
	if err != nil {
		lc.Error("Failed to update last connected value for device: " + device.Name)
	}
}
