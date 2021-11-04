package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"Lab2_Distribuidos/squid_game/lider"
)

// func check_error(place string, err * Error){
// 	if err != nil {
// 			log.Fatalf("Error in %s: %v", err)
// 		}
// }

func main(){
	fmt.Println("Lider node")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	grpcServer := grpc.NewServer()
	lider.RegisterPlayerServiceServer(grpcServer, &lider.Server{})
	err = grpcServer.Serve(lis) // bind server

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

}