FROM golang:1.24-alpine AS dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY . ./
RUN go mod download

CMD ["air", "-c", "api.air.toml"]
