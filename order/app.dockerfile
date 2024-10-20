FROM golang:1.22.5-alpine3.11 as build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/bensmile/microservice-grpc-graphql
COPY go.mod go.sum ./
COPY vendor vendor
COPY order order
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./order/cmd/order

FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD [ "./app" ]