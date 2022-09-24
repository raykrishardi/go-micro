package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name of the exchange
		"topic",      // type
		true,         //durable?
		false,        //auto delete when unused?
		false,        //used only internally (noLocal)? false because using between microservices
		false,        // no-wait, trevor said it's not important so just put false
		nil,          // no args for this
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name of the queue, if empty then random name
		false, // durable?
		false, // auto delete when unused?
		true,  // exclusive channel for current operations, don't share it around
		false, // no-wait?
		nil,   // no arguments
	)
}
