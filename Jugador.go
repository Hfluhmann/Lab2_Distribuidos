package main

import (
	"fmt"
	"math/rand"
	"time"
)

//Jugar primera fase para un bot
func Fase1() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var suma int = 0

	for j := 0; j < 4; j++ {

		//generar seed cambiantes para los random

		//Generar numero entre el (0 y el 9) + 1
		var num_jefe int = 7 //numero generico para considerar que elige el Jefe
		var actual int = r1.Intn(10) + 1

		fmt.Println(suma, actual)

		if actual < num_jefe {
			// SOBREVIVE LA ELECCION
			suma += actual
		} else {
			// NO SOBREVIVE
		}

		if suma < 21 {
			// NO SOBREVIVE
		}

	}

}

func main() {

	Fase1()

}
