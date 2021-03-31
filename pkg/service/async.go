// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018-2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	commonDTO "github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/cache"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/transformer"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
	dsModels "github.com/tuya/tuya-edge-driver-sdk-go/pkg/models"
)

// processAsyncResults processes readings that are pushed from
// a DS implementation. Each is reading is optionally transformed
// before being pushed to Core Data.
// In this function, AsyncBufferSize is used to create a buffer for
// processing AsyncValues concurrently, so that events may arrive
// out-of-order in core-data / app service when AsyncBufferSize value
// is greater than or equal to two. Alternatively, we can process
// AsyncValues one by one in the same order by changing the AsyncBufferSize
// value to one.
func (s *DeviceService) processAsyncResults(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		wg.Done()
	}()

	working := make(chan bool, s.config.Service.AsyncBufferSize)
	for {
		select {
		case <-ctx.Done():
			return
		case acv := <-s.asyncCh:
			go s.sendAsyncValues(acv, working)
		}
	}
}

// sendAsyncValues convert AsyncValues to event and send the event to CoreData
func (s *DeviceService) sendAsyncValues(acv *dsModels.AsyncValues, working chan bool) {
	working <- true
	defer func() {
		<-working
	}()
	readings := make([]models.Reading, 0, len(acv.CommandValues))

	device, ok := cache.Devices().ForName(acv.DeviceName)
	if !ok {
		s.LoggingClient.Error(fmt.Sprintf("processAsyncResults - recieved Device %s not found in cache", acv.DeviceName))
		return
	}

	for _, cv := range acv.CommandValues {
		// get the device resource associated with the rsp.RO
		dr, ok := cache.Profiles().DeviceResource(device.ProfileName, cv.DeviceResourceName)
		if !ok {
			s.LoggingClient.Error(fmt.Sprintf("processAsyncResults - Device Resource %s not found in Device %s", cv.DeviceResourceName, acv.DeviceName))
			continue
		}

		// device resourse property转换
		if s.config.Device.DataTransform {
			err := transformer.TransformReadResult(cv, dr.Properties, s.LoggingClient)
			if err != nil {
				s.LoggingClient.Error(fmt.Sprintf("processAsyncResults - CommandValue (%s) transformed failed: %v", cv.String(), err))

				if errors.As(err, &transformer.OverflowError{}) {
					cv = dsModels.NewStringValue(cv.DeviceResourceName, cv.Origin, transformer.Overflow)
				} else if errors.As(err, &transformer.NaNError{}) {
					cv = dsModels.NewStringValue(cv.DeviceResourceName, cv.Origin, transformer.NaN)
				} else {
					cv = dsModels.NewStringValue(cv.DeviceResourceName, cv.Origin, fmt.Sprintf("Transformation failed for device resource, with value: %s, property value: %v, and error: %v", cv.String(), dr.Properties, err))
				}
			}
		}

		err := transformer.CheckAssertion(cv, dr.Properties.Assertion, &device, s.LoggingClient, s.tedgeClients.DeviceClient)
		if err != nil {
			s.LoggingClient.Error(fmt.Sprintf("processAsyncResults - Assertion failed for device resource: %s, with value: %s and assertion: %s, %v", cv.DeviceResourceName, cv.String(), dr.Properties.Assertion, err))
			cv = dsModels.NewStringValue(cv.DeviceResourceName, cv.Origin, fmt.Sprintf("Assertion failed for device resource, with value: %s and assertion: %s", cv.String(), dr.Properties.Assertion))
		}

		ro, err := cache.Profiles().ResourceOperation(device.ProfileName, cv.DeviceResourceName, common.GetCmdMethod)
		if err != nil {
			s.LoggingClient.Debug(fmt.Sprintf("processAsyncResults - getting resource operation failed: %s", err.Error()))
		} else if len(ro.Mappings) > 0 {
			newCV, ok := transformer.MapCommandValue(cv, ro.Mappings)
			if ok {
				cv = newCV
			} else { // TODO
				s.LoggingClient.Warn(fmt.Sprintf("processAsyncResults - Mapping failed for Device Resource Operation: %s, with value: %s, %v", ro, cv.String(), err))
			}
		}

		// TODO float encoding
		// TODO 直接创建dtos reading
		reading := common.CommandValueToReading(cv, device.Name, device.ProfileName, dr.Properties.MediaType, "")
		readings = append(readings, reading)
	}

	// push to Core Data
	event := models.Event{
		Id:          uuid.NewString(),
		DeviceName:  device.Name,
		ProfileName: device.ProfileName,
		Created:     time.Now().UnixNano() / 1e6,
		Origin:      common.GetUniqueOrigin(),
		Readings:    readings,
	}

	common.SendEvent(dtos.FromEventModelToDTO(event), s.LoggingClient, s.tedgeClients.EventClient)
}

