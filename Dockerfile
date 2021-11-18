FROM golang:1.17 AS build

# Build a statically linked binary
ENV CGO_ENABLED=0
# Targetting linux/docker.
ENV GOOS=linux

WORKDIR /go/src/github.com/meyersh/tcrp/
COPY main.go go.mod ./
RUN go build

FROM scratch
COPY --from=build /go/src/github.com/meyersh/tcrp/tcrp /
CMD ["/tcrp"]
