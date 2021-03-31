// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"os"

	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/flags"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/handlers"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/interfaces"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/gorilla/mux"

	"github.com/tuya/tuya-edge-driver-sdk-go/internal/autodiscovery"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/clients"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/container"
)

func Main(serviceName string, serviceVersion string, proto interface{}, ctx context.Context, cancel context.CancelFunc, router *mux.Router) {
	startupTimer := startup.NewStartUpTimer(serviceName)
	sdkFlags := flags.New()
	sdkFlags.Parse(os.Args[1:])

	ds = &DeviceService{}
	ds.Initialize(serviceName, serviceVersion, proto)

	dic := di.NewContainer(di.ServiceConstructorMap{
		container.ConfigurationName: func(get di.Get) interface{} {
			return ds.config
		},
	})

	httpServer := handlers.NewHttpServer(router, true)

	bootstrap.Run(
		ctx,
		cancel,
		sdkFlags,
		ds.ServiceName,
		common.ConfigStemDevice+common.ConfigMajorVersion,
		ds.config,
		startupTimer,
		dic,
		[]interfaces.BootstrapHandler{
			httpServer.BootstrapHandler,
			clients.NewClients().BootstrapHandler,
			NewBootstrap(router).BootstrapHandler,
			autodiscovery.BootstrapHandler,
			handlers.NewStartMessage(serviceName, serviceVersion).BootstrapHandler,
		})

	ds.Stop(false)
}
