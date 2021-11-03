package lider

import (
	"log"
	"strconv"
	"os"

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
	//err = os.WriteFile("players.txt", []bytes(strconv.FormatInt(con, 10)), 0644)
	_, err = f.WriteString(fmt.Sprintf("%d\n", con))

	return &Interaction{Player: con}, nil
}

func (s * Server) GetPrizePool(ctx context.Context, in *Interaction) (*Interaction, error) {
	log.Printf("Lider received action. Player: %d, Play: %s", in.Player, in.Play)
	return &Interaction{Response: "Lider server"}, nil
}