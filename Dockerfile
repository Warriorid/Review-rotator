# Stage 1: Modules caching
FROM golang:1.24-alpine AS modules
RUN apk add --no-cache git
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Stage 2: Builder
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/Review-rotator ./cmd/main.go

# Stage 3: Final
FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /bin/Review-rotator /app/Review-rotator
COPY --from=builder /app/configs /app/configs

WORKDIR /app

EXPOSE 8080
CMD ["/app/Review-rotator"]