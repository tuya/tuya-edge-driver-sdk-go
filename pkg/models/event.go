// -- Mode: Go; indent-tabs-mode: t --
//
// Copyright (C) 2019 Intel Ltd
//
// SPDX-License-Identifier: Apache-2.0

package models

import (
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
)

// Event is a wrapper of contract.Event to provide more Binary related operation in Device Service.
type Event struct {
	models.Event
	EncodedEvent []byte
}

// HasBinaryValue confirms whether an event contains one or more
// readings populated with a BinaryValue payload.
func (e Event) HasBinaryValue() bool {
	if len(e.Readings) > 0 {
		for r := range e.Readings {
			if e.Readings[r].GetBaseReading().ValueType == contracts.ValueTypeBinary {
				return true
			}
		}
	}
	return false
}
