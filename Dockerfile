# syntax=docker/dockerfile:1

#############################################################################################
# Build stage
#############################################################################################
FROM golang:1.23 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest
COPY ./ ./
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o /recsis

#############################################################################################
# Dev stage
#############################################################################################
FROM alpine AS dev-stage
WORKDIR /app
COPY --from=build-stage /recsis /recsis
COPY --from=build-stage /app/*templ.go /app
EXPOSE 8000
ENTRYPOINT ["/recsis"]

#############################################################################################
# Deploy stage
#############################################################################################
FROM gcr.io/distroless/base-debian11 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /recsis /recsis
EXPOSE 8000
USER nonroot:nonroot
ENTRYPOINT ["/recsis"]
