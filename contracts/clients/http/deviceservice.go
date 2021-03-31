package http

import (
	"context"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/schema"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/http/utils"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/interfaces"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/responses"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
)

type DeviceServiceClient struct {
	baseUrl string
}

// NewDeviceServiceClient creates an instance of DeviceServiceClient
func NewDeviceServiceClient(baseUrl string) interfaces.DeviceServiceClient {
	return &DeviceServiceClient{
		baseUrl: baseUrl,
	}
}

func (dsc *DeviceServiceClient) Add(ctx context.Context, reqs []requests.AddDeviceServiceRequest) (
	res []common.BaseWithIdResponse, err errors.EdgeX) {
	err = utils.PostRequest(ctx, &res, dsc.baseUrl+contracts.ApiDeviceServiceRoute, reqs)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dsc *DeviceServiceClient) Update(ctx context.Context, reqs []requests.UpdateDeviceServiceRequest) (
	res []common.BaseResponse, err errors.EdgeX) {
	err = utils.PatchRequest(ctx, &res, dsc.baseUrl+contracts.ApiDeviceServiceRoute, reqs)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dsc *DeviceServiceClient) AllDeviceServices(ctx context.Context, labels []string, offset int, limit int) (
	res responses.MultiDeviceServicesResponse, err errors.EdgeX) {
	requestParams := url.Values{}
	if len(labels) > 0 {
		requestParams.Set(contracts.Labels, strings.Join(labels, contracts.CommaSeparator))
	}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	err = utils.GetRequest(ctx, &res, dsc.baseUrl, contracts.ApiAllDeviceServiceRoute, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dsc *DeviceServiceClient) DeviceServiceByName(ctx context.Context, name string) (
	res responses.DeviceServiceResponse, err errors.EdgeX) {
	path := path.Join(contracts.ApiDeviceServiceRoute, contracts.Name, url.QueryEscape(name))
	err = utils.GetRequest(ctx, &res, dsc.baseUrl, path, nil)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dsc *DeviceServiceClient) DeleteByName(ctx context.Context, name string) (
	res common.BaseResponse, err errors.EdgeX) {
	path := path.Join(contracts.ApiDeviceServiceRoute, contracts.Name, url.QueryEscape(name))
	err = utils.DeleteRequest(ctx, &res, dsc.baseUrl, path)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dsc *DeviceServiceClient) DeviceServiceByID(ctx context.Context, id string) (
	res responses.DeviceServiceResponse, err errors.EdgeX) {
	path := path.Join(contracts.ApiDeviceServiceRoute, contracts.Id, url.QueryEscape(id))
	err = utils.GetRequest(ctx, &res, dsc.baseUrl, path, nil)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dsc *DeviceServiceClient) DeleteByID(ctx context.Context, id string) (
	res common.BaseResponse, err errors.EdgeX) {
	path := path.Join(contracts.ApiDeviceServiceRoute, contracts.Id, url.QueryEscape(id))
	err = utils.DeleteRequest(ctx, &res, dsc.baseUrl, path)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dsc *DeviceServiceClient) DeviceServicesSearch(ctx context.Context, offset int, limit int, req requests.DeviceServiceSearchQueryRequest) (res responses.MultiDeviceServicesResponse, edgexError errors.EdgeX) {
	requestParams := url.Values{}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	var encoder = schema.NewEncoder()
	err := encoder.Encode(req, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}

	err = utils.GetRequest(ctx, &res, dsc.baseUrl, contracts.ApiDeviceServiceSearchRoute, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}
