

# build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd

# run stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]