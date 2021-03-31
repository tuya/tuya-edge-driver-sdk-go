// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/interfaces"
)

// v2版本相关依赖客户端
var (
	CommonClientName                        = di.TypeInstanceToName((*interfaces.CommonClient)(nil))
	MetadataDeviceClientName                = di.TypeInstanceToName((*interfaces.DeviceClient)(nil))
	MetadataDeviceProfileClientName         = di.TypeInstanceToName((*interfaces.DeviceProfileClient)(nil))
	MetadataDeviceServiceClientName         = di.TypeInstanceToName((*interfaces.DeviceServiceClient)(nil))
	MetadataDeviceServiceCallbackClientName = di.TypeInstanceToName((*interfaces.DeviceServiceCallbackClient)(nil))
	MetadataProvisionWatcherClientName      = di.TypeInstanceToName((*interfaces.ProvisionWatcherClient)(nil))
	CoredataEventClientName                 = di.TypeInstanceToName((*interfaces.EventClient)(nil))
)

func CommonClientFrom(get di.Get) interfaces.CommonClient {
	return get(CommonClientName).(interfaces.CommonClient)
}

func MetadataDeviceClientFrom(get di.Get) interfaces.DeviceClient {
	return get(MetadataDeviceClientName).(interfaces.DeviceClient)
}

func MetadataDeviceProfileClientFrom(get di.Get) interfaces.DeviceProfileClient {
	return get(MetadataDeviceProfileClientName).(interfaces.DeviceProfileClient)
}

func MetadataDeviceServiceClientFrom(get di.Get) interfaces.DeviceServiceClient {
	return get(MetadataDeviceServiceClientName).(interfaces.DeviceServiceClient)
}

func MetadataDeviceServiceCallbackClientFrom(get di.Get) interfaces.DeviceServiceCallbackClient {
	return get(MetadataDeviceServiceCallbackClientName).(interfaces.DeviceServiceCallbackClient)
}

func MetadataProvisionWatcherClientFrom(get di.Get) interfaces.ProvisionWatcherClient {
	return get(MetadataProvisionWatcherClientName).(interfaces.ProvisionWatcherClient)
}

func CoredataEventClientFrom(get di.Get) interfaces.EventClient {
	return get(CoredataEventClientName).(interfaces.EventClient)
}
