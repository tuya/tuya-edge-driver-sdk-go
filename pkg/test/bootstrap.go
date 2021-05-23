package test

import (
	"bou.ke/monkey"
	"context"
	"github.com/gorilla/mux"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/clients/http"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/clients"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/pkg/service"

	"reflect"
)

const (
	deviceservice_id    = "./data/device/service.json"
	device_service_name = "./data/device/devices.json"
	deviceprofile_name  = "./data/deviceprofile/"
)

func Bootstrap(serviceName string, serviceVersion string, driver interface{}) {
	InitStub()
	ctx, cancel := context.WithCancel(context.Background())
	service.Main(serviceName, serviceVersion, driver, ctx, cancel, mux.NewRouter())
}

func InitStub() {
	monkey.Patch(clients.CheckDependencyServices, CheckDependencyServicesStub)

	monkey.PatchInstanceMethod(reflect.TypeOf(&http.DeviceServiceClient{}), "DeviceServiceByID", DeviceServiceByIDStub)
	monkey.PatchInstanceMethod(reflect.TypeOf(&http.DeviceServiceClient{}), "Update", UpdateStub)
	monkey.PatchInstanceMethod(reflect.TypeOf(&http.DeviceClient{}), "DevicesByServiceName", DevicesByServiceNameStub)
	monkey.PatchInstanceMethod(reflect.TypeOf(&http.DeviceProfileClient{}), "DeviceProfileByName", DeviceProfileByNameStub)
	monkey.PatchInstanceMethod(reflect.TypeOf(&http.ProvisionWatcherClient{}), "ProvisionWatchersByServiceName", ProvisionWatchersByServiceNameStub)

	monkey.Patch(common.SendEvent, SendEventStub)
}
