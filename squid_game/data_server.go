package main

import (
	"fmt"
	"log"
	// "math/rand"
	"net"
	// "time"

	"Lab2_Distribuidos/squid_game/data"

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
	fmt.Println("------------- Iniciando Data Node ---------------")
	fmt.Println("-------------------------------------------------\n")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9004))
	check_error(err, "Error al escuchar en el puerto 9004")

	server := &data.Server{}

	grpcServer := grpc.NewServer()
	data.RegisterDataServiceServer(grpcServer, server)
	err = grpcServer.Serve(lis) // bind server

	check_error(err, "Error al iniciar el servidor de Datos")
}
