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
2. ./fj/compiler/Dockerfile を編集してください。
  自分の使う環境に合わせてDockerfileを編集してください。
3. ./fj/compiler/に移動して、localbuild.sh を実行してください。
  実行権限がない場合は、chmod +x ./fj/compiler/localbuild.shを実行してください。
4. 同様に ./fj/worker/に移動して、localbuild.sh を実行してください。
