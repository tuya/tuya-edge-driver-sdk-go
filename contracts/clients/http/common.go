//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/http/utils"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/interfaces"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
)

type commonClient struct {
	baseUrl string
}

// NewCommonClient creates an instance of CommonClient
func NewCommonClient(baseUrl string) interfaces.CommonClient {
	return &commonClient{
		baseUrl: baseUrl,
	}
}

func (cc *commonClient) Configuration(ctx context.Context) (common.ConfigResponse, errors.EdgeX) {
	cr := common.ConfigResponse{}
	err := utils.GetRequest(ctx, &cr, cc.baseUrl, contracts.ApiConfigRoute, nil)
	if err != nil {
		return cr, errors.NewCommonEdgeXWrapper(err)
	}
	return cr, nil
}

func (cc *commonClient) Metrics(ctx context.Context) (common.MetricsResponse, errors.EdgeX) {
	mr := common.MetricsResponse{}
	err := utils.GetRequest(ctx, &mr, cc.baseUrl, contracts.ApiMetricsRoute, nil)
	if err != nil {
		return mr, errors.NewCommonEdgeXWrapper(err)
	}
	return mr, nil
}

func (cc *commonClient) Ping(ctx context.Context) (common.PingResponse, errors.EdgeX) {
	pr := common.PingResponse{}
	err := utils.GetRequest(ctx, &pr, cc.baseUrl, contracts.ApiPingRoute, nil)
	if err != nil {
		return pr, errors.NewCommonEdgeXWrapper(err)
	}
	return pr, nil
}

func (cc *commonClient) Version(ctx context.Context) (common.VersionResponse, errors.EdgeX) {
	vr := common.VersionResponse{}
	err := utils.GetRequest(ctx, &vr, cc.baseUrl, contracts.ApiVersionRoute, nil)
	if err != nil {
		return vr, errors.NewCommonEdgeXWrapper(err)
	}
	return vr, nil
}
