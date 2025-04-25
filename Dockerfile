FROM golang:1.24.2 AS builder

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main .

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/main ./

ENTRYPOINT ["./main"]