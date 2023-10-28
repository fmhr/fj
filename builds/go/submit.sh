# Cloud Buildをトリガーし、その出力を変数に保存
OUTPUT=$(gcloud builds submit --config cloudbuild.yaml .)

# 出力からBUILD_IDを抽出
BUILD_ID=$(echo "$OUTPUT" | grep "^ID" -A 1 | tail -n1 | awk '{print $1}')
# BUILD_IDを表示
echo $BUILD_ID

