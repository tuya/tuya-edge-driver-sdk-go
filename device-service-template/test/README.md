
# 本地单机测试方法

可以在不启动边缘网关的情况下，通过读取配置文件的方式启动项目进行测试。你也可以通过此方法熟悉项目流程，但是此方法只支持 `x86` 或者 `x86-64` 架构平台。

1. 在项目文件下中创建test目录，test目录结构如下,[参考这里](./)。
```
└── test
    ├── README.md
    ├── data
    │   ├── device
    │   │   ├── devices.json
    │   │   └── service.json
    │   └── deviceprofile
    │       └── 产品1.json
    │     
    ├── main.go
    └── res
        ├── Simple-Driver.yaml
        └── configuration.toml
```

2. 根据在data目录下配置服务、设备、产品的`.json`文件，[参考这里](./data)。注意：**data目录下目录结构必须如此不可更改**。


3. 调用`test.Bootstrap`方法，单机启动程序，[参考这里](./main.go)。