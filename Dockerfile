FROM scratch

ARG TARGETOS
ARG TARGETARCH

COPY target/builds/gcp-kafka-auth-handler-$TARGETOS-$TARGETARCH /usr/local/bin/gcp-kafka-auth-handler

ENTRYPOINT [ "/usr/local/bin/gcp-kafka-auth-handler", "serve" ]
