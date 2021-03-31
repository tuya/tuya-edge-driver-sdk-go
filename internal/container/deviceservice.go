// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"

	contract "github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	"github.com/tuya/tuya-edge-driver-sdk-go/pkg/models"
)

var (
	DeviceServiceName     = di.TypeInstanceToName(contract.DeviceService{})
	ProtocolDiscoveryName = di.TypeInstanceToName((*models.ProtocolDiscovery)(nil))
	ProtocolDriverName    = di.TypeInstanceToName((*models.ProtocolDriver)(nil))
)

func DeviceServiceFrom(get di.Get) contract.DeviceService {
	return get(DeviceServiceName).(contract.DeviceService)
}

func ProtocolDiscoveryFrom(get di.Get) models.ProtocolDiscovery {
	casted, ok := get(ProtocolDiscoveryName).(models.ProtocolDiscovery)
	if ok {
		return casted
	}
	return nil
}

func ProtocolDriverFrom(get di.Get) models.ProtocolDriver {
	return get(ProtocolDriverName).(models.ProtocolDriver)
}
