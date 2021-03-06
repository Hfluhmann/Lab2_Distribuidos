package lider

import (
	"fmt"
	"io"
	"log"

	"os"
	"time"
	//"fmt"
	"golang.org/x/net/context"
	"github.com/streadway/amqp"
	"Lab2_Distribuidos/squid_game/name"
	"google.golang.org/grpc"
	"github.com/joho/godotenv"
)

func get_env_var(key string) string {

	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }
  

type Connection struct {
	Id     int
	Active bool
	Jugada int
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
	Writed		      bool
	R1				  int
	R2				  int
	R3				  int
	R4				  int
	Randoms           []int
	JugadoresFase2    []int
	JugadoresFase3    []int
	RespuestasFase3   [16]int
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
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
			Id:     player_id,
			Active: true,
			Jugada: 0,
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

	max_players := s.Max_players
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
func depositar(player int32, ronda int32) {
	ip := get_env_var("IP_POZO")
	conn, err := amqp.Dial("amqp://"+get_env_var("RM_USER")+":"+get_env_var("RM_PASS")+"@"+ip+":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
	"hello", // name
	false,   // durable
	false,   // delete when unused
	false,   // exclusive
	false,   // no-wait
	nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := fmt.Sprintf("%d %d", player, ronda)
	err = ch.Publish(
	"",     // exchange
	q.Name, // routing key
	false,  // mandatory
	false,  // immediate
	amqp.Publishing {
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	failOnError(err, "Failed to publish a message")

}
func (s *Server) Fase1P1(stream PlayerService_Fase1P1Server) error {

	log.Printf("Fase 1 Iniciada")

	//receive player request
	req, err := stream.Recv()
	round := req.Round
	// print round
	log.Printf("Round: %d R1: %d R2: %d R3: %d Connected %d", round, s.R1, s.R2, s.R3, s.Connected_players)
	if round+1 == 2 {
		for {
			if s.R1 >= s.Connected_players {
				break
			}
			// sleep 1 sec
			time.Sleep(1 * time.Second)
		}
	} else if round+1 == 3 {
		for {
			if s.R2 >= s.Connected_players {
				break
			}
			// sleep 1 sec
			time.Sleep(1 * time.Second)
		}
	} else if round+1 == 4 {
		for {
			if s.R3 >= s.Connected_players {
				break
			}
			// sleep 1 sec
			time.Sleep(1 * time.Second)
		}
	}


	// El lider selecciona un valor entre 6 y 10
	valor_lider := s.Randoms[round]

	req, err = stream.Recv()
	if err == io.EOF {
		return err
	}
	player_id := req.Player
	log.Printf("Mi valor %d es %d", round+1, valor_lider)

	check_error(err, "Error al recibir jugada")
	log.Printf("Jugador %d. Resp: %d", player_id, req.Play)

	s.Connection[req.Player-1].Jugada = int(req.Play)
	if s.Connection[req.Player-1].Active == true {

		var movimiento bool = comparar(int(req.Play), valor_lider)
		if movimiento {
			s.Connection[req.Player-1].Jugada += int(req.Play) 
			
			//save play in player data
			// notificar al player que sobrevivio
			resp := PlayerResponse{Type: 1, Response: 1}
			err := stream.Send(&resp)
			check_error(err, "Error al notificar al player que sobrevivio la ronda")
		} else {
			s.Connection[req.Player-1].Active = false
			s.Connected_players -= 1
			// notificar al player que murio
			resp := PlayerResponse{Type: 1, Response: 0}
			err := stream.Send(&resp)
			check_error(err, "Error al notificar al player que murio en la ronda")
			depositar(player_id, round+1)
		}
		s.Contestados += 1
	}

	if round == 0 {
		s.R1 += 1
	} else if round == 1 {
		s.R2 += 1
	} else if round == 2 {
		s.R3 += 1
	} else if round == 3 {
		s.R4 += 1
	}

	return nil
}
func (s *Server) SaveJugadasRonda1(stream PlayerService_SaveJugadasRonda1Server) error {
	log.Printf("Guardando jugadas")
	//receive player request
	req, err := stream.Recv()

	ip := get_env_var("IP_NAMENODE")
	conn_name, err := grpc.Dial(ip+":9003", grpc.WithInsecure())
	check_error(err, "Error al conectar con el servidor")
	defer conn_name.Close()

	c_name := name.NewNameServiceClient(conn_name)
	stream_name, err := c_name.Registrar(context.Background())
	check_error(err, "Error al crear el stream de nombres")

	req_name := &name.NameRequest{
		Player: int32(req.Player),
		Ronda: 1,
		Jugadas: req.Jugadas,
	}
	
	// send req to name server
	err = stream_name.Send(req_name)
	check_error(err, "Error al enviar el request")

	resp, err := stream_name.Recv()
	if check_error(err, "Error al recibir respuesta del servidor"){
		return nil
	}
	log.Printf("Recibido: %d", resp.Response)

	//send response to stream
	err = stream.Send(&PlayerResponse{Type: 1, Response: 1})
	check_error(err, "Error al enviar respuesta al player")
	return nil
}
func (s *Server) Fase1P2(stream PlayerService_Fase1P2Server) error {
	//Eliminar jugadores que no lograron el puntaje
	req, err := stream.Recv()
	player_id := req.Player
	check_error(err, "No se pudo determinar el numero de jugador")
	log.Printf("Jugador %d. Resp: %d. Active: %d", player_id, s.Connection[player_id-1].Jugada, s.Connection[player_id-1].Active)
	if s.Connection[player_id-1].Active == true {
		if req.Total < 21 {
			//print player Jugada
			log.Printf("Jugador %d. Resp: %d", player_id, s.Connection[player_id-1].Jugada)
			s.Connection[player_id-1].Active = false
			s.Connected_players -= 1
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

	if s.Change_fase == true {
		s.Change_fase = false
	}

	return nil
}
func (s *Server) Fase2(stream PlayerService_Fase2Server) error {

	//wait until change fase is true
	for {
		if s.Change_fase == true {
			break
		}
		time.Sleep(1 * time.Second)
	}
	log.Printf("VOY A RECIBIR 111")
	Check, err := stream.Recv()
	log.Printf("YA RECIBI")
	if s.Connection[Check.Player].Active {
		che := PlayerResponse{Type: 2, Response: int32(1)}
		err = stream.Send(&che)
		check_error(err, "Error al notificar al player que puede empezar el juego")
	} else {
		che := PlayerResponse{Type: 2, Response: int32(0)}
		err = stream.Send(&che)
		check_error(err, "Error al notificar al player que murio antes de iniciar el juego")
	}
	log.Printf("AUXILIO")
	valor_lider := s.Randoms[4]

	//ver que onda con los equipos y con conseguir los valores

	log.Printf("111")
	req, err := stream.Recv()
	log.Printf("222")
	player_id := req.Player
	s.JugadoresFase2 = append(s.JugadoresFase2, int(player_id))
	check_error(err, "aaaa")
	log.Printf("333")

	// wait until jugadoresfase2 have Connected_players
	for {
		if len(s.JugadoresFase2) == s.Connected_players {
			break
		}
		time.Sleep(1 * time.Second)
	}


	//sacar suma de resultados de cada team
	for i := 0; i < s.Jugadores2; i++ {
		if s.JugadoresFase2[i] == int(req.Player) && i < s.Jugadores2/2 {
			s.Team1 += int(req.Play)
		} else if s.JugadoresFase2[i] == int(req.Player) && i >= s.Jugadores2/2 {
			s.Team2 += int(req.Play)
		}
	}

	log.Printf("000")
	if valor_lider%2 != s.Team1%2 && valor_lider%2 != s.Team2%2 {
		//Matar equipo random

		perdedor := s.Randoms[7]

		if perdedor == 0 {
			//Matar equipo 1
			log.Printf("MATADO EL EQUIPO 1")

			for i := 0; i < s.Jugadores2; i++ {

				if s.JugadoresFase2[i] == int(req.Player) && i < s.Jugadores2/2 {

					s.Connection[req.Player-1].Active = false
					s.Connected_players -= 1
					// notificar al player que murio
					resp := PlayerResponse{Type: 1, Response: 0}
					err := stream.Send(&resp)
					check_error(err, "Error al notificar al player que murio en la ronda")

				}

			}
			log.Printf("444")

		} else {

			log.Printf("MATADO EL EQUIPO 2")

			for i := 0; i < s.Jugadores2; i++ {

				if s.JugadoresFase2[i] == int(req.Player) && i >= s.Jugadores2/2 {

					s.Connection[req.Player-1].Active = false
					s.Connected_players -= 1
					// notificar al player que murio
					resp := PlayerResponse{Type: 1, Response: 0}
					err := stream.Send(&resp)
					check_error(err, "Error al notificar al player que murio en la ronda")

				}

			}
		}
		
		log.Printf("555")
	} else if valor_lider%2 != s.Team1%2 && valor_lider%2 == s.Team2%2 {
		//Matar equipo 1
		log.Printf("MATADO EL EQUIPO 1")

		for i := 0; i < s.Jugadores2; i++ {

			if s.JugadoresFase2[i] == int(req.Player) && i < s.Jugadores2/2 {

				s.Connection[req.Player-1].Active = false
				s.Connected_players -= 1
				// notificar al player que murio
				resp := PlayerResponse{Type: 1, Response: 0}
				err := stream.Send(&resp)
				check_error(err, "Error al notificar al player que murio en la ronda")

			}

		}
		log.Printf("666")
	} else if valor_lider%2 == s.Team1%2 && valor_lider%2 != s.Team2%2 {
		//Matar equipo 2
		log.Printf("MATADO EL EQUIPO 2")

		for i := 0; i < s.Jugadores2; i++ {

			if s.JugadoresFase2[i] == int(req.Player) && i >= s.Jugadores2/2 {

				s.Connection[req.Player-1].Active = false
				s.Connected_players -= 1
				// notificar al player que murio
				resp := PlayerResponse{Type: 1, Response: 0}
				err := stream.Send(&resp)
				check_error(err, "Error al notificar al player que murio en la ronda")

			}

		}
		log.Printf("777")
	}

	resp := PlayerResponse{Type: 1, Response: 1}
	err = stream.Send(&resp)
	check_error(err, "Error al notificar al player que sobrevivio la ronda")
	log.Printf("888")
	return nil
}
func (s *Server) Fase3(stream PlayerService_Fase3Server) error {

	Check, err := stream.Recv()
	if s.Connection[Check.Player].Active {
		che := PlayerResponse{Type: 2, Response: int32(1)}
		err = stream.Send(&che)
	} else {
		che := PlayerResponse{Type: 2, Response: int32(0)}
		err = stream.Send(&che)
	}

	valor_lider := s.Randoms[5]

	//ver que onda con los equipos y con conseguir los valores

	req, err := stream.Recv()
	s.JugadoresFase3 = append(s.JugadoresFase2, int(req.Player))
	check_error(err, "aaaa de la fase 3")

	//POR HACER pedir valores y guardarlos en arreglo

	//ver diferencia de valores

	for i := 0; i < s.Jugadores3; i++ {
		if s.JugadoresFase3[i] == int(req.Player) {
			s.RespuestasFase3[i] = Abs(int(req.Play) - valor_lider)
		}
	}

	for i := 0; i < s.Jugadores3/2; i++ {
		if s.RespuestasFase3[2*i] == s.RespuestasFase3[(2*i)+1] {
			fmt.Println("AMBOS VIVEN")

			if s.JugadoresFase3[2*i] == int(req.Player) || s.JugadoresFase3[(2*i)+1] == int(req.Player) {

				// ambos viven
				resp := PlayerResponse{Type: 1, Response: 1}
				err = stream.Send(&resp)
				check_error(err, "Error al notificar al player que sobrevivio la ronda")

			}
		} else if s.RespuestasFase3[2*i] < s.RespuestasFase3[(2*i)+1] {

			if s.JugadoresFase3[2*i] == int(req.Player) {

				//gana participante 1
				resp := PlayerResponse{Type: 1, Response: 1}
				err = stream.Send(&resp)
				check_error(err, "Error al notificar al player que sobrevivio la ronda")

			}
			if s.JugadoresFase3[(2*i)+1] == int(req.Player) {
				// muere participante 2

				s.Connection[req.Player-1].Active = false
				s.Connected_players -= 1
				// notificar al player que murio
				resp := PlayerResponse{Type: 1, Response: 0}
				err := stream.Send(&resp)
				check_error(err, "Error al notificar al player que murio en la ronda")

			}

		} else {

			if s.JugadoresFase3[(2*i)+1] == int(req.Player) {

				// gana participante 2
				resp := PlayerResponse{Type: 1, Response: 1}
				err = stream.Send(&resp)
				check_error(err, "Error al notificar al player que sobrevivio la ronda")

			}
			if s.JugadoresFase3[2*i] == int(req.Player) {
				// muere participante 1

				s.Connection[req.Player-1].Active = false
				s.Connected_players -= 1
				// notificar al player que murio
				resp := PlayerResponse{Type: 1, Response: 0}
				err := stream.Send(&resp)
				check_error(err, "Error al notificar al player que murio en la ronda")

			}

		}

	}

	return nil
}
