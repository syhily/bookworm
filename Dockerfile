# Generate
FROM ghcr.io/a-h/templ:latest AS generate-stage
COPY --chown=65532:65532 . /build
WORKDIR /build
RUN ["templ", "generate"]

# Build
FROM golang:latest AS build-stage
COPY --from=generate-stage /build /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags=all="-l -B" --ldflags="-s -w" -buildvcs=false -o /build/bookworm

# Deploy
FROM gcr.io/distroless/base-debian12 AS deploy-stage
WORKDIR /
COPY --from=build-stage /build/bookworm /bookworm
EXPOSE 4321
USER nonroot:nonroot
ENTRYPOINT ["/bookworm", "--port", "4321"]
