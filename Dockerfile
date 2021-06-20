FROM alpine

WORKDIR /app

COPY bin/shrtnr .
COPY templates templates

CMD ["/app/shrtnr", ""]
