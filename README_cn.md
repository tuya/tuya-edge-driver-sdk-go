# Tuya-edge-driver-sdk-go

## Overview

该sdk用于开发接入tuya 边缘网关的南向设备服务。

## Usage

开发人员需要在自己的服务中实现 [`ProtocolDriver`](./pkg/models/protocoldriver.go) 接口，并在`main`函数中启动设备服务。可以使用 [`startup`](./pkg/startup/bootstrap.go) 包来实现`main`函数。

请参考sdk中提供的南向设备服务开发模板[device-service-template](./device-srvice-template) ,建议使用该模版的格式来开发南向设备服务。

## 技术支持

Tuya IoT 开发者平台: https://developer.tuya.com/en

Tuya 开发者帮助中心: https://support.tuya.com/en/help

Tuya 工单系统: https://service.console.tuya.com/

## License

[Apache-2.0](LICENSE)

持续开发中...