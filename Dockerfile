FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum /build/
RUN set -ex && \
    go mod download

COPY . /build/

RUN set -ex && \
    go build -o dist/kubersctl .


FROM alpine:latest

ENV KUBERSCTL_BIN=/usr/local/bin/kubersctl \
    USER_NAME=kubers

RUN mkdir /licenses
COPY LICENSE /licenses

RUN set -ex && \
    addgroup -S ${USER_NAME} && \
    adduser -S --no-create-home -G ${USER_NAME} --disabled-password ${USER_NAME}

# switch to non-root user
USER ${USER_NAME}

COPY --chown=${USER_NAME}:${USER_NAME} --from=builder /build/dist/kubersctl ${KUBERSCTL_BIN}

ENTRYPOINT ["/usr/local/bin/kubersctl"]
