// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2021 IOTech Ltd
// Copyright (c) 2019 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/gorilla/mux"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	sdkCommon "github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/controller/correlation"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
)

type RestController struct {
	LoggingClient  logger.LoggingClient
	router         *mux.Router
	reservedRoutes map[string]bool
	httpController *HttpController
	dic            *di.Container
}

func NewRestController(r *mux.Router, dic *di.Container) *RestController {
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)
	return &RestController{
		LoggingClient:  lc,
		router:         r,
		reservedRoutes: make(map[string]bool),
		httpController: NewHttpController(dic),
		dic:            dic,
	}
}

func (c *RestController) InitRestRoutes() {
	c.LoggingClient.Info("Registering v2 routes...")

	c.addReservedRoute(contracts.ApiPingRoute, c.httpController.Ping).Methods(http.MethodGet)
	c.addReservedRoute(contracts.ApiVersionRoute, c.httpController.Version).Methods(http.MethodGet)
	c.addReservedRoute(contracts.ApiConfigRoute, c.httpController.Config).Methods(http.MethodGet)
	c.addReservedRoute(contracts.ApiMetricsRoute, c.httpController.Metrics).Methods(http.MethodGet)

	c.addReservedRoute(sdkCommon.APIV2SecretRoute, c.httpController.Secret).Methods(http.MethodPost)

	c.addReservedRoute(contracts.ApiDiscoveryRoute, c.httpController.Discovery).Methods(http.MethodPost)

	c.addReservedRoute(contracts.ApiDeviceNameCommandNameRoute, c.httpController.Command).Methods(http.MethodPut, http.MethodGet)

	c.addReservedRoute(contracts.ApiDeviceCallbackRoute, c.httpController.AddDevice).Methods(http.MethodPost)
	c.addReservedRoute(contracts.ApiDeviceCallbackRoute, c.httpController.UpdateDevice).Methods(http.MethodPut)
	c.addReservedRoute(contracts.ApiDeviceCallbackNameRoute, c.httpController.DeleteDevice).Methods(http.MethodDelete)
	c.addReservedRoute(contracts.ApiProfileCallbackRoute, c.httpController.UpdateProfile).Methods(http.MethodPut)
	c.addReservedRoute(contracts.ApiProvisionWatcherRoute, c.httpController.AddProvisionWatcher).Methods(http.MethodPost)
	c.addReservedRoute(contracts.ApiProvisionWatcherRoute, c.httpController.UpdateProvisionWatcher).Methods(http.MethodPut)
	c.addReservedRoute(contracts.ApiProvisionWatcherByNameRoute, c.httpController.DeleteProvisionWatcher).Methods(http.MethodDelete)

	c.router.Use(correlation.ManageHeader)
	c.router.Use(correlation.OnResponseComplete)
	c.router.Use(correlation.OnRequestBegin)
}

func (c *RestController) addReservedRoute(route string, handler func(http.ResponseWriter, *http.Request)) *mux.Route {
	c.reservedRoutes[route] = true
	return c.router.HandleFunc(
		route,
		func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), bootstrapContainer.LoggingClientInterfaceName, c.LoggingClient)
			handler(
				w,
				r.WithContext(ctx))
		})
}

func (c *RestController) AddRoute(route string, handler func(http.ResponseWriter, *http.Request), methods ...string) error {
	if c.reservedRoutes[route] {
		return errors.New("route is reserved")
	}

	c.router.HandleFunc(
		route,
		func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), bootstrapContainer.LoggingClientInterfaceName, c.LoggingClient)
			handler(
				w,
				r.WithContext(ctx))
		}).Methods(methods...)
	c.LoggingClient.Debug("Route added", "route", route, "methods", fmt.Sprintf("%v", methods))

	return nil
}

func (c *RestController) Router() *mux.Router {
	return c.router
}
