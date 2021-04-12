// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v2/config"
)

// ConfigurationStruct contains the configuration properties for the device service.
type ConfigurationStruct struct {
	// WritableInfo contains configuration settings that can be changed in the Registry .
	Writable WritableInfo
	// Clients is a map of services used by a DS.
	Clients map[string]config.ClientInfo
	// Registry contains registry-specific settings.
	Registry config.RegistryInfo
	// Service contains DeviceService-specific settings.
	Service ServiceInfo
	// Device contains device-specific configuration settings.
	Device DeviceInfo
	// DeviceList is the list of pre-define Devices
	DeviceList []DeviceConfig `consul:"-"`
	// Driver is a string map contains customized configuration for the protocol driver implemented based on Device SDK
	Driver map[string]interface{}
	// SecretStore contains information for connecting to the secure SecretStore (Vault) to retrieve or store secrets
	SecretStore config.SecretStoreInfo
}

// UpdateFromRaw converts configuration received from the registry to a service-specific configuration struct which is
// then used to overwrite the service's existing configuration struct.
func (c *ConfigurationStruct) UpdateFromRaw(rawConfig interface{}) bool {
	configuration, ok := rawConfig.(*ConfigurationStruct)
	if ok {
		// Check that information was successfully read from Registry
		if configuration.Service.Port == 0 {
			return false
		}
		*c = *configuration
	}
	return ok
}

// EmptyWritablePtr returns a pointer to a service-specific empty WritableInfo struct.  It is used by the bootstrap to
// provide the appropriate structure to registry.Client's WatchForChanges().
func (c *ConfigurationStruct) EmptyWritablePtr() interface{} {
	return &WritableInfo{}
}

// UpdateWritableFromRaw converts configuration received from the registry to a service-specific WritableInfo struct
// which is then used to overwrite the service's existing configuration's WritableInfo struct.
func (c *ConfigurationStruct) UpdateWritableFromRaw(rawWritable interface{}) bool {
	writable, ok := rawWritable.(*WritableInfo)
	if ok {
		c.Writable = *writable
	}
	return ok
}

// GetBootstrap returns the configuration elements required by the bootstrap.  Currently, a copy of the configuration
// data is returned.  This is intended to be temporary -- since ConfigurationStruct drives the configuration.toml's
// structure -- until we can make backwards-breaking configuration.toml changes (which would consolidate these fields
// into an config.BootstrapConfiguration struct contained within ConfigurationStruct).
func (c *ConfigurationStruct) GetBootstrap() config.BootstrapConfiguration {
	return config.BootstrapConfiguration{
		Clients:     c.Clients,
		Service:     c.Service.GetBootstrapServiceInfo(),
		Registry:    c.Registry,
		SecretStore: c.SecretStore,
	}
}

// GetLogLevel returns the current ConfigurationStruct's log level.
func (c *ConfigurationStruct) GetLogLevel() string {
	return c.Writable.LogLevel
}

// GetRegistryInfo gets the config.RegistryInfo field from the ConfigurationStruct.
func (c *ConfigurationStruct) GetRegistryInfo() config.RegistryInfo {
	return c.Registry
}

// GetInsecureSecrets returns the service's InsecureSecrets.
func (c *ConfigurationStruct) GetInsecureSecrets() config.InsecureSecrets {
	return c.Writable.InsecureSecrets
}
