FROM golang:1.21-bullseye as build

WORKDIR /go/src/app

ADD . /go/src/app

RUN go get -d -v ./...

RUN CGO_ENABLED=0 go build -o /go/bin/ccliwrapper

FROM scratch

COPY --from=build /go/bin/ccliwrapper /

ENTRYPOINT ["/ccliwrapper"]
