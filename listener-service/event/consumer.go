package event

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.starlark.net/lib/json"
)

// This struct is used to RECEIVE events FROM QUEUE (i.e. listener service as consumer while broker service as producer)
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// Producer struct which is used to SEND TO the QUEUE
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Might listen to many topics
func (consumer *Consumer) Listen(topics []string) error {
	// Get reference to channel to get reference to queue
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Get reference to the queue using the channel
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// Loop through each topics, and bind exchange with queue and topic name
	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,
			s,
			"logs_topic", // IMPORTANT: Must match what's declared in the exchange
			false,        // no-wait
			nil,          // no args passed
		)

		if err != nil {
			return err
		}
	}

	// Get messages
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Consume messages forever until the program exits
	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	log.Printf("waiting for message on [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever // wait forever (BLOCKING)

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// logic for auth

	case "default":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	// Marshall/convert to json the payload
	jsonData, err := json.MarshalIndent(&entry, "", "\t")
	if err != nil {
		return err
	}

	// Call the logger microservice
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		// Something went wrong on the server side
		return err
	}

	return nil
}
