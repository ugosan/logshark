FROM --platform=$BUILDPLATFORM golang:1.19 as builder
ENV TERM "xterm-256color"
ARG TARGETARCH

WORKDIR /build
COPY . /build
RUN GOOS=linux GOARCH=$TARGETARCH go build -buildmode exe -ldflags="-w -s" -o ./logshark ./cmd

FROM alpine:latest
WORKDIR /root
COPY --from=builder /build/logshark /usr/local/bin/logshark
ENTRYPOINT ["logshark"]



