//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package contracts

// Miscellaneous constants
const (
	ClientMonitorDefault = 15000              // Defaults the interval at which a given service client will refresh its endpoint from the Registry, if used
	CorrelationHeader    = "X-Correlation-ID" // Sets the key of the Correlation ID HTTP header
)

// Constants related to the possible content types supported by the APIs
const (
	ContentType     = "Content-Type"
	ContentTypeCBOR = "application/cbor"
	ContentTypeJSON = "application/json"
	ContentTypeYAML = "application/x-yaml"
	ContentTypeText = "text/plain"
	ContentTypeXML  = "application/xml"
)

// Constants related to defined routes in the v2 service APIs
const (
	ApiVersion = "v2"
	ApiBase    = "/api/v2"

	ApiEventRoute                      = ApiBase + "/event"
	ApiEventProfileNameDeviceNameRoute = ApiEventRoute + "/{" + ProfileName + "}" + "/{" + DeviceName + "}"
	ApiAllEventRoute                   = ApiEventRoute + "/" + All
	ApiEventIdRoute                    = ApiEventRoute + "/" + Id + "/{" + Id + "}"
	ApiEventCountRoute                 = ApiEventRoute + "/" + Count
	ApiEventCountByDeviceNameRoute     = ApiEventCountRoute + "/" + Device + "/" + Name + "/{" + Name + "}"
	ApiEventByDeviceNameRoute          = ApiEventRoute + "/" + Device + "/" + Name + "/{" + Name + "}"
	ApiEventByTimeRangeRoute           = ApiEventRoute + "/" + Start + "/{" + Start + "}/" + End + "/{" + End + "}"
	ApiEventByAgeRoute                 = ApiEventRoute + "/" + Age + "/{" + Age + "}"

	ApiReadingRoute                            = ApiBase + "/reading"
	ApiAllReadingRoute                         = ApiReadingRoute + "/" + All
	ApiReadingCountRoute                       = ApiReadingRoute + "/" + Count
	ApiReadingCountByDeviceNameRoute           = ApiReadingCountRoute + "/" + Device + "/" + Name + "/{" + Name + "}"
	ApiReadingByDeviceNameRoute                = ApiReadingRoute + "/" + Device + "/" + Name + "/{" + Name + "}"
	ApiReadingByResourceNameRoute              = ApiReadingRoute + "/" + ResourceName + "/{" + ResourceName + "}"
	ApiReadingByTimeRangeRoute                 = ApiReadingRoute + "/" + Start + "/{" + Start + "}/" + End + "/{" + End + "}"
	ApiReadingByDeviceNameAndResourceNameRoute = ApiReadingRoute + "/" + Device + "/" + Name + "/{" + DeviceName + "}/" + ResourceName + "/{" + ResourceName + "}"

	ApiDeviceProfileRoute                       = ApiBase + "/deviceprofile"
	ApiDeviceProfileUploadFileRoute             = ApiDeviceProfileRoute + "/uploadfile"
	ApiDeviceProfileByNameRoute                 = ApiDeviceProfileRoute + "/" + Name + "/{" + Name + "}"
	ApiDeviceProfileByIdRoute                   = ApiDeviceProfileRoute + "/" + Id + "/{" + Id + "}"
	ApiAllDeviceProfileRoute                    = ApiDeviceProfileRoute + "/" + All
	ApiDeviceProfileByManufacturerRoute         = ApiDeviceProfileRoute + "/" + Manufacturer + "/{" + Manufacturer + "}"
	ApiDeviceProfileByModelRoute                = ApiDeviceProfileRoute + "/" + Model + "/{" + Model + "}"
	ApiDeviceProfileByManufacturerAndModelRoute = ApiDeviceProfileRoute + "/" + Manufacturer + "/{" + Manufacturer + "}" + "/" + Model + "/{" + Model + "}"
	ApiDeviceProfileCountRoute                  = ApiDeviceProfileRoute + "/" + Count
	ApiDeviceProfileSearchRoute                 = ApiDeviceProfileRoute + "/" + Search

	ApiDeviceServiceRoute          = ApiBase + "/deviceservice"
	ApiAllDeviceServiceRoute       = ApiDeviceServiceRoute + "/" + All
	ApiDeviceServiceByNameRoute    = ApiDeviceServiceRoute + "/" + Name + "/{" + Name + "}"
	ApiDeviceServiceByIdRoute      = ApiDeviceServiceRoute + "/" + Id + "/{" + Id + "}"
	ApiDeviceServiceRunStatusRoute = ApiDeviceServiceRoute + "/" + "run_status"
	ApiDeviceServiceSearchRoute    = ApiDeviceServiceRoute + "/" + Search

	ApiDeviceRoute                = ApiBase + "/device"
	ApiAllDeviceRoute             = ApiDeviceRoute + "/" + All
	ApiDeviceIdExistsRoute        = ApiDeviceRoute + "/" + Check + "/" + Id + "/{" + Id + "}"
	ApiDeviceNameExistsRoute      = ApiDeviceRoute + "/" + Check + "/" + Name + "/{" + Name + "}"
	ApiDeviceByIdRoute            = ApiDeviceRoute + "/" + Id + "/{" + Id + "}"
	ApiDeviceByNameRoute          = ApiDeviceRoute + "/" + Name + "/{" + Name + "}"
	ApiDeviceByProfileIdRoute     = ApiDeviceRoute + "/" + Profile + "/" + Id + "/{" + Id + "}"
	ApiDeviceByProfileNameRoute   = ApiDeviceRoute + "/" + Profile + "/" + Name + "/{" + Name + "}"
	ApiDeviceByServiceIdRoute     = ApiDeviceRoute + "/" + Service + "/" + Id + "/{" + Id + "}"
	ApiDeviceByServiceNameRoute   = ApiDeviceRoute + "/" + Service + "/" + Name + "/{" + Name + "}"
	ApiDeviceNameCommandNameRoute = ApiDeviceByNameRoute + "/{" + Command + "}"
	ApiDeviceCountRoute           = ApiDeviceRoute + "/" + Count
	ApiDeviceActiveRoute          = ApiDeviceRoute + "/active"
	ApiDeviceSearchRoute          = ApiDeviceRoute + "/" + Search

	ApiAuthCodeRoute     = ApiBase + "/authcode"
	ApiAuthCodeInfoRoute = ApiBase + ApiAuthCodeRoute + "/info"

	ApiProvisionWatcherRoute              = ApiBase + "/provisionwatcher"
	ApiAllProvisionWatcherRoute           = ApiProvisionWatcherRoute + "/" + All
	ApiProvisionWatcherByIdRoute          = ApiProvisionWatcherRoute + "/" + Id + "/{" + Id + "}"
	ApiProvisionWatcherByNameRoute        = ApiProvisionWatcherRoute + "/" + Name + "/{" + Name + "}"
	ApiProvisionWatcherByProfileNameRoute = ApiProvisionWatcherRoute + "/" + Profile + "/" + Name + "/{" + Name + "}"
	ApiProvisionWatcherByServiceNameRoute = ApiProvisionWatcherRoute + "/" + Service + "/" + Name + "/{" + Name + "}"

	ApiConfigRoute  = ApiBase + "/config"
	ApiMetricsRoute = ApiBase + "/metrics"
	ApiPingRoute    = ApiBase + "/ping"
	ApiVersionRoute = ApiBase + "/version"

	ApiDeviceCallbackRoute      = ApiBase + "/callback/device"
	ApiDeviceCallbackNameRoute  = ApiBase + "/callback/device/name/{name}"
	ApiProfileCallbackRoute     = ApiBase + "/callback/profile"
	ApiProfileCallbackNameRoute = ApiBase + "/callback/profile/name/{name}"
	ApiWatcherCallbackRoute     = ApiBase + "/callback/watcher"
	ApiWatcherCallbackNameRoute = ApiBase + "/callback/watcher/name/{name}"
	ApiServiceCallbackRoute     = ApiBase + "/callback/service"
	ApiDiscoveryRoute           = ApiBase + "/discovery"

	//功能点
	ApiFuncPointRoute       = ApiDeviceProfileRoute + "/{" + DeviceProfileId + "}" + "/func_point"
	ApiFuncPointByIdRoute   = ApiFuncPointRoute + "/" + Id + "/{" + Id + "}"
	ApiAllFuncPointRoute    = ApiFuncPointRoute + "/" + All
	ApiFuncPointCountRoute  = ApiFuncPointRoute + "/" + Count
	ApiFuncPointSearchRoute = ApiFuncPointRoute + "/" + Search

	ApiIntervalRoute       = ApiBase + "/interval"
	ApiAllIntervalRoute    = ApiIntervalRoute + "/" + All
	ApiIntervalByNameRoute = ApiIntervalRoute + "/" + Name + "/{" + Name + "}"

	ApiDeviceLibraryRoute       = ApiBase + "/devicelibrary"
	ApiAllDeviceLibraryRoute    = ApiDeviceLibraryRoute + "/" + All
	ApiDeviceLibraryByIdRoute   = ApiDeviceLibraryRoute + "/" + Id + "/{" + Id + "}"
	ApiDeviceLibrarySearchRoute = ApiDeviceLibraryRoute + "/" + Search

	ApiFileRoute = ApiBase + "/file"

	// 北向服务
	ApiEdgeHubDeviceActiveRoute = "/api/v1/app/device/tuya"
	ApiEdgeHubConfigRoute       = "/api/v1/app/config/tuya"
)

