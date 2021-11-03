
package main

import (
	"log"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"Lab2_Distribuidos/Lider_node/lider"
)

func check_error(e error) {
    if e != nil {
        panic(e)
    }
}

func print_options(){
	fmt.Println("1. Solicitar unirme")
	fmt.Println("2. Enviar Jugada")
	fmt.Println("3. Ver Monto del pozo")

}

func main() {

	//var conn *grpc.ClientConn
	bool_connected := false
	var player_id int64;
	for {
		print_options()
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
				}
			}
		} else {
			log.Printf("No es posible pedir conexión. Eres el jugador: %d", player_id)
		}

	}

	// for{
	// 	response, err := c.SendPlay(context.Background(), &lider.Interaction{Play: "the big play"})
	
	// 	if err != nil {
	// 		log.Fatalf("Error when calling SayHello: %s", err)
	// 	}
	// 	log.Printf("Response from server: %s", response.Response)
	// }

}