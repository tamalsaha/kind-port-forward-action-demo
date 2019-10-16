###################################
#Build stage
FROM golang:1.12-alpine3.10 AS build

ARG GOPROXY
ENV GOPROXY ${GOPROXY:-direct}

ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GOFLAGS -mod=vendor

ADD . /go/src/app
WORKDIR /go/src/app
RUN go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/static
COPY --from=build /go/bin/app /
CMD ["/app"]
