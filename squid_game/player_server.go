
package main

import (
	"log"
	"fmt"
	//"time"
	//"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"Lab2_Distribuidos/squid_game/lider"
)

func check_error(e error, msg string) {
    if e != nil {
		log.Printf("%s", msg)
        panic(e)
    }
}

func print_options(flag bool){
	fmt.Println("1. Solicitar unirme")
	if flag{
		fmt.Println("2. Enviar Jugada")
		fmt.Println("3. Ver Monto del pozo")
	}
}

func main() {

	//var conn *grpc.ClientConn
	bool_connected := false
	var player_id int32;
	for {
		print_options(false)
		// read a int from stdin
		fmt.Print("Ingresar opción: ")
		var option int
		fmt.Scanf("%d", &option)

		if (option == 1 && !bool_connected){
			conn, err := grpc.Dial("172.17.0.4:9000", grpc.WithInsecure())
			check_error(err, "Error al conectar con el servidor")
			defer conn.Close()
		
			c := lider.NewPlayerServiceClient(conn)

			stream, err := c.PlayerHandler(context.Background())
			check_error(err, "Error al crear el stream")

			//done := make(chan bool)
			
			//fase := 0;
			go func() {
				req := &lider.PlayerRequest{Type: 0}
				// send reto to stream
				err := stream.Send(req)
				check_error(err, "Error al enviar la solicitud")
				//receive response from stream
				res, err := stream.Recv()
				check_error(err, "Error al recibir la respuesta")

				if err == nil {
					// log response
					player_id = res.Player
					log.Printf("Conectado. Eres el jugador: %d", player_id)
				}
				
			}()

		} else {
			log.Printf("No es posible pedir conexión.")
			//exit the program
			return
		}
	}
}