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


## Google Cloud Runを使った並列テスト
#### 1. fj setupCloudを実行する。
1. ```fj setupCloud ```をコンテストのディレクトリで実行してcompileコンテナと実際にテストするworkerコンテナ必要なDockerfileとシェルスクリプトを生成します。
2. ./fj/compiler/ の中のファイルを使用する言語にあわせて修正します

ConcurrentRequestsは、並列実行するときのリクエスト数です。ローカルではコンテナを１つしか立ち上げないので、1にしてください。
設定後に、```fj test --cloud```を実行して、正しく動作するか確認してください。

## Docker Desktopをインストールする。
[https://docs.docker.com/engine/install/]を参照してインストールしてください。
## Google Cloud 
0. SDKをインストールする。
1. 新しいプロジェクトを作成します。　プロジェクト名を適当に　ahc027-test とします。
端末を開いて、
```gcloud auth login```を実行して、ログインします
```gcloud projects PROJECT_NAME ``` でプロジェクトを作成します.
Google Cloud　のコンソールページから作成したプロジェクトを選択して、メニューの「お支払い」から「請求先アカウントにリンク」を選択します。
Cloud Build APIの有効化:　```gcloud services enable cloudbuild.googleapis.com --project=YOUR_PROJECT_ID```
Cloud Storage APIの有効化: ```gcloud services enable storage.googleapis.com --project=YOUR_PROJECT_ID```
3. fj/worker/gcloudbuild.sh と fj/compiler/gcloudbuile.sh のGCLOUD_PROJECTを自分のプロジェクトIDに変更します。
4. 途中に、google cloud run のAPIを有効にする質問がくるので y を推します
5. fj/compiler/gcloudbuild.sh のIMAGE_NAMEを変更します。
6. fj/compiler/gcloudbuild.sh を実行して、コンパイル用のコンテナをビルドします。
スクリプトが成功すると以下の様なメッセージが表示されます。
```
(略)
DONE
Done.
Service [go-compiler] revision [go-compiler-00002-ushi] has been deployed and is serving 100 percent of traffic.
Service URL: https://go-compiler-xjfasdfae-an.a.run.app
```
5. Service URLをコピーして、fj/config.tomlのCompilerURLを変更します。
6. 同様に　fj/worker/gcloudbuild.sh を実行して、URLをfj/config.tomlのWorkerURLを変更します。
7. WebブラウザでGoogle Cloud Runのコンソールを開きます。ジャッジコンテナ(worker)を選んで、YAML のタブを開いて、autoscaling.knative.dev/maxScale: '100'　の'100'を'1000'に変更します。(現在設定できる最大値です。)
gcloud storage buckets create [BUCKET_NAME] --location=[LOCATION]でバケットを作成します。
指定したバケット名をfj/config.tomlのbucketに設定します。
