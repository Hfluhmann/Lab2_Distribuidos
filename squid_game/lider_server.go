package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"Lab2_Distribuidos/squid_game/lider"
)

func check_error(e error, msg string) bool {
    if e != nil {
		log.Printf("%s", msg)
        log.Printf("Error: %v", e)
		return true
    }
	return false
}

func print_lider_options(s * lider.Server) {
	log.Println("\n---------------------------------------------")
	log.Printf("Estamos en la etapa: %d", s.Fase)
	log.Println("1. Iniciar siguiente fase")
	log.Println("2. Pedir jugadas de un jugador")
}

func main(){
	fmt.Println("---------------------------------------------")
	fmt.Println("------------- Iniciando Lider ---------------")
	fmt.Println("---------------------------------------------\n")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	var connections []*lider.Connection

	server := &lider.Server{connections, 0, 2, 0, false}
	go func(){
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
		if !server.Change_fase{
			//read int from stdin
			print_lider_options(server)
			var input int
			fmt.Scanf("%d", &input)
			if input == 1 {
				server.Fase++
				server.Change_fase = true
				log.Printf("Comenzando Fase %d", server.Fase)
			} else if input == 2 {
				log.Println("\n---------------------------------------------")
				log.Println("Las jugadas del jugador")
			}
		}
	}

}