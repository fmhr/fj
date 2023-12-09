# fj
```fj```コマンドはAtCoder Heuristic Contestの問題を解くことを助けるツールです。 このコマンドは、テストの実行を自動化します。
(注意)現在、AHC023,AHC022,を対象に作成しています。
## インストール
最新のGo(1.21)が導入されている環境であれば、次のコマンドだけでインストールできます。
```
go install github.com/fmhr/fj/cmd/fj@latest
```
すでにインストールしていて、アップデートをする場合は、次のコマンドを実行してください。
または＠latestを具体的なバージョンに置き換えてください。
```
go clean -modcache
go install github.com/fmhr/fj/cmd/fj@latest
```
OSにはLinux(Windows Subsystem for Linux)かmacOSを推奨しますが、Windows上でも動作します。
詳細な手順については(準備中)を読んでください。

# Features
- テストの実行
- リアクティブ問題に対応
- 並列実行
- 実行時に"N=3"のような標準エラー出力をすると、自動で収集してテストケースごとに表示する
- Googole Cloud Runを使った並列テスト


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
fj tests 100
```


設定ファイル(fj/config.toml)の設定例
```
Language = 'Go'
Cmd = './bin/a.out'
Args = []
Reactive = true
TesterPath = 'tools/target/release/tester'
VisPath = 'tools/target/release/vis'
GenPath = 'tools/target/release/gen'
InfilePath = 'tools/in/'
OutfilePath = 'out/'
Jobs = 3
Cloud = false
CloudURL = 'http://localhost:8888'
CompilerURL = 'http://localhost:8080/compiler'
WorkerURL = 'http://localhost:8081/worker'
Source = 'src/main.go'
CompileCmd = 'go build -o bin/main src/main.go'
Binary = 'bin/main'
ConcurrentRequests = 1000
TimeLimitMS = 10000

```
## Google Cloud Runを使った並列テスト

### ローカルでテストする。
Gcloudで動かす前に、ローカルでテストすることをおすすめします。
#### 1. Docker Desktopをインストールする。
[https://docs.docker.com/engine/install/]を参照してインストールしてください。
#### 2. fj setupCloudを実行する。
1. ```fj setupCloud ```を実行して必要なDockerfileを生成します。
2. ./fj/compiler/Dockerfile を自身の使う言語に合わせて編集してください。
3. ./fj/compiler/に移動して、sh localbuild.sh を実行してください。
4. 同様に ./fj/worker/に移動して、localbuild.sh を実行してください。
5.　設定ファイル(fj/config.toml)を自分の環境に合わせて編集します。
```
設定例
CompilerURL = 'http://localhost:8080/compiler'
Source = 'src/main.go'
CompileCmd = 'go build -o bin/main src/main.go'
Binary = 'bin/main'
WorkerURL = 'http://localhost:8081/worker'
ConcurrentRequests = 1
```
ConcurrentRequestsは、並列実行するときのリクエスト数です。ローカルではコンテナを１つしか立ち上げないので、1にしてください。
設定後に、```fj test --cloud```を実行して、正しく動作するか確認してください。