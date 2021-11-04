package lider

import (
	"log"
	//"strconv"
	//"os"
	//"time"

	//"fmt"

	//"golang.org/x/net/context"
)

type Connection struct {
	stream PlayerService_PlayerHandlerServer
	id     int
	active bool
	error  chan error
}

type Server struct {
	Connection []*Connection
}

func check_error(e error, msg string) {
    if e != nil {
		log.Printf("%s", msg)
        panic(e)
    }
}


// Send Play
// func (s * Server) SendPlay(ctx context.Context, in *Interaction) (*Interaction, error) {
// 	log.Printf("Lider received action. Player: %d, Play: %s", in.Player, in.Play)
// 	return &Interaction{Response: strconv.FormatInt(in.Player, 10)}, nil
// }

// func (s * Server) GetConnection(ctx context.Context, in *Interaction) (*Interaction, error) {
// 	// var con int64 = 0;
// 	// file, err := os.Open("players.txt")
// 	// check_error(err)
// 	// _, err = fmt.Fscanf(file, "%d\n", &con)
// 	// check_error(err)

// 	// log.Printf("Jugador solicitando conexion. %d/16", con)
// 	// if (con <= 16) {
// 	// 	con += 1
// 	// 	log.Printf("Jugador conectado. %d/16", con)
// 	// } else{
// 	// 	log.Printf("Rechazando conexion. 16/16")
// 	// }
	
// 	// err = os.Remove("players.txt")
// 	// f, _ := os.Create("players.txt")
// 	// defer f.Close()
// 	// _, err = f.WriteString(fmt.Sprintf("%d\n", con))

// 	return &Interaction{Player: 1/*con*/}, nil
// }

// func (s * Server) GetPrizePool(ctx context.Context, in *Interaction) (*Interaction, error) {
// 	log.Printf("Lider received action. Player: %d, Play: %s", in.Player, in.Play)
// 	return &Interaction{Response: "Lider server"}, nil
// }

// func (s Server) SubscribePlayer(src PlayerService_SubscribePlayerServer) error {

// 	log.Println("Player connected")

// 	resp := PlayerResponse{Type: 0, Response: 1}
// 	if err := src.Send(&resp); err != nil {
// 		log.Printf("send error %v", err)
// 	}

// 	// print the received requests each 5 seconds
// 	for {
// 		req, err := src.Recv()
// 		if err != nil {
// 			log.Printf("receive error %v", err)
// 			return err
// 		}
// 		log.Printf("Type %d: Play: %d", req.Type, req.Play)
// 		// sleep 5 seconds
// 		time.Sleep(5 * time.Second)
// 	}

// 	return nil
// }

// func (s Server) CreateConnection(stream PlayerService_CreateConnectionClient) error {
	
// 	req, err := stream.Recv()
// 	if err != nil {
// 		log.Printf("receive error %v", err)
// 		return err
// 	}

// 	if req.Type == 0 && len(s.Connection) <= 16 {
// 		player_id := len(s.Connection)+1
// 		conn := &Connection{
// 			stream: stream,
// 			id: player_id,
// 			active: true,
// 			error:  make(chan error),
// 		}

// 		resp := PlayerRequest{Type: 0, Player: int32(player_id)}
// 		if err := stream.Send(&resp); err != nil {
// 			log.Printf("send error %v", err)
// 		}

// 		s.Connection = append(s.Connection, conn)
// 		return <-conn.error
// 	} else{
// 		log.Printf("Error al conectar al jugador")
// 		return nil
// 	}
// }

func (s *Server) PlayerHandler(stream PlayerService_PlayerHandlerServer) error {
	req, err := stream.Recv()
	check_error(err, "Error when reading from stream")
	
	if req.Type == 0 { // player trying to connect
		if len(s.Connection) <= 16 {
			player_id := len(s.Connection)+1
			conn := &Connection{
				stream: stream,
				id: player_id,
				active: true,
				error:  make(chan error),
			}

			resp := PlayerResponse{Type: 0, Player: int32(player_id)}
			if err := stream.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}

			s.Connection = append(s.Connection, conn)
			log.Printf("Players %v", s.Connection)

			log.Printf("Player connected. %d/16", player_id)
			return <-conn.error
		} else{
			log.Printf("Error al conectar al jugador")
			return nil
		}
	} else{
		return nil
	}
	return nil

}
