#!/bin/sh

# .gcloudignore　必要に応じて追加
# woekerのイメージはコンテストに応じて変更する

set -e #コマンドが失敗したらそこで終了する
set -x #実行したコマンドを表示する

if [ -z "$GCLOUD_PROJECT" ] || [ -z "$GCLOUD_IMAGE_TAG" ]; then
    echo "GCLOUD_PROJECT or GCLOUD_IMAGE_TAG is not set."
    exit 1
fi

gcloud builds submit --tag gcr.io/$GCLOUD_PROJECT/$GCLOUD_IMAGE_TAG .
gcloud run deploy worker\
    --image gcr.io/$GCLOUD_PROJECT/$GCLOUD_IMAGE_TAG \
    --platform managed --allow-unauthenticated \
    --concurrency 1
