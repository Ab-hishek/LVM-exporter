FROM golang:alpine

RUN apk update \
    && apk add lvm2

WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -o lvm-exporter
EXPOSE 9101
CMD [ "./lvm-exporter" ]