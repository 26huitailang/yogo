FROM debian:12-slim as base

FROM base as build-stage

# 设置环境变量，Go 官方仓库通常建议设置这些变量
ENV GOLANG_VERSION 1.21.7
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# 更新系统并安装必要的包以支持下载和解压 Go 安装程序
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

# 下载并安装 Go
RUN curl -fsSL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz -o go.tar.gz && \
    tar -C /usr/local -xzvf go.tar.gz && \
    rm go.tar.gz
RUN go version

# 设置工作目录
WORKDIR /app

# 前端工具安装
COPY scripts ./scripts
RUN chmod +x ./scripts/*.sh
RUN /usr/bin/bash ./scripts/install-front-dev.sh

# 在这里添加任何额外的构建或运行时所需的命令

# 清理缓存及无用文件
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

FROM build-stage as dev-stage

ARG BUILD_ENVIRONMENT=local
ARG APP_HOME=/app
ENV GOPROXY https://goproxy.io,direct
ENV BUILD_ENV ${BUILD_ENVIRONMENT}
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Create devcontainer user and add it to sudoers
RUN mkdir -p /etc/sudoers.d \
  && groupadd --gid 1000 dev-user \
  && useradd --uid 1000 --gid dev-user --shell /bin/bash --create-home dev-user \
  && echo dev-user ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/dev-user \
  && chmod 0440 /etc/sudoers.d/dev-user
RUN mkdir -p /go && chown -R dev-user:dev-user /go

# devcontainer dependencies and utils
RUN apt-get update && apt-get install --no-install-recommends -y \
  sudo git bash-completion vim ssh gcc g++ net-tools \
  # cleaning up unused files
  && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false \
  && rm -rf /var/lib/apt/lists/*
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
#RUN apk --no-cache add ca-certificates sqlite

# dev scripts
COPY scripts ./scripts
RUN chmod +x ./scripts/*.sh
RUN /usr/bin/bash ./scripts/install-front-dev.sh

COPY ./compose/production/web/entrypoint /entrypoint
RUN sed -i 's/\r$//g' /entrypoint
RUN chmod +x /entrypoint

COPY ./compose/local/web/start /start
RUN sed -i 's/\r$//g' /start
RUN chmod +x /start

WORKDIR ${APP_HOME}
USER dev-user

ENTRYPOINT ["/entrypoint"]
