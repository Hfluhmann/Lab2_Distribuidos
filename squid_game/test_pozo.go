package main

import (
	"fmt"
	// "io"
	"log"
	// "math/rand"
	// "time"

	"golang.org/x/net/context"
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

func depositar(player int, ronda int) {
	conn, err := amqp.Dial("amqp://client:1234@"+get_env_var("IP_POZO")+":5672/")
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

	body := fmt.Sprintf("%d %d", player, ronda)
	err = ch.Publish(
	"",     // exchange
	q.Name, // routing key
	false,  // mandatory
	false,  // immediate
	amqp.Publishing {
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	failOnError(err, "Failed to publish a message")

}

func main() {
	ip := get_env_var("IP_LIDER")

	//read a number from stdin
	var num int
	fmt.Printf("1. Consultar el pozo\n2.Depositar: ")
	fmt.Scanf("%d", &num)
	if num  == 1 {
		conn, err := grpc.Dial(ip+":9001", grpc.WithInsecure())
		check_error(err, "Error al conectar con el servidor del pozo")
		defer conn.Close()
	
		c := pozo.NewPozoServiceClient(conn)
		stream, err := c.Consultar(context.Background())
		check_error(err, "Error al crear el stream")
	
		resp, err := stream.Recv()
		if check_error(err, "Error al recibir respuesta del servidor"){
			return
		}
		log.Printf("Recibido: %d", resp.Response)
	} else {
		depositar(2, 3)
	}



}