# syntax=docker/dockerfile:1

# Build
FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server

# Deploy
FROM scratch

COPY --from=build /bin/server /bin/server

EXPOSE 8080

CMD [ "/bin/server" ]
