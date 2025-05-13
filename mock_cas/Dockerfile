# syntax=docker/dockerfile:1

#############################################################################################
# Build stage
#############################################################################################
FROM golang:1.23 AS build-stage
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /mockcas

#############################################################################################
# Dev stage
#############################################################################################
FROM alpine AS dev-stage
WORKDIR /app
COPY --from=build-stage /mockcas /mockcas
EXPOSE 8001
ENTRYPOINT ["/mockcas"]