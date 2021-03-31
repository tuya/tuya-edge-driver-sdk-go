// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"context"
	"encoding/json"
	"path/filepath"
	"runtime"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/responses"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
)

var (
	ValidBooleanWatcher         = dtos.ProvisionWatcher{}
	ValidIntegerWatcher         = dtos.ProvisionWatcher{}
	ValidUnsignedIntegerWatcher = dtos.ProvisionWatcher{}
	ValidFloatWatcher           = dtos.ProvisionWatcher{}
	DuplicateFloatWatcher       = dtos.ProvisionWatcher{}
	NewProvisionWatcher         = dtos.ProvisionWatcher{}
)

type ProvisionWatcherClientMock struct {
}

// Add adds a new provision watcher.
func (ProvisionWatcherClientMock) Add(ctx context.Context, reqs []requests.AddProvisionWatcherRequest) ([]common.BaseWithIdResponse, errors.EdgeX) {
	panic("implement me")
}

// Update updates provision watchers.
func (ProvisionWatcherClientMock) Update(ctx context.Context, reqs []requests.UpdateProvisionWatcherRequest) ([]common.BaseResponse, errors.EdgeX) {
	panic("implement me")
}

// AllProvisionWatchers returns all provision watchers. ProvisionWatchers can also be filtered by labels.
// The result can be limited in a certain range by specifying the offset and limit parameters.
// offset: The number of items to skip before starting to collect the result set. Default is 0.
// limit: The number of items to return. Specify -1 will return all remaining items after offset. The maximum will be the MaxResultCount as defined in the configuration of service. Default is 20.
func (ProvisionWatcherClientMock) AllProvisionWatchers(ctx context.Context, labels []string, offset int, limit int) (responses.MultiProvisionWatchersResponse, errors.EdgeX) {
	panic("implement me")
}

// ProvisionWatcherByName returns a provision watcher by name.
func (ProvisionWatcherClientMock) ProvisionWatcherByName(ctx context.Context, name string) (responses.ProvisionWatcherResponse, errors.EdgeX) {
	panic("implement me")
}

// DeleteProvisionWatcherByName deletes a provision watcher by name.
func (ProvisionWatcherClientMock) DeleteProvisionWatcherByName(ctx context.Context, name string) (common.BaseResponse, errors.EdgeX) {
	panic("implement me")
}

// ProvisionWatchersByProfileName returns provision watchers associated with the specified device profile name.
// The result can be limited in a certain range by specifying the offset and limit parameters.
// offset: The number of items to skip before starting to collect the result set. Default is 0.
// limit: The number of items to return. Specify -1 will return all remaining items after offset. The maximum will be the MaxResultCount as defined in the configuration of service. Default is 20.
func (ProvisionWatcherClientMock) ProvisionWatchersByProfileName(ctx context.Context, name string, offset int, limit int) (responses.MultiProvisionWatchersResponse, errors.EdgeX) {
	panic("implement me")
}

// ProvisionWatchersByServiceName returns provision watchers associated with the specified device service name.
// The result can be limited in a certain range by specifying the offset and limit parameters.
// offset: The number of items to skip before starting to collect the result set. Default is 0.
// limit: The number of items to return. Specify -1 will return all remaining items after offset. The maximum will be the MaxResultCount as defined in the configuration of service. Default is 20.
func (ProvisionWatcherClientMock) ProvisionWatchersByServiceName(ctx context.Context, name string, offset int, limit int) (responses.MultiProvisionWatchersResponse, errors.EdgeX) {
	err := populateProvisionWatcherMock()
	if err != nil {
		return responses.MultiProvisionWatchersResponse{}, errors.NewCommonEdgeXWrapper(err)
	}
	return responses.MultiProvisionWatchersResponse{
		BaseResponse: common.NewBaseResponse("", "", 0),
		ProvisionWatchers: []dtos.ProvisionWatcher{
			ValidBooleanWatcher,
			ValidIntegerWatcher,
			ValidUnsignedIntegerWatcher,
			ValidFloatWatcher,
		},
	}, nil
}

func populateProvisionWatcherMock() error {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	watchers, err := loadData(basepath + "/data/provisionwatcher")
	if err != nil {
		return err
	}

	_ = json.Unmarshal(watchers[WatcherBool], &ValidBooleanWatcher)
	_ = json.Unmarshal(watchers[WatcherInt], &ValidIntegerWatcher)
	_ = json.Unmarshal(watchers[WatcherUint], &ValidUnsignedIntegerWatcher)
	_ = json.Unmarshal(watchers[WatcherFloat], &ValidFloatWatcher)
	_ = json.Unmarshal(watchers[WatcherFloat], &DuplicateFloatWatcher)
	_ = json.Unmarshal(watchers[WatcherNew], &NewProvisionWatcher)

	return nil
}
