docker build -t go-compiler-test .
docker run -p 8080:8080 go-compiler-test