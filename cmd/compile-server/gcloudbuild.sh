#.gcloudignore を忘れないように。

gcloud builds submit --tag gcr.io/fj-test-399812/go-compiler:latest .
gcloud run deploy go-compiler\
    --image gcr.io/fj-test-399812/go-compiler:latest \
    --platform managed --allow-unauthenticated