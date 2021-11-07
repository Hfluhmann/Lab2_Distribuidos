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
func conexion(ip string) (fase int32, player_id int32) {

	var option int = 0
	for option != 1 && option != 2 {
		log.Printf("\n1. Solitar conexion\n2. Salir")
		fmt.Scanf("%d", &option)
	}
	if option == 2 {
		return
	}

	// Conexion y espera de jugadores
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
	return fase, player_id
}

func fase1(fase int32, player_id int32, ip string, human bool) bool {

	log.Println("\n---------------------------------------------")
	log.Printf("Bienvenido al primer juego, jugador %d\nEspera a que te demos la orden para comenzar...", player_id)

	conn, err := grpc.Dial(ip+":9000", grpc.WithInsecure())
	check_error(err, "Error al conectar con el servidor")
	defer conn.Close()

	c := lider.NewPlayerServiceClient(conn)
	total := 0
	for i := 0; i < 4; i++ {
		stream, err := c.Fase1P1(context.Background())
		if !check_error(err, "Error al crear el stream fase 1") {
			// block until we get a response equal to 0

			// send req to server
			req := &lider.PlayerRequest{Round: int32(i)}
			err = stream.Send(req)
			check_error(err, "Error al enviar la solicitud")

			check_error(err, "Error al recibir inicio de ronda")
			log.Println("\n---------------------------------------------")
			log.Printf("Iniciando Ronda %d...", i+1)

			var value int
			if !human {
				s1 := rand.NewSource(time.Now().UnixNano() * int64(player_id))
				r1 := rand.New(s1)
				value = int(r1.Intn(10) + 1)
			} else {
				fmt.Printf("\nIngrese un valor (1-5): ")
				fmt.Scanf("%d", &value)
			}

			total += value

			log.Printf("Enviando valor: %d", value)
			// send player request to stream
			req = &lider.PlayerRequest{Type: 1, Player: player_id, Play: int32(value)}
			err = stream.Send(req)
			check_error(err, "Error al enviar la jugada")

			res, err := stream.Recv()
			check_error(err, "Error al recibir la respuesta de la jugada en fase 1")
			if res.Type == 1 {
				if res.Response == 0 {
					log.Printf("Has muerto R.I.P.")
					return false
				} else if res.Response == 1 {
					log.Printf("Has sobrevivido a la ronda")
				}
			}
		}
	}
	log.Println("\n---------------------------------------------")
	stream_f2, err := c.Fase1P2(context.Background())
	req := &lider.PlayerRequest{Type: 1, Player: player_id, Total: int32(total)}
	err = stream_f2.Send(req)
	check_error(err, "Error al consultar resultado del juego")

	res, err := stream_f2.Recv()
	check_error(err, "Error al recibir el resultado del juego")
	if res.Type == 1 {
		if res.Response == 0 {
			log.Printf("Has muerto R.I.P.")
		} else {
			log.Printf("Has sobrevivido al juego")
		}
	}
	return true
}

func fase2(fase int32, player_id int32, ip string) {
	//POR HACER checkear que no fue eliminado por sobrar al comenzar la ronda

	log.Println("\n---------------------------------------------")
	log.Printf("Bienvenido al segundo juego, jugador %d\nEspera a que te demos la orden para comenzar...", player_id)

	conn, err := grpc.Dial(ip+":9000", grpc.WithInsecure())
	check_error(err, "Error al conectar con el servidor")
	defer conn.Close()

	c := lider.NewPlayerServiceClient(conn)
	stream, err := c.Fase2(context.Background())
	if !check_error(err, "Error al crear el stream fase 1") {

		check := &lider.PlayerRequest{Type: 2, Player: player_id}
		err = stream.Send(check)

		res, err := stream.Recv()

		check_error(err, "")

		if res.Response == 0 {
			return
		}

		s1 := rand.NewSource(time.Now().UnixNano() * int64(player_id))
		r1 := rand.New(s1)
		var value int = int(r1.Intn(4) + 1)

		log.Printf("Enviando valor: %d", value)
		// send player request to stream
		req := &lider.PlayerRequest{Type: 1, Player: player_id, Play: int32(value)}
		err = stream.Send(req)

		res, err = stream.Recv()

		check_error(err, "Error al recibir la respuesta de la jugada en fase 2")
		if res.Type == 1 {
			if res.Response == 0 {
				log.Printf("Has muerto R.I.P.")
				return
			} else if res.Response == 1 {
				log.Printf("Has sobrevivido a la ronda")
			}
		}

	}

	return
}

func fase3(fase int32, player_id int32, ip string) {
	//POR HACER checkear que no fue eliminado por sobrar al comenzar la ronda

	log.Println("\n---------------------------------------------")
	log.Printf("Bienvenido al tercer juego, jugador %d\nEspera a que te demos la orden para comenzar...", player_id)

	conn, err := grpc.Dial(ip+":9000", grpc.WithInsecure())
	check_error(err, "Error al conectar con el servidor")
	defer conn.Close()

	c := lider.NewPlayerServiceClient(conn)
	stream, err := c.Fase3(context.Background())
	if !check_error(err, "Error al crear el stream fase 1") {

		check := &lider.PlayerRequest{Type: 2, Player: player_id}
		err = stream.Send(check)

		res, err := stream.Recv()

		check_error(err, "")

		if res.Response == 0 {
			return
		}

		s1 := rand.NewSource(time.Now().UnixNano() * int64(player_id))
		r1 := rand.New(s1)
		var value int = int(r1.Intn(10) + 1)

		log.Printf("Enviando valor: %d", value)
		// send player request to stream
		req := &lider.PlayerRequest{Type: 1, Player: player_id, Play: int32(value)}
		err = stream.Send(req)

		res, err = stream.Recv()

		check_error(err, "Error al recibir la respuesta de la jugada en fase 3")
		if res.Type == 1 {
			if res.Response == 0 {
				log.Printf("Has muerto R.I.P.")
				return
			} else if res.Response == 1 {
				log.Printf("Has sobrevivido a la ronda y ganado el juego")
			}
		}

	}

	return
}

func main() {
	ip := "172.17.0.4"

	fase, player_id := conexion(ip)
	
	
	// read in from stdin
	var option int
	fmt.Println("\n1. Jugador Bot\n2. Jugador Humano")
	fmt.Scanf("%d", &option)

	if option == 1 {
		flag := fase1(fase, player_id, ip, false)
		if flag {
			fase2(fase, player_id, ip)
		}
		// fase3(fase, player_id, ip)

	} else if option == 2 {
		flag := fase1(fase, player_id, ip, true)
		if flag {
			fase2(fase, player_id, ip)
		}
		// fase3(fase, player_id, ip)
	}

}
