# alpine 1.11.6
FROM alpine@sha256:39eda93d15866957feaee28f8fc5adb545276a64147445c64992ef69804dbf01 as builder

ARG VERSION

WORKDIR /build

RUN set -ex && \
    apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates

ENV USER=kubers
ENV UID=10001

RUN set -ex && \
    adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

RUN set -ex && \
    wget https://github.com/jacops/kubers/releases/download/v${VERSION#?}/kubers_${VERSION#?}_Linux_x86_64.tar.gz

RUN set -ex && \
    tar -xvf kubers_${VERSION#?}_Linux_x86_64.tar.gz

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

USER kubers:kubers

COPY --from=builder /build/kubers-agent /usr/local/bin/kubers-agent

ENTRYPOINT ["/usr/local/bin/kubers-agent"]
