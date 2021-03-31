// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/mock"
)

var pws []models.ProvisionWatcher

func init() {
	serviceName := "watcher-cache-test"
	pwc := &mock.ProvisionWatcherClientMock{}
	ctx := context.WithValue(context.Background(), common.CorrelationHeader, uuid.New().String())
	resp, _ := pwc.ProvisionWatchersByServiceName(ctx, serviceName, 0, -1)
	for i := range resp.ProvisionWatchers {
		pws = append(pws, dtos.ToProvisionWatcherModel(resp.ProvisionWatchers[i]))
	}
}

func TestNewProvisionWatcherCache(t *testing.T) {
	pwc := newProvisionWatcherCache([]models.ProvisionWatcher{})
	if _, ok := pwc.(ProvisionWatcherCache); !ok {
		t.Errorf("the newProvisionWatcherCacher function supposed to return a value which holds the ProvisionWatcherCache type")
	}
}

func TestProvisionWatcherCache_ForName(t *testing.T) {
	pwc := newProvisionWatcherCache(pws)

	if w, found := pwc.ForName(mock.ValidBooleanWatcher.Name); !found {
		t.Error("supposed to find a matching watcher in cache by a valid watcher name")
	} else {
		assert.Equal(t, mock.ValidBooleanWatcher, w)
	}

	if _, found := pwc.ForName(mock.NewProvisionWatcher.Name); found {
		t.Error("not supposed to find a matching watcher by an invalid watcher name")
	}
}

func TestProvisionWatcherCache_ForId(t *testing.T) {
	pwc := newProvisionWatcherCache(pws)

	if w, found := pwc.ForId(mock.ValidBooleanWatcher.Id); !found {
		t.Error("supposed to find a matching watcher in cache by a valid watcher id")
	} else {
		assert.Equal(t, mock.ValidBooleanWatcher, w)
	}

	if _, found := pwc.ForId(mock.NewProvisionWatcher.Id); found {
		t.Error("not supposed to find a matching watcher by an invalid watcher id")
	}
}

func TestProvisionWatcherCache_All(t *testing.T) {
	pwc := newProvisionWatcherCache(pws)
	pwsFromCache := pwc.All()

	for _, pwFromCache := range pwsFromCache {
		for _, w := range pws {
			if pwFromCache.Id == w.Id {
				assert.Equal(t, w, pwFromCache)
				continue
			}
		}
	}
}

func TestProvisionWatcherCache_Add(t *testing.T) {
	pwc := newProvisionWatcherCache(pws)

	if err := pwc.Add(dtos.ToProvisionWatcherModel(mock.NewProvisionWatcher)); err != nil {
		t.Error("failed to add a new watcher to cache")
	}

	if w, found := pwc.ForId(mock.NewProvisionWatcher.Id); !found {
		t.Error("unable to find the watcher which just be added to cache")
	} else {
		assert.Equal(t, mock.NewProvisionWatcher, w)
	}

	if err := pwc.Add(dtos.ToProvisionWatcherModel(mock.DuplicateFloatWatcher)); err == nil {
		t.Error("supposed to get an error when adding a duplicate watcher to cache")
	}
}

func TestProvisionWatcherCache_RemoveByName(t *testing.T) {
	pwc := newProvisionWatcherCache(pws)

	if err := pwc.RemoveByName(mock.NewProvisionWatcher.Name); err == nil {
		t.Error("supposed to get an error when removing a watcher doesn't exist in cache")

		if err := pwc.RemoveByName(mock.ValidBooleanWatcher.Name); err != nil {
			t.Error("failed to remove watcher from cache by name")
		}

		if _, found := pwc.ForName(mock.ValidBooleanWatcher.Name); found {
			t.Error("unable to remove watcher from cache by name")
		}
	}
}

func TestProvisionWatcherCache_Remove(t *testing.T) {
	pwc := newProvisionWatcherCache(pws)

	if err := pwc.RemoveById(mock.NewProvisionWatcher.Id); err == nil {
		t.Error("supposed to get an error when removing a watcher which doesn't exist in cache")
	}

	if err := pwc.RemoveById(mock.ValidBooleanWatcher.Id); err != nil {
		t.Error("failed to remove watcher from cache by id")
	}

	if _, found := pwc.ForId(mock.ValidBooleanWatcher.Id); found {
		t.Error("unable to remove watcher from cache by id")
	}
}

func TestProvisionWatcherCache_Update(t *testing.T) {
	pwc := newProvisionWatcherCache(pws)

	if err := pwc.Update(dtos.ToProvisionWatcherModel(mock.NewProvisionWatcher)); err == nil {
		t.Error("supposed to get an error when updating a watcher which doesn't exist in cache")
	}

	mock.ValidBooleanWatcher.AdminState = models.Locked
	if err := pwc.Update(dtos.ToProvisionWatcherModel(mock.ValidBooleanWatcher)); err != nil {
		t.Error("failed to update watcher in cache")
	}

	if ud0, found := pwc.ForId(mock.ValidBooleanWatcher.Id); !found {
		t.Error("unable to find the watcher in cache after updating it")
	} else {
		assert.Equal(t, mock.ValidBooleanWatcher, ud0)
	}
}

func TestProvisionWatcherCache_UpdateAdminState(t *testing.T) {
	pwc := newProvisionWatcherCache(pws)

	if err := pwc.UpdateAdminState(mock.NewProvisionWatcher.Id, models.Locked); err == nil {
		t.Error("supposed to get an error when updating AdminState of the watcher which doesn't exist in cache")
	}
	if err := pwc.UpdateAdminState(mock.ValidBooleanWatcher.Id, models.Locked); err != nil {
		t.Error("failed to update AdminState")
	}
	if ud0, _ := pwc.ForId(mock.ValidBooleanWatcher.Id); ud0.AdminState != models.Locked {
		t.Error("succeeded in executing UpdateAdminState, but the value of AdminState was not updated")
	}
}
