FROM docker.io/golang:1.11 AS development
MAINTAINER Tim Jones <timniverse@gmail.com>

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/TimJones/godotty
