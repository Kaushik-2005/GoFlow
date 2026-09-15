FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=docker -X main.commit=unknown" \
    -o /app/goflow ./cmd/goflow

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/goflow /app/goflow
COPY --from=builder /app/migrations /app/migrations

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/goflow"]
