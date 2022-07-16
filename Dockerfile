FROM golang:1.18 as builder

WORKDIR /build
COPY . .
RUN env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o local-persist

# generate clean, final image for end users
FROM alpine

COPY --from=builder /build/local-persist local-persist

RUN mkdir -p /run/docker/plugins /state /docker-data

# executable
ENTRYPOINT [ "/local-persist" ]