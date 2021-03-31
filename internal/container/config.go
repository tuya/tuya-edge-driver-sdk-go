// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"

	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
)

var ConfigurationName = di.TypeInstanceToName(common.ConfigurationStruct{})

// ConfigurationFrom helper function queries the DIC and returns device service's ConfigurationStruct implementation.
func ConfigurationFrom(get di.Get) *common.ConfigurationStruct {
	return get(ConfigurationName).(*common.ConfigurationStruct)
}
