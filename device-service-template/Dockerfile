# 南向设备服务dockerfile模板，推荐使用多阶段构建镜像
# 开发时请根据情况修改
ARG BASE=golang:1.15-alpine3.12
FROM ${BASE} AS builder

ARG MAKE='make build'

WORKDIR /device-service-template

# 换源加速
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk add --update --no-cache make git

COPY /device-service-template .

RUN ${MAKE}

FROM alpine:3.12

# 服务运行时监听的端口，开发时请根据情况修改
ENV APP_PORT=49990
EXPOSE $APP_PORT

WORKDIR /
COPY --from=builder /device-service-template/cmd/device-service-template /bin/device-service-template
COPY --from=builder /device-service-template/cmd/res/configuration.toml /etc/device-service-template/res/configuration.toml

ENTRYPOINT ["/bin/device-service-template"]
CMD ["--confdir=/etc/device-service-template/res"]
