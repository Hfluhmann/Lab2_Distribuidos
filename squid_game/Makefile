proto:
	@export PATH="$PATH:$(go env GOPATH)/bin"
	@protoc --go_out=plugins=grpc:lider --go_opt=paths=source_relative lider.proto

proto_pozo:
	@export PATH="$PATH:$(go env GOPATH)/bin"
	@protoc --go_out=plugins=grpc:pozo --go_opt=paths=source_relative pozo.proto

proto_name:
	@export PATH="$PATH:$(go env GOPATH)/bin"
	@protoc --go_out=plugins=grpc:name --go_opt=paths=source_relative name.proto

proto_data:
	@export PATH="$PATH:$(go env GOPATH)/bin"
	@protoc --go_out=plugins=grpc:data --go_opt=paths=source_relative data.proto

lider: proto
	@go run lider_server.go

player: proto
	@go run player_server.go

pozo: proto_pozo
	@go run pozo_server.go

name: proto_name
	@go run name_server.go

data: proto_name proto_data
	@go run data_server.go

test_pozo: proto_pozo
	@go run test_pozo.go