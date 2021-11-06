package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"Lab2_Distribuidos/squid_game/lider"

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

	var connections []*lider.Connection
	var playersData [16]*lider.Player
	var randoms []int
	s1 := rand.NewSource(time.Now().UnixNano() * 100)
	r1 := rand.New(s1)
	for i := 0; i < 9; i++ {

		if i < 4 { //numeros primer juego
			randoms = append(randoms, r1.Intn(5)+6)
		} else if i == 4 { //numero segunda ronda
			randoms = append(randoms, r1.Intn(4)+1)
		} else if i == 5 { //numero 3ra ronda
			randoms = append(randoms, r1.Intn(10)+1)
		} else if i == 6 { //si sobra un jugador para el segundo juego
			randoms = append(randoms, r1.Intn(16)+1)
		} else if i == 7 { //si hay que matar un equipo aleatorio en el 2do juego
			randoms = append(randoms, r1.Intn(2))
		} else if i == 8 { //si sobra un jugador para el tercer juego
			randoms = append(randoms, r1.Intn(16)+1)
		}

	}

	server := &lider.Server{
		Connection:        connections,
		Fase:              0,
		Max_players:       2,
		Connected_players: 0,
		Change_fase:       false,
		Players_data:      playersData,
		Randoms:           randoms,
		Change_round:      false,
		Round:             0,
		Contestados:       0,
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
				server.Change_fase = true
				log.Printf("Comenzando Fase %d", server.Fase)

				if server.Fase == 1 {
					for i := 0; i < 4; i++ {
						if server.Connected_players == 0 {
							log.Printf("No quedan jugadores\nHa terminado esta version de Squid Game")
							return
						}
						log.Printf("Ronda: %d", i)
						log.Printf("1. Comenzar ronda")
						fmt.Scanf("%d", &input)
						if input == 1 {
							server.Round = i
							server.Change_round = true
						}
						for server.Contestados < server.Connected_players {
							time.Sleep(2 * time.Second)
						}
						server.Change_round = false
						server.Contestados = 0
					}
				} else if server.Fase == 2 {
					continue
				} else if server.Fase == 3 {
					continue
				} else {
					log.Printf("Ha terminado esta version de Squid Game")
					return
				}

			} else if input == 2 {
				log.Println("\n---------------------------------------------")
				log.Println("Las jugadas del jugador")
			}
		}
	}

}
