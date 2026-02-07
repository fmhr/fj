# fj
```fj```コマンドはAtCoder Heuristic Contestの問題を解くことを助けるツールです。 このコマンドは、テストの実行を自動化します。
## インストール
Go(1.21~)が導入されている環境であれば、次のコマンドだけでインストールできます。
```
go install github.com/fmhr/fj/cmd/fj@latest
```
すでにインストールしていて、アップデートをする場合は、次のコマンドを実行してください。
または＠latestを具体的なバージョンに置き換えてください。
```
go clean -modcache
go install github.com/fmhr/fj/cmd/fj@latest
```
macOS上で開発しています。

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

fj/ ディレクトリにfj_config.tomlが生成されるので、エディタで開いて設定を変更してください。

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
1. 新しいプロジェクトを作成します。

gcloud Cloud にログインします。

```gcloud auth login```

## 1. プロジェクトを作成して選択
```gcloud projects create PROJECT_NAME && gcloud config set project PROJECT_NAME```

## 2. APIを有効化

```
gcloud services enable \
  cloudbuild.googleapis.com \
  storage.googleapis.com \
  artifactregistry.googleapis.com
```

必要に応じてGoogle Cloud　のコンソールページから作成したプロジェクトを選択して、メニューの「お支払い」から「請求先アカウントにリンク」を選択します。

## 3. DockerImageを保存するためのレジストリを作成

```gcloud artifacts repositories create REPOSITORY_NAME --repository-format=docker --location=asia-northeast1```

## 4. コンパイル後の実行ファイルを保存するためのバケットを作成

```gcloud storage buckets create gs://YOUR_BUCKET_NAME --location=asia-northeast1```

## 5. Cloud Build サービスアカウントに Artifact Registry へのアクセス権限を付与 

````
gcloud projects add-iam-policy-binding PROJECT_ID \
    --member="serviceAccount:PROJECT_NUMBER@cloudbuild.gserviceaccount.com" \
    --role="roles/artifactregistry.writer"
````
PROJECT_ID: あなたの Google Cloud プロジェクト ID を指定します。
PROJECT_NUMBER: あなたの Google Cloud プロジェクト番号を指定します。


ローカルのfj/config.tomlのbucketに指定したYOUR_BUCKET_NAMEを設定


## コンテナのビルド
``` fj setupCloud```を実行すると、fj/compiler/とfj/worker/にDockerfileとgcloudbuild.shが生成されます。

### コンパイラコンテナのビルド
1. fj/compiler/gcloudbuild.sh のGCLOUD_PROJECTとIMAGE_NAMEを変更します。
2. fj/compiler/Dockerfileをコンパイルする言語に合わせて修正します。
3. fj/compiler/に移動して、gcloudbuild.sh を実行して、コンパイル用のコンテナをビルドします。

```Do you want enable these APIs to continue (this will take a few minutes)? (y/N)?  ```と聞かれたら、```y```を入力してください。

スクリプトが成功すると以下の様なメッセージが表示されます。
```
(略)
DONE
Done.
Service [go-compiler] revision [go-compiler-00002-ushi] has been deployed and is serving 100 percent of traffic.
Service URL: https://go-compiler-abcdefghij-an.a.run.app
```
4. Service URLをコピーして、fj/config.tomlのCompilerURLに先ほどのURLの末尾に/compilerを追加して設定します。
### ジャッジコンテナのビルド
1. AtCoderが配布しているローカルtoolsを解答して、toolsディレクトリをfj/worker/にコピーします。
2. fj/worker/gcloudbuild.sh のGCLOUD_PROJECTとIMAGE_NAMEを変更します。
3. fj/worker/に移動して gcloudbuild.sh を実行して、コンパイラコンテナ同様 Service URL: をfj/config.tomlのWorkerURLの末尾に/workerを追加して設定します。
7. WebブラウザでGoogle Cloud Runのコンソールを開きます。ジャッジコンテナ(worker)を選んで、YAML のタブを開いて、autoscaling.knative.dev/maxScale: '100'　の'100'を'1000'に変更します。(現在設定できる最大値です。)
