package pozo

import (

	"io"
	"log"
	"os"
	// "time"
	"fmt"
	//"golang.org/x/net/context"
)

func read_pozo() (player int, ronda int, monto int) {
	log.Printf("Reading pozo")
	file, err := os.Open("pozo/pozo.txt")
	check_error(err, "Error al abrir el archivo")
	defer file.Close()
  
	var p int
	var r int
	var m int

	_, err = fmt.Fscanf(file, "Jugador_%d Ronda_%d %d\n", &p, &r, &m)
	for {
		player = p
		ronda = r
		monto = m
		_, err := fmt.Fscanf(file, "Jugador_%d Ronda_%d %d\n", &p, &r, &m)
		if err == io.EOF || err != nil {
			break
		}
	}
	return player, ronda, monto
}

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

func (s *Server) Consultar(stream PozoService_ConsultarServer) error {
	_, _, monto := read_pozo()
	//send monto to stream
	err := stream.Send(&PozoResponse{
		Type: 1,
		Response: int32(monto),
	})
	check_error(err, "Error al enviar el monto")
	return nil
}