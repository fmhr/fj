
# 2. 公式ツールのビルド
FROM rust:buster as toolsbuilder
WORKDIR /work

# IMPORTANT: 問題によってテスターは異なるので、適宜URLを変更すること
# ------------------------------------------------------------------
# RUN wget -q https://img.atcoder.jp/ahc024/AU5KcDyn.zip -O tools.zip 
RUN wget -q https://img.atcoder.jp/ahc023/asprocon10_5XtyY7PgVPtBANvV.zip -O tools.zip \
    && unzip tools.zip \
    && (cd tools && cargo build --release --quiet)


FROM golang:1.21 as workerbuild
WORKDIR /work
RUN git clone -b make-builds  https://github.com/fmhr/fj.git \
    && cd fj/cmd/worker \
    && CGO_ENABLED=0 go build -o /work/worker



# 3. ビルド用のイメージ
FROM debian:buster
WORKDIR /work
COPY --from=workerbuild /work/worker /work/worker
COPY --from=toolsbuilder /work/tools/target/release/* /work/tools/target/release/

EXPOSE $PORT
CMD ["/work/worker"]