//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/http/utils"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/interfaces"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/responses"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
)

type DeviceProfileClient struct {
	baseUrl string
}

// NewDeviceProfileClient creates an instance of DeviceProfileClient
func NewDeviceProfileClient(baseUrl string) interfaces.DeviceProfileClient {
	return &DeviceProfileClient{
		baseUrl: baseUrl,
	}
}

func (client *DeviceProfileClient) Add(ctx context.Context, reqs []requests.DeviceProfileRequest) ([]common.BaseWithIdResponse, errors.EdgeX) {
	var responses []common.BaseWithIdResponse
	err := utils.PostRequest(ctx, &responses, client.baseUrl+contracts.ApiDeviceProfileRoute, reqs)
	if err != nil {
		return responses, errors.NewCommonEdgeXWrapper(err)
	}
	return responses, nil
}

func (client *DeviceProfileClient) Update(ctx context.Context, reqs []requests.DeviceProfileRequest) ([]common.BaseResponse, errors.EdgeX) {
	var responses []common.BaseResponse
	err := utils.PutRequest(ctx, &responses, client.baseUrl+contracts.ApiDeviceProfileRoute, reqs)
	if err != nil {
		return responses, errors.NewCommonEdgeXWrapper(err)
	}
	return responses, nil
}

func (client *DeviceProfileClient) AddByYaml(ctx context.Context, yamlFilePath string) (common.BaseWithIdResponse, errors.EdgeX) {
	var responses common.BaseWithIdResponse
	err := utils.PostByFileRequest(ctx, &responses, client.baseUrl+contracts.ApiDeviceProfileUploadFileRoute, yamlFilePath)
	if err != nil {
		return responses, errors.NewCommonEdgeXWrapper(err)
	}
	return responses, nil
}

func (client *DeviceProfileClient) UpdateByYaml(ctx context.Context, yamlFilePath string) (common.BaseResponse, errors.EdgeX) {
	var responses common.BaseResponse
	err := utils.PutByFileRequest(ctx, &responses, client.baseUrl+contracts.ApiDeviceProfileUploadFileRoute, yamlFilePath)
	if err != nil {
		return responses, errors.NewCommonEdgeXWrapper(err)
	}
	return responses, nil
}

func (client *DeviceProfileClient) DeleteByName(ctx context.Context, name string) (common.BaseResponse, errors.EdgeX) {
	var response common.BaseResponse
	requestPath := path.Join(contracts.ApiDeviceProfileRoute, contracts.Name, url.QueryEscape(name))
	err := utils.DeleteRequest(ctx, &response, client.baseUrl, requestPath)
	if err != nil {
		return response, errors.NewCommonEdgeXWrapper(err)
	}
	return response, nil
}

func (client *DeviceProfileClient) DeviceProfileByName(ctx context.Context, name string) (res responses.DeviceProfileResponse, edgexError errors.EdgeX) {
	requestPath := path.Join(contracts.ApiDeviceProfileRoute, contracts.Name, url.QueryEscape(name))
	err := utils.GetRequest(ctx, &res, client.baseUrl, requestPath, nil)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (client *DeviceProfileClient) AllDeviceProfiles(ctx context.Context, labels []string, offset int, limit int) (res responses.MultiDeviceProfilesResponse, edgexError errors.EdgeX) {
	requestParams := url.Values{}
	if len(labels) > 0 {
		requestParams.Set(contracts.Labels, strings.Join(labels, contracts.CommaSeparator))
	}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	err := utils.GetRequest(ctx, &res, client.baseUrl, contracts.ApiAllDeviceProfileRoute, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (client *DeviceProfileClient) DeviceProfilesByModel(ctx context.Context, model string, offset int, limit int) (res responses.MultiDeviceProfilesResponse, edgexError errors.EdgeX) {
	requestPath := path.Join(contracts.ApiDeviceProfileRoute, contracts.Model, url.QueryEscape(model))
	requestParams := url.Values{}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	err := utils.GetRequest(ctx, &res, client.baseUrl, requestPath, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (client *DeviceProfileClient) DeviceProfilesByManufacturer(ctx context.Context, manufacturer string, offset int, limit int) (res responses.MultiDeviceProfilesResponse, edgexError errors.EdgeX) {
	requestPath := path.Join(contracts.ApiDeviceProfileRoute, contracts.Manufacturer, url.QueryEscape(manufacturer))
	requestParams := url.Values{}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	err := utils.GetRequest(ctx, &res, client.baseUrl, requestPath, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (client *DeviceProfileClient) DeviceProfilesByManufacturerAndModel(ctx context.Context, manufacturer string, model string, offset int, limit int) (res responses.MultiDeviceProfilesResponse, edgexError errors.EdgeX) {
	requestPath := path.Join(contracts.ApiDeviceProfileRoute, contracts.Manufacturer, url.QueryEscape(manufacturer), contracts.Model, url.QueryEscape(model))
	requestParams := url.Values{}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	err := utils.GetRequest(ctx, &res, client.baseUrl, requestPath, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}
