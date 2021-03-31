package interfaces

import (
	"context"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/requests"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/responses"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
)

// DeviceServiceClient defines the interface for interactions with the Device Service endpoint on the EdgeX Foundry core-metadata service.
type DeviceServiceClient interface {
	// Add adds new device services.
	Add(ctx context.Context, reqs []requests.AddDeviceServiceRequest) ([]common.BaseWithIdResponse, errors.EdgeX)
	// Update updates device services.
	Update(ctx context.Context, reqs []requests.UpdateDeviceServiceRequest) ([]common.BaseResponse, errors.EdgeX)
	// AllDeviceServices returns all device services. Device services can also be filtered by labels.
	// The result can be limited in a certain range by specifying the offset and limit parameters.
	// offset: The number of items to skip before starting to collect the result set. Default is 0.
	// limit: The number of items to return. Specify -1 will return all remaining items after offset. The maximum will be the MaxResultCount as defined in the configuration of service. Default is 20.
	AllDeviceServices(ctx context.Context, labels []string, offset int, limit int) (responses.MultiDeviceServicesResponse, errors.EdgeX)
	// DeviceServiceByName returns a device service by name.
	DeviceServiceByName(ctx context.Context, name string) (responses.DeviceServiceResponse, errors.EdgeX)
	// DeviceServiceByID returns a device service by id.
	DeviceServiceByID(ctx context.Context, id string) (responses.DeviceServiceResponse, errors.EdgeX)
	// DeleteByName deletes a device service by name.
	DeleteByName(ctx context.Context, name string) (common.BaseResponse, errors.EdgeX)
	// DeleteByID deletes a device service by id.
	DeleteByID(ctx context.Context, id string) (common.BaseResponse, errors.EdgeX)
	DeviceServicesSearch(ctx context.Context, offset int, limit int, req requests.DeviceServiceSearchQueryRequest) (responses.MultiDeviceServicesResponse, errors.EdgeX)
}
