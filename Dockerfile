FROM alpine:3.9

ARG VERSION

ENV LOG_LEVEL ERROR
ENV LISTEN_ADDRESS 0.0.0.0:8500 

RUN apk add curl ca-certificates

WORKDIR /bin

# Download binary
RUN curl -L https://github.com/biensupernice/krane/releases/download/${VERSION}/krane_${VERSION}_linux_386.tar.gz | tar xz && chmod +x krane 

EXPOSE 8500
VOLUME "/var/run/docker.sock"

ENTRYPOINT krane
