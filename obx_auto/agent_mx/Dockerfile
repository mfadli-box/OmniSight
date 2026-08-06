# Stage 1: Build Go binary
FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o agent_mx .

# Stage 2: Minimal runtime image
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/agent_mx .

ENV API_URL=http://localhost:3666
ENV KEY_ROBOT=
ENV IDA_HOST=
ENV IPP_SYSLOG=0.0.0.0:514
ENV DIR_SYSLOG=./syslog_buffer
ENV DIR_DECRYPT=./decrypt
ENV LOG_LEVEL=info
ENV IVL_POOL=30
ENV IVL_HEARTBEAT=60
ENV TZ=UTC

CMD ["./agent_mx"]
