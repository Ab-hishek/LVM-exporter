FROM golang:1.13

RUN apt-get update || exit 0
RUN apt-get install -y lvm2

WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -o lvm-exporter
EXPOSE 9101
CMD [ "./lvm-exporter" ]