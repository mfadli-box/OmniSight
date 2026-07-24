# Stage 1: Mengompilasi Kode Binary Go
FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ict_log_rotate .

# Stage 2: Runtime Image Ringan
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/ict_log_rotate .

ENV PG_HOST=erp_pool
ENV PG_PORT=5432
ENV PG_USER=dbe
ENV PG_PASS=rahasia
ENV PG_DATA=erp
ENV IS_POOL=true
ENV RE_PATH=/var/log/ict_log_rotate
ENV RE_NORMAL=7
ENV RE_ATTACK=90
CMD ["./ict_log_rotate"]
