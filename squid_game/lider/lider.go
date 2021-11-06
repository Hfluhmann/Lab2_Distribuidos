package lider

import (
	"fmt"
	"io"
	"log"

	//"os"
	"time"
	//"fmt"
	//"golang.org/x/net/context"
)

type Round struct {
	Plays []int
}

type Player struct {
	Round1 *Round
	Round2 *Round
	Round3 *Round
}

type Connection struct {
	id     int
	active bool
	jugada int
	error  chan error
}

type Server struct {
	Connection        []*Connection
	Fase              int
	Max_players       int
	Connected_players int
	Round             int
	Contestados       int
	Team1             int
	Team2             int
	Jugadores2        int
	Jugadores3        int
	Change_fase       bool
	Change_round      bool
	Players_data      [16]*Player
	Randoms           []int
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

		player_id := len(s.Connection) + 1
		conn := &Connection{
			id:     player_id,
			active: true,
			jugada: 0,
			error:  make(chan error),
		}

		round1 := &Round{}
		round2 := &Round{}
		round3 := &Round{}

		player := &Player{
			Round1: round1,
			Round2: round2,
			Round3: round3,
		}

		s.Players_data[player_id-1] = player

		resp := PlayerResponse{Type: 0, Player: int32(player_id)}
		if err := stream.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}

		s.Connection = append(s.Connection, conn)
		s.Connected_players += 1

		log.Printf("Player connected. %d/16", player_id)
		return nil
	} else {
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
	if len(s.Connection) < max_players {
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
		// sleep for 5 seconds
		time.Sleep(5 * time.Second)
	}
	// send response to player
	resp := PlayerResponse{Type: 0, Response: 1}
	if err := stream.Send(&resp); err != nil {
		log.Printf("send error %v", err)
	}
	return nil
}

//Compara los puntajes seleccionados de un participante con el lider en el primer juego, y decide si eliminarlo
func comparar(valor_jugador int, valor_lider int) bool {
	if valor_jugador < valor_lider {

		//El jugador sigue vivo
		return true

	} else {

		//El jugador debe ser eliminado
		return false
	}

}

func (s *Server) Fase1P1(stream PlayerService_Fase1P1Server) error {
	log.Printf("Fase 1 Iniciada")
	log.Printf("Ronda Server: %d. Change Round: %d", s.Round, s.Change_round)

	for !s.Change_round {
		//log.Printf("Esperando cambio de ronda")
		time.Sleep(2 * time.Second)
		resp := PlayerResponse{Type: 1, Response: 1, Round: int32(s.Round)}
		err := stream.Send(&resp)
		check_error(err, "")
	}
	// resp := PlayerResponse{Type: 1, Response: 1, Round: int32(s.Round)}
	// err := stream.Send(&resp)

	// El lider selecciona un valor entre 6 y 10
	valor_lider := s.Randoms[s.Round]

	req, err := stream.Recv()
	if err == io.EOF {
		return err
	}
	player_id := req.Player
	log.Printf("Mi valor %d es %d", s.Round+1, valor_lider)

	s.Players_data[player_id-1].Round1.Plays = append(s.Players_data[player_id-1].Round1.Plays, int(req.Play))
	check_error(err, "Error al recibir jugada")
	log.Printf("Jugador %d. Resp: %d", player_id, req.Play)

	s.Connection[req.Player-1].jugada = int(req.Play)

	if s.Connection[req.Player-1].active == true {

		var movimiento bool = comparar(int(req.Play), valor_lider) //POR HACER se esta usando un valor generico 7 como valor del juagdor
		if movimiento {
			s.Connection[req.Player-1].jugada += int(req.Play) //POR HACER se suma el puntaje que puso el jugador

			// notificar al player que sobrevivio
			resp := PlayerResponse{Type: 1, Response: 1}
			err := stream.Send(&resp)
			check_error(err, "Error al notificar al player que sobrevivio la ronda")
		} else {
			s.Connection[req.Player-1].active = false
			s.Connected_players -= 1
			// notificar al player que murio
			resp := PlayerResponse{Type: 1, Response: 0}
			err := stream.Send(&resp)
			check_error(err, "Error al notificar al player que murio en la ronda")
		}
		s.Contestados += 1
	}

	return nil
}

