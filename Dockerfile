FROM golang:1.22.0 AS builder
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# 设置 Go 模块代理，确保构建时能下载依赖
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /workspace
ENV GO111MODULE=auto

# 禁用CGO并启用静态编译
ENV CGO_ENABLED=0

COPY . .

# 整理模块依赖
RUN go mod tidy

# 完全静态编译 - 使用 CGO_DISABLED=0 并强制静态链接库
RUN CGO_ENABLED=0  go build -a -o algo-wizard  main.go

# 调试：列出生成的文件及其权限
RUN ls -l /workspace

FROM alpine:3.1
USER 0
WORKDIR /

# 从构建阶段复制生成的静态二进制文件
COPY --from=builder /workspace/algo-wizard .
COPY --from=builder /workspace/project/ project/
COPY --from=builder /workspace/types/ types/

# 确保二进制文件具有执行权限
RUN chmod +x /algo-wizard && \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 保持容器运行
#CMD ["tail", "-f", "/dev/null"]

ENTRYPOINT ["/algo-wizard"]


