ARG DBFILE='/app/data/scheduler.db'
ARG PORT=7540
ARG TODO_PASSWORD

FROM golang:1.22-alpine3.19 as builder

WORKDIR /app
COPY . .

# Install gcc
RUN apk add build-base
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -a -o todoshka ./cmd/main.go

FROM alpine:3.19

ARG DBFILE
ARG PORT
ARG TODO_PASSWORD

VOLUME /app/data

EXPOSE $PORT
ENV TODO_DBFILE=${DBFILE}
ENV TODO_PORT=${PORT}
ENV TODO_PASSWORD=${TODO_PASSWORD}

# place all necessary executables and other files into /app directory
WORKDIR /app/
COPY --from=builder app/todoshka todoshka
COPY --from=builder app/web web

CMD ["./todoshka"]