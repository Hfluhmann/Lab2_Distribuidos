package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"time"
)

//Jugar primera fase para un bot
func Fase1(player_number int64) int {

	//generar seed cambiantes para los random
	s1 := rand.NewSource(time.Now().UnixNano() * player_number)
	r1 := rand.New(s1)

	return r1.Intn(10) + 1

}

//Jugar segunda fase para un bot
func Fase2(player_number int64) int {

	s2 := rand.NewSource(time.Now().UnixNano() * player_number)
	r2 := rand.New(s2)

	return r2.Intn(4) + 1

	l := list.New() // Initialize an empty list
	l.PushFront(10)
	fmt.Println(l.Front()) // &{{0x43e280 0x43e280 <nil> <nil>} 0}
}

//Jugar tercera fase para un bot
func Fase3(player_number int64) int {

	s3 := rand.NewSource(time.Now().UnixNano() * player_number)
	r3 := rand.New(s3)

	return r3.Intn(10) + 1

}

func delete() {

}
func main() {

	Fase2(1)

}
