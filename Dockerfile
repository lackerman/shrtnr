FROM golang

WORKDIR /go/src/github.com/lackerman/shrtnr

COPY . .

RUN go get -d -v golang.org/x/net/html \
	&& go get ./... \
	&& go generate ./... \
	&& CGO_ENABLED=0 GOOS=linux go build -o app cmd/main.go

FROM alpine

COPY --from=0 /go/src/github.com/lackerman/shrtnr/app .
COPY --from=0 /go/src/github.com/lackerman/shrtnr/public public
COPY --from=0 /go/src/github.com/lackerman/shrtnr/templates templates

CMD ["./app"]
