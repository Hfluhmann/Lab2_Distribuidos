package main

import (
	"log"

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

func main() {

	conn, err := grpc.Dial("172.17.0.2:9000", grpc.WithInsecure())
	check_error(err, "Error al conectar con el servidor")
	defer conn.Close()

	stream, err := c.PlayerHandler(context.Background())

}
