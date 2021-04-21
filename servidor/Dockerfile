FROM golang
WORKDIR /
COPY . .
ENV HOST=0.0.0.0:8081
RUN go mod download
EXPOSE 8081
CMD ["go", "run", "servidor.go"]