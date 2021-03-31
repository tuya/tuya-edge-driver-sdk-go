// This package provides a simple device-service-template of a device service.
package main

import (
	"device-service-template/internal/driver"

	device "github.com/tuya/tuya-edge-driver-sdk-go"
	"github.com/tuya/tuya-edge-driver-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-service-template"
)

func main() {
	sd := driver.SimpleDriver{}
	startup.Bootstrap(serviceName, device.Version, &sd)
}
