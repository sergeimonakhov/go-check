FROM alpine

COPY bin/go-check /usr/local/bin/go-check
COPY ./docker-entrypoint.sh /

RUN apk --no-cache add --update \
      openssl \
      ca-certificates \
    && chmod +x /usr/local/bin/go-check \
         /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
