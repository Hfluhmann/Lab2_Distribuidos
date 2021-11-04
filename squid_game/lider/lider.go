package lider

import (
	"log"
	"strconv"
	"os"
	"time"

	"fmt"

	"golang.org/x/net/context"
)

type Server struct {
}

func check_error(e error) {
    if e != nil {
        panic(e)
    }
}


// Send Play
func (s * Server) SendPlay(ctx context.Context, in *Interaction) (*Interaction, error) {
	log.Printf("Lider received action. Player: %d, Play: %s", in.Player, in.Play)
	return &Interaction{Response: strconv.FormatInt(in.Player, 10)}, nil
}

func (s * Server) GetConnection(ctx context.Context, in *Interaction) (*Interaction, error) {
	var con int64 = 0;
	file, err := os.Open("players.txt")
	check_error(err)
	_, err = fmt.Fscanf(file, "%d\n", &con)
	check_error(err)

	log.Printf("Jugador solicitando conexion. %d/16", con)
	if (con <= 16) {
		con += 1
		log.Printf("Jugador conectado. %d/16", con)
	} else{
		log.Printf("Rechazando conexion. 16/16")
	}
	
	err = os.Remove("players.txt")
	f, _ := os.Create("players.txt")
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%d\n", con))

	return &Interaction{Player: con}, nil
}

func (s * Server) GetPrizePool(ctx context.Context, in *Interaction) (*Interaction, error) {
	log.Printf("Lider received action. Player: %d, Play: %s", in.Player, in.Play)
	return &Interaction{Response: "Lider server"}, nil
}

func (s Server) SubscribePlayer(src PlayerService_SubscribePlayerServer) error {

	log.Println("Player connected")

	resp := PlayerResponse{Type: 0, Response: 1}
	if err := src.Send(&resp); err != nil {
		log.Printf("send error %v", err)
	}

	// print the received requests each 5 seconds
	for {
		req, err := src.Recv()
		if err != nil {
			log.Printf("receive error %v", err)
			return err
		}
		log.Printf("Type %d: Play: %d", req.Type, req.Play)
		// sleep 5 seconds
		time.Sleep(5 * time.Second)
	}


	// ctx := src.Context()

	// for {

	// 	// exit if context is done
	// 	// or continue
	// 	select {
	// 	case <-ctx.Done():
	// 		return ctx.Err()
	// 	default:
	// 	}

	// 	// receive data from stream
	// 	req, err := srv.Recv()
	// 	if err == io.EOF {
	// 		// return will close stream from server side
	// 		log.Println("exit")
	// 		return nil
	// 	}
	// 	if err != nil {
	// 		log.Printf("receive error %v", err)
	// 		continue
	// 	}

	// 	// continue if number reveived from stream
	// 	// less than max
	// 	if req.Type < 1 || req.Type > 4 {
	// 		continue
	// 	}

	// 	if req.Type == 1 {
	// 		resp := PlayerResponse{Response: "Moriste"}
	// 		if err := src.Send(&resp); err != nil {
	// 			log.Printf("send error %v", err)
	// 		}
	// 	}
		
	// 	log.Printf("El jugador Murió")
	// }

	return nil
}