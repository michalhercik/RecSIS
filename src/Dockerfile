# syntax=docker/dockerfile:1

#############################################################################################
# Build stage
#############################################################################################
FROM golang:1.23 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@v0.2.793
COPY ./ ./
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o /recsis

#############################################################################################
# Dev stage
#############################################################################################
FROM alpine AS dev-stage
WORKDIR /app
COPY --from=build-stage /recsis ./recsis
COPY ./static ./static
COPY ./config.prod.toml ./config.prod.toml
EXPOSE 8000
EXPOSE 80
ENTRYPOINT ["./recsis", "--config", "/app/config.prod.toml"]

#############################################################################################
# Deploy stage
#############################################################################################
FROM gcr.io/distroless/base-debian11 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /recsis /recsis
COPY ./static ./static
COPY ./config.prod.toml ./config.prod.toml
EXPOSE 8000
EXPOSE 80
USER nonroot:nonroot
ENTRYPOINT ["/recsis --config /app/config.prod.toml"]
