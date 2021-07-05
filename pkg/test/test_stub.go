package test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/google/uuid"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/http"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/interfaces"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	v2common "github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/responses"
	eErr "github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/container"
	context2 "github.com/tuya/tuya-edge-driver-sdk-go/internal/context"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
)

func CheckDependencyServicesStub(ctx context.Context, startupTimer startup.Timer, dic *di.Container) bool {
	var dependencyList = []string{common.ClientData, common.ClientMetadata}
	var waitGroup sync.WaitGroup
	checkingErr := true

	dependencyCount := len(dependencyList)
	waitGroup.Add(dependencyCount)

	for i := 0; i < dependencyCount; i++ {
		go func(wg *sync.WaitGroup, serviceName string) {
			defer wg.Done()
			configuration := container.ConfigurationFrom(dic.Get)
			lc := bootstrapContainer.LoggingClientFrom(dic.Get)
			lc.Infof("Check %v service's status by ping...", serviceName)
			addr := configuration.Clients[serviceName].Url()
			lc.Infof("Check %v service's , addr: %s", serviceName, addr)
		}(&waitGroup, dependencyList[i])
	}
	waitGroup.Wait()

	return checkingErr
}

func DeviceServiceByIDStub(dsc *http.DeviceServiceClient, ctx context.Context, id string) (
	res responses.DeviceServiceResponse, err eErr.EdgeX) {
	bytes, err1 := ioutil.ReadFile(deviceservice_id)
	dsr := &responses.DeviceServiceResponse{}
	json.Unmarshal(bytes, dsr)

	res = *dsr
	if err1 != nil {
		return res, eErr.NewCommonEdgeXWrapper(err1)
	}
	return res, nil
}

func UpdateStub(dsc *http.DeviceServiceClient, ctx context.Context, reqs []requests.UpdateDeviceServiceRequest) (
	res []v2common.BaseResponse, err eErr.EdgeX) {
	//err = utils.PatchRequest(ctx, &res, dsc.baseUrl+v2.ApiDeviceServiceRoute, reqs)
	//if err != nil {
	//	return res, errors.NewCommonEdgeXWrapper(err)
	//}
	return res, nil
}

func DevicesByServiceNameStub(dc *http.DeviceClient, ctx context.Context, name string, offset int, limit int) (res responses.MultiDevicesResponse, err eErr.EdgeX) {
	//requestPath := path.Join(v2.ApiDeviceRoute, v2.Service, v2.Name, url.QueryEscape(name))
	//requestParams := url.Values{}
	//requestParams.Set(v2.Offset, strconv.Itoa(offset))
	//requestParams.Set(v2.Limit, strconv.Itoa(limit))

	bytes, err1 := ioutil.ReadFile(device_service_name)

	err1 = json.Unmarshal(bytes, &res)
	//err = utils.GetRequest(ctx, &res, dc.baseUrl, requestPath, requestParams)
	if err1 != nil {
		return res, eErr.NewCommonEdgeXWrapper(err1)
	}
	return res, nil
}

func DeviceProfileByNameStub(client *http.DeviceProfileClient, ctx context.Context, name string) (res responses.DeviceProfileResponse, edgexError eErr.EdgeX) {

	bytes, err1 := ioutil.ReadFile(deviceprofile_name + name + ".json")

	err1 = json.Unmarshal(bytes, &res)
	//err = utils.GetRequest(ctx, &res, dc.baseUrl, requestPath, requestParams)
	if err1 != nil {
		return res, eErr.NewCommonEdgeXWrapper(err1)
	}
	return res, nil

	//requestPath := path.Join(v2.ApiDeviceProfileRoute, v2.Name, url.QueryEscape(name))
	//err := utils.GetRequest(ctx, &res, client.baseUrl, requestPath, nil)
	//if err != nil {
	//	return res, errors.NewCommonEdgeXWrapper(err)
	//}
	//return res, nil
}

//Get "http://localhost:48081/api/v2/provisionwatcher/service/name/
func ProvisionWatchersByServiceNameStub(pwc *http.ProvisionWatcherClient, ctx context.Context, name string, offset int, limit int) (res responses.MultiProvisionWatchersResponse, err eErr.EdgeX) {
	//requestPath := path.Join(v2.ApiProvisionWatcherRoute, v2.Service, v2.Name, url.QueryEscape(name))
	//requestParams := url.Values{}
	//requestParams.Set(v2.Offset, strconv.Itoa(offset))
	//requestParams.Set(v2.Limit, strconv.Itoa(limit))
	//err = utils.GetRequest(ctx, &res, pwc.baseUrl, requestPath, requestParams)
	//if err != nil {
	//	return res, eErr.NewCommonEdgeXWrapper(err)
	//}

	return
}

func SendEventStub(event dtos.Event, lc logger.LoggingClient, ec interfaces.EventClient) {
	correlation := uuid.New().String()
	ctx := context.WithValue(context.Background(), common.CorrelationHeader, correlation)
	/*
		if event.HasBinaryValue() {
			ctx = context.WithValue(ctx, clients.ContentType, clients.ContentTypeCBOR)
		} else {
			ctx = context.WithValue(ctx, clients.ContentType, clients.ContentTypeJSON)
		}
		// Call MarshalEvent to encode as byte array whether event contains binary or JSON readings
		var err error
		if len(event.EncodedEvent) <= 0 {
			event.EncodedEvent, err = ec.MarshalEvent(event.Event)
			if err != nil {
				lc.Error("SendEvent: Error encoding event", "device", event.Device, clients.CorrelationHeader, correlation, "error", err)
			} else {
				lc.Debug("SendEvent: EventClient.MarshalEvent encoded event", clients.CorrelationHeader, correlation)
			}
		} else {
			lc.Debug("SendEvent: EventClient.MarshalEvent passed through encoded event", clients.CorrelationHeader, correlation)
		}
	*/
	//req := requests.AddEventRequest{
	//	BaseRequest: v2common.NewBaseRequest(),
	//	Event:       event,
	//}
	// Call Add to post event to core data
	//responseBody, errPost := ec.Add(ctx, req)
	//if errPost != nil {
	//	lc.Error("SendEvent Failed to push event", "device", event.DeviceName, "response", responseBody, "error", errPost)
	//} else {
	lc.Debug("SendEvent: Pushed event to core data", contracts.ContentType, context2.FromContext(ctx, contracts.ContentType), contracts.CorrelationHeader, correlation)
	lc.Trace("SendEvent: Pushed this event to core data", contracts.ContentType, context2.FromContext(ctx, contracts.ContentType), contracts.CorrelationHeader, correlation, "event", event)
	//}
}
