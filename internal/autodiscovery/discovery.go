// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package autodiscovery

import (
	"fmt"
	"sync"

	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
	dsModels "github.com/tuya/tuya-edge-driver-sdk-go/pkg/models"
)

type discoveryLocker struct {
	busy bool
	mux  sync.Mutex
}

var locker discoveryLocker

func DiscoveryWrapper(discovery dsModels.ProtocolDiscovery, lc logger.LoggingClient) {
	locker.mux.Lock()
	if locker.busy {
		lc.Info("another device discovery process is currently running")
		locker.mux.Unlock()
		return
	}
	locker.busy = true
	locker.mux.Unlock()

	lc.Debug(fmt.Sprintf("protocol discovery triggered"))
	discovery.Discover()

	// ReleaseLock
	locker.mux.Lock()
	locker.busy = false
	locker.mux.Unlock()
}
