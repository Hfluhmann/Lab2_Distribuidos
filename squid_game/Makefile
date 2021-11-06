proto:
	@protoc --go_out=plugins=grpc:lider --go_opt=paths=source_relative lider.proto

proto_pozo:
	@protoc --go_out=plugins=grpc:pozo --go_opt=paths=source_relative pozo.proto

lider: proto
	@go run lider_server.go

player: proto
	@go run player_server.go

pozo: proto_pozo
	@go run pozo_server.go

test_pozo: proto_pozo
	@go run test_pozo.go