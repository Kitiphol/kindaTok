FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o video-convertor ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/video-convertor .
COPY machinery.yaml .


RUN apk add --no-cache ffmpeg curl



EXPOSE 8081

CMD ["./video-convertor"]