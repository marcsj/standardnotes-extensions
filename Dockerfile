FROM golang:1.16.3-alpine AS build

RUN apk add --no-cache ca-certificates
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix nocgo -o app cmd/server/main.go


FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /src/standardnotes-extensions /
COPY --from=build /src/definitions /definitions
ENTRYPOINT ["/app"]
