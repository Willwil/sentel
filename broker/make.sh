GOARCH=amd64 GOOS=linux go build
docker build -t sentel/broker .
