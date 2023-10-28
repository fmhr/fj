#!/bin/sh

GCLOUD_PROJECT="project2323"
GCLOUD_IMAGE_TAG="fj-worker"

set -e #コマンドが失敗したらそこで終了する
set -x #実行したコマンドを表示する

gcloud builds submit --tag gcr.io/$GCLOUD_PROJECT/$GCLOUD_IMAGE_TAG .
gcloud run deploy worker\
    --image gcr.io/$GCLOUD_PROJECT/$GCLOUD_IMAGE_TAG \
    --platform managed --allow-unauthenticated \
    --concurrency 1
