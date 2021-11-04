package main

import (
	"fmt"
	"math/rand"
	"time"
)

func kill_player(player_number int) {
	//POR HACER le dice al jugador numero player_number que se murio y se debe desconectar
	//POR HACER le dice al pozo que el jugador numero player_number murio y que actualice el pozo y el txt
	return
}

//Compara los puntajes seleccionados de un participante con el lider en el primer juego, y decide si eliminarlo
func comparar(valor_jugador int, valor_lider int) bool {
	if valor_jugador < valor_lider {

		//El jugador sigue vivo
		return true

	} else {

		//El jugador debe ser eliminado
		return false
	}

}

//Analizis del Lider sobre el primer juego
func Jugar_F1() {
	//puntaje acumulado que lleva un jugador en esta ronda (solo suma si sigue vivo)
	puntajes := [16]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	//0 si el jugador de la posicion correspondiente esta vivo, 1 si esta eliminado
	eliminados := [16]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	//for que lleva cada ronda del primer juego
	for i := 0; i < 4; i++ {

		//generar seed cambiantes para los random
		s1 := rand.NewSource(time.Now().UnixNano() * 100)
		r1 := rand.New(s1)
		// El lider selecciona un valor entre 6 y 10
		valor_lider := r1.Intn(5) + 6

		fmt.Println(valor_lider)
		for j := 0; j < 16; j++ {
			if eliminados[j] == 0 {
				//POR HACER recibir numero de cada jugador
				var movimiento bool = comparar(7, valor_lider) //POR HACER se esta usando un valor generico 7 como valor del juagdor
				if movimiento {
					puntajes[j] += 7 //POR HACER se suma el puntaje que puso el jugador
				} else {
					kill_player(j + 1)
					eliminados[j] = 1
				}
			}
		}
	}

	for j := 0; j < 16; j++ {
		if eliminados[j] == 0 {
			if puntajes[j] < 21 {
				kill_player(j + 1)
				eliminados[j] = 1
			}
		}
	}

	fmt.Println(puntajes)
	fmt.Println(eliminados)
	return

}

//Analizis del Lider sobre el segundo juego
func Jugar_F2(player_number int64) int {

	s2 := rand.NewSource(time.Now().UnixNano() * player_number)
	r2 := rand.New(s2)

	return r2.Intn(4) + 1

}

//Jugar tercera fase para un bot
func Jugar_F3(player_number int64) int {

	s3 := rand.NewSource(time.Now().UnixNano() * player_number)
	r3 := rand.New(s3)

	return r3.Intn(10) + 1

}

func main() {

	Jugar_F1()
	return
}
