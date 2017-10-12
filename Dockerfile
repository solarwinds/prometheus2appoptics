FROM alpine:latest
RUN apk add --no-cache ca-certificates && \
    update-ca-certificates
WORKDIR /p2l
COPY p2l /p2l/p2l
COPY run.sh /p2l/run.sh
EXPOSE 4567
RUN ["chmod", "+x", "run.sh"]
CMD "./run.sh"

