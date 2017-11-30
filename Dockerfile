FROM alpine:latest
RUN apk add --no-cache ca-certificates && \
    update-ca-certificates
WORKDIR /prometheus2appoptics
COPY prometheus2appoptics /prometheus2appoptics/prometheus2appoptics
COPY run.sh /prometheus2appoptics
EXPOSE 4567
RUN ["chmod", "+x", "run.sh"]
CMD "./run.sh"

