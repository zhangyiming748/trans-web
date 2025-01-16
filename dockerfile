FROM golang:1.23.4-alpine3.21 AS builder
# 将文件复制和工作目录设置合并到一个 WORKDIR 中，WORKDIR 会自动创建目录
WORKDIR /app
COPY . /app
# 合并多个 Go 环境设置和构建命令
RUN go env -w GO111MODULE=on && \
    # go env -w GOPROXY=https://goproxy.cn,direct && \
    go vet && go mod tidy && go mod vendor && go build -o gin main.go

FROM alpine:3.21

RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#http://mirrors4.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories
RUN apk add translate-shell

# 复制文件
COPY --from=builder /app/gin /usr/bin/gin

ENTRYPOINT ["/usr/bin/gin"]