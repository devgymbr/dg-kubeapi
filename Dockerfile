FROM golang:1.21.1 as base
RUN apk update 
WORKDIR /src/devgym
COPY go.mod go.sum ./
COPY . . 
RUN go build -o kube main.go

FROM alpine:3 as binary
COPY --from=base /src/devgym/kube .
EXPOSE 8080
CMD ["./kube"]