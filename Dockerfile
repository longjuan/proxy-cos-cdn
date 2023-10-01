FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go build -o proxy-cos-cdn

FROM debian:stable-slim
WORKDIR /app
COPY --from=builder /app/proxy-cos-cdn /app/proxy-cos-cdn
RUN apt-get update && apt-get install -y ca-certificates
EXPOSE 3321
CMD ["/app/proxy-cos-cdn"]