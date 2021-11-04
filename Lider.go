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

	//Eliminar jugadores que no lograron el puntaje
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
func Jugar_F2(cant_jugadores int) {

	var sobrevivientes = make([]int, cant_jugadores) //ARREGLO TEMPORAL CON LOS EQUIPOS (1ra mitad: equipo 1/ 2da mitad equipo 2)
	for i := 0; i < cant_jugadores; i++ {
		sobrevivientes[i] += (i + 1)
	}

	//crear valor aleatorio del lider
	s2 := rand.NewSource(time.Now().UnixNano() * 100)
	r2 := rand.New(s2)
	valor_lider := r2.Intn(4) + 1

	//guardar numeros de cada particupante
	var teams = make([]int, cant_jugadores)
	for i := 0; i < cant_jugadores; i++ {
		teams[i] += i //POR HACER sumar valor selecionado por cada jugador
	}

	var team1 int = 0
	var team2 int = 0

	//sacar suma de resultados de cada team
	for i := 0; i < cant_jugadores; i++ {
		if i < cant_jugadores/2 {
			team1 += teams[i]
		} else {
			team2 += teams[i]
		}
	}

	//ver cual es el equipo perdedor y eliminarlo
	if valor_lider%2 != team1%2 && valor_lider%2 != team2%2 {
		//Matar equipo random
		s := rand.NewSource(time.Now().UnixNano() * 27)
		r := rand.New(s)
		perdedor := r.Intn(2)
		if perdedor == 0 {
			//Matar equipo 1
			fmt.Println("MATADO EL EQUIPO 1")
			for i := 0; i < cant_jugadores/2; i++ {
				//fmt.Println(sobrevivientes[i])
				kill_player(sobrevivientes[i])
			}
		} else {

			fmt.Println("MATADO EL EQUIPO 2")
			for i := cant_jugadores / 2; i < cant_jugadores; i++ {
				//fmt.Println(sobrevivientes[i])
				kill_player(sobrevivientes[i])
			}
		}

	} else if valor_lider%2 != team1%2 && valor_lider%2 == team2%2 {
		//Matar equipo 1
		fmt.Println("MATADO EL EQUIPO 1")
		for i := 0; i < cant_jugadores/2; i++ {
			//fmt.Println(sobrevivientes[i])
			kill_player(sobrevivientes[i])
		}

	} else if valor_lider%2 == team1%2 && valor_lider%2 != team2%2 {
		//Matar equipo 2
		fmt.Println("MATADO EL EQUIPO 2")
		for i := cant_jugadores / 2; i < cant_jugadores; i++ {
			//fmt.Println(sobrevivientes[i])
			kill_player(sobrevivientes[i])
		}
	}

	fmt.Println(valor_lider)
	fmt.Println(team1)
	fmt.Println(team2)

	return

}

//Analizis del Lider sobre el tercer juego
func Jugar_F3() int {

	s3 := rand.NewSource(time.Now().UnixNano() * 100)
	r3 := rand.New(s3)

	return r3.Intn(10) + 1

}

func main() {

	Jugar_F2(8)
	return
}
