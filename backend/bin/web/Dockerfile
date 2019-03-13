FROM alpine:latest

RUN apk add --no-cache ca-certificates && \
    adduser -h /home/ashirt -S -D ashirt

USER ashirt
WORKDIR /home/ashirt

COPY web /home/ashirt/public-api
COPY migrations /migrations

CMD /home/ashirt/public-api
