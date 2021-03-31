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

type DeviceClient struct {
	baseUrl string
}

// NewDeviceClient creates an instance of DeviceClient
func NewDeviceClient(baseUrl string) interfaces.DeviceClient {
	return &DeviceClient{
		baseUrl: baseUrl,
	}
}

func (dc DeviceClient) Add(ctx context.Context, reqs []requests.AddDeviceRequest) (res []common.BaseWithIdResponse, err errors.EdgeX) {
	err = utils.PostRequest(ctx, &res, dc.baseUrl+contracts.ApiDeviceRoute, reqs)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dc *DeviceClient) Update(ctx context.Context, reqs []requests.UpdateDeviceRequest) (res []common.BaseResponse, err errors.EdgeX) {
	err = utils.PatchRequest(ctx, &res, dc.baseUrl+contracts.ApiDeviceRoute, reqs)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dc *DeviceClient) AllDevices(ctx context.Context, labels []string, offset int, limit int) (res responses.MultiDevicesResponse, err errors.EdgeX) {
	requestParams := url.Values{}
	if len(labels) > 0 {
		requestParams.Set(contracts.Labels, strings.Join(labels, contracts.CommaSeparator))
	}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	err = utils.GetRequest(ctx, &res, dc.baseUrl, contracts.ApiAllDeviceRoute, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dc *DeviceClient) DeviceNameExists(ctx context.Context, name string) (res common.BaseResponse, err errors.EdgeX) {
	path := path.Join(contracts.ApiDeviceRoute, contracts.Check, contracts.Name, url.QueryEscape(name))
	err = utils.GetRequest(ctx, &res, dc.baseUrl, path, nil)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dc *DeviceClient) DeviceByName(ctx context.Context, name string) (res responses.DeviceResponse, err errors.EdgeX) {
	path := path.Join(contracts.ApiDeviceRoute, contracts.Name, url.QueryEscape(name))
	err = utils.GetRequest(ctx, &res, dc.baseUrl, path, nil)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dc *DeviceClient) DeleteDeviceByName(ctx context.Context, name string) (res common.BaseResponse, err errors.EdgeX) {
	path := path.Join(contracts.ApiDeviceRoute, contracts.Name, url.QueryEscape(name))
	err = utils.DeleteRequest(ctx, &res, dc.baseUrl, path)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dc *DeviceClient) DevicesByProfileName(ctx context.Context, name string, offset int, limit int) (res responses.MultiDevicesResponse, err errors.EdgeX) {
	requestPath := path.Join(contracts.ApiDeviceRoute, contracts.Profile, contracts.Name, url.QueryEscape(name))
	requestParams := url.Values{}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	err = utils.GetRequest(ctx, &res, dc.baseUrl, requestPath, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dc *DeviceClient) DevicesByServiceName(ctx context.Context, name string, offset int, limit int) (res responses.MultiDevicesResponse, err errors.EdgeX) {
	requestPath := path.Join(contracts.ApiDeviceRoute, contracts.Service, contracts.Name, url.QueryEscape(name))
	requestParams := url.Values{}
	requestParams.Set(contracts.Offset, strconv.Itoa(offset))
	requestParams.Set(contracts.Limit, strconv.Itoa(limit))
	err = utils.GetRequest(ctx, &res, dc.baseUrl, requestPath, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dc *DeviceClient) DeviceById(ctx context.Context, id string) (res responses.DeviceResponse, err errors.EdgeX) {
	path := path.Join(contracts.ApiDeviceRoute, contracts.Id, url.QueryEscape(id))
	err = utils.GetRequest(ctx, &res, dc.baseUrl, path, nil)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

func (dc *DeviceClient) DeleteDeviceById(ctx context.Context, id string) (res common.BaseResponse, err errors.EdgeX) {
	path := path.Join(contracts.ApiDeviceRoute, contracts.Id, url.QueryEscape(id))
	err = utils.DeleteRequest(ctx, &res, dc.baseUrl, path)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}
