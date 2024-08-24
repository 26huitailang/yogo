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
  make \
  vim \
  gcc \
  g++ \
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

# 清理缓存及无用文件
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /app
COPY . /app
RUN make yogo

FROM base as runtime-stage

# 设置工作目录
WORKDIR /app

ARG BUILD_ENVIRONMENT=production
ARG APP_HOME=/app
ENV BUILD_ENV ${BUILD_ENVIRONMENT}

COPY --from=build-stage /app/bin /app/bin
RUN chmod 0755 /app/bin/*
# Create devcontainer user and add it to sudoers
RUN mkdir -p /etc/sudoers.d \
  && groupadd --gid 1000 dev-user \
  && useradd --uid 1000 --gid dev-user --shell /bin/bash --create-home dev-user \
  && echo dev-user ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/dev-user \
  && chmod 0440 /etc/sudoers.d/dev-user

COPY ./compose/production/web/entrypoint /entrypoint
RUN sed -i 's/\r$//g' /entrypoint
RUN chmod +x /entrypoint && chown dev-user:dev-user /entrypoint

COPY ./compose/production/web/start /start
RUN sed -i 's/\r$//g' /start
RUN chmod +x /start && chown dev-user:dev-user /start

RUN chown dev-user:dev-user /app
WORKDIR ${APP_HOME}
USER dev-user

ENTRYPOINT ["/entrypoint"]
