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

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/http/utils"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/interfaces"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/responses"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
)

type eventClient struct {
	baseUrl string
}

// NewEventClient creates an instance of EventClient
func NewEventClient(baseUrl string) interfaces.EventClient {
	return &eventClient{
		baseUrl: baseUrl,
	}
}

func (ec *eventClient) Add(ctx context.Context, req requests.AddEventRequest) (
	common.BaseWithIdResponse, errors.EdgeX) {
	path := path.Join(contracts.ApiEventRoute, url.QueryEscape(req.Event.ProfileName), url.QueryEscape(req.Event.DeviceName))
	var br common.BaseWithIdResponse
	err := utils.PostRequest(ctx, &br, ec.baseUrl+path, req)
	if err != nil {
		return br, errors.NewCommonEdgeXWrapper(err)
	}
	return br, nil
}

func (ec *eventClient) AllEvents(ctx context.Context, offset, limit int) (responses.MultiEventsResponse, errors.EdgeX) {
	requestParams := url.Values{}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	res := responses.MultiEventsResponse{}
	err := utils.GetRequest(ctx, &res, ec.baseUrl, contracts.ApiAllEventRoute, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (ec *eventClient) EventsByDeviceName(ctx context.Context, name string, offset, limit int) (
	responses.MultiEventsResponse, errors.EdgeX) {
	requestPath := path.Join(contracts.ApiEventRoute, contracts.Device, contracts.Name, url.QueryEscape(name))
	requestParams := url.Values{}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	res := responses.MultiEventsResponse{}
	err := utils.GetRequest(ctx, &res, ec.baseUrl, requestPath, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (ec *eventClient) DeleteByDeviceName(ctx context.Context, name string) (common.BaseResponse, errors.EdgeX) {
	path := path.Join(contracts.ApiEventRoute, contracts.Device, contracts.Name, url.QueryEscape(name))
	res := common.BaseResponse{}
	err := utils.DeleteRequest(ctx, &res, ec.baseUrl, path)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (ec *eventClient) EventsByTimeRange(ctx context.Context, start, end, offset, limit int) (
	responses.MultiEventsResponse, errors.EdgeX) {
	requestPath := path.Join(contracts.ApiEventRoute, contracts.Start, strconv.Itoa(start), contracts.End, strconv.Itoa(end))
	requestParams := url.Values{}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	res := responses.MultiEventsResponse{}
	err := utils.GetRequest(ctx, &res, ec.baseUrl, requestPath, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (ec *eventClient) DeleteByAge(ctx context.Context, age int) (common.BaseResponse, errors.EdgeX) {
	path := path.Join(contracts.ApiEventRoute, contracts.Age, strconv.Itoa(age))
	res := common.BaseResponse{}
	err := utils.DeleteRequest(ctx, &res, ec.baseUrl, path)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}
