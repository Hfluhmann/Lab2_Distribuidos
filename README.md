# Lab2_Distribuidos

# Install Dependencies
```bash
apt update
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go get -u github.com/golang/protobuf/protoc-gen-go
apt install -y protobuf-compiler
export PATH="$PATH:$(go env GOPATH)/bin"
```
## Run Compile
```bash
protoc --go_out=plugins=grpc:<DIROUT> --go_opt=paths=source_relative <FILENAME>.proto
```

# Install RabbitMQ
```bash
apt update
apt install rabbitmq-server
service rabbitmq-server start
rabbitmqctl add_user 'client' '1234'
rabbitmqctl set_permissions -p "name" "client" ".*" ".*" ".*"
```
