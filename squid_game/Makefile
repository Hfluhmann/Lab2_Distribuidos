proto:
	protoc --go_out=plugins=grpc:lider --go_opt=paths=source_relative lider.proto

lider: proto
	go run lider_server.go

player: proto
	go run player_server.go