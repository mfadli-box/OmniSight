# Stage 1: Build Go binary
FROM golang:1.26.4-alpine AS builder
RUN apk add --no-cache git ca-certificates

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o agent_vm .

# Stage 2: Minimal runtime image
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata docker-cli docker-cli-compose openssl git

WORKDIR /root/

COPY --from=builder /app/agent_vm .

ENV API_URL=http://localhost:3666
ENV KEY_ROBOT=
ENV IDA_HOST=
ENV IVL_POOL=30
ENV IVL_HEARTBEAT=60
ENV SOC_DOCKER=/var/run/docker.sock
ENV DIR_NGINX=/etc/nginx
ENV DIR_LEGO=/etc/lego
ENV DIR_TRIVY=/tmp/trivy
ENV LOG_LEVEL=info
ENV TZ=UTC

CMD ["./agent_vm"]
