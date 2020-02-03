FROM alpine:latest
RUN apk add --no-cache ca-certificates && \
    update-ca-certificates
WORKDIR /prometheus2appoptics
COPY prometheus2appoptics /prometheus2appoptics/prometheus2appoptics
COPY run.sh /prometheus2appoptics
EXPOSE 4567
ENV SEND_STATS false
ENV ACCESS_TOKEN abc123fake

CMD ["./prometheus2appoptics", "--access-token=$ACCESS_TOKEN", "--SEND_STATS=$SEND_STATS"]

