// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package callback

import (
	"context"
	"fmt"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/autoevent"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/cache"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/container"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
)

func UpdateProfile(profileRequest requests.DeviceProfileRequest, lc logger.LoggingClient) errors.EdgeX {
	_, ok := cache.Profiles().ForName(profileRequest.Profile.Name)
	if !ok {
		errMsg := fmt.Sprintf("failed to find profile %s", profileRequest.Profile.Name)
		return errors.NewCommonEdgeX(errors.KindInvalidId, errMsg, nil)
	}

	err := cache.Profiles().Update(dtos.ToDeviceProfileModel(profileRequest.Profile))
	if err != nil {
		errMsg := fmt.Sprintf("failed to update profile %s", profileRequest.Profile.Name)
		return errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
	}

	lc.Debug(fmt.Sprintf("profile %s updated", profileRequest.Profile.Name))
	return nil
}

func AddDevice(addDeviceRequest requests.AddDeviceRequest, dic *di.Container) errors.EdgeX {
	device := dtos.ToDeviceModel(addDeviceRequest.Device)
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)

	edgexErr := updateAssociatedProfile(device.ProfileName, dic)
	if edgexErr != nil {
		errMsg := fmt.Sprintf("failed to update device profile %s", device.ProfileName)
		return errors.NewCommonEdgeX(errors.KindServerError, errMsg, edgexErr)
	}

	edgexErr = cache.Devices().Add(device)
	if edgexErr != nil {
		errMsg := fmt.Sprintf("failed to add device %s", device.Name)
		return errors.NewCommonEdgeX(errors.KindServerError, errMsg, edgexErr)
	}
	lc.Debug(fmt.Sprintf("device %s added", device.Name))

	driver := container.ProtocolDriverFrom(dic.Get)
	err := driver.AddDevice(device.Name, device.Protocols, device.AdminState)
	if err == nil {
		lc.Debug(fmt.Sprintf("Invoked driver.AddDevice callback for %s", device.Name))
	} else {
		errMsg := fmt.Sprintf("driver.AddDevice callback failed for %s", device.Name)
		return errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
	}

	lc.Debug(fmt.Sprintf("Handler - starting AutoEvents for device %s", device.Name))
	autoevent.GetManager().RestartForDevice(device.Name, dic)
	return nil
}

func UpdateDevice(updateDeviceRequest requests.UpdateDeviceRequest, dic *di.Container) errors.EdgeX {
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)

	device, ok := cache.Devices().ForName(*updateDeviceRequest.Device.Name)
	if !ok {
		errMsg := fmt.Sprintf("failed to find device %s", *updateDeviceRequest.Device.Name)
		return errors.NewCommonEdgeX(errors.KindInvalidId, errMsg, nil)
	}

	requests.ReplaceDeviceModelFieldsWithDTO(&device, updateDeviceRequest.Device)
	// TODO: uncomment when core-contracts v2 client is ready.
	edgexErr := updateAssociatedProfile(device.ProfileName, dic)
	if edgexErr != nil {
		errMsg := fmt.Sprintf("failed to update device profile %s", device.ProfileName)
		return errors.NewCommonEdgeX(errors.KindServerError, errMsg, edgexErr)
	}

	edgexErr = cache.Devices().Update(device)
	if edgexErr != nil {
		errMsg := fmt.Sprintf("failed to update device %s", device.Name)
		return errors.NewCommonEdgeX(errors.KindServerError, errMsg, edgexErr)
	}
	lc.Debug(fmt.Sprintf("device %s updated", device.Name))

	driver := container.ProtocolDriverFrom(dic.Get)
	err := driver.UpdateDevice(device.Name, device.Protocols, device.AdminState)
	if err == nil {
		lc.Debug(fmt.Sprintf("Invoked driver.UpdateDevice callback for %s", device.Name))
	} else {
		errMsg := fmt.Sprintf("driver.UpdateDevice callback failed for %s", device.Name)
		return errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
	}

	lc.Debug(fmt.Sprintf("Handler - starting AutoEvents for device %s", device.Name))
	autoevent.GetManager().RestartForDevice(device.Name, dic)
	return nil
}