// Constants related to defined url path names and parameters in the v2 service APIs
const (
	All             = "all"
	Id              = "id"
	Created         = "created"
	Modified        = "modified"
	Pushed          = "pushed"
	Count           = "count"
	Device          = "device"
	DeviceId        = "deviceId"
	DeviceName      = "deviceName"
	Check           = "check"
	Profile         = "profile"
	Service         = "service"
	Command         = "command"
	ProfileName     = "profileName"
	ServiceName     = "serviceName"
	ResourceName    = "resourceName"
	Start           = "start"
	End             = "end"
	Age             = "age"
	Scrub           = "scrub"
	Type            = "type"
	Name            = "name"
	Label           = "label"
	Manufacturer    = "manufacturer"
	Model           = "model"
	ValueType       = "valueType"
	Offset          = "offset"         //query string to specify the number of items to skip before starting to collect the result set.
	Limit           = "limit"          //query string to specify the numbers of items to return
	Labels          = "labels"         //query string to specify associated user-defined labels for querying a given object. More than one label may be specified via a comma-delimited list
	PushEvent       = "ds-pushevent"   //query string to specify if an event should be pushed to the EdgeX system
	ReturnEvent     = "ds-returnevent" //query string to specify if an event should be returned from device service
	Search          = "search"
	DeviceProfileId = "deviceProfileId"
	MarkCode        = "markCode" //标示符
)

