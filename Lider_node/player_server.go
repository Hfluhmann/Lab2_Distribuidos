
package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"Lab2_Distribuidos/Lider_node/lider"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("172.17.0.4:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := lider.NewInteractionServiceClient(conn)

	response, err := c.SendPlay(context.Background()})
	if err != nil {
		log.Fatalf("Error when getting connection: %s", err)
	}
	log.Printf("Eres el jugador: %s", response.Response)
	player_id := response.Response

	// for{
	// 	response, err := c.SendPlay(context.Background(), &lider.Interaction{Play: "the big play"})
	
	// 	if err != nil {
	// 		log.Fatalf("Error when calling SayHello: %s", err)
	// 	}
	// 	log.Printf("Response from server: %s", response.Response)
	// }

}