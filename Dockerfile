FROM golang:1.14.3 AS builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y upx

WORKDIR /build

COPY go.mod go.sum /build/

RUN go mod download
RUN go mod verify

COPY . /build/

ENV LD_FLAGS="-w"
ENV CGO_ENABLED=0

RUN go install -v -tags netgo -ldflags "${LD_FLAGS}" .
RUN upx -9 /go/bin/linky

FROM busybox
LABEL maintainer="Robert Jacob <xperimental@solidproject.de>"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/linky /bin/linky

USER nobody
ENTRYPOINT ["/bin/linky"]
