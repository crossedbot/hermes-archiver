ARG OS=debian:bullseye-slim
ARG GOLANG_VERSION=1.19-bullseye
ARG CGO=0
ARG GOOS=linux
ARG GOARCH=amd64

#-------------------------------------------------------------------------------
FROM golang:${GOLANG_VERSION} AS gobuilder

ARG CGO
ARG GOOS
ARG GOARCH

RUN go version
WORKDIR /go/src/
COPY . .
RUN cd cmd/indexer && \
    CGO_ENABLED='${CGO}' GOOS='${GOOS}' GOARCH='${GOARCH}' \
    make -f /go/src/Makefile build

#-------------------------------------------------------------------------------
FROM ${OS}

ENV HERMES_HOME /usr/local/hermes
ENV PATH ${HERMES_HOME}/bin:$PATH
RUN mkdir -vp ${HERMES_HOME}
WORKDIR ${HERMES_HOME}

COPY --from=gobuilder /go/bin/indexer ./bin/indexer
COPY ./scripts/indexer/run.bash ./bin/run-indexer

EXPOSE 8888
ENTRYPOINT [ "run-indexer", "-d", "redis:6379", "-i", "/dns4/ipfs/tcp/5001" ]
