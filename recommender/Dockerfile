# syntax=docker/dockerfile:1

#############################################################################################
# Build stage
#############################################################################################
FROM golang:1.23 AS build-stage
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /recommender

#############################################################################################
# Dev stage
#############################################################################################
FROM alpine AS dev-stage
WORKDIR /app
COPY --from=build-stage /recommender ./recommender
COPY ./config.dev.toml ./config.dev.toml
EXPOSE 8002
ENTRYPOINT ["./recommender", "--config", "/app/config.dev.toml"]