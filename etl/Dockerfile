# syntax=docker/dockerfile:1

#############################################################################################
# Build stage
#############################################################################################
FROM golang:1.23 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /etl

#############################################################################################
# Dev stage
#############################################################################################
FROM alpine AS dev-stage
WORKDIR /app
COPY --from=build-stage /etl /etl
COPY ./docker-cmd.sh ./docker-cmd.sh
COPY ./config.docker.toml ./config.docker.toml
# RUN apk add --update --no-cache openssh
# RUN apk add --update --no-cache autossh
RUN apk add openssh
RUN apk add autossh
ENTRYPOINT ["./docker-cmd.sh"]

#############################################################################################
# Deploy stage
#############################################################################################
FROM gcr.io/distroless/base-debian11 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /etl /etl
COPY ./config.prod.toml ./config.prod.toml
USER nonroot:nonroot
ENTRYPOINT ["/etl", "--config", "/app/config.prod.toml"]
