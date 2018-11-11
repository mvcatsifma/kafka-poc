package main

import (
	"bytes"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	consumer1, err := consumer.ConsumePartition("Products", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer1.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumedProducts := 0

	wg.Add(1)
	go func() {
	ConsumeProductsLoop:
		for {
			select {
			case msg := <-consumer1.Messages():
				result := &Product{}
				json.NewDecoder(bytes.NewReader(msg.Value)).Decode(result)
				log.Printf("Consumed message offset %d\n", msg.Offset)
				log.Printf("Msg %+v\n", result)
				consumedProducts++
			case <-signals:
				break ConsumeProductsLoop
			}
		}
	}()

	consumer2, err := consumer.ConsumePartition("Status", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer2.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	consumedStatus := 0
	wg.Add(1)
	go func() {
	ConsumeStatusLoop:
		for {
			select {
			case msg := <-consumer2.Messages():
				result := &Status{}
				json.NewDecoder(bytes.NewReader(msg.Value)).Decode(result)
				log.Printf("Consumed message offset %d\n", msg.Offset)
				log.Printf("Msg %+v\n", result)
				consumedStatus++
			case <-signals:
				break ConsumeStatusLoop
			}
		}
	}()

	wg.Wait()

	log.Printf("Consumed products: %d\n", consumedProducts)
	log.Printf("Consumed status: %d\n", consumedStatus)
}

type Status struct {
	UUID      uuid.UUID
	MachineId string
	Status    string
	Since     time.Time
}

type Product struct {
	UUID      uuid.UUID
	MachineId string
	Ok        bool
	Produced  time.Time
}
