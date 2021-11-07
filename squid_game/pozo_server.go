package main

import (

	"log"
	"fmt"
	"net"
	"io"
	"os"

	// "golang.org/x/net/context"
	"google.golang.org/grpc"
	"Lab2_Distribuidos/squid_game/pozo"
	"github.com/streadway/amqp" //rabbitmq
)

func check_error(e error, msg string) bool {
	if e != nil {
		log.Printf("%s", msg)
		log.Printf("Error: %v", e)
		return true
	}
	return false
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
}

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

  
func main() {

	fmt.Println("---------------------------------------------")
	fmt.Println("------------- Iniciando Pozo ---------------")
	fmt.Println("---------------------------------------------\n")
	
	go func() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9001))
	if check_error(err, "Error al iniciar el servidor") {
		return
	}
		server := &pozo.Server{}
		grpcServer := grpc.NewServer()
		pozo.RegisterPozoServiceServer(grpcServer, server)
		err = grpcServer.Serve(lis) // bind server
		check_error(err, "Error al iniciar el servidor de registro de jugadores")
	}()

	//------------------------------------------------------
	//----------------- RabbitMQ ---------------------------
	conn, err := amqp.Dial("amqp://client:1234@172.17.0.5:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// read player and round from string
			_, _, monto := read_pozo()
			log.Printf("monto: %d", monto)
			var player, round int
			fmt.Sscanf(string(d.Body), "%d %d", &player, &round)
			log.Printf("Player: %d Ronda: %d monto: %d", player, round, monto)
			line := fmt.Sprintf("\nJugador_%d Ronda_%d %d", player, round, monto+100000000)
			// append line file pozo/pozo.txt
			file, err := os.OpenFile("pozo/pozo.txt", os.O_APPEND|os.O_WRONLY, 0600)
			check_error(err, "Error al abrir el archivo")
			defer file.Close()
			_, err = file.WriteString(line)
			check_error(err, "Error al escribir en el archivo")
		
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
