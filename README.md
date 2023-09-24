# fj
fj コマンドはAtCoder Heuristic Contestの問題を解くことを助けるツールです。 このコマンドは、テストの実行を自動化します。
# Features
- テストの実行
- リアクティブ問題に対応
- 並列実行
- 実行時に"N=3"のような標準エラー出力をすると、自動で収集してテストケースごとに表示する
## In progress
- Google Cloud による並列実行
- 実行時間制限
- 
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
設定ファイル(fj_config.toml)の設定例
```
Cmd = './bin/main'  実行時のコマンド
Args = []           必要に応じて設定
Reactive = false    リアクティブ問題のときはtrue
TesterPath = 'tools/target/release/tester'   
VisPath = 'tools/target/release/vis'
GenPath = 'tools/target/release/gen'
InfilePath = 'tools/in/'
OutfilePath = 'out/'
Jobs = 7             並列実行数 CPUコア数以下がいい
Cloud = false        準備中
CloudURL = 'http://localhost:8888' 準備中
TimeLimitMS = 0      準備中
```

テストケース、seed=0を実行
```
fj -mode test
```
seed=777を実行
```
fj -mode test -seed 777
```

テストケース、seed=３〜99を実行
```
fj -mode test -start 3 -end 100
```

# Example
# FAQ
