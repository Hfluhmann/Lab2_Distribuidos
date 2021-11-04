package main

import (
	"fmt"
	"math/rand"
	"time"
)

func kill_player(player_number int) {
	//le dice al jugador numero player_number que se murio y se debe desconectar
	//le dice al pozo que el jugador numero player_number murio y que actualice el pozo y el txt
	return
}

//Jugar primera fase para un bot
func comparar(valor_jugador int, valor_lider int) bool {
	if valor_jugador < valor_lider {

		//El jugador sigue vivo
		return true

	} else {

		//El jugador debe ser eliminado
		return false
	}

}

func final_F1() {
	var a [2]int

	fmt.Println(a)
	fmt.Println(a[0])
	fmt.Println(a[1])
	return
}
func Jugar_F1() {
	puntajes := [16]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	eliminados := [16]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	//generar seed cambiantes para los random
	for i := 0; i < 4; i++ {

		s1 := rand.NewSource(time.Now().UnixNano() * 100)
		r1 := rand.New(s1)
		valor_lider := r1.Intn(5) + 6

		fmt.Println(valor_lider)
		for j := 0; j < 16; j++ {
			if eliminados[j] == 0 {
				//recibir numero de cada jugador
				var movimiento bool = comparar(7, valor_lider) //se esta usando un valor generico 7 como valor del juagdor
				if movimiento {
					puntajes[j] += 7 //se suma el puntaje que puso el jugador
				} else {
					kill_player(j)
					eliminados[j] = 1
				}
			}
		}
	}
	fmt.Println(puntajes)
	fmt.Println(eliminados)
	return

}

//Jugar segunda fase para un bot
func Fase2(player_number int64) int {

	s2 := rand.NewSource(time.Now().UnixNano() * player_number)
	r2 := rand.New(s2)

	return r2.Intn(4) + 1

}

//Jugar tercera fase para un bot
func Fase3(player_number int64) int {

	s3 := rand.NewSource(time.Now().UnixNano() * player_number)
	r3 := rand.New(s3)

	return r3.Intn(10) + 1

}

func main() {

	Jugar_F1()
	return
}
