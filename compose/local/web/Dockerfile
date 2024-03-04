FROM golang:1.22-bookworm as golang

FROM golang as golang-build-stage

ARG BUILD_ENVIRONMENT=local
ARG APP_HOME=/app
ENV GOPROXY https://goproxy.io,direct

#RUN sed -i "s@http://deb.debian.org@http://mirrors.aliyun.com@g" /etc/apt/sources.list && rm -rf /var/lib/apt/lists/* && apt-get update
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apt-get update && apt-get install --no-install-recommends -y \
  gcc \
  g++

WORKDIR ${APP_HOME}

RUN go version
COPY go.mod .
COPY go.sum .
RUN go mod vendor

# alpine/scratch/busybox choose one
FROM golang as golang-run-stage

ARG BUILD_ENVIRONMENT=local
ARG APP_HOME=/app

ENV BUILD_ENV ${BUILD_ENVIRONMENT}
ENV GOPROXY https://goproxy.io,direct

WORKDIR ${APP_HOME}

# Create devcontainer user and add it to sudoers
RUN mkdir -p /etc/sudoers.d \
  && groupadd --gid 1000 dev-user \
  && useradd --uid 1000 --gid dev-user --shell /bin/bash --create-home dev-user \
  && echo dev-user ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/dev-user \
  && chmod 0440 /etc/sudoers.d/dev-user

# devcontainer dependencies and utils
RUN apt-get update && apt-get install --no-install-recommends -y \
  sudo git bash-completion vim ssh
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
#RUN apk --no-cache add ca-certificates sqlite

COPY ./compose/production/web/entrypoint /entrypoint
RUN sed -i 's/\r$//g' /entrypoint
RUN chmod +x /entrypoint

COPY ./compose/local/web/start /start
RUN sed -i 's/\r$//g' /start
RUN chmod +x /start

ENTRYPOINT ["/entrypoint"]

# Python 'run' stage
#FROM python as python-run-stage
#
#
#
#




## Install required system dependencies
#RUN apt-get update && apt-get install --no-install-recommends -y \
#  # psycopg2 dependencies
#  libpq-dev \
#  # Translations dependencies
#  gettext \
#  # cleaning up unused files
#  && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false \
#  && rm -rf /var/lib/apt/lists/*
