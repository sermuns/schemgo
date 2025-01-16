FROM golang:1.23-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -ldflags "-s -w" .

FROM scratch
WORKDIR /app
COPY --from=builder /build/schemgo /bin/
ENTRYPOINT [ "schemgo" ]
