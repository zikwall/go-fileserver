FROM golang:alpine as app-builder
RUN apk update
RUN apk add --no-cache bash
RUN mkdir -p /go/tmp/app/{tmp}
WORKDIR /go/tmp/app
COPY . .
RUN go get .
RUN go list -m all
RUN CGO_ENABLED=0 go test -v ./...
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -tags timetzdata -o main /go/tmp/app .

FROM scratch
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=app-builder /go/tmp/app/tmp /go/src/app/tmp
COPY --from=app-builder /go/tmp/app/main /go/src/app/
WORKDIR /go/src/app/
CMD ["/go/src/app/main"]