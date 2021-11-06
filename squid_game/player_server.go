package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

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

func print_options(flag bool) {
	fmt.Println("\n----------------------------")
	if flag {
		fmt.Println("1. Enviar Jugada")
		fmt.Println("2. Ver Monto del pozo")
	} else {
		fmt.Println("1. Solicitar unirme")
	}
}

func main() {
	ip := "172.17.0.2"
	//var conn *grpc.ClientConn
	bool_connected := false
	var player_id int32
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
		if option == 1 && !bool_connected {

			conn, err := grpc.Dial(ip+":9000", grpc.WithInsecure())
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
			conn, err := grpc.Dial(ip+":9000", grpc.WithInsecure())
			check_error(err, "Error al conectar con el servidor")
			defer conn.Close()

			if fase == 1 {
				for i := 0; i < 4; i++ {
					c := lider.NewPlayerServiceClient(conn)
					stream, err := c.Fase1P1(context.Background())
					if !check_error(err, "Error al crear el stream fase 1") {
						res, err := stream.Recv()
						for i != int(res.Round) {
							time.Sleep(2 * time.Second)
						}
						check_error(err, "Error al recibir inicio de ronda")
						log.Println("\n---------------------------------------------")
						log.Printf("Pedir jugada")

						s1 := rand.NewSource(time.Now().UnixNano() * int64(player_id))
						r1 := rand.New(s1)
						var value int = int(r1.Intn(5) + 1)
						// send player request to stream
						req := &lider.PlayerRequest{Type: 1, Player: player_id, Play: int32(value)}
						err = stream.Send(req)
						check_error(err, "Error al enviar la jugada")

						res, err = stream.Recv()
						check_error(err, "Error al recibir la respuesta de la jugada en fase 1")
						if res.Type == 1 {
							if res.Response == 0 {
								log.Printf("Has muerto R.I.P.")
								return
							} else if res.Response == 1 {
								log.Printf("Has sobrevivido a la ronda")
							}
						}
					}
				}
				c := lider.NewPlayerServiceClient(conn)
				stream, err := c.Fase1P2(context.Background())
				req := &lider.PlayerRequest{Type: 1, Player: player_id}
				err = stream.Send(req)
				check_error(err, "Error al consultar resultado del juego")

				res, err := stream.Recv()
				check_error(err, "Error al recibir el resultado del juego")
				if res.Type == 1 {
					if res.Response == 0 {
						log.Printf("Has muerto R.I.P.")
					} else {
						log.Printf("Has sobrevivido al juego")
					}
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
