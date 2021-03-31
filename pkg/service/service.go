// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2020 IOTech Ltd
// Copyright (c) 2019 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a basic EdgeX Foundry device service implementation
// meant to be embedded in an command, similar in approach to the builtin
// net/http package.
package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/config"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	eErr "github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/autoevent"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/clients"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/container"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/controller"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
	dsModels "github.com/tuya/tuya-edge-driver-sdk-go/pkg/models"
)

var (
	ds *DeviceService
)

type DeviceService struct {
	ServiceName   string
	LoggingClient logger.LoggingClient
	tedgeClients  clients.TedgeClients
	controller    *controller.RestController
	config        *common.ConfigurationStruct
	deviceService models.DeviceService
	driver        dsModels.ProtocolDriver
	discovery     dsModels.ProtocolDiscovery
	asyncCh       chan *dsModels.AsyncValues
	deviceCh      chan []dsModels.DiscoveredDevice
	initialized   bool
}

func (s *DeviceService) Initialize(serviceName, serviceVersion string, proto interface{}) {
	if serviceName == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Please specify device service name")
		os.Exit(1)
	}
	s.ServiceName = serviceName

	if serviceVersion == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Please specify device service version")
		os.Exit(1)
	}
	common.ServiceVersion = serviceVersion

	if driver, ok := proto.(dsModels.ProtocolDriver); ok {
		s.driver = driver
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Please implement and specify the protocoldriver")
		os.Exit(1)
	}

	if discovery, ok := proto.(dsModels.ProtocolDiscovery); ok {
		s.discovery = discovery
	} else {
		s.discovery = nil
	}

	s.config = &common.ConfigurationStruct{}
}

func (s *DeviceService) UpdateFromContainer(r *mux.Router, dic *di.Container) {
	s.LoggingClient = bootstrapContainer.LoggingClientFrom(dic.Get)
	// v2
	s.tedgeClients.CommonClient = container.CommonClientFrom(dic.Get)
	s.tedgeClients.DeviceClient = container.MetadataDeviceClientFrom(dic.Get)
	s.tedgeClients.DeviceProfileClient = container.MetadataDeviceProfileClientFrom(dic.Get)
	s.tedgeClients.DeviceServiceClient = container.MetadataDeviceServiceClientFrom(dic.Get)
	s.tedgeClients.CallbackClient = container.MetadataDeviceServiceCallbackClientFrom(dic.Get)
	s.tedgeClients.ProvisionWatcherClient = container.MetadataProvisionWatcherClientFrom(dic.Get)

	s.config = container.ConfigurationFrom(dic.Get)
	s.controller = controller.NewRestController(r, dic)
}

// Name returns the name of this Device Service
func (s *DeviceService) Name() string {
	return s.ServiceName
}

// Version returns the version number of this Device Service
func (s *DeviceService) Version() string {
	return common.ServiceVersion
}

// AsyncReadings returns a bool value to indicate whether the asynchronous reading is enabled.
func (s *DeviceService) AsyncReadings() bool {
	return s.config.Service.EnableAsyncReadings
}

func (s *DeviceService) DeviceDiscovery() bool {
	return s.config.Device.Discovery.Enabled
}

// AddRoute allows leveraging the existing internal web server to add routes specific to Device Service.
func (s *DeviceService) AddRoute(route string, handler func(http.ResponseWriter, *http.Request), methods ...string) error {
	return s.controller.AddRoute(route, handler, methods...)
}

// Stop shuts down the Service
func (s *DeviceService) Stop(force bool) {
	if s.initialized {
		_ = s.driver.Stop(false)
	}
	autoevent.GetManager().StopAutoEvents()
}