func DeleteDevice(name string, dic *di.Container) errors.EdgeX {
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)
	// check the device exist and stop its autoevents
	device, ok := cache.Devices().ForName(name)
	if ok {
		lc.Debugf("Handler - stopping AutoEvents for device %s", device.Name)
		autoevent.GetManager().StopForDevice(device.Name)
	} else {
		errMsg := fmt.Sprintf("failed to find device %s", name)
		return errors.NewCommonEdgeX(errors.KindInvalidId, errMsg, nil)
	}

	// remove the device in cache
	edgexErr := cache.Devices().RemoveByName(name)
	if edgexErr != nil {
		errMsg := fmt.Sprintf("failed to remove device %s", device.Name)
		return errors.NewCommonEdgeX(errors.KindServerError, errMsg, edgexErr)
	}
	lc.Debugf("Removed device: %s", device.Name)

	driver := container.ProtocolDriverFrom(dic.Get)
	err := driver.RemoveDevice(device.Name, device.Protocols)
	if err == nil {
		lc.Debugf("Invoked driver.RemoveDevice callback for %s", device.Name)
	} else {
		errMsg := fmt.Sprintf("driver.RemoveDevice callback failed for %s", device.Name)
		return errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
	}

	// a special case in which user updates the device profile after deleting all
	// devices in metadata, the profile won't be correctly updated because metadata
	// does not know which device service callback it needs to call. Remove the unused
	// device profile in cache so that if it is updated in metadata, next time the
	// device using it is added/updated, the cache can receive the updated one as well.
	if cache.CheckProfileNotUsed(device.ProfileName) {
		edgexErr = cache.Profiles().RemoveByName(device.ProfileName)
		if edgexErr != nil {
			lc.Warn("failed to remove unused profile", edgexErr.DebugMessages())
		}
	}

	return nil
}

func AddProvisionWatcher(addProvisionWatcherRequest requests.AddProvisionWatcherRequest, lc logger.LoggingClient) errors.EdgeX {
	provisionWatcher := dtos.ToProvisionWatcherModel(addProvisionWatcherRequest.ProvisionWatcher)

	edgexErr := cache.ProvisionWatchers().Add(provisionWatcher)
	if edgexErr != nil {
		errMsg := fmt.Sprintf("failed to add provision watcher %s", provisionWatcher.Name)
		return errors.NewCommonEdgeX(errors.KindServerError, errMsg, edgexErr)
	}

	lc.Debugf("provision watcher %s added", provisionWatcher.Name)
	return nil
}

func UpdateProvisionWatcher(updateProvisionWatcherRequest requests.UpdateProvisionWatcherRequest, lc logger.LoggingClient) errors.EdgeX {
	provisionWatcher, ok := cache.ProvisionWatchers().ForName(*updateProvisionWatcherRequest.ProvisionWatcher.Name)
	if !ok {
		errMsg := fmt.Sprintf("failed to find provision watcher %s", *updateProvisionWatcherRequest.ProvisionWatcher.Name)
		return errors.NewCommonEdgeX(errors.KindInvalidId, errMsg, nil)
	}

	requests.ReplaceProvisionWatcherModelFieldsWithDTO(&provisionWatcher, updateProvisionWatcherRequest.ProvisionWatcher)

	edgexErr := cache.ProvisionWatchers().Update(provisionWatcher)
	if edgexErr != nil {
		errMsg := fmt.Sprintf("failed to update provision watcher %s", provisionWatcher.Name)
		return errors.NewCommonEdgeX(errors.KindServerError, errMsg, edgexErr)
	}

	lc.Debugf("provision watcher %s updated", provisionWatcher.Name)
	return nil
}

func DeleteProvisionWatcher(name string, lc logger.LoggingClient) errors.EdgeX {
	err := cache.ProvisionWatchers().RemoveByName(name)
	if err != nil {
		errMsg := fmt.Sprintf("failed to remove provision watcher %s", name)
		return errors.NewCommonEdgeX(errors.KindInvalidId, errMsg, err)
	}

	lc.Debugf("removed provision watcher %s", name)
	return nil
}

// updateAssociatedProfile updates the profile specified in AddDeviceRequest or UpdateDeviceRequest
// to stay consistent with core metadata.
func updateAssociatedProfile(profileName string, dic *di.Container) errors.EdgeX {
	// TODO: uncomment when core-contracts v2 clients are ready.
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)
	//gc := container.GeneralClientFrom(dic.Get)
	//vdc := container.CoredataValueDescriptorClientFrom(dic.Get)
	dpc := container.MetadataDeviceProfileClientFrom(dic.Get)
	lc.Debugf("get profile: ", profileName)

	resp, err := dpc.DeviceProfileByName(context.Background(), profileName)
	if err != nil {
		errMsg := fmt.Sprintf("failed to find profile %s in metadata", profileName)
		return errors.NewCommonEdgeX(errors.KindInvalidId, errMsg, nil)
	}
	fmt.Printf("%+v\n", resp)
	_, exist := cache.Profiles().ForName(profileName)
	if exist == false {
		err = cache.Profiles().Add(dtos.ToDeviceProfileModel(resp.Profile))
		if err == nil {
			//provision.CreateDescriptorsFromProfile(&profile, lc, gc, vdc)
			//lc.Info(fmt.Sprintf("Added device profile: %s", profileName))
		} else {
			errMsg := fmt.Sprintf("failed to add profile %s", profileName)
			return errors.NewCommonEdgeX(errors.KindServerError, errMsg, err)
		}
	} else {
		err := cache.Profiles().Update(dtos.ToDeviceProfileModel(resp.Profile))
		if err != nil {
			lc.Warn(fmt.Sprintf("failed to to update profile %s in cache, using the original one", profileName))
		}
	}

	return nil
}
