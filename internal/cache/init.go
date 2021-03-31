// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/interfaces"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
)

var (
	initOnce sync.Once
)

// Init basic state for cache
func InitCache(serviceName string,
	lc logger.LoggingClient,
	dp interfaces.DeviceProfileClient,
	dc interfaces.DeviceClient,
	pwc interfaces.ProvisionWatcherClient) {
	initOnce.Do(func() {
		ctx := context.WithValue(context.Background(), common.CorrelationHeader, uuid.New().String())
		mdr, err := dc.DevicesByServiceName(ctx, serviceName, 0, -1)
		if err != nil {
			lc.Error("get device list error", err)
		}
		var dcs []models.Device
		for i := range mdr.Devices {
			dcs = append(dcs, dtos.ToDeviceModel(mdr.Devices[i]))
		}
		newDeviceCache(dcs)

		var (
			dps   []models.DeviceProfile
			dpMap = make(map[string]struct{})
		)
		for i := range dcs {
			if _, ok := dpMap[dcs[i].ProfileName]; ok {
				continue
			}
			dpr, err := dp.DeviceProfileByName(ctx, dcs[i].ProfileName)
			if err != nil {
				lc.Error(fmt.Sprintf("get device profile(%s) error: %+v", dcs[i].ProfileName, err))
				continue
			}
			dpMap[dcs[i].Name] = struct{}{}
			dps = append(dps, dtos.ToDeviceProfileModel(dpr.Profile))
		}
		newProfileCache(dps)

		pwr, err := pwc.ProvisionWatchersByServiceName(ctx, serviceName, 0, -1)
		if err != nil {
			lc.Error(fmt.Sprintf("get device profile(%s) error: %+v", serviceName, err))
		}
		var pws []models.ProvisionWatcher
		for i := range pwr.ProvisionWatchers {
			pws = append(pws, dtos.ToProvisionWatcherModel(pwr.ProvisionWatchers[i]))
		}
		newProvisionWatcherCache(pws)
	})
}
