//
// Copyright (C) 2021 IOTech Ltd
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

type ProvisionWatcherClient struct {
	baseUrl string
}

// NewProvisionWatcherClient creates an instance of ProvisionWatcherClient
func NewProvisionWatcherClient(baseUrl string) interfaces.ProvisionWatcherClient {
	return &ProvisionWatcherClient{
		baseUrl: baseUrl,
	}
}

func (pwc *ProvisionWatcherClient) Add(ctx context.Context, reqs []requests.AddProvisionWatcherRequest) (res []common.BaseWithIdResponse, err errors.EdgeX) {
	err = utils.PostRequest(ctx, &res, pwc.baseUrl+contracts.ApiProvisionWatcherRoute, reqs)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}

	return
}

func (pwc *ProvisionWatcherClient) Update(ctx context.Context, reqs []requests.UpdateProvisionWatcherRequest) (res []common.BaseResponse, err errors.EdgeX) {
	err = utils.PatchRequest(ctx, &res, pwc.baseUrl+contracts.ApiProvisionWatcherRoute, reqs)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}

	return
}

func (pwc *ProvisionWatcherClient) AllProvisionWatchers(ctx context.Context, labels []string, offset int, limit int) (res responses.MultiProvisionWatchersResponse, err errors.EdgeX) {
	requestParams := url.Values{}
	if len(labels) > 0 {
		requestParams.Set(contracts.Labels, strings.Join(labels, contracts.CommaSeparator))
	}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	err = utils.GetRequest(ctx, &res, pwc.baseUrl, contracts.ApiAllProvisionWatcherRoute, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}

	return
}

func (pwc *ProvisionWatcherClient) ProvisionWatcherByName(ctx context.Context, name string) (res responses.ProvisionWatcherResponse, err errors.EdgeX) {
	path := path.Join(contracts.ApiProvisionWatcherRoute, contracts.Name, url.QueryEscape(name))
	err = utils.GetRequest(ctx, &res, pwc.baseUrl, path, nil)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}

	return
}

func (pwc *ProvisionWatcherClient) DeleteProvisionWatcherByName(ctx context.Context, name string) (res common.BaseResponse, err errors.EdgeX) {
	path := path.Join(contracts.ApiProvisionWatcherRoute, contracts.Name, url.QueryEscape(name))
	err = utils.DeleteRequest(ctx, &res, pwc.baseUrl, path)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}

	return
}

func (pwc *ProvisionWatcherClient) ProvisionWatchersByProfileName(ctx context.Context, name string, offset int, limit int) (res responses.MultiProvisionWatchersResponse, err errors.EdgeX) {
	requestPath := path.Join(contracts.ApiProvisionWatcherRoute, contracts.Profile, contracts.Name, url.QueryEscape(name))
	requestParams := url.Values{}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	err = utils.GetRequest(ctx, &res, pwc.baseUrl, requestPath, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}

	return
}

func (pwc *ProvisionWatcherClient) ProvisionWatchersByServiceName(ctx context.Context, name string, offset int, limit int) (res responses.MultiProvisionWatchersResponse, err errors.EdgeX) {
	requestPath := path.Join(contracts.ApiProvisionWatcherRoute, contracts.Service, contracts.Name, url.QueryEscape(name))
	requestParams := url.Values{}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	err = utils.GetRequest(ctx, &res, pwc.baseUrl, requestPath, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}

	return
}
