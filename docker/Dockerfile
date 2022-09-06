FROM golang:1.18-rc as builder

COPY . /yogo
WORKDIR /yogo
# 这里如果1不是静态编译，放到scratch镜像中runtime有问题
RUN go build .

FROM alpine:3.15 as runtime

WORKDIR /app
COPY --from=builder /code/yogo /app/
CMD ["echo", "done"]
