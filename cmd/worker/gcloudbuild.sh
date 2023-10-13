#.gcloudignore を忘れないように。
set -e #コマンドが失敗したらそこで終了する
set -x #実行したコマンドを表示する
gcloud builds submit --tag gcr.io/fj-test-399812/worker-image:latest .
gcloud run deploy worker\
    --image gcr.io/fj-test-399812/worker-image:latest \
    --platform managed --allow-unauthenticated \
    --concurrency 1
