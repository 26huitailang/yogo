FROM teamcity-agent-go:1.17 as builder

COPY . /code
WORKDIR /code
# 这里如果1不是静态编译，放到scratch镜像中runtime有问题
RUN CGO_ENABLED=0 go build .

FROM scratch as runner

COPY --from=builder /code/yogo /app/
WORKDIR /app

CMD ["./yogo"]
