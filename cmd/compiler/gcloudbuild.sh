#.gcloudignore を忘れないように。
set -e #コマンドが失敗したらそこで終了する
set -x #実行したコマンドを表示する
gcloud builds submit --tag gcr.io/fj-test-399812/go-compiler:latest .
gcloud run deploy go-compiler\
    --image gcr.io/fj-test-399812/go-compiler:latest \
    --platform managed --allow-unauthenticated