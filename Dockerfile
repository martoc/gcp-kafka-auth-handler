FROM golang:1.25.5 AS builder

FROM scratch

ARG TARGETOS
ARG TARGETARCH

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/passwd
COPY target/builds/kafka-auth-handler-$TARGETOS-$TARGETARCH /usr/local/bin/kafka-auth-handler

USER nobody

ENTRYPOINT [ "/usr/local/bin/kafka-auth-handler", "serve" ]
