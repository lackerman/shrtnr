FROM alpine

WORKDIR /app

COPY shrtnr .
COPY templates templates

CMD ["/app/shrtnr", ""]
