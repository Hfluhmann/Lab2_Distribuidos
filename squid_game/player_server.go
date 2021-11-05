
package main

import (
	"log"
	"fmt"
	//"time"
	"io"

	"golang.org/x/net/context"
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

func print_options(flag bool){
	fmt.Println("\n----------------------------")
	if flag{
		fmt.Println("1. Enviar Jugada")
		fmt.Println("2. Ver Monto del pozo")
	} else {
		fmt.Println("1. Solicitar unirme")
	}
}

func main() {

	//var conn *grpc.ClientConn
	bool_connected := false
	var player_id int32;
	var fase int32 = 0
	for {
		if fase > 0 {
			log.Println("\n---------------------------------------------")
			log.Printf("Estamos en la fase: %d", fase)
		} else {
			log.Println("\n---------------------------------------------")
			log.Printf("Estamos en la fase de preparación")
		}
		print_options(bool_connected)
		// read a int from stdin
		fmt.Print("Ingresar opción: ")
		var option int
		fmt.Scanf("%d", &option)

		// conexion de los jugadores y comienzo de la primera fase
		if (option == 1 && !bool_connected){

			conn, err := grpc.Dial("172.17.0.4:9000", grpc.WithInsecure())
			check_error(err, "Error al conectar con el servidor")
			defer conn.Close()
		
			c := lider.NewPlayerServiceClient(conn)
			stream, err := c.PlayerHandler(context.Background())
			check_error(err, "Error al crear el stream")
			
			ctx := stream.Context()
			done := make(chan bool)
			go func() {
				<-ctx.Done()
				if err := ctx.Err(); err != nil {
					log.Println(err)
				}
				close(done)
			}()


			req := &lider.PlayerRequest{Type: 0}
			// send reto to stream
			err = stream.Send(req)
			check_error(err, "Error al enviar la solicitud")
			
			//receive response from stream
			res, err := stream.Recv()
			check_error(err, "Error al recibir la respuesta")
			
			if err == io.EOF {
				log.Println("No se aceptan más jugadores")
				return
			}

			if err == nil {
				if err := stream.CloseSend(); err != nil {
					log.Println(err)
				}

				// log response
				player_id = res.Player
				log.Printf("Conectado. Eres el jugador: %d", player_id)
				bool_connected = true

				stream, err := c.WaitingRoom(context.Background())
				check_error(err, "Error al conectar a la sala de espera")

				//send player data to waiting room
				req = &lider.PlayerRequest{Type: 0, Player: player_id}
				err = stream.Send(req)
				check_error(err, "Error al enviar datos a la sala de espera")

				for {
					res, err := stream.Recv()
					check_error(err, "Error al recibir la respuesta de la sala de espera")
					if res.Type == 0 && res.Response == 0 {
						log.Printf("Esperando jugadores")
					} else if res.Type == 0 && res.Response > 0 {
						log.Println("\n---------------------------------------------")
						log.Printf("Iniciando la fase: %d", res.Response)
						fase = res.Response
						break
					} else {
						log.Printf("Esperando cambio de fase")
					}
				}
			}
			

		} else if bool_connected && option == 1 && fase > 0 {
			// se esta jugar en la fase
			conn, err := grpc.Dial("172.17.0.4:9000", grpc.WithInsecure())
			check_error(err, "Error al conectar con el servidor")
			defer conn.Close()
			
			if fase == 1 {
				c := lider.NewPlayerServiceClient(conn)
				stream, err := c.Fase1(context.Background())
				if !check_error(err, "Error al crear el stream fase 1") {
					log.Println("\n---------------------------------------------")
					log.Printf("Pedir jugada")

					// send player request to stream
					req := &lider.PlayerRequest{Type: 1, Player: player_id, Play: 1}
					err = stream.Send(req)
					check_error(err, "Error al enviar la jugada")
				}
			// } else if fase == 2 {
			// 	c := lider.NewPlayerServiceClient(conn)
			// 	stream, err := c.Fase2(context.Background())
			// 	if !check_error(err, "Error al crear el stream fase 2") {
			// 		log.Println("\n---------------------------------------------")
			// 		log.Printf("Pedir jugada")
			// 	}
			// } else if fase == 3 {
			// 	c := lider.NewPlayerServiceClient(conn)
			// 	stream, err := c.Fase3(context.Background())
			// 	if !check_error(err, "Error al crear el stream fase 3") {
			// 		log.Println("\n---------------------------------------------")
			// 		log.Printf("Pedir jugada")
			// 	}
			}
		} else {
			log.Printf("No es posible pedir conexión.")
			//exit the program
			return
		}

	}
}