# ローカルで動かすための手順
#  docker build --no-cache -t fj-server .# ローカルで動かすための手順
#  docker run -p 8888:8080 --rm fj-server
#  docker exec -it container_name sh

FROM rust:buster as rust-builder
WORKDIR /work
COPY /tools .
# コンテストによって、以下のコメントアウトを切り替える
# gen
RUN cargo build --release --bin=gen
RUN cp target/release/gen .
RUN chmod +x /work/gen
# vis リアクティブ問題のときコメントアウト
RUN cargo build --release --bin=vis
RUN cp target/release/vis .
RUN chmod +x /work/vis
# tester リアクティブな問題の時だけ使う
#RUN cargo build --release --bin=tester
#RUN cp target/release/tester .
#RUN chmod +x /work/tester

# fj/server
FROM golang:latest as golang-builder
RUN git clone -b severBuild https://github.com/fmhr/fj.git /go/src/github.com/fmhr/fj
WORKDIR /go/src/github.com/fmhr/fj/cmd/server
RUN CGO_ENABLED=0 go build -o fj-server
#COPY /Users/fumihiro/Desktop/fj /go/src/github.com/fmhr/fj
#WORKDIR /go/src/github.com/fmhr/fj/cmd/server
#RUN CGO_ENABLED=0 go build -o fj-server

# gここは提出言語によって変える
FROM golang:1.20 as golang-builder2
WORKDIR /work
COPY src/ .
RUN go mod init main
RUN CGO_ENABLED=0 go build -o main



# 新しいステージを開始し、最小限のイメージをベースにする
# gcloudにログインしてないと失敗する $ gcloud auth login
# プロジェクトをたててGoogle Container Registry APIを有効にする
#FROM gcr.io/google-cloud-builders/debian-slim

FROM debian:buster
WORKDIR /app
# バイナリをコピー
#COPY --from=rust-builder /work/tester /app/tester
COPY --from=rust-builder /work/gen /app/gen
COPY --from=rust-builder /work/vis /app/vis
COPY --from=golang-builder /go/src/github.com/fmhr/fj/cmd/server/fj-server /app/fj-server
COPY --from=golang-builder2 /work/main /app/main

# ローカル環境で必要な時だけ使う
#COPY ahc022-d9fa70931824.json /app/service_account_key.json
#ENV GOOGLE_APPLICATION_CREDENTIALS /app/service_account_key.json


RUN echo -n "./main" > cmd.txt
ENV PATH /app:$PATH

# ポート8080を公開
EXPOSE 8080

# バイナリを実行
CMD ["/app/fj-server"]
