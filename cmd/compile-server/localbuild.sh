set -e
docker build -t go-compiler-test -f Dockerfile.golang .
docker run -p 8080:8080 go-compiler-test