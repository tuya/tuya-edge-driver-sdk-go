## **[English](README.md) | [中文](README_cn.md)**

# Device-service-template

基于`tedge-driver-sdk-go`开发南向设备服务的模版。

## Protocol Driver

开发者需要实现`ProtocolDriver`接口，接口定义如下：

```go
type ProtocolDriver interface {
    // 南向设备服务初始化函数
    // asyncChan: 用于异步上报设备上报上来的数据的channel
    // deviceCh: 自动发现设备，目前不支持
    Initialize(lc logger.LoggingClient, asyncCh chan<- *AsyncValues, deviceCh chan<- []DiscoveredDevice) error

    // 控制设备的读指令会通过该函数到达业务代码层，开发者需要在该函数内实现设备控制流程。
    // deviceName: 要控制的设备的名字
    // protocols: 控制设备时所需要的一些参数，由对应的开发者定义，例如要开发一款modbus协议的南向设备服务，
    // 则开发者需要定义设备的通信地址，读写数据的功能码，读取的数据的长度等等参数。
    // reqs: 设备读控制指令的请求
    HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []CommandRequest) ([]*CommandValue, error)

    // 控制设备的写指令会通过该函数到达业务代码层，开发者需要在该函数内实现设备控制流程。
    // deviceName: 要控制的设备的名字
    // protocols: 控制设备时所需要的一些参数，由开发者定义，例如要开发一款modbus协议的南向设备服务，
    // 则开发者需要定义设备的通信地址，读写数据的功能码，读取的数据的长度等等参数。
    // reqs: 设备写控制指令的请求
    // params: 设备写控制指令时需要的一些参数
    HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []CommandRequest, params []*CommandValue) error

    // 服务停止时会调用该函数，业务层的退出清理工作可以在此函数内实现
  	// force: 目前默认为false
    Stop(force bool) error

    // sdk在接收到添加设备请求时会调用该接口通知业务层
    // deviceName: 新添加设备的名字
    // protocols: 控制设备时所需要的一些参数，由开发者定义
  	// adminState: 管理员角度对设备状态的定义，locked，unlocked。当设备状态被设置为locked时，设备将不可以通过tedge进行读写操作。
    AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error

    // sdk在接收到更新设备请求时会调用该接口通知业务层
    // deviceName: 要更新设备的名字
    // protocols: 控制设备时所需要的一些参数，由开发者定义
  	// adminState: 管理员角度对设备状态的定义，locked，unlocked。当设备状态被设置为locked时，设备将不可以通过tedge进行读写操作。
    UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error

    // sdk在接收到删除设备请求时会调用该接口通知业务层
    // deviceName: 要删除设备的名字
    // protocols: 控制设备时所需要的一些参数，由开发者定义
    RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error
}
```

