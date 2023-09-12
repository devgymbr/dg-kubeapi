FROM golang:1.21.1 as base
WORKDIR /src/devgym
COPY go.mod go.sum ./
COPY . . 
RUN go build -o kube main.go
EXPOSE 8080
CMD ["./kube"]