func (s *Server) Fase1P2(stream PlayerService_Fase1P2Server) error {
	//Eliminar jugadores que no lograron el puntaje
	req, err := stream.Recv()
	player_id := req.Player
	check_error(err, "No se pudo determinar el numero de jugador")
	if s.Connection[player_id-1].active == true {
		if s.Connection[player_id-1].jugada < 21 {
			s.Connection[player_id].active = false

			//notificar al player que murio
			resp := PlayerResponse{Type: 1, Response: 0}
			err := stream.Send(&resp)
			check_error(err, "Error al notificar al player que murio en la ronda")
		} else {
			resp := PlayerResponse{Type: 1, Response: 2}
			err := stream.Send(&resp)
			check_error(err, "Error al notificar al player que sobrevivio al juego")
		}
	}

	return nil
}
func (s *Server) Fase2(stream PlayerService_Fase2Server) error {

	valor_lider := s.Randoms[4]

	//ver que onda con los equipos y con conseguir los valores

	//sacar suma de resultados de cada team
	/*
		for i := 0; i < cant_jugadores; i++ {
			if i < cant_jugadores/2 {
				team1 += teams[i]
			} else {
				team2 += teams[i]
			}
		}*/

	if valor_lider%2 != s.Team1%2 && valor_lider%2 != s.Team2%2 {
		//Matar equipo random

		perdedor := s.Randoms[7]
		if perdedor == 0 {
			//Matar equipo 1
			log.Printf("MATADO EL EQUIPO 1")
			/*
				for i := 0; i < cant_jugadores/2; i++ {
					kill_player(sobrevivientes[i])
				}*/
		} else {

			log.Printf("MATADO EL EQUIPO 2")
			/*
				for i := cant_jugadores / 2; i < cant_jugadores; i++ {
					kill_player(sobrevivientes[i])
				}*/
		}

	} else if valor_lider%2 != s.Team1%2 && valor_lider%2 == s.Team2%2 {
		//Matar equipo 1
		log.Printf("MATADO EL EQUIPO 1")
		/*
			for i := 0; i < cant_jugadores/2; i++ {
				kill_player(sobrevivientes[i])
			}*/

	} else if valor_lider%2 == s.Team1%2 && valor_lider%2 != s.Team2%2 {
		//Matar equipo 2
		log.Printf("MATADO EL EQUIPO 2")
		/*
			for i := cant_jugadores / 2; i < cant_jugadores; i++ {
				kill_player(sobrevivientes[i])
			}*/
	}

	return nil
}
func (s *Server) Fase3(stream PlayerService_Fase3Server) error {

	valor_lider := s.Randoms[5]

	//POR HACER pedir valores y guardarlos en arreglo

	//ver diferencia de valores
	for i := 0; i < s.Jugadores3; i++ {
		teams[i] = Abs(teams[i] - valor_lider)
	}

	for i := 0; i < s.Jugadores3/2; i++ {
		if teams[2*i] == teams[(2*i)+1] {
			fmt.Println("AMBOS VIVEN")
			// ambos viven
			/*
				winner_player(sobrevivientes[2*i])
				winner_player(sobrevivientes[(2*i)+1])*/
		} else if teams[2*i] < teams[(2*i)+1] {
			// gana participante 1
			// muere participante 2
			/*
				kill_player(sobrevivientes[(2*i)+1])
				winner_player(sobrevivientes[2*i])*/
		} else {
			// gana participante 2
			// muere participante 1
			/*
				kill_player(sobrevivientes[2*i])
				winner_player(sobrevivientes[(2*i)+1])*/
		}

	}

	return nil
}
