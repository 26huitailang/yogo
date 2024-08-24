FROM debian:12-slim as base

FROM base as build-stage

ARG GO_VERSION=1.22.6

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install --no-install-recommends -y \
  sudo \
  git \
  bash-completion \
  vim \
  ssh \
  gcc \
  g++ \
  net-tools \
  make \
  htop \
  # cleaning up unused files
  && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false \
  && rm -rf /var/lib/apt/lists/*

# 下载并安装 Go
RUN export ARCH=`dpkg --print-architecture` && curl -fsSL -v https://golang.org/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz -o go.tar.gz && \
  tar -C /usr/local -xzvf go.tar.gz && \
  rm go.tar.gz

# 前端工具安装
RUN curl -sL https://deb.nodesource.com/setup_20.x | bash - && \
  apt-get install nodejs -y && \
  node --version && \
  npm install -g pnpm && \
  pnpm --version

# 在这里添加任何额外的构建或运行时所需的命令

# 清理缓存及无用文件
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

FROM build-stage as dev-stage

# 设置工作目录
WORKDIR /app

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
  git \
  bash-completion \
  vim \
  htop \
  # cleaning up unused files
  && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false \
  && rm -rf /var/lib/apt/lists/*

COPY ./compose/production/web/entrypoint /entrypoint
RUN sed -i 's/\r$//g' /entrypoint
RUN chmod +x /entrypoint && chown dev-user:dev-user /entrypoint

COPY ./compose/local/web/start /start
RUN sed -i 's/\r$//g' /start
RUN chmod +x /start && chown dev-user:dev-user /start

WORKDIR ${APP_HOME}
USER dev-user

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/air-verse/air@latest

ENTRYPOINT ["/entrypoint"]
