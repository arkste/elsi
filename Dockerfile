FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
ENV GO111MODULE=on
ENV CGO_ENABLED=0
COPY . $GOPATH/src/arkste/elsi/
WORKDIR $GOPATH/src/arkste/elsi/
RUN go get -d -v
RUN go build -o /go/bin/elsi

FROM scratch
COPY --from=builder /go/bin/elsi /go/bin/elsi
ENTRYPOINT ["/go/bin/elsi"]