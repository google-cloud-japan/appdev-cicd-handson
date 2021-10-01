FROM golang:1.17-alpine3.14 AS build
RUN apk --no-cache add gcc musl-dev git
COPY . /go/src/
WORKDIR /go/src/
RUN go build -a -installsuffix cgo -ldflags "-s -w" -o /app

FROM alpine:3.14
RUN apk add --no-cache ca-certificates
COPY --from=build /app /app
ENTRYPOINT ["/app"]
