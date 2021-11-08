package main

import (
	"fmt"
	"io"
	"os"
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/joho/godotenv"

	"Lab2_Distribuidos/squid_game/lider"
	"Lab2_Distribuidos/squid_game/pozo"
)

func get_env_var(key string) string {

	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

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

func save_jugadas(jugadas []int32, player_id int32, ip string){ 
	conn, err := grpc.Dial(ip+":9000", grpc.WithInsecure())
	check_error(err, "Error al conectar con el servidor")
	defer conn.Close()

	c := lider.NewPlayerServiceClient(conn)
	stream_save, err := c.SaveJugadasRonda1(context.Background())
	check_error(err, "Error al crear el stream de guardado")

	//send jugadas to stream
	req := &lider.PlayerRequest{Type: 1, Player: player_id, Jugadas: jugadas} 
	err = stream_save.Send(req)

	//receive response from stream
	_, err = stream_save.Recv()
	check_error(err, "Error al recibir la respuesta")
}

func fase1(fase int32, player_id int32, ip string, human bool) bool {

	log.Println("\n---------------------------------------------")
	log.Printf("Bienvenido al primer juego, jugador %d\nEspera a que te demos la orden para comenzar...", player_id)

	conn, err := grpc.Dial(ip+":9000", grpc.WithInsecure())
	check_error(err, "Error al conectar con el servidor")
	defer conn.Close()

	c := lider.NewPlayerServiceClient(conn)
	total := 0
	var jugadas []int32
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
				fmt.Printf("\nIngrese un valor (1-10): ")
				fmt.Scanf("%d", &value)
			}
			value  = 6
			total += value
			jugadas = append(jugadas, int32(value))

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
					save_jugadas(jugadas, player_id, ip)
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
			save_jugadas(jugadas, player_id, ip)
			return false
		} else {
			log.Printf("Has sobrevivido al juego")
		}
	}

	save_jugadas(jugadas, player_id, ip)
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
		check_error(err, "Error al enviar los datos del jugador a fase 2")

		res, err := stream.Recv()
		check_error(err, "Error al recibir la confirmacion de inicio de fase 2")
		log.Printf("99999")

		if res.Response == 0 {
			log.Printf("He muerto porque sobraba un jugador para la segunda fase")
			return
		}
		log.Printf("111")
		s1 := rand.NewSource(time.Now().UnixNano() * int64(player_id))
		r1 := rand.New(s1)
		log.Printf("222")
		var value int = int(r1.Intn(4) + 1)
		log.Printf("333")
		log.Printf("Enviando valor: %d", value)
		// send player request to stream
		req := &lider.PlayerRequest{Type: 1, Player: player_id, Play: int32(value)}
		err = stream.Send(req)
		log.Printf("444")
		res, err = stream.Recv()
		log.Printf("555")
		check_error(err, "Error al recibir la respuesta de la jugada en fase 2")
		if res.Type == 1 {
			log.Printf("666")
			if res.Response == 0 {
				log.Printf("777")
				log.Printf("Has muerto R.I.P.")
				return
			} else if res.Response == 1 {
				log.Printf("888")
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
	ip := get_env_var("IP_LIDER")
	ip_pozo := get_env_var("IP_POZO")

	fase, player_id := conexion(ip)
	
	
	fmt.Println("\n---------------------------------------------")
	fmt.Println("1. Jugar\n2.Consultar pozo\n3.Salir")
	var opcion int
	fmt.Scanf("%d", &opcion)
	if opcion == 1 {
		fmt.Println("\n1. Jugador Bot\n2. Jugador Humano")
		fmt.Scanf("%d", &opcion)

		if opcion == 1 {
			fase1(fase, player_id, ip, false)
			// flag := fase1(fase, player_id, ip, false)
			// if flag {
			// 	fase2(fase, player_id, ip)
			// }
			// fase3(fase, player_id, ip)

		} else if opcion == 2 {
			fase1(fase, player_id, ip, true)
			// flag := fase1(fase, player_id, ip, true)
			// if flag {
			// 	fase2(fase, player_id, ip)
			// }
			// fase3(fase, player_id, ip)
		}
	} else if opcion == 2 {
		conn, err := grpc.Dial(ip_pozo+":9002", grpc.WithInsecure())
		check_error(err, "Error al conectar con el servidor del pozo")
		defer conn.Close()
	
		c := pozo.NewPozoServiceClient(conn)
		stream, err := c.Consultar(context.Background())
		check_error(err, "Error al crear el stream")
	
		resp, err := stream.Recv()
		if check_error(err, "Error al recibir respuesta del servidor"){
			return
		}
		log.Printf("Recibido: %d", resp.Response)
	} else if opcion == 3 {
		return
	}
	

}
