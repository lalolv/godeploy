# 打包依赖阶段使用golang作为基础镜像
FROM golang:1.17 as builder

# 启用go module
ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

# CGO_ENABLED禁用cgo 然后指定OS等，并go build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -v -o /app/godeploy


# 运行阶段指定scratch作为基础镜像
FROM alpine
LABEL maintainer="work@lalo.im" 
LABEL version="0.2"

# 添加 bash
RUN apk add --no-cache bash

WORKDIR /app

# 复制应用文件和配置文件
COPY --from=builder /app/godeploy .
COPY --from=builder /app/app.conf .
# 为了防止代码中请求https链接报错，我们需要将证书纳入到scratch中
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/cert

# 指定运行时环境变量
ENV GIN_MODE=release PORT=8080

EXPOSE 8080
ENTRYPOINT ["./godeploy"]
