package lider

import (
	"log"
	//"strconv"
	//"os"
	"time"

	//"fmt"

	//"golang.org/x/net/context"
)

type Connection struct {
	id     int
	active bool
	error  chan error
}

type Server struct {
	Connection []*Connection
	Fase 	 int
	Max_players int
	Connected_players int
	Change_fase bool
}

func check_error(e error, msg string) bool {
    if e != nil {
		log.Printf("%s", msg)
        log.Printf("Error: %v", e)
		return true
    }

	return false
}

func ConnectPlayer(req *PlayerRequest, stream PlayerService_PlayerHandlerServer, s *Server) error {
	if len(s.Connection) <= 16 {

		player_id := len(s.Connection)+1
		conn := &Connection{
			id: player_id,
			active: true,
			error:  make(chan error),
		}

		resp := PlayerResponse{Type: 0, Player: int32(player_id)}
		if err := stream.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}

		s.Connection = append(s.Connection, conn)
		s.Connected_players += 1

		log.Printf("Player connected. %d/16", player_id)
		return nil
	} else{
		log.Printf("Error al conectar al jugador")
		resp := PlayerResponse{Type: 0, Player: -1}
		if err := stream.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
		return nil
	}
	return nil
}

func (s *Server) PlayerHandler(stream PlayerService_PlayerHandlerServer) error {
	
	ctx := stream.Context()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	max_players := 2
	if len(s.Connection) < max_players{
		log.Println("Esperando a los jugadores")
		req, err := stream.Recv()
		if check_error(err, "Error when reading from stream") {
			return err
		}
		ConnectPlayer(req, stream, s)
	}
	return nil
}

func (s *Server) WaitingRoom(stream PlayerService_WaitingRoomServer) error {
	//receive player request
	req, err := stream.Recv()
	if check_error(err, "Error when reading from stream") {
		return err
	}
	if req.Type == 0 && req.Player > 0 {
		log.Printf("Jugador %d conectado a la sala de espera", req.Player)
		for s.Connected_players < s.Max_players {
		
			// send response to player
			resp := PlayerResponse{Type: 0, Response: 0}
			if err := stream.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
			// sleep for 2 seconds
			time.Sleep(5 * time.Second)
		}
	}

	for !s.Change_fase {
		// send response to player
		resp := PlayerResponse{Type: 0, Response: -1}
		if err := stream.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
		// sleep for 2 seconds
		time.Sleep(5 * time.Second)
	}
	// send response to player
	resp := PlayerResponse{Type: 0, Response: 1}
	if err := stream.Send(&resp); err != nil {
		log.Printf("send error %v", err)
	}
	return nil
}

func (s *Server) Fase1(stream PlayerService_Fase1Server) error {
	log.Printf("Fase 1 Iniciada")
	return nil
}
func (s *Server) Fase2(stream PlayerService_Fase2Server) error {
	return nil
}
func (s *Server) Fase3(stream PlayerService_Fase3Server) error {
	return nil
}

