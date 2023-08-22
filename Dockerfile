FROM golang:1.18
ADD ./ /go/src/
WORKDIR /go/src
RUN go env -w GOPROXY=https://proxy.golang.com.cn,https://goproxy.cn,direct
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM alpine:latest
MAINTAINER nick
WORKDIR /app/
COPY --from=0 /go/src/app ./app
VOLUME ["/app/config/conf","/app/config/secret"]
EXPOSE 8080
ENTRYPOINT ["./app"]
