FROM golang:1.16-alpine
WORKDIR /app
COPY . .
RUN \
    go mod init server && \
    go get && \
    go build
CMD ["./server"]
