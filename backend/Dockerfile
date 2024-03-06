FROM golang:1.20-alpine as builder
ENV CGO_ENABLED=0
WORKDIR /app
RUN go install github.com/githubnemo/CompileDaemon@latest
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -C . -o ./backend

FROM alpine
WORKDIR /app
COPY .default.env /app/.default.env
COPY --from=builder /app/backend /app/backend
COPY --from=builder /go/bin/CompileDaemon /bin/CompileDaemon
CMD ["CompileDaemon", "-command=./app/backend"]
