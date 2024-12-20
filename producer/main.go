package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

type Temperature struct {
	TimeStamp int64 `json:"timestamp"`
	Degree    int   `json:"degree"`
}

func main() {
	log.Println("RabbitMQ producer running")

	conn, err := amqp.Dial("amqp://localhost:5672")
	Panic(err)
	defer conn.Close()

	channel, err := conn.Channel()
	Panic(err)
	defer channel.Close()

	queueName := "temperature"
	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	Panic(err)

	http.HandleFunc("/temperature", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		temperature := new(Temperature)

		if err := json.NewDecoder(r.Body).Decode(temperature); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		temperatureInBytes, err := json.Marshal(temperature)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		msg := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(temperatureInBytes),
		}

		err = channel.Publish("", queueName, false, false, msg)
		if err != nil {
			log.Println("Failed to publish message: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("temperature is tored in RabbitMQ channel temperature, %v\n", temperature)

		response := map[string]interface{}{
			"success": true,
			"msg":     fmt.Sprintf("Temperature for timecode %d is added successfully! ", temperature.TimeStamp),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, "Error posting temperature", http.StatusInternalServerError)
		}
	})

	log.Println("listening on 3030")
	log.Fatal(http.ListenAndServe(":3030", nil))
}
