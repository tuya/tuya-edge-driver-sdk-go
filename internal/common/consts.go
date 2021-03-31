// -*- mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2020 IOTech Ltd
// Copyright (c) 2019 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
)

const (
	ClientData     = "Data"
	ClientMetadata = "Metadata"

	EnvInstanceName = "EDGEX_INSTANCE_NAME"

	Colon      = ":"
	HttpScheme = "http://"
	HttpProto  = "HTTP"

	ConfigStemDevice   = "edgex/devices/"
	ConfigMajorVersion = "1.0/"

	APICallbackRoute        = contracts.ApiCallbackRoute
	APIValueDescriptorRoute = contracts.ApiValueDescriptorRoute
	APIPingRoute            = contracts.ApiPingRoute
	APIVersionRoute         = contracts.ApiVersionRoute
	APIMetricsRoute         = contracts.ApiMetricsRoute
	APIConfigRoute          = contracts.ApiConfigRoute
	APIAllCommandRoute      = contracts.ApiDeviceRoute + "/all/{command}"
	APIIdCommandRoute       = contracts.ApiDeviceRoute + "/{id}/{command}"
	APINameCommandRoute     = contracts.ApiDeviceRoute + "/name/{name}/{command}"
	APIDiscoveryRoute       = contracts.ApiBase + "/discovery"
	APITransformRoute       = contracts.ApiBase + "/debug/transformData/{transformData}"

	APIV2SecretRoute = contracts.ApiBase + "/secret"

	IdVar        string = "id"
	NameVar      string = "name"
	CommandVar   string = "command"
	GetCmdMethod string = "get"
	SetCmdMethod string = "set"

	DeviceResourceReadOnly  string = "R"
	DeviceResourceWriteOnly string = "W"

	CorrelationHeader = contracts.CorrelationHeader
	URLRawQuery       = "urlRawQuery"
	SDKReservedPrefix = "ds-"
)

// SDKVersion indicates the version of the SDK - will be overwritten by build
var SDKVersion string = "0.0.0"

// ServiceVersion indicates the version of the device service itself, not the SDK - will be overwritten by build
var ServiceVersion string = "0.0.0"
