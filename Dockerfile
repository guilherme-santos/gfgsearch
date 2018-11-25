# Builder

FROM golang:1.11-alpine as builder

RUN apk update \
    && apk upgrade \
    && apk add --no-cache git bash make \
    && go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/guilherme-santos/gfgsearch

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only
COPY . ./
RUN make build-static

# Final docker image

FROM alpine:3.7

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/guilherme-santos/gfgsearch/gfgsearch /bin/

CMD ["gfgsearch"]
