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
}

// models to dtos
func SendEvent(event dtos.Event, lc logger.LoggingClient, ec interfaces.EventClient) {
	correlation := uuid.New().String()
	ctx := context.WithValue(context.Background(), CorrelationHeader, correlation)
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
