ARG GOLANG_VERSION=1.14-buster
ARG CGO=0
ARG OS=linux
ARG ARCH=amd64

FROM golang:${GOLANG_VERSION} AS gobuilder
RUN go version
WORKDIR /go/src/
COPY . .
RUN CGO_ENABLED=${CGO} GOOS=${OS} GOARCH=${ARCH} go build \
    -a -installsuffix cgo \
    -o /go/bin/indexer ./cmd/indexer
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -o /go/bin/replayer ./cmd/replayer