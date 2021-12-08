FROM teamcity-agent-go:1.17.4 as builder

COPY . /code
WORKDIR /code
# 这里如果1不是静态编译，放到scratch镜像中runtime有问题
RUN CGO_ENABLED=1 go build .

FROM alpine:3.15 as runtime

ENV PATH /app:$PATH
RUN apk add --no-cache gcompat
COPY --from=builder /code/yogo /app/
WORKDIR /app

CMD ["yogo"]
