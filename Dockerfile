# Build stage
FROM golang:1.25 AS build-stage
WORKDIR /build
COPY . .
RUN go get -d -v ./... && \
  CGO_ENABLED=0 GOOS=linux go build -o app cmd/main.go

# Release stage
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build-stage /build/app /app

USER nonroot:nonroot

EXPOSE 8080
CMD ["/app"]