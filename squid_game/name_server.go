package main

import (
	"fmt"
	"log"
	// "math/rand"
	// "time"
	"net"

	"Lab2_Distribuidos/squid_game/name"

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

func main() {
	fmt.Println("-------------------------------------------------")
	fmt.Println("------------- Iniciando Name Node ---------------")
	fmt.Println("-------------------------------------------------\n")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9003))
	check_error(err, "Error al escuchar en el puerto 9003")

	ips := [3]string{"172.17.0.3", "172.17.0.4", "172.17.0.5"}
	server := &name.Server{Ips: ips}

	grpcServer := grpc.NewServer()
	name.RegisterNameServiceServer(grpcServer, server)
	err = grpcServer.Serve(lis) // bind server

	check_error(err, "Error al iniciar el servidor de nombre")
}
