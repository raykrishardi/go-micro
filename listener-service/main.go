package main

import (
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Panic(err)
	}
	defer rabbitConn.Close()
	log.Println("connected to rabbitmq")

	// start listening for messages (subscribe to the topic for event based consumption instead of polling)

	// create consumer

	// watch consumer and consume events/topics
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		// TODO: change to docker service name
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			log.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			log.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		// continue // optional
	}

	return connection, nil
}
