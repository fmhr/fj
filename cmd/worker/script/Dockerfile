FROM rust:latest as toolsbuilder
WORKDIR /work

# 直接ダウンロードするとき　例)AHC027
#RUN wget -q https://img.atcoder.jp/ahc027/aPdjCUIZ_v2.zip -O tools.zip \
#    && unzip tools.zip 

# ローカルのtoolsをコピーするときは tools/　をDockerfileと同じ階層に置く
COPY ./tools /work/tools
RUN cd tools && cargo build --release --quiet


FROM golang:1.25 as workerbuild
WORKDIR /work
RUN git clone https://github.com/fmhr/fj.git

WORKDIR /work/fj
RUN go mod tidy

WORKDIR /work/fj/cmd/worker 
RUN CGO_ENABLED=0 go build -o /work/worker


# 3. 実行時の環境
FROM ubuntu:24.04
WORKDIR /judge

RUN apt-get update \
    && apt-get install -y ca-certificates openssl time \
    && rm -rf /var/lib/apt/lists/*

COPY --from=workerbuild /work/worker /judge/worker
COPY --from=toolsbuilder /work/tools/target/release/* /judge/tools/target/release/

RUN mkdir -p /judge/bin

CMD ["/judge/worker"]