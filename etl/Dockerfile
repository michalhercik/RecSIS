# syntax=docker/dockerfile:1

#############################################################################################
# Build stage
#############################################################################################
FROM golang:1.23 AS go-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /etl

#############################################################################################
# Dev stage
#############################################################################################
FROM alpine AS ssh-stage
WORKDIR /app
COPY --from=go-stage /etl /etl
COPY ./docker-cmd.sh ./docker-cmd.sh
RUN apk add --update --no-cache openssh
RUN apk add --update --no-cache autossh
CMD ./docker-cmd.sh
# CMD sshpass -p "$SSH_PASS" ssh -T -L 10502:mffout.karlov.mff.cuni.cz:10502 hercik@acheron.ms.mff.cuni.cz -p 42049 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null && /etl
# CMD /bin/sh

#############################################################################################
# Deploy stage
#############################################################################################
FROM gcr.io/distroless/base-debian11 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /recsis /recsis
COPY ./static ./static
EXPOSE 8000
USER nonroot:nonroot
ENTRYPOINT ["/recsis"]
