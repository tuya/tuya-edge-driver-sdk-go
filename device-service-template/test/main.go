// This package provides a simple device-service-template of a device service.
package main

import (
	"device-service-template/internal/driver"
	"github.com/tuya/tuya-edge-driver-sdk-go/pkg/test"

	device "github.com/tuya/tuya-edge-driver-sdk-go"
)

const (
	serviceName string = "device-service-template"
)

func main() {
	sd := driver.SimpleDriver{}
	test.Bootstrap(serviceName, device.Version, &sd)
}
