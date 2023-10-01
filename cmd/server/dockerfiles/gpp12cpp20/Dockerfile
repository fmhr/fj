# 1. サーバーのビルド
FROM golang:latest as serverbuilder
RUN CGO_ENABLED=0 go install github.com/fmhr/fj/cmd/server@v0.0.19

# 2. 公式ツールのビルド
FROM rust:buster as toolsbuilder
WORKDIR /work

# IMPORTANT: 問題によってテスターは異なるので、適宜URLを変更すること
# ------------------------------------------------------------------
RUN wget -q https://img.atcoder.jp/ahc024/AU5KcDyn.zip -O tools.zip 
# ------------------------------------------------------------------
RUN unzip tools.zip 

RUN (cd tools && cargo build --release --quiet)

# 3. ソースコードのビルド
# INPORTANT: 言語によって環境が異なるので、適宜変更すること
# ------------------------------------------------------------------
FROM gcc:12.3 as builder

WORKDIR /work

# AtCoder Library (ACL)のダウンロードと配置
RUN wget -q https://github.com/atcoder/ac-library/releases/download/v1.5.1/ac-library.zip -O acl.zip \
    && unzip acl.zip \
    && rm -rf acl acl.zip
 
# ソースファイルのコピー
COPY Main.cpp /work

# ソースファイルのコンパイル
# -static　オプションを追加していることに。注意してください。　軽量なイメージを作るために必要です。
RUN g++-12 -std=gnu++20 -O2 -DONLINE_JUDGE -DATCODER -Wall -Wextra \
    -mtune=native -march=native -fconstexpr-depth=2147483647 -fconstexpr-loop-limit=2147483647 \
    -fconstexpr-ops-limit=2147483647 -I/opt/ac-library -I/opt/boost/gcc/include -L/opt/boost/gcc/lib \
    -o a.out Main.cpp -lgmpxx -lgmp -I/usr/include/eigen3 -static
# -------------------------------------------------------------------------

# 4. 実行環境の構築
FROM debian:buster

WORKDIR /work

COPY --from=serverbuilder /go/bin/server /usr/local/bin/server

COPY --from=toolsbuilder /work/tools/target/release/* /work/tools/target/release/

# IMPORTANT: 言語によって実行ファイルが異なるので、適宜変更すること
COPY --from=builder /work/a.out /work/a.out
RUN echo -n "./a.out" > cmd.txt
ENV EXECUTE_COMMAND "./a.out"

EXPOSE $PORT
CMD ["server"]


# メモ
# docker build -t gpp-compiler .
# docker run -it  gpp-compiler /bin/bash
# docker run -it --rm -p 8080:8080 gpp-compiler
