FROM golang:1.20.5 AS builder
ARG APP_NAME=commonapi
ENV GO111MODULE=on

WORKDIR /build

# 复制构建应用程序所需的代码
# 可能需要更改下边的命令，只复制您实际需要的内容。
COPY . .

# 构建应用程序
RUN GOOS=linux GOARCH=amd64 go build -mod vendor -o ${APP_NAME} cmd/server/server.go

# 我们创建一个 /dist 目录， 仅包含运行时必须的文件
# 然后，他会被复制到输出镜像的 / （根目录）
WORKDIR /dist
RUN cp /build/${APP_NAME} ./${APP_NAME}

# 可选项:如果您的应用程序使用动态链接(通常情况下使用 CGO)，
# 这将收集相关库，以便稍后将它们复制到最终镜像
# 注意: 确保您遵守您复制和分发的库的许可条款
RUN ldd ${APP_NAME} | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'
# RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/

# 在运行时复制或创建您的应用程序需要的其他目录/文件。
# 例如，本例使用 /data 作为工作目录，在正常运行容器时，该目录可能绑定到永久目录
RUN mkdir /data

# 构建最小运行时镜像
FROM scratch

COPY --chown=0:0 --from=builder /dist /
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 设置应用程序以 /data 文件夹中的非 root 用户身份运
# User ID 65534 通常是 'nobody' 用户.
# 映像的执行者仍应在安装过程中指定一个用户。
COPY --chown=65534:0 --from=builder /data /data
USER 65534
WORKDIR /data