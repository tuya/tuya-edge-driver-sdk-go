// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package startup

import (
	"context"

	"github.com/gorilla/mux"

	"github.com/tuya/tuya-edge-driver-sdk-go/pkg/service"
)

// TODO change error code

func Bootstrap(serviceName string, serviceVersion string, driver interface{}) {
	ctx, cancel := context.WithCancel(context.Background())
	service.Main(serviceName, serviceVersion, driver, ctx, cancel, mux.NewRouter())
}
