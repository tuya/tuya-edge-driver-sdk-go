# Tuya-edge-driver-sdk-go

### **[English](README.md) | [中文版](README_cn.md)**

## Overview

The SDK is used to develop southbound device services connected to tuya edge gateway.

## Usage

Developers need to implement the [`ProtocolDriver`](./pkg/models/protocoldriver.go) interface in their services and start the device service in the `main` function. You can use the [`startup`](./pkg/startup/bootstrap.go) package to implement the `main` function.

Please refer to the southbound device service development template **device-service-template** provided in the SDK. It is recommended to use this template format to develop southbound device services.

## Technical Support

Tuya IoT Developer Platform: https://developer.tuya.com/en/

Tuya Developer Help Center: https://support.tuya.com/en/help

Tuya Work Order System: https://service.console.tuya.com/

## License

[Apache-2.0](LICENSE)

Continuous development...