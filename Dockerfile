FROM golang:1.13.7 AS builder
COPY . /go/src/github.com/solarwinds/prometheus2appoptics
WORKDIR /go/src/github.com/solarwinds/prometheus2appoptics
RUN make build

FROM alpine:latest
RUN apk add --no-cache ca-certificates && \
    update-ca-certificates
COPY --from=builder /go/src/github.com/solarwinds/prometheus2appoptics/prometheus2appoptics .
RUN chmod +x prometheus2appoptics
ENV SEND_STATS false
ENV APPOPTICS_TOKEN abc123fake
CMD ["./prometheus2appoptics"]

