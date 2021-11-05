package lider

import (
	"log"
	//"strconv"
	//"os"
	//"time"

	//"fmt"

	//"golang.org/x/net/context"
)

type Connection struct {
	stream PlayerService_PlayerHandlerServer
	id     int
	active bool
	error  chan error
}

type Server struct {
	Connection []*Connection
	Fase 	 int
}

func check_error(e error, msg string) bool {
    if e != nil {
		log.Printf("%s", msg)
        panic(e)
		return true
    }

	return false
}

func ConnectPlayer(req *PlayerRequest, stream PlayerService_PlayerHandlerServer, s *Server) error {
	if len(s.Connection) <= 16 {

		player_id := len(s.Connection)+1
		conn := &Connection{
			stream: stream,
			id: player_id,
			active: true,
			error:  make(chan error),
		}

		resp := PlayerResponse{Type: 0, Player: int32(player_id)}
		if err := stream.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}

		s.Connection = append(s.Connection, conn)

		log.Printf("Player connected. %d/16", player_id)
		return <-conn.error
	} else{
		log.Printf("Error al conectar al jugador")
		return nil
	}
	return nil
}

// func PlayersConnections(stream PlayerService_PlayerHandlerServer, s *Server) {
// 	for {
// 		req, err := stream.Recv()
// 		if check_error(err, "Error when reading from stream") {
// 			return
// 		}

// 		if req.Type == 1 { // player trying to connect
// 			ConnectPlayer(req, stream, s)
// 		} else{
// 			return
// 		}
// 	}
// }

func (s *Server) PlayerHandler(stream PlayerService_PlayerHandlerServer) error {
	req, err := stream.Recv()
	if check_error(err, "Error when reading from stream") {
		return err
	}
	
	if req.Type == 0 { // player trying to connect
		ConnectPlayer(req, stream, s)
	} else{
		return nil
	}
	return nil

}
