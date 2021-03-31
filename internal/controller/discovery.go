// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"net/http"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	edgexErr "github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/autodiscovery"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/container"
)

func (c *HttpController) Discovery(writer http.ResponseWriter, request *http.Request) {
	ds := container.DeviceServiceFrom(c.dic.Get)
	if ds.AdminState == models.Locked {
		err := edgexErr.NewCommonEdgeX(edgexErr.KindServiceLocked, "service locked", nil)
		c.sendEdgexError(writer, request, err, contracts.ApiDiscoveryRoute)
		return
	}

	configuration := container.ConfigurationFrom(c.dic.Get)
	if !configuration.Device.Discovery.Enabled {
		err := edgexErr.NewCommonEdgeX(edgexErr.KindServiceUnavailable, "device discovery disabled", nil)
		c.sendEdgexError(writer, request, err, contracts.ApiDiscoveryRoute)
		return
	}

	discovery := container.ProtocolDiscoveryFrom(c.dic.Get)
	if discovery == nil {
		err := edgexErr.NewCommonEdgeX(edgexErr.KindNotImplemented, "protocolDiscovery not implemented", nil)
		c.sendEdgexError(writer, request, err, contracts.ApiDiscoveryRoute)
		return
	}

	go autodiscovery.DiscoveryWrapper(discovery, c.lc)
	c.sendResponse(writer, request, contracts.ApiDiscoveryRoute, nil, http.StatusAccepted)
}
