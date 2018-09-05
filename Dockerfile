FROM alpine

WORKDIR /app

COPY bin/shrtnr .
COPY templates templates

CMD ["./shrtnr"]
