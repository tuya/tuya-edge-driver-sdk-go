// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package autoevent

import (
	"testing"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
)

func TestCompareReadings(t *testing.T) {
	readings := make([]dtos.BaseReading, 4)
	readings[0] = dtos.BaseReading{
		ResourceName:  "Temperature",
		ValueType:     contracts.ValueTypeString,
		SimpleReading: dtos.SimpleReading{Value: "10"},
	}
	readings[1] = dtos.BaseReading{
		ResourceName:  "Humidity",
		ValueType:     contracts.ValueTypeString,
		SimpleReading: dtos.SimpleReading{Value: "50"},
	}
	readings[2] = dtos.BaseReading{
		ResourceName:  "Pressure",
		ValueType:     contracts.ValueTypeString,
		SimpleReading: dtos.SimpleReading{Value: "3"},
	}
	readings[3] = dtos.BaseReading{
		ResourceName:  "Image",
		ValueType:     contracts.ValueTypeBinary,
		BinaryReading: dtos.BinaryReading{BinaryValue: []byte("This is a image")},
	}

	lc := logger.NewMockClient()
	autoEvent := models.AutoEvent{Frequency: "500ms"}
	e, err := NewExecutor("hasBinaryTrue", autoEvent)
	if err != nil {
		t.Errorf("Autoevent executor creation failed: %v", err)
	}
	resultFalse := compareReadings(e, readings, lc)
	if resultFalse {
		t.Error("compare readings with cache failed, the result should be false in the first place")
	}

	readings[1] = dtos.BaseReading{
		ResourceName:  "Humidity",
		ValueType:     contracts.ValueTypeString,
		SimpleReading: dtos.SimpleReading{Value: "51"},
	}
	resultFalse = compareReadings(e, readings, lc)
	if resultFalse {
		t.Error("compare readings with cache failed, the result should be false")
	}

	readings[3] = dtos.BaseReading{
		ResourceName:  "Image",
		ValueType:     contracts.ValueTypeBinary,
		BinaryReading: dtos.BinaryReading{BinaryValue: []byte("This is not a image")},
	}
	resultFalse = compareReadings(e, readings, lc)
	if resultFalse {
		t.Error("compare readings with cache failed, the result should be false")
	}

	resultTrue := compareReadings(e, readings, lc)
	if !resultTrue {
		t.Error("compare readings with cache failed, the result should be true with unchanged readings")
	}

	e, err = NewExecutor("hasBinaryFalse", autoEvent)
	if err != nil {
		t.Errorf("Autoevent executor creation failed: %v", err)
	}
	// This scenario should not happen in real case
	resultFalse = compareReadings(e, readings, lc)
	if resultFalse {
		t.Error("compare readings with cache failed, the result should be false in the first place")
	}

	readings[0] = dtos.BaseReading{
		ResourceName:  "Temperature",
		ValueType:     contracts.ValueTypeString,
		SimpleReading: dtos.SimpleReading{Value: "20"},
	}
	resultFalse = compareReadings(e, readings, lc)
	if resultFalse {
		t.Error("compare readings with cache failed, the result should be false")
	}

	readings[3] = dtos.BaseReading{
		ResourceName:  "Image",
		ValueType:     contracts.ValueTypeBinary,
		BinaryReading: dtos.BinaryReading{BinaryValue: []byte("This is a image")},
	}
	resultTrue = compareReadings(e, readings, lc)
	if !resultTrue {
		t.Error("compare readings with cache failed, the result should always be true in such scenario")
	}

	resultTrue = compareReadings(e, readings, lc)
	if !resultTrue {
		t.Error("compare readings with cache failed, the result should be true with unchanged readings")
	}
}
