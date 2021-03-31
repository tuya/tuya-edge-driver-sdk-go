// -*- mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package clients

import "github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/interfaces"

type TedgeClients struct {
	CommonClient           interfaces.CommonClient
	DeviceClient           interfaces.DeviceClient
	DeviceServiceClient    interfaces.DeviceServiceClient
	DeviceProfileClient    interfaces.DeviceProfileClient
	CallbackClient         interfaces.DeviceServiceCallbackClient
	ProvisionWatcherClient interfaces.ProvisionWatcherClient
	EventClient            interfaces.EventClient
}