/*
func (s *DeviceService) selfRegister() eErr.EdgeX {
	// 1 search
	ctx := context.WithValue(context.Background(), common.CorrelationHeader, uuid.New().String())
	s.LoggingClient.Debug("Trying to find  DeviceService: " + s.ServiceName)
	dsr, err := s.tedgeClients.DeviceServiceClient.DeviceServiceByName(ctx, s.ServiceName)
	if err != nil { // 3. not exists add new
		if errsc, ok := err.(eErr.EdgeX); ok && (errsc.Code() == http.StatusNotFound) {
			s.LoggingClient.Info(fmt.Sprintf(" DeviceService %s doesn't exist, creating a new one", s.ServiceName))
			ba := config.ClientInfo{
				Host:     s.config.Service.Host,
				Port:     s.config.Service.Port,
				Protocol: s.config.Service.Protocol,
			}
			newDeviceService := dtos.DeviceService{
				//Name: s.ServiceName,
				Labels:          s.config.Service.Labels,
				BaseAddress:     ba.Url(),
				AdminState:      models.Unlocked,
				DeviceLibraryId: s.config.Service.DeviceLibraryId,
				Config: models.DeviceServiceConfig{
					Server: models.DeviceServiceConfigServer{
						Protocol: s.config.Service.Protocol,
						Port:     s.config.Service.Port,
					},
				},
			}
			req := requests.AddDeviceServiceRequest{
				Service: newDeviceService,
			}
			reqs := make([]requests.AddDeviceServiceRequest, 0, 1)
			reqs = append(reqs, req)
			resp, err := s.tedgeClients.DeviceServiceClient.Add(ctx, reqs)
			if err != nil {
				s.LoggingClient.Error(fmt.Sprintf("Failed to add  Deviceservice %s: %v", s.ServiceName, err))
				return err
			}
			if err := common.VerifyIdFormat(resp[0].Id, "Device Service"); err != nil {
				return err
			}
			// NOTE - this differs from Addressable and Device Resources,
			// neither of which require the '.Service'prefix
			newDeviceService.Id = resp[0].Id
			// TODO device service by name again
			s.deviceService = newDeviceService
			s.LoggingClient.Debug("New DeviceService Id: " + newDeviceService.Id)
		} else {
			s.LoggingClient.Error(fmt.Sprintf("DeviceServicForName failed: %v", err))
			return err
		}
	} else { // 2. exists update
		s.LoggingClient.Info(fmt.Sprintf("DeviceService %s exists, updating it", dsr.Service.Name))
		id := dsr.Service.Id
		name := s.ServiceName
		unlock := models.Unlocked
		ba := s.config.Clients[common.ClientMetadata].Url()
		update := dtos.UpdateDeviceService{
			Id:          &id,
			Name:        &name,
			Labels:      s.config.Service.Labels,
			BaseAddress: &ba,
			AdminState:  &unlock,
		}
		req := requests.UpdateDeviceServiceRequest{
			Service: update,
		}
		reqs := make([]requests.UpdateDeviceServiceRequest, 0, 1)
		reqs = append(reqs, req)
		resp, err := s.tedgeClients.DeviceServiceClient.Update(ctx, reqs)
		if err != nil {
			s.LoggingClient.Error(fmt.Sprintf("Failed to update  DeviceService %s: %v, %v", s.deviceService.Id, err, resp))
			return err
		}
		s.deviceService.Labels = s.config.Service.Labels
		s.deviceService.BaseAddress = ba
		s.deviceService.AdminState = unlock
	}
	return nil
}
*/

func (s *DeviceService) updateService() eErr.EdgeX {
	// 获取驱动实例ID
	id := s.config.Service.ID
	if id == "" {
		s.LoggingClient.Error("device service instance id is required")
		return eErr.NewCommonEdgeXWrapper(errors.New("device service instance id is required"))
	}
	// 获取配置
	ctx := context.WithValue(context.Background(), common.CorrelationHeader, uuid.New().String())
	s.LoggingClient.Debug("Trying to find  DeviceService with instance id: " + id)
	dsr, err := s.tedgeClients.DeviceServiceClient.DeviceServiceByID(ctx, id)
	if err != nil { // not exists, then return
		s.LoggingClient.Error(fmt.Sprintf("DeviceServic instance by id find failed: %v", err))
		return err
	} else { // exists update
		svc := dtos.ToDeviceServiceModel(dsr.Service)
		s.LoggingClient.Info(fmt.Sprintf("DeviceService instance with id: %s exists, updating it", dsr.Service.Id))

		unlock := models.Unlocked
		baInfo := config.ClientInfo{
			Host:     s.config.Service.Host,
			Port:     s.config.Service.Port,
			Protocol: s.config.Service.Protocol,
		}
		ba := baInfo.Url()
		update := dtos.UpdateDeviceService{
			Id:          &svc.Id,
			Name:        &svc.Name,
			Labels:      s.config.Service.Labels,
			BaseAddress: &ba,
			AdminState:  &unlock,
		}
		req := requests.UpdateDeviceServiceRequest{
			Service: update,
		}
		reqs := make([]requests.UpdateDeviceServiceRequest, 0, 1)
		reqs = append(reqs, req)
		resp, err := s.tedgeClients.DeviceServiceClient.Update(ctx, reqs)
		if err != nil {
			s.LoggingClient.Error(fmt.Sprintf("Failed to update  DeviceService %s: %v, %v", s.deviceService.Id, err, resp))
			return err
		}
		// update local config
		svc.Labels = s.config.Service.Labels
		svc.BaseAddress = ba
		svc.AdminState = models.Unlocked
		s.deviceService = svc
	}
	return nil
}

func RunningService() *DeviceService {
	return ds
}

// DriverConfigs retrieves the driver specific configuration
func DriverConfigs() map[string]string {
	return ds.config.Driver
}
