# Multi-stage build to generate custom k6 with extension
FROM golang:1.20 as builder
WORKDIR $GOPATH/src/go.k6.io/k6
ADD . .
# RUN apk --no-cache add build-base git
RUN go install go.k6.io/xk6/cmd/xk6@latest
RUN CGO_ENABLED=1 xk6 build \
    --with github.com/aseara/xk6-m3u8=. \
    --output /tmp/k6

# Create image for running your customized k6
FROM ubuntu:23.04
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && apt-get autoremove -y && apt-get autoclean -y \
    && useradd -m k6 && usermod -a -G k6 k6
COPY --from=builder /tmp/k6 /usr/bin/k6
COPY --from=redboxoss/scuttle:latest scuttle /bin/scuttle

USER k6
WORKDIR /home/k6

ENTRYPOINT ["scuttle","k6"]
