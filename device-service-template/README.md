## **[English](README.md) | [中文](README_cn.md)**

# Device-service-template

`tedge-driver-sdk-go` based development template.

## Protocol Driver

Developers need to implement the `ProtocolDriver` interface，Interface definition:

```go
type ProtocolDriver interface {
	// Initialize performs protocol-specific initialization for the device service.
	// The given *AsyncValues channel can be used to push asynchronous events and
	// readings to Core Data. The given []DiscoveredDevice channel is used to send
	// discovered devices that will be filtered and added to Core Metadata asynchronously.
	Initialize(lc logger.LoggingClient, asyncCh chan<- *AsyncValues, deviceCh chan<- []DiscoveredDevice) error

	// HandleReadCommands passes a slice of CommandRequest struct each representing
	// a ResourceOperation for a specific device resource.
	HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []CommandRequest) ([]*CommandValue, error)

	// HandleWriteCommands passes a slice of CommandRequest struct each representing
	// a ResourceOperation for a specific device resource.
	// Since the commands are actuation commands, params provide parameters for the individual
	// command.
	HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []CommandRequest, params []*CommandValue) error

	// Stop instructs the protocol-specific DS code to shutdown gracefully, or
	// if the force parameter is 'true', immediately. The driver is responsible
	// for closing any in-use channels, including the channel used to send async
	// readings (if supported).
	Stop(force bool) error

	// AddDevice is a callback function that is invoked
	// when a new Device associated with this Device Service is added
	AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error

	// UpdateDevice is a callback function that is invoked
	// when a Device associated with this Device Service is updated
	UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error

	// RemoveDevice is a callback function that is invoked
	// when a Device associated with this Device Service is removed
	RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error
}

```