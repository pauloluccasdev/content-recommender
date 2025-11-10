FROM golang:1.23.4-bullseye AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/server/main.go

FROM gcr.io/distroless/static-debian12
WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
ENV GIN_MODE=release

CMD ["/app/main"]