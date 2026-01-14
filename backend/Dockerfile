# syntax=docker/dockerfile:1

# Build
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server

# Deploy - use alpine for healthcheck tools
FROM alpine:3.19

# Install curl for healthcheck and ca-certificates for HTTPS
RUN apk add --no-cache curl ca-certificates tzdata

COPY --from=build /bin/server /bin/server

EXPOSE 8080

CMD [ "/bin/server" ]
