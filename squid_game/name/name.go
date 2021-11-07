package name

import (

	"io"
	"log"
	"os"
	// "math/rand"
	// "time"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"Lab2_Distribuidos/squid_game/data"
)


func check_error(e error, msg string) bool {
	if e != nil {
		log.Printf("%s", msg)
		log.Printf("Error: %v", e)
		return true
	}
	return false
}

type Server struct {
	Ips [3]string
}

func append_jugada(ip string, player int32, ronda int32){
	//open file
	f, err := os.OpenFile("name/jugadas.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if check_error(err, "Error opening file") {
		return
	}
	defer f.Close()

	//write to file
	_, err = f.WriteString(fmt.Sprintf("Jugador_%d Ronda_%d %s\n", int(player), ronda, ip))
	if check_error(err, "Error writing to file") {
		return
	}
}

func find_ip(player int32, ronda int32) string {
	//open file
	f, err := os.Open("name/jugadas.txt")
	if check_error(err, "Error opening file") {
		return ""
	}
	defer f.Close()

	//read from file
	var ip string
	var player_aux int32
	var ronda_aux int32
	for {
		_, err = fmt.Fscanf(f, "Jugador_%d Ronda_%d %s\n", &player_aux, &ronda_aux, &ip)
		if err == io.EOF {
			break
		}
		if check_error(err, "Error reading from file") {
			return ""
		}
		if player_aux == player && ronda_aux == ronda {
			return ip
		}
	}
	return ""
}

func (s *Server) Registrar(stream NameService_RegistrarServer) error {
	// get a random number between 0 and 2
	// s1 := rand.NewSource(time.Now().UnixNano() * 100)
	// r1 := rand.New(s1)
	i := 2//r1.Intn(3)
	
	
	//receive the request
	req, err := stream.Recv() // req tiene, jugador, ronda y array con las jugadas
	log.Printf("Registrando jugadas de Jugador %d", req.Player)
	if check_error(err, "Error receiving request") {
		return err
	}

	append_jugada(s.Ips[i], req.Player, req.Ronda)

	// connect to data server
	conn, err := grpc.Dial(s.Ips[i]+":9004", grpc.WithInsecure())
	if check_error(err, "Error connecting to data server") {
		return err
	}
	defer conn.Close()

	c := data.NewDataServiceClient(conn)
	stream_data, err := c.RegistrarJugadas(context.Background())
	if check_error(err, "Error connecting to data server") {
		return err
	}
	// send the request to the data server
	req_data := data.DataRequest{
		Player: req.Player,
		Ronda: req.Ronda,
		Jugadas: req.Jugadas,
	}
	err = stream_data.Send(&req_data)
	if check_error(err, "Error sending request to data server") {
		return err
	}
	// receive the response from the data server
	_, err = stream_data.Recv()
	if check_error(err, "Error receiving response from data server") {
		return err
	}

	// send response to stream
	err = stream.Send(&NameResponse{
		Type: 0,
		Response: 0,
	})
	if check_error(err, "Error sending response to lider stream") {
		return err
	}

	return nil
}

func (s *Server) ObtenerJugadas(stream NameService_ObtenerJugadasServer) error {
	req, err := stream.Recv() // req tiene, jugador, ronda y array con las jugadas
	log.Printf("Solicitando jugadas de Jugador %d", req.Player)
	if check_error(err, "Error receiving request") {
		return err
	}

	ip := find_ip(req.Player, req.Ronda)
	if ip == "" {
		return fmt.Errorf("Error finding ip")
	}

	// connect to data server
	conn, err := grpc.Dial(ip+":9004", grpc.WithInsecure())
	if check_error(err, "Error connecting to data server") {
		return err
	}
	defer conn.Close()

	c := data.NewDataServiceClient(conn)
	stream_data, err := c.ObtenerJugadas(context.Background())
	if check_error(err, "Error connecting to data server") {
		return err
	}
	// send the request to the data server
	req_data := data.DataRequest{
		Player: req.Player,
		Ronda: req.Ronda,
	}
	err = stream_data.Send(&req_data)
	if check_error(err, "Error sending request to data server") {
		return err
	}
	// receive the response from the data server
	resp_data, err := stream_data.Recv()
	if check_error(err, "Error receiving response from data server") {
		return err
	}

	// send response to stream
	err = stream.Send(&NameResponse{
		Type: 0,
		Jugadas: resp_data.Jugadas,
	})
	if check_error(err, "Error sending response to stream") {
		return err
	}
	
	return nil
}