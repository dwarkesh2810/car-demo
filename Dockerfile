# Build Golang binary
FROM golang:1.21.5 AS build-golang

WORKDIR /app

COPY . .
RUN go build -o /car_demo

EXPOSE 8080
EXPOSE 8000
CMD ["/car_demo"]