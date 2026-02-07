#!/bin/bash

set -e

AC_GO_VERSION=1.25.1

cd /tmp                                    
wget -q https://go.dev/dl/go${AC_GO_VERSION}.linux-amd64.tar.gz
tar -C /opt -xf go${AC_GO_VERSION}.linux-amd64.tar.gz
export PATH=$PATH:/opt/go/bin 

mkdir -p /go/src/atcoder.jp/golang
cd /go/src/atcoder.jp/golang

go mod init atcoder.jp/golang

go get -u \
  github.com/emirpasic/gods \
  gonum.org/v1/gonum \
  github.com/liyue201/gostl \
  github.com/benbjohnson/immutable \
  golang.org/x/exp \
  github.com/monkukui/ac-library-go