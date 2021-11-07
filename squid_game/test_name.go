package main

import (
	"fmt"
	// "io"
	"log"
	// "math/rand"
	// "time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"Lab2_Distribuidos/squid_game/name"
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
	ip := "172.17.0.4"
	
	//read option from stdin
	var option int
	log.Printf("\n1. Registrar jugadas\n2.Pedir Jugadas")
	fmt.Scanf("%d", &option)

	if option == 1 {
		//register

		conn, err := grpc.Dial(ip+":9003", grpc.WithInsecure())
		check_error(err, "Error al conectar con el servidor del pozo")
		defer conn.Close()

		c := name.NewNameServiceClient(conn)
		stream, err := c.Registrar(context.Background())
		check_error(err, "Error al crear el stream")

		// send req to name server
		jugadas := []int32{11,12,13,14,15}
		req := &name.NameRequest{
			Player: 3,
			Ronda: 12,
			Jugadas: jugadas,
		}

		// send req to name server
		err = stream.Send(req)
		check_error(err, "Error al enviar el request")

		resp, err := stream.Recv()
		if check_error(err, "Error al recibir respuesta del servidor"){
			return
		}
		log.Printf("Recibido: %d", resp.Response)
	} else if option == 2 {
		conn, err := grpc.Dial(ip+":9003", grpc.WithInsecure())
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