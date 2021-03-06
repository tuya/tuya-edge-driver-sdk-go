// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package correlation

import (
	"context"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
)

func IdFromContext(ctx context.Context) string {
	hdr, ok := ctx.Value(contracts.CorrelationHeader).(string)
	if !ok {
		hdr = ""
	}
	return hdr
}

func LoggingClientFromContext(ctx context.Context) logger.LoggingClient {
	lc, ok := ctx.Value(bootstrapContainer.LoggingClientInterfaceName).(logger.LoggingClient)
	if !ok {
		lc = nil
	}
	return lc
}
