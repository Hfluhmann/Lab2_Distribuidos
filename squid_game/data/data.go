package data

import (

	// "io"
	"log"
	"os"
	// "math/rand"
	// "time"
	"fmt"
	//"golang.org/x/net/context"
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
}

func write_jugadas(req * DataRequest) {
	player := req.Player
	ronda := req.Ronda
	jugadas := req.Jugadas
	filename := fmt.Sprintf("data/jugador_%d_ronda_%d.txt", player, ronda)

	//write one play per line in file filename
	file, err := os.Create(filename)
	if check_error(err, "Error creating file") {
		return
	}
	defer file.Close()

	for _, jugada := range jugadas {
		log.Printf("Jugada: %d", jugada)
		_, err = file.WriteString(fmt.Sprintf("%d\n",jugada))
		if check_error(err, "Error writing to file") {
			return
		}
	}
}

func read_jugadas(player int32, ronda int32) []int32 {
	filename := fmt.Sprintf("data/jugador_%d_ronda_%d.txt", player, ronda)
	file, err := os.Open(filename)
	if check_error(err, "Error opening file") {
		return nil
	}
	defer file.Close()

	var jugadas []int32
	for {
		var jugada int32
		_, err = fmt.Fscanln(file, &jugada)
		if err != nil {
			break
		}
		jugadas = append(jugadas, jugada)
	}
	return jugadas
}

func (s *Server) RegistrarJugadas(stream DataService_RegistrarJugadasServer) error {
	// receive req from stream
	req, err := stream.Recv()
	if check_error(err, "Error receiving request") {
		return err
	}
	log.Printf("Registrando jugadas de jugador %d", req.Player)

	//write jugadas to file
	write_jugadas(req)

	//send response
	err = stream.Send(&DataResponse{
		Type: 0,
		Response: 0,
	})
	
	return nil
}

func (s *Server) ObtenerJugadas(stream DataService_ObtenerJugadasServer) error {
	// receive req from stream
	req, err := stream.Recv()
	if check_error(err, "Error receiving request") {
		return err
	}
	log.Printf("Obteniendo jugadas de jugador %d", req.Player)

	//read jugadas from file
	jugadas := read_jugadas(req.Player, req.Ronda)

	//send response
	err = stream.Send(&DataResponse{
		Type: 1,
		Jugadas: jugadas,
	})

	return nil
}