FROM golang:1.21 as serverbuild

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o server

# ここからターゲットのコンパイル環境をつくる
FROM golang:1.20 as gocompile

WORKDIR /work

ENV PATH="/opt/go/bin:${PATH}"

#RUN go install github.com/fmorenovr/gods/...
#RUN go install gonum.org/v1/gonum/...
#RUN go install github.com/liyue201/gostl/...
#RUN go install golang.org/x/exp/...

COPY --from=serverbuild /app/server /work

RUN chmod +x /work/server

EXPOSE 8080

CMD ["/work/server"]