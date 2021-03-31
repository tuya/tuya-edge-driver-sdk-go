// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/mock"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
)

func TestInitCache(t *testing.T) {
	serviceName := "init-cache-test"
	lc := logger.NewMockClient()
	//vdc := &mock.ValueDescriptorMock{}
	dp := &mock.DeviceProfileClientMock{}
	dc := &mock.DeviceClientMock{}
	pwc := &mock.ProvisionWatcherClientMock{}
	InitCache(serviceName, lc, dp, dc, pwc)

	ctx := context.WithValue(context.Background(), common.CorrelationHeader, uuid.New().String())

	/*
		vdsBeforeAddingToCache, _ := vdc.ValueDescriptors(ctx)
		if vl := len(ValueDescriptors().All()); vl != len(vdsBeforeAddingToCache) {
			t.Errorf("the expected number of valuedescriptors in cache is %d but got: %d:", len(vdsBeforeAddingToCache), vl)
		}
	*/

	dsBeforeAddingToCache, err := dc.DevicesByServiceName(ctx, serviceName, 0, -1)
	if err != nil {
		t.Error("get device by service name error: ", err)
	}
	if dl := len(Devices().All()); dl != len(dsBeforeAddingToCache.Devices) {
		t.Errorf("the expected number of devices in cache is %d but got: %d:", len(dsBeforeAddingToCache.Devices), dl)
	}

	pMap := make(map[string]models.DeviceProfile, len(dsBeforeAddingToCache.Devices)*2)
	for _, d := range dsBeforeAddingToCache.Devices {
		dp, b := Profiles().ForId(d.ProfileId)
		if !b {
			t.Error("device profile not exists in local cache")
			continue
		}
		pMap[d.ProfileName] = dp
	}
	if pl := len(Profiles().All()); pl != len(pMap) {
		t.Errorf("the expected number of device profiles in cache is %d but got: %d:", len(pMap), pl)
	} else {
		psFromCache := Profiles().All()
		for _, p := range psFromCache {
			assert.Equal(t, pMap[p.Name], p)
		}
	}
}
