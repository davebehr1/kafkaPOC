package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/davebehr1/saramaExample/consumers"
	"github.com/davebehr1/saramaExample/models"
	"github.com/davebehr1/saramaExample/producers"
)

const topic = "sample-topic"

func main() {
	producer, err := producers.NewProducer()
	if err != nil {
		fmt.Println("Could not create producer: ", err)
	}

	consumer, err := sarama.NewConsumer(producers.Brokers, nil)
	if err != nil {
		fmt.Println("Could not create consumer: ", err)
	}

	consumers.Subscribe(topic, consumer)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "Hello Sarama!") })

	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		r.ParseForm()
		msg := producers.PrepareMessage(topic, r.FormValue("q"))
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Fprintf(w, "%s error occured.", err.Error())
		} else {
			fmt.Fprintf(w, "Message was saved to partion: %d.\nMessage offset is: %d.\n", partition, offset)
		}
	})

	http.HandleFunc("/retrieve", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, html.EscapeString(models.GetMessage())) })

	log.Fatal(http.ListenAndServe(":8081", nil))
}
