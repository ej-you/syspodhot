# ----- #
# BUILD #
# ----- #

FROM golang:1.27.3-alpine3.24 AS build

# set up workdir
WORKDIR /go/src

# install dependences
COPY ./go.mod .
COPY ./go.sum .

# compile app
COPY ./main.go ./main.go
COPY ./config ./config
COPY ./internal ./internal
RUN go build ./main.go

# --- #
# RUN #
# --- #

FROM alpine:3.21 AS run

WORKDIR /app

# copy compiled app
COPY --from=build /go/src/main .
# copy templates and static files
COPY ./web ./web

# run app
CMD ["/bin/sh", "-c", "/app/main"]
