# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Cert stage
FROM alpine:latest as certs
RUN apk --update add ca-certificates

# Builder stage
FROM golang:1.13.4-alpine3.10 as builder

# Output dir
RUN mkdir -p /build

# Set the Current Working Directory inside the container
WORKDIR /build

# Copy mod file inside the container
COPY go.mod .

# Copy sum file inside the contaner
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy source inside the container
COPY . .

# Compile output
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix cgo -o /bin/go-appcenter cmd/appcenter/main.go

# Thin stage
FROM scratch
ENV PATH=/bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /bin/go-appcenter /bin/go-appcenter

CMD ["/bin/go-appcenter"]
ENTRYPOINT ["/bin/go-appcenter"]

