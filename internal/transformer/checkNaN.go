// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package transformer

import (
	"math"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/pkg/models"
)

// NaNError is used to throw the NaN error for the floating-point value
type NaNError struct{}

func (e NaNError) Error() string {
	return "not a valid float value NaN"
}

func isNaN(cv *models.CommandValue) (bool, error) {
	switch cv.Type {
	case contracts.ValueTypeFloat32:
		v, err := cv.Float32Value()
		if err != nil {
			return false, err
		}
		if math.IsNaN(float64(v)) {
			return true, nil
		}
	case contracts.ValueTypeFloat64:
		v, err := cv.Float64Value()
		if err != nil {
			return false, err
		}
		if math.IsNaN(v) {
			return true, nil
		}
	}
	return false, nil
}
