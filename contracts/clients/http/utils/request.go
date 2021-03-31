//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
)

// Helper method to make the get request and return the body
func GetRequest(ctx context.Context, returnValuePointer interface{}, baseUrl string, requestPath string, requestParams url.Values) errors.EdgeX {
	req, err := createRequest(ctx, http.MethodGet, baseUrl, requestPath, requestParams)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	res, err := sendRequest(ctx, req)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to parse the response body", err)
	}
	return nil
}

func GetRequest2(ctx context.Context, baseUrl, requestPath, requestParams string) ([]byte, errors.EdgeX) {
	var (
		err error
		req *http.Request
		res []byte
	)

	req, err = createRequest2(ctx, http.MethodGet, baseUrl, requestPath, requestParams)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	res, err = sendRequest(ctx, req)
	if err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

// Helper method to make the post JSON request and return the body
func PostRequest(
	ctx context.Context,
	returnValuePointer interface{},
	url string,
	data interface{}) errors.EdgeX {

	req, err := createRequestWithRawData(ctx, http.MethodPost, url, data)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	res, err := sendRequest(ctx, req)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to parse the response body", err)
	}
	return nil
}

// Helper method to make the put JSON request and return the body
func PutRequest(
	ctx context.Context,
	returnValuePointer interface{},
	url string,
	data interface{}) errors.EdgeX {

	req, err := createRequestWithRawData(ctx, http.MethodPut, url, data)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	res, err := sendRequest(ctx, req)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to parse the response body", err)
	}
	return nil
}

func PutRequest2(ctx context.Context, url string, body []byte) ([]byte, errors.EdgeX) {
	var (
		err error
		req *http.Request
		res []byte
	)
	if req, err = createRequestWithMarshalData(ctx, http.MethodPut, url, body); err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	if res, err = sendRequest(ctx, req); err != nil {
		return res, errors.NewCommonEdgeXWrapper(err)
	}
	return res, nil
}

// PatchRequest makes a PATCH request and unmarshals the response to the returnValuePointer
func PatchRequest(
	ctx context.Context,
	returnValuePointer interface{},
	url string,
	data interface{}) errors.EdgeX {

	req, err := createRequestWithRawData(ctx, http.MethodPatch, url, data)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	res, err := sendRequest(ctx, req)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to parse the response body", err)
	}
	return nil
}

// Helper method to make the post file request and return the body
func PostByFileRequest(
	ctx context.Context,
	returnValuePointer interface{},
	url string,
	filePath string) errors.EdgeX {

	req, err := createRequestFromFilePath(ctx, http.MethodPost, url, filePath)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	res, err := sendRequest(ctx, req)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to parse the response body", err)
	}
	return nil
}

// Helper method to make the put file request and return the body
func PutByFileRequest(
	ctx context.Context,
	returnValuePointer interface{},
	url string,
	filePath string) errors.EdgeX {

	req, err := createRequestFromFilePath(ctx, http.MethodPut, url, filePath)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	res, err := sendRequest(ctx, req)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to parse the response body", err)
	}
	return nil
}

// Helper method to make the delete request and return the body
func DeleteRequest(ctx context.Context, returnValuePointer interface{}, baseUrl string, requestPath string) errors.EdgeX {
	req, err := createRequest(ctx, http.MethodDelete, baseUrl, requestPath, nil)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	res, err := sendRequest(ctx, req)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to parse the response body", err)
	}
	return nil
}
