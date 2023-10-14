set -xe
docker build -t worker-local .
docker run -p 8080:8080 worker-local