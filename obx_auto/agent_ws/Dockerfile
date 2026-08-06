# Stage 1: Mengompilasi Kode Binary Go
FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o agent_ws .

# Stage 2: Runtime Image Ringan
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/agent_ws .

ENV PG_HOST=localhost
ENV PG_PORT=5432
ENV PG_USER=postgres
ENV PG_PASS=rahasia
ENV PG_DATA=ict
ENV IS_POOL=true
ENV ES_LINK=http://elasticsearch:9200
ENV FT_HTTP=150
ENV RE_PATH=/var/log/agent_ws
ENV RE_NORMAL=7
ENV RE_ATTACK=90
ENV TZ=UTC
CMD ["./agent_ws"]
