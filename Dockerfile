FROM  golang:1.19-alpine3.18 as builder

RUN apk add --no-cache make ca-certificates gcc musl-dev linux-headers git jq bash

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

RUN go mod download

# build Nekoswap_runes with the shared go.mod & go.sum files
COPY . /app/DRNG

WORKDIR /app/DRNG

RUN make build

FROM alpine:3.18

COPY --from=builder /app/DRNG/bin/DRNG-service /usr/local/bin
COPY --from=builder /app/DRNG/tests /app/DRNG/tests

WORKDIR /app/DRNG
