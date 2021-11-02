package main

import (
	"fmt"
	"os"
	"strconv"
)

var pozo int = 0

func morision(jugador int, ronda int) int {
	pozo += 100000000

	file, err := os.OpenFile("pozo.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {

		fmt.Println(err)

	} else {
		file.WriteString("Jugador_" + strconv.Itoa(jugador) + " Ronda_" + strconv.Itoa(ronda) + " " + strconv.Itoa(pozo))
		file.WriteString("\n")
	}

	file.Close()
	return pozo
}

func main() {

	for j := 7; j <= 9; j++ {
		morision(j, j*3)
	}

}
