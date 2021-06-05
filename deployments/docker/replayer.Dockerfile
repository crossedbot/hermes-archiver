ARG OS_NICKNAME=buster
ARG OS=debian:buster-slim
ARG ARCH=x64

FROM ${OS}

ENV HERMES_HOME /usr/local/hermes
ENV PATH ${HERMES_HOME}/bin:$PATH
RUN mkdir -vp ${HERMES_HOME}
WORKDIR ${HERMES_HOME}

COPY --from=hermes-builder /go/bin/replayer ./bin/replayer
COPY ./scripts/replayer/run.bash ./bin/run-replayer

EXPOSE 8989
ENTRYPOINT [ "run-replayer", "-d", "redis:6379", "-i", "/dns4/ipfs/tcp/5001" ]