# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /k-means-MR
COPY go.mod /k-means-MR
CMD [ "sh" ]