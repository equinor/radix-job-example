FROM golang:alpine3.10 as builder

RUN apk update && \
    apk add ca-certificates curl git  && \
    apk add --no-cache gcc musl-dev && \
    go get -u golang.org/x/lint/golint github.com/frapposelli/wwhrd

RUN addgroup -S -g 1001 radix-non-root-group && \
    adduser -S -u 1001 -G radix-non-root-group radix-non-root-user

WORKDIR ..

# Install project dependencies
COPY ./app ./app
WORKDIR ./app
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -installsuffix cgo -o /usr/local/bin/radix-job-example

USER 1001
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/radix-job-example"]