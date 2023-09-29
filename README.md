# fj
```fj```コマンドはAtCoder Heuristic Contestの問題を解くことを助けるツールです。 このコマンドは、テストの実行を自動化します。
(注意)現在、AHC023,AHC022,を対象に作成しています。
## インストール
最新のGo(1.21)が導入されている環境であれば、次のコマンドだけでインストールできます。
```
go install github.com/fmhr/fj/cmd/fj@latest
```
OSにはLinux(Windows Subsystem for Linux)かmacOSを推奨しますが、Windows上でも動作します。
詳細な手順については(準備中)を読んでください。

# Features
- テストの実行
- リアクティブ問題に対応
- 並列実行
- 実行時に"N=3"のような標準エラー出力をすると、自動で収集してテストケースごとに表示する


## 使用方法
1. AtCoderの公式サイトからローカルテスターをダウンロードします。
2. ダウンロードしたテスターをコンテスト用のフォルダに解凍します。
3. tools/のREADMEを参照し、テストケースを生成します。
  - リアクティブ問題の場合、testerを実行してください。
  - それ以外の場合は、visを実行して実行ファイルを生成してください。

## 設定ファイル(fj_config.toml)を生成
コンテスト用のフォルダに移動して.以下のコマンドを実行する。
```
fj inti
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
TimeLimitMS = 5000     タイムアウト時間。コンテストのTLEではなく無限ループなどで終了しない時用
```


## サンプル 

テストケース、seed=0を実行
```
fj test
```
seed=777をテスト
```
fj test 777
```
seed(3~99)をテスト
```
fj tests -start 3 -end 100
```
seed(0~99)をテスト
```
fj tests -e 100
```
