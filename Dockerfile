FROM golang:1.12-alpine3.10 AS builder
WORKDIR /go/src/github.com/iavael/test-project
COPY . .
RUN go build -v .

FROM alpine:3.16.2
EXPOSE 8080
VOLUME /var/lib/metrics/
COPY --from=builder /go/src/github.com/iavael/test-project/test-project /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/test-project", "-s=/var/lib/metrics/storage.txt"]
CMD ["0.0.0.0:8080"]
