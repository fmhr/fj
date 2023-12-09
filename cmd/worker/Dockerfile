
# 2. 公式ツールのビルド
FROM rust:buster as toolsbuilder
WORKDIR /work

# IMPORTANT: 問題によってテスターは異なるので、適宜URLを変更すること
# ------------------------------------------------------------------
# RUN wget -q https://img.atcoder.jp/ahc024/AU5KcDyn.zip -O tools.zip 
# RUN wget -q https://img.atcoder.jp/ahc023/asprocon10_5XtyY7PgVPtBANvV.zip -O tools.zip \
# AHC027
RUN wget -q https://img.atcoder.jp/ahc027/aPdjCUIZ_v2.zip -O tools.zip \
    && unzip tools.zip \
    && (cd tools && cargo build --release --quiet)


FROM golang:1.21 as workerbuild
WORKDIR /work
RUN git clone -b debug-for-ahc027 https://github.com/fmhr/fj.git \
    && cd fj/cmd/worker \
    && CGO_ENABLED=0 go build -o /work/worker


# 3. 実行時の環境
FROM ubuntu:22.10
WORKDIR /judge
COPY --from=workerbuild /work/worker /judge/worker
COPY --from=toolsbuilder /work/tools/target/release/* /judge/tools/target/release/
RUN mkdir -p /judge/bin

EXPOSE $PORT
CMD ["/judge/worker"]