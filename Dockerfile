ARG ARCH=
FROM ${ARCH}golang:alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN apk add --update-cache build-base \
 && go build

FROM ${ARCH}alpine:latest AS container

COPY --from=builder /go/src/app/journalist /usr/bin/journalist

CMD ["journalist"]
