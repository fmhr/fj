# *プロトタイプ
# fj
fj コマンドはAtCoder Heuristic Contestの問題を解くことを助けるツールです。 このコマンドは、テストの実行を自動化します。
# Features
- テストの実行
## In progress
- リアクティブ問題に対応する
# How to install
最新のGoをインストール

[https://go.dev/doc/install](https://go.dev/install/)
```
go intall github.con/fmhr/fj
```
# How to use
- AtCoder公式からローカルテスターをダウンロードして、コンテスト用のフォルダに。解凍する。
- READMEを参照してテストケースを生成して、必要に応じてtester,visをコンパイルする。
## 設定ファイル(fj_config.toml)を生成
コンテスト用のフォルダに移動して.以下のコマンドを実行する。
```
fj -mode inti
```
設定ファイルfj_config.tomlのCmdに実行コマンドを入れる

テストケース、seed=1を実行
```
fj -mode run
```
テストケース、seed=３〜99を実行
```
fj -mode run -start 3 -end 100
```

# Example
# FAQ