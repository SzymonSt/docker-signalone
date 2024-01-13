FROM golang:1.20-alpine as builder
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -C ./cmd -o ../backend

FROM alpine
WORKDIR /app
COPY .default.env /app/.default.env
COPY --from=builder /app/backend /app/backend
CMD ["/app/backend"]
