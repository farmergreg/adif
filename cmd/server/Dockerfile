# Stage: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN go mod init server
RUN go get github.com/hamradiolog-net/adif@latest

COPY . .
RUN go build -o server

# Stage: Build Production Container
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]
CMD ["-addr", "0.0.0.0:8080"]