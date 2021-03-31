// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	commonDTO "github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/callback"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
)

func (c *HttpController) DeleteDevice(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars[common.NameVar]

	err := callback.DeleteDevice(name, c.dic)
	if err == nil {
		res := commonDTO.NewBaseResponse("", "", http.StatusOK)
		c.sendResponse(writer, request, contracts.ApiDeviceCallbackNameRoute, res, http.StatusOK)
	} else {
		c.sendEdgexError(writer, request, err, contracts.ApiDeviceCallbackNameRoute)
	}
}

func (c *HttpController) AddDevice(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var addDeviceRequest requests.AddDeviceRequest

	err := json.NewDecoder(request.Body).Decode(&addDeviceRequest)
	if err != nil {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "failed to decode JSON", err)
		c.sendEdgexError(writer, request, edgexErr, contracts.ApiDeviceCallbackRoute)
		return
	}

	edgexErr := callback.AddDevice(addDeviceRequest, c.dic)
	if edgexErr == nil {
		res := commonDTO.NewBaseResponse(addDeviceRequest.RequestId, "", http.StatusOK)
		c.sendResponse(writer, request, contracts.ApiDeviceCallbackRoute, res, http.StatusOK)
	} else {
		c.sendEdgexError(writer, request, edgexErr, contracts.ApiDeviceCallbackRoute)
	}
}

func (c *HttpController) UpdateDevice(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var updateDeviceRequest requests.UpdateDeviceRequest

	err := json.NewDecoder(request.Body).Decode(&updateDeviceRequest)
	if err != nil {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "failed to decode JSON", err)
		c.sendEdgexError(writer, request, edgexErr, contracts.ApiDeviceCallbackRoute)
		return
	}

	edgexErr := callback.UpdateDevice(updateDeviceRequest, c.dic)
	if edgexErr == nil {
		res := commonDTO.NewBaseResponse(updateDeviceRequest.RequestId, "", http.StatusOK)
		c.sendResponse(writer, request, contracts.ApiDeviceCallbackRoute, res, http.StatusOK)
	} else {
		c.sendEdgexError(writer, request, edgexErr, contracts.ApiDeviceCallbackRoute)
	}
}

func (c *HttpController) UpdateProfile(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var edgexErr errors.EdgeX
	var profileRequest requests.DeviceProfileRequest

	err := json.NewDecoder(request.Body).Decode(&profileRequest)
	if err != nil {
		edgexErr = errors.NewCommonEdgeX(errors.KindServerError, "failed to decode JSON", err)
		c.sendEdgexError(writer, request, edgexErr, contracts.ApiProfileCallbackRoute)
		return
	}

	edgexErr = callback.UpdateProfile(profileRequest, c.lc)
	if edgexErr == nil {
		res := commonDTO.NewBaseResponse(profileRequest.RequestId, "", http.StatusOK)
		c.sendResponse(writer, request, contracts.ApiProfileCallbackRoute, res, http.StatusOK)
	} else {
		c.sendEdgexError(writer, request, edgexErr, contracts.ApiProfileCallbackRoute)
	}
}

func (c *HttpController) DeleteProvisionWatcher(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars[common.NameVar]

	err := callback.DeleteProvisionWatcher(name, c.lc)
	if err == nil {
		res := commonDTO.NewBaseResponse("", "", http.StatusOK)
		c.sendResponse(writer, request, contracts.ApiWatcherCallbackNameRoute, res, http.StatusOK)
	} else {
		c.sendEdgexError(writer, request, err, contracts.ApiWatcherCallbackNameRoute)
	}
}

func (c *HttpController) AddProvisionWatcher(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var addProvisionWatcherRequest requests.AddProvisionWatcherRequest

	err := json.NewDecoder(request.Body).Decode(&addProvisionWatcherRequest)
	if err != nil {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "failed to decode JSON", err)
		c.sendEdgexError(writer, request, edgexErr, contracts.ApiWatcherCallbackRoute)
		return
	}

	edgexErr := callback.AddProvisionWatcher(addProvisionWatcherRequest, c.lc)
	if edgexErr == nil {
		res := commonDTO.NewBaseResponse("", "", http.StatusOK)
		c.sendResponse(writer, request, contracts.ApiWatcherCallbackRoute, res, http.StatusOK)
	} else {
		c.sendEdgexError(writer, request, edgexErr, contracts.ApiWatcherCallbackRoute)
	}
}

func (c *HttpController) UpdateProvisionWatcher(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var updateProvisionWatcherRequest requests.UpdateProvisionWatcherRequest

	err := json.NewDecoder(request.Body).Decode(&updateProvisionWatcherRequest)
	if err != nil {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "failed to decode JSON", err)
		c.sendEdgexError(writer, request, edgexErr, contracts.ApiWatcherCallbackRoute)
		return
	}

	edgexErr := callback.UpdateProvisionWatcher(updateProvisionWatcherRequest, c.lc)
	if edgexErr == nil {
		res := commonDTO.NewBaseResponse("", "", http.StatusOK)
		c.sendResponse(writer, request, contracts.ApiWatcherCallbackRoute, res, http.StatusOK)
	} else {
		c.sendEdgexError(writer, request, edgexErr, contracts.ApiWatcherCallbackRoute)
	}
}
