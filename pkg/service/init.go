// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/gorilla/mux"

	"github.com/tuya/tuya-edge-driver-sdk-go/internal/autoevent"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/cache"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/container"
	"github.com/tuya/tuya-edge-driver-sdk-go/pkg/models"
)

// Bootstrap contains references to dependencies required by the BootstrapHandler.
type Bootstrap struct {
	router *mux.Router
}

// NewBootstrap is a factory method that returns an initialized Bootstrap receiver struct.
func NewBootstrap(router *mux.Router) *Bootstrap {
	return &Bootstrap{
		router: router,
	}
}

func (b *Bootstrap) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup, startupTimer startup.Timer, dic *di.Container) (success bool) {
	ds.UpdateFromContainer(b.router, dic)
	autoevent.NewManager(ctx, wg, ds.config.Service.AsyncBufferSize, dic)

	// 由服务自注册改为通过ID更新服务配置
	if err := ds.updateService(); err != nil {
		ds.LoggingClient.Error(fmt.Sprintf("update device service instance config error: %v\n", err))
		return false
	}

	// initialize devices, deviceResources, provisionWatchers & profiles cache
	cache.InitCache(
		ds.deviceService.Name,
		ds.LoggingClient,
		container.MetadataDeviceProfileClientFrom(dic.Get),
		container.MetadataDeviceClientFrom(dic.Get),
		container.MetadataProvisionWatcherClientFrom(dic.Get))

	if ds.AsyncReadings() {
		ds.asyncCh = make(chan *models.AsyncValues, ds.config.Service.AsyncBufferSize)
		go ds.processAsyncResults(ctx, wg)
	}
	if ds.DeviceDiscovery() {
		ds.deviceCh = make(chan []models.DiscoveredDevice, 1)
		go ds.processAsyncFilterAndAdd(ctx, wg)
	}

	err := ds.driver.Initialize(ds.LoggingClient, ds.asyncCh, ds.deviceCh)
	if err != nil {
		ds.LoggingClient.Error(fmt.Sprintf("Driver.Initialize failed: %v\n", err))
		return false
	}
	ds.initialized = true

	dic.Update(di.ServiceConstructorMap{
		container.ProtocolDiscoveryName: func(get di.Get) interface{} {
			return ds.discovery
		},
		container.ProtocolDriverName: func(get di.Get) interface{} {
			return ds.driver
		},
		//v2
		container.DeviceServiceName: func(get di.Get) interface{} {
			return ds.deviceService
		},
		container.ProtocolDiscoveryName: func(get di.Get) interface{} {
			return ds.discovery
		},
		container.ProtocolDriverName: func(get di.Get) interface{} {
			return ds.driver
		},
	})

	ds.controller.InitRestRoutes()

	autoevent.GetManager().StartAutoEvents(dic)
	http.TimeoutHandler(nil, time.Millisecond*time.Duration(ds.config.Service.Timeout), "Request timed out")

	return true
}
