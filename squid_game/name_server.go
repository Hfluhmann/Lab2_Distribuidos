package main

import (
	"fmt"
	"log"
	// "math/rand"
	// "time"
	"net"
	"os"

	"Lab2_Distribuidos/squid_game/name"
	"github.com/joho/godotenv"
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

func get_env_var(key string) string {

	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

func main() {
	fmt.Println("-------------------------------------------------")
	fmt.Println("------------- Iniciando Name Node ---------------")
	fmt.Println("-------------------------------------------------\n")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9003))
	check_error(err, "Error al escuchar en el puerto 9003")

	ips := [3]string{get_env_var("IP_JUGADORES"), get_env_var("IP_POZO"), get_env_var("IP_POZO")}
	server := &name.Server{Ips: ips}

	grpcServer := grpc.NewServer()
	name.RegisterNameServiceServer(grpcServer, server)
	err = grpcServer.Serve(lis) // bind server

	check_error(err, "Error al iniciar el servidor de nombre")
}
