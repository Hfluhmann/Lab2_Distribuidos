package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"Lab2_Distribuidos/squid_game/lider"
	"Lab2_Distribuidos/squid_game/name"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func check_error(e error, msg string) bool {
	if e != nil {
		log.Printf("%s", msg)
		log.Printf("Error: %v", e)
		return true
	}
	return false
}

func print_lider_options(s *lider.Server) {
	log.Println("\n---------------------------------------------")
	log.Printf("Estamos en la etapa: %d", s.Fase)
	log.Println("1. Iniciar siguiente fase")
	log.Println("2. Pedir jugadas de un jugador")
}

func main() {
	fmt.Println("---------------------------------------------")
	fmt.Println("------------- Iniciando Lider ---------------")
	fmt.Println("---------------------------------------------\n")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// ip del name node
	var ip_name string = "172.17.0.5"
	max := 3
	var connections []*lider.Connection
	var randoms []int
	s1 := rand.NewSource(time.Now().UnixNano() * 100)
	r1 := rand.New(s1)
	for i := 0; i < 9; i++ {

		if i < 4 { //numeros primer juego
			randoms = append(randoms, r1.Intn(5)+11)
		} else if i == 4 { //numero segunda ronda
			randoms = append(randoms, r1.Intn(4)+1)
		} else if i == 5 { //numero 3ra ronda
			randoms = append(randoms, r1.Intn(10)+1)
		} else if i == 6 { //si sobra un jugador para el segundo juego
			randoms = append(randoms, r1.Intn(max)+1)
		} else if i == 7 { //si hay que matar un equipo aleatorio en el 2do juego
			randoms = append(randoms, r1.Intn(2))
		} else if i == 8 { //si sobra un jugador para el tercer juego
			randoms = append(randoms, r1.Intn(max)+1)
		}

	}

	var jugadoresFase2 []int
	var jugadoresFase3 []int
	var respuestasFase3 [16]int

	server := &lider.Server{
		Connection:        connections,
		Fase:              0,
		Max_players:       max,
		Connected_players: 0,
		Change_fase:       false,
		Team1:             0,
		Team2:             0,
		Jugadores2:        0,
		Jugadores3:        0,
		Randoms:           randoms,
		Change_round:      false,
		Round:             0,
		Contestados:       0,
		JugadoresFase2:    jugadoresFase2,
		JugadoresFase3:    jugadoresFase3,
		RespuestasFase3:   respuestasFase3,
		R1:				   0,
		R2:				   0,
		R3:				   0,
		R4:				   0,
		Writed:			   false,
	}

	go func() {
		grpcServer := grpc.NewServer()
		lider.RegisterPlayerServiceServer(grpcServer, server)
		err = grpcServer.Serve(lis) // bind server

		check_error(err, "Error al iniciar el servidor de registro de jugadores")
	}()

	// consola del lider
	for server.Connected_players < server.Max_players {
		log.Printf("Esperando Jugadores")
		//sleep for 2 seconds
		time.Sleep(5 * time.Second)
	}

	for {
		if !server.Change_fase {
			//read int from stdin
			print_lider_options(server)
			var input int
			fmt.Scanf("%d", &input)
			if input == 1 {
				server.Fase++

				if server.Fase == 2 {
					log.Printf("-------------------------------  %d %d", server.Randoms[6], len(server.Connection))
					for server.Connected_players%2 != 0 {
						if server.Connection[server.Randoms[6]-1].Active == true {
							// matar jugador de la posicion
							server.Connection[server.Randoms[6]-1].Active = false
							server.Connected_players -= 1
						} else {
							if server.Randoms[6] == max {
								server.Randoms[6] = 1
							} else {
								server.Randoms[6] += 1
							}
						}
					}
				}
				if server.Fase == 3 {
					
					if server.Connection[server.Randoms[8]-1].Active == true {
						// matar jugador de la posicion
						server.Connection[server.Randoms[8]-1].Active = false
						server.Connected_players -= 1
					} else {
						if server.Randoms[8] == 16 {
							server.Randoms[8] = 1
						} else {
							server.Randoms[8] += 1
						}
					}
				}

				server.Change_fase = true
				log.Printf("Comenzando Fase %d", server.Fase)

			} else if input == 2 {
				conn, err := grpc.Dial(ip_name+":9003", grpc.WithInsecure())
				check_error(err, "Error al conectar con el servidor del pozo")
				defer conn.Close()

				c := name.NewNameServiceClient(conn)
				stream, err := c.ObtenerJugadas(context.Background())

				req := &name.NameRequest{
					Player: 3,
					Ronda: 12,
				}
				err = stream.Send(req)
				check_error(err, "Error al enviar el request")

				resp, err := stream.Recv()
				if check_error(err, "Error al recibir respuesta del servidor"){
					return
				}
				// print all the jugadas
				for _, jugada := range resp.Jugadas {
					log.Printf("Jugada: %d", jugada)
				}
			}
		}
	}

}
