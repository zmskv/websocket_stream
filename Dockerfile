FROM golang:1.22 as builder

WORKDIR /app/websocket_stream


COPY websocket_stream/go.mod websocket_stream/go.sum ./
RUN go mod download

COPY ./websocket_stream .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

FROM debian:bullseye-slim

COPY --from=builder /app/websocket_stream/server /app/websocket_stream/server

WORKDIR /app/websocket_stream


EXPOSE 8080
EXPOSE 9092


CMD ["./server"]
