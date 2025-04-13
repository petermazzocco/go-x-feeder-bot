FROM golang:1.23.2 AS build
WORKDIR /go/src/app
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOPROXY=direct
RUN go build -v -o app .


FROM scratch
COPY --from=build /go/src/app/app /go/bin/app
ENTRYPOINT ["/go/bin/app"]