// processAsyncFilterAndAdd filter and add devices discovered by
// device service protocol discovery.
func (s *DeviceService) processAsyncFilterAndAdd(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case devices := <-s.deviceCh:
			newDevReqs := make([]requests.AddDeviceRequest, 0, len(devices))
			ctx := context.Background()
			pws := cache.ProvisionWatchers().All()
			for _, d := range devices {
				for _, pw := range pws {
					if whitelistPass(d, pw, s.LoggingClient) && blacklistPass(d, pw, s.LoggingClient) {
						if _, ok := cache.Devices().ForName(d.Name); ok {
							s.LoggingClient.Debug(fmt.Sprintf("Candidate discovered device %s already existed", d.Name))
							break
						}
						s.LoggingClient.Info(fmt.Sprintf("Adding discovered device %s to Edgex", d.Name))
						millis := time.Now().UnixNano() / 1e6
						device := models.Device{
							Name:           d.Name,
							ProfileName:    pw.ProfileName,
							Protocols:      d.Protocols,
							Labels:         d.Labels,
							ServiceName:    pw.ServiceName,
							AdminState:     pw.AdminState,
							OperatingState: models.Up,
							AutoEvents:     nil,
						}
						device.Created = millis
						device.Description = d.Description
						// models to dtos
						newDevReqs = append(newDevReqs, requests.AddDeviceRequest{
							BaseRequest: commonDTO.NewBaseRequest(),
							Device:      dtos.FromDeviceModelToDTO(device),
						})
					}
				}
				// add here
				_, err := s.tedgeClients.DeviceClient.Add(ctx, newDevReqs)
				if err != nil {
					s.LoggingClient.Error(fmt.Sprintf("failed to create discovered device %v", err))
				} else {
					break
				}
			}
			s.LoggingClient.Debug("Filtered device addition finished")
		}
	}
}

func whitelistPass(d dsModels.DiscoveredDevice, pw models.ProvisionWatcher, lc logger.LoggingClient) bool {
	// ignore the device protocol properties name
	for _, protocol := range d.Protocols {
		matchedCount := 0
		for name, regex := range pw.Identifiers {
			if value, ok := protocol[name]; ok {
				matched, err := regexp.MatchString(regex, value)
				if !matched || err != nil {
					lc.Debug(fmt.Sprintf("Device %s's %s value %s did not match PW identifier: %s", d.Name, name, value, regex))
					break
				}
				matchedCount += 1
			}
		}
		// match succeed on all identifiers
		if matchedCount == len(pw.Identifiers) {
			return true
		}
	}
	return false
}

func blacklistPass(d dsModels.DiscoveredDevice, pw models.ProvisionWatcher, lc logger.LoggingClient) bool {
	// a candidate should match none of the blocking identifiers
	for name, blacklist := range pw.BlockingIdentifiers {
		// ignore the device protocol properties name
		for _, protocol := range d.Protocols {
			if value, ok := protocol[name]; ok {
				for _, v := range blacklist {
					if value == v {
						lc.Debug(fmt.Sprintf("Discovered Device %s's %s should not be %s", d.Name, name, value))
						return false
					}
				}
			}
		}
	}
	return true
}
