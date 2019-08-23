# Development image


FROM golang:1.12-alpine as builder

RUN apk add git shadow ca-certificates

# Prepare directory for source code and empty directory, which we copy
# to scratch image
RUN mkdir -p /usr/src/watchcluster/tmp

ADD ./go.mod ./go.sum /usr/src/watchcluster/
WORKDIR /usr/src/watchcluster
RUN go mod download

# Add source code
ADD . /usr/src/watchcluster


# Build binary without linking to glibc, so we can use scratch image
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o wc .

FROM alpine:3.7
# Copy executable
COPY --from=builder /usr/src/watchcluster/tmp /watchcluster
COPY --from=builder /usr/src/watchcluster/wc /watchcluster/wc

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Required for running as nobody
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Required for writing logs
COPY --from=builder --chown=nobody:nobody /usr/src/watchcluster/tmp /tmp
USER nobody
WORKDIR /watchcluster
ENTRYPOINT ["./wc"]







# FROM golang:alpine AS BUILD-ENV

# ARG GOOS_VAL 
# ARG GOARCH_VAL

# RUN mkdir -p /go/src/github.com/ipochi/watchcluster/vendor

# COPY vendor/ /go/src/github.com/ipochi/watchcluster/vendor
# COPY cmd/ /go/src/github.com/ipochi/watchcluster/cmd
# COPY pkg/ /go/src/github.com/ipochi/watchcluster/pkg

# RUN cd /go/src/github.com/ipochi/watchcluster/cmd/watchcluster && \
#   GOOS=${GOOS_VAL} GOARCH=${GOARCH_VAL} go build -o /go/bin/watchcluster

# ENV KUBE_LATEST_VERSION="v1.13.0"

# RUN apk add --no-cache ca-certificates bash git \
#   && wget -q https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/amd64/kubectl -O /usr/local/bin/kubectl \
#   && chmod +x /usr/local/bin/kubectl

# # Production image
# FROM alpine

# COPY --from=BUILD-ENV /go/bin/watchcluster /go/bin/watchcluster
# COPY --from=BUILD-ENV /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# COPY --from=BUILD-ENV /usr/local/bin/kubectl /usr/local/bin/kubectl

# ENTRYPOINT /go/bin/watchcluster
