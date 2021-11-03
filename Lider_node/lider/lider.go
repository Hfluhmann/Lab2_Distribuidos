package lider

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

// Send Play
func (s * Server) SendPlay(ctx context.Context, in *Interaction) (*Interaction, error) {
	log.Printf("Lider receibed action. Player: %d, Play: %s", in.Type, in.Play)
	return &Interaction{Response: "Lider server"}, nil
}