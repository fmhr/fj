#!/bin/sh

GCLOUD_PROJECT="ahc-contests"
REPOSITORY="my-app-images"
IMAGE_NAME=""
REGION="asia-northeast1"

set -e #コマンドが失敗したらそこで終了する
set -x #実行したコマンドを表示する

if [ -z "$IMAGE_NAME" ]; then
  echo "IMAGE_NAME is not set. e.g., IMAGE_NAME=ahc001a"
  exit 1    
fi

gcloud builds submit --tag ${REGION}-docker.pkg.dev/${GCLOUD_PROJECT}/${REPOSITORY}/${IMAGE_NAME}:latest .

gcloud run deploy ${IMAGE_NAME}\
    --image ${REGION}-docker.pkg.dev/${GCLOUD_PROJECT}/${REPOSITORY}/${IMAGE_NAME}:latest \
    --platform managed --allow-unauthenticated \
    --concurrency 1