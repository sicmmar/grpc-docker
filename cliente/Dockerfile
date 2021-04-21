FROM golang
WORKDIR /
COPY . .
RUN go mod download
EXPOSE 8080
CMD ["go", "run", "cliente.go"]