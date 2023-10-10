FROM golang:1.21 as workerbuild
WORKDIR /work
RUN git clone -b upgrade  https://github.com/fmhr/fj.git \
    cd fj/cmd/worker \
    && CGO_ENABLED go build -o /work/worker


FROM debian:buster
WORKDIR /app
COPY --from=workerbuild /work/worker /app/worker
EXPOSE 8080
CMD ["/app/worker"]