FROM golang:1.16.4-alpine3.13 as build

COPY . /src/go

WORKDIR /src/go

ENV PORT 5000

# install all dependencies
RUN go get ./...

# build the binary
RUN go build

# Put back once we have an application
CMD ["go-heroku-server"]

EXPOSE 5000