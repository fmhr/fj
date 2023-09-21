# Dockerfile は問題フォルダの直下に置く前提
# Rustのツールをビルドするための新しいステージ
FROM rust:latest as rust-builder
WORKDIR /rust-app
COPY rust-app/ ./tools
RUN cargo update
RUN cargo build --release


## ここは提出するコードに依存する
#FROM golang:1.21 as go-builder
## 作業ディレクトリの設定
#WORKDIR /app
## Goのソースコードをコピー
#COPY . .
## Goアプリのビルド
#RUN go build -o main .


# 実行用のイメージの作成
FROM debian:buster-slim
COPY --from=go-builder /app/main /app/main
COPY --from=rust-builder /rust-app/target/release/rust-app /app/rust-app

# ポートを設定
EXPOSE 8080

# 実行コマンド
CMD ["/app/main"]
