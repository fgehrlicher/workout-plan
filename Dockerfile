FROM golang:1.12-alpine as build-env
RUN apk add git make

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make

FROM alpine:3.9

WORKDIR /usr/local/bin
COPY --from=build-env /app/bin/server ./server
COPY --from=build-env /app/bin/claim-generator ./claim-generator
COPY --from=build-env /app/bin/validator ./validator

ENTRYPOINT ["server"]

