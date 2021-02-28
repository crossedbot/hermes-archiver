ARG OS_NICKNAME=buster
ARG OS=debian:buster-slim
ARG ARCH=x64

FROM ${OS}

ENV HERMES_HOME /usr/local/hermes
ENV PATH ${HERMES_HOME}/bin:$PATH
RUN mkdir -vp ${HERMES_HOME}
WORKDIR ${HERMES_HOME}

COPY --from=hermes-builder /go/bin/indexer ./bin/indexer
COPY ./scripts/indexer/run.bash ./bin/run-indexer

EXPOSE 8888
ENTRYPOINT [ "run-indexer", "-d", "redis:6379", "-i", "/dns4/ipfs/tcp/5001" ]