// Constants related to the default value of query strings in the v2 service APIs
const (
	DefaultOffset  = 0
	DefaultLimit   = 20
	CommaSeparator = ","
	ValueYes       = "yes"
	ValueNo        = "no"
)

// Constants related to Reading ValueTypes
const (
	ValueTypeBool         = "Bool"
	ValueTypeString       = "String"
	ValueTypeUint8        = "Uint8"
	ValueTypeUint16       = "Uint16"
	ValueTypeUint32       = "Uint32"
	ValueTypeUint64       = "Uint64"
	ValueTypeInt8         = "Int8"
	ValueTypeInt16        = "Int16"
	ValueTypeInt32        = "Int32"
	ValueTypeInt64        = "Int64"
	ValueTypeFloat32      = "Float32"
	ValueTypeFloat64      = "Float64"
	ValueTypeBinary       = "Binary"
	ValueTypeBoolArray    = "BoolArray"
	ValueTypeStringArray  = "StringArray"
	ValueTypeUint8Array   = "Uint8Array"
	ValueTypeUint16Array  = "Uint16Array"
	ValueTypeUint32Array  = "Uint32Array"
	ValueTypeUint64Array  = "Uint64Array"
	ValueTypeInt8Array    = "Int8Array"
	ValueTypeInt16Array   = "Int16Array"
	ValueTypeInt32Array   = "Int32Array"
	ValueTypeInt64Array   = "Int64Array"
	ValueTypeFloat32Array = "Float32Array"
	ValueTypeFloat64Array = "Float64Array"
)

// service key
const (
	TedgeEdgeHubServiceKey      = "tedge-edge-hub"
	TedgeDeviceModuleServiceKey = "tedge-device-module"
	TedgeDriverModuleServiceKey = "tedge-driver-module"
	TedgeManageServerServiceKey = "tedge-manage-server"
)

// dirver server
const (
	// 驱动服务运行状态 启动、停止、启动中
	StartedRunStatus  = 1
	StopedRunStatus   = 2
	StartingRunStatus = 3
)

// device
const (
	// 设备激活状态 未激活、已激活、激活失败
	DeviceActiveStatusInactivated = "inactivated"
	DeviceActiveStatusActivated   = "activated"
	DeviceActiveStatusActiveFail  = "activeFail"
)

// Constants related to defined routes in the service APIs
const (
	ApiAddressableRoute        = "/api/v1/addressable"
	ApiCallbackRoute           = "/api/v1/callback"
	ApiCommandRoute            = "/api/v1/command"
	ApiHealthRoute             = "/api/v1/health"
	ApiLoggingRoute            = "/api/v1/logs"
	ApiNotificationRoute       = "/api/v1/notification"
	ApiNotifyRegistrationRoute = "/api/v1/notify/registrations"
	ApiOperationRoute          = "/api/v1/operation"
	ApiRegistrationRoute       = "/api/v1/registration"
	ApiRegistrationByNameRoute = ApiRegistrationRoute + "/name"
	ApiSubscriptionRoute       = "/api/v1/subscription"
	ApiTransmissionRoute       = "/api/v1/transmission"
	ApiValueDescriptorRoute    = "/api/v1/valuedescriptor"
	ApiIntervalActionRoute     = "/api/v1/intervalaction"
)

// Constants related to how services identify themselves in the Service Registry
const (
	ServiceKeyPrefix                    = "edgex-"
	ConfigSeedServiceKey                = "edgex-config-seed"
	CoreCommandServiceKey               = "edgex-core-command"
	CoreDataServiceKey                  = "edgex-core-data"
	CoreMetaDataServiceKey              = "edgex-core-metadata"
	SupportLoggingServiceKey            = "edgex-support-logging"
	SupportNotificationsServiceKey      = "edgex-support-notifications"
	SystemManagementAgentServiceKey     = "edgex-sys-mgmt-agent"
	SupportSchedulerServiceKey          = "edgex-support-scheduler"
	SecuritySecretStoreSetupServiceKey  = "edgex-security-secretstore-setup"
	SecuritySecretsSetupServiceKey      = "edgex-security-secrets-setup"
	SecurityProxySetupServiceKey        = "edgex-security-proxy-setup"
	SecurityFileTokenProviderServiceKey = "edgex-security-file-token-provider"
	SecurityBootstrapRedisKey           = "edgex-security-bootstrap-redis"
)
