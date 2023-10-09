FROM golang:1.21
WORKDIR /app
COPY . ./
RUN  go mod download
RUN CGO_ENABLED=0 go build -o server
EXPOSE 8080
CMD ["/app/server"]