FROM golang:1.19 AS build
RUN mkdir /app
WORKDIR /app

ADD go.mod /app
ADD go.sum /app
ADD golang-server /app/golang-server
RUN go mod download

RUN mv /etc/localtime /etc/localtime-backup
RUN ln -s /usr/share/zoneinfo/Europe/Madrid /etc/localtime
RUN mkdir bin

ADD cmd /app/cmd
# ADD contracts /app/contracts
ADD pkg /app/pkg
ADD go.mod /app/go.mod
ADD go.sum /app/go.sum

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/obac-idaas-server cmd/rest_server/rest_server.go

FROM alpine:3.16.2
WORKDIR /app
COPY resources /app/resources
COPY --from=build /app/bin .

CMD [ "./obac-idaas-server" ]
