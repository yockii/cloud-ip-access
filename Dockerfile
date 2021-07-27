FROM golang:alpine as builder
WORKDIR /usr/src/app
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add --no-cache upx tzdata

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w" -o server && upx --best server -o _upx_server && mv -f _upx_server server

FROM scratch as runner
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /usr/src/app/server /opt/app/
COPY --from=builder /usr/src/app/conf /opt/app/conf/
EXPOSE 80
WORKDIR /opt/app
CMD ["/opt/app/server"]
