#!/bin/sh

GCLOUD_PROJECT="ahc-contests"
REGION="asia-northeast1"
REPOSITORY="my-app-images"
IMAGE_NAME="compiler-go"
IMAGE_URL="${REGION}-docker.pkg.dev/${GCLOUD_PROJECT}/${REPOSITORY}/${IMAGE_NAME}:latest"

set -eu #コマンドが失敗したらそこで終了する
set -x #実行したコマンドを表示する

gcloud builds submit --tag ${IMAGE_URL} .

gcloud run deploy ${IMAGE_NAME}\
    --region ${REGION} \
    --image ${IMAGE_URL} \
    --platform managed \
    --allow-unauthenticated 
