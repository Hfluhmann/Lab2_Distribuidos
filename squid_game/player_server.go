
package main

import (
	"log"
	"fmt"
	"time"
	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"Lab2_Distribuidos/squid_game/lider"
)

func check_error(e error) {
    if e != nil {
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
	var player_id int64;
	for {
		print_options(false)
		// read a int from stdin
		fmt.Print("Ingresar opción: ")
		var option int
		fmt.Scanf("%d", &option)

		if (option == 1 && !bool_connected){
			conn, err := grpc.Dial("172.17.0.4:9000", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %s", err)
			}
			defer conn.Close()
		
			c := lider.NewPlayerServiceClient(conn)
			response, err := c.GetConnection(context.Background(), &lider.Interaction{})
			if err != nil {
				log.Fatalf("Error when getting connection: %s", err)
			} else{
				player_id = response.Player
				if player_id == -1 {
					fmt.Println("No se pudo conectar")
				} else{
					log.Printf("Eres el jugador: %d", player_id)
					bool_connected = true
					
					stream, err := c.SubscribePlayer(context.Background())
					if err != nil {
						log.Fatalf("Error when creating stream: %s", err)
					}
				
					//ctx := stream.Context()
					done := make(chan bool)
				
					fase := 0;
					//play := -1;
				
					// listen responses from Lider
					go func() {
						for {
							res, err := stream.Recv()
							if err == io.EOF {
								close(done)
								return
							}
							if err != nil {
								log.Fatalf("can not receive %v", err)
							}
				
							if res.Type == 0 {
								log.Printf("Fase Changed: %d", res.Response)
								fase = int(res.Response)

								for{
									player_req := lider.PlayerRequest{Type: 1, Play: 2}
									if err := stream.Send(&player_req); err != nil {
										log.Fatalf("can not send %v", err)
									}
									//sleep 5 seconds
									time.Sleep(5 * time.Second)
								}
							}
						}
					}()



				}
			}
		} else {
			log.Printf("No es posible pedir conexión. Eres el jugador: %d", player_id)
			//exit the program
			return
		}
	}

	// create bidirectional stream connection
	// stream, err := c.SubscribePlayer(context.Background())
	// if err != nil {
	// 	log.Fatalf("Error when creating stream: %s", err)
	// }

	// //ctx := stream.Context()
	// done := make(chan bool)

	// fase := 0;
	// //play := -1;

	// // listen responses from Lider
	// go func() {
	// 	for {
	// 		res, err := stream.Recv()
	// 		if err == io.EOF {
	// 			close(done)
	// 			return
	// 		}
	// 		if err != nil {
	// 			log.Fatalf("can not receive %v", err)
	// 		}

	// 		if res.Type == 0 {
	// 			fase = int(res.Response)
	// 		}
	// 	}
	// }()

	// // send request to Lider
	// go func() {
	// 	if fase == 0 {
	// 		// sleep 2 seconds
	// 		log.Printf("Game dont start yet")
	// 		time.Sleep(3 * time.Second)
	// 	} else if fase == 1 {
	// 		log.Printf("Fase 1")
	// 	} else if fase == 2 {
	// 		log.Printf("Fase 2")
	// 	} else if fase == 3 {
	// 		log.Printf("Fase 3")
	// 	} else {
	// 		log.Printf("Error")
	// 	}
	// 	// for i := 1; i <= 10; i++ {
	// 	// 	// generates random number and sends it to stream
	// 	// 	rnd := int32(rand.Intn(i))
	// 	// 	req := pb.Request{Num: rnd}
	// 	// 	if err := stream.Send(&req); err != nil {
	// 	// 		log.Fatalf("can not send %v", err)
	// 	// 	}
	// 	// 	log.Printf("%d sent", req.Num)
	// 	// 	time.Sleep(time.Millisecond * 200)
	// 	// }
	// 	// if err := stream.CloseSend(); err != nil {
	// 	// 	log.Println(err)
	// 	// }
	// }()

}