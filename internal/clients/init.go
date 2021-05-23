// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2020 IOTech Ltd
// Copyright (c) 2019 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	tdHttp "github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/http"
	dtCommon "github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/container"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
)

// Clients contains references to dependencies required by the Clients bootstrap implementation.
type Clients struct {
}

// NewClients create a new instance of Clients
func NewClients() *Clients {
	return &Clients{}
}

func (_ *Clients) BootstrapHandler(
	ctx context.Context,
	wg *sync.WaitGroup,
	startupTimer startup.Timer,
	dic *di.Container) bool {
	return InitDependencyClients(ctx, startupTimer, dic)
}

// InitDependencyClients triggers Service Client Initializer to establish connection to Metadata and Core Data Services
// through Metadata Client and Core Data Client.
// Service Client Initializer also needs to check the service status of Metadata and Core Data Services,
// because they are important dependencies of Device Service.
// The initialization process should be pending until Metadata Service and Core Data Service are both available.
func InitDependencyClients(ctx context.Context, startupTimer startup.Timer, dic *di.Container) bool {
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)

	if err := validateClientConfig(container.ConfigurationFrom(dic.Get)); err != nil {
		lc.Error(err.Error())
		return false
	}

	if CheckDependencyServices(ctx, startupTimer, dic) == false {
		return false
	}

	initializeClientsClients(dic)

	lc.Info("Service clients initialize successful.")
	return true
}

func validateClientConfig(configuration *common.ConfigurationStruct) error {

	if len(configuration.Clients[common.ClientMetadata].Host) == 0 {
		return fmt.Errorf("fatal error; Host setting for Core Metadata client not configured")
	}

	if configuration.Clients[common.ClientMetadata].Port == 0 {
		return fmt.Errorf("fatal error; Port setting for Core Metadata client not configured")
	}

	if len(configuration.Clients[common.ClientData].Host) == 0 {
		return fmt.Errorf("fatal error; Host setting for Core Data client not configured")
	}

	if configuration.Clients[common.ClientData].Port == 0 {
		return fmt.Errorf("fatal error; Port setting for Core Ddata client not configured")
	}

	// TODO: validate other settings for sanity: maxcmdops, ...

	return nil
}

func CheckDependencyServices(ctx context.Context, startupTimer startup.Timer, dic *di.Container) bool {
	var dependencyList = []string{common.ClientData, common.ClientMetadata}
	var waitGroup sync.WaitGroup
	checkingErr := true

	dependencyCount := len(dependencyList)
	waitGroup.Add(dependencyCount)

	for i := 0; i < dependencyCount; i++ {
		go func(wg *sync.WaitGroup, serviceName string) {
			defer wg.Done()
			if checkServiceAvailable(ctx, serviceName, startupTimer, dic) == false {
				checkingErr = false
			}
		}(&waitGroup, dependencyList[i])
	}
	waitGroup.Wait()

	return checkingErr
}

// ping检测
func checkServiceAvailable(ctx context.Context, serviceId string, startupTimer startup.Timer, dic *di.Container) bool {
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)

	for startupTimer.HasNotElapsed() {
		select {
		case <-ctx.Done():
			return false
		default:
			configuration := container.ConfigurationFrom(dic.Get)
			if checkServiceAvailableByPing(serviceId, configuration, lc) == nil {
				return true
			}
			startupTimer.SleepForInterval()
		}
	}

	lc.Error(fmt.Sprintf("dependency %s service checking time out", serviceId))
	return false
}

func checkServiceAvailableByPing(serviceId string, configuration *common.ConfigurationStruct, lc logger.LoggingClient) error {
	lc.Info(fmt.Sprintf("Check %v service's status by ping...", serviceId))
	addr := configuration.Clients[serviceId].Url()
	timeout := int64(configuration.Service.Timeout) * int64(time.Millisecond)

	client := http.Client{
		Timeout: time.Duration(timeout),
	}

	resp, err := client.Get(addr + contracts.ApiPingRoute)
	if err != nil {
		lc.Error(err.Error())
		return err
	}
	defer resp.Body.Close()
	var (
		body  []byte
		pResp dtCommon.PingResponse
	)
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		lc.Error("read response body error", err.Error())
		return err
	}
	if err = json.Unmarshal(body, &pResp); err != nil {
		lc.Error("unmarshal response body error", err.Error())
		return err
	}
	lc.Info(fmt.Sprintf("Check %v service's response: %+v", serviceId, pResp))
	return err
}

// 初始化v2版本需要的客户端
func initializeClientsClients(dic *di.Container) {
	configuration := container.ConfigurationFrom(dic.Get)
	cdBaseUrl := configuration.Clients[common.ClientMetadata].Url()
	dBaseUrl := configuration.Clients[common.ClientData].Url()

	cc := tdHttp.NewCommonClient(cdBaseUrl)
	dcV2 := tdHttp.NewDeviceClient(cdBaseUrl)
	dpcV2 := tdHttp.NewDeviceProfileClient(cdBaseUrl)
	dscV2 := tdHttp.NewDeviceServiceClient(cdBaseUrl)
	dsccV2 := tdHttp.NewDeviceServiceCallbackClient(cdBaseUrl)
	pwcV2 := tdHttp.NewProvisionWatcherClient(cdBaseUrl)

	ecV2 := tdHttp.NewEventClient(dBaseUrl)

	dic.Update(di.ServiceConstructorMap{
		container.CommonClientName: func(get di.Get) interface{} {
			return cc
		},
		container.MetadataDeviceClientName: func(get di.Get) interface{} {
			return dcV2
		},
		container.MetadataDeviceProfileClientName: func(get di.Get) interface{} {
			return dpcV2
		},
		container.MetadataDeviceServiceClientName: func(get di.Get) interface{} {
			return dscV2
		},
		container.MetadataDeviceServiceCallbackClientName: func(get di.Get) interface{} {
			return dsccV2
		},
		container.MetadataProvisionWatcherClientName: func(get di.Get) interface{} {
			return pwcV2
		},
		container.CoredataEventClientName: func(get di.Get) interface{} {
			return ecV2
		},
	})
}
