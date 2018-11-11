package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin/json"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

var machines = []string{"machine-1", "machine-2", "machine-3", "machine-4", "machine-5"}

var status = []string{"RUNNING", "STOPPED", "SERVICING", "RESETTING"}

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		wg                          sync.WaitGroup
		enqueued, successes, errors int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			fmt.Println("success")
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println(err)
			errors++
		}
	}()

	for _, machine := range machines {
		wg.Add(1)
		go func(machine string) {
			defer wg.Done()
		ProducerLoop:
			for {
				time.Sleep(time.Second * time.Duration(rand.Int63n(2)))

				product := &Product{
					UUID:      uuid.Must(uuid.NewV4()),
					Produced:  time.Now(),
					Ok:        rand.Float32() < 0.9,
					MachineId: machine,
				}

				message := &sarama.ProducerMessage{Topic: "Products", Value: product}
				select {
				case producer.Input() <- message:
					enqueued++
				case <-signals:
					producer.AsyncClose() // Trigger a shutdown of the producer.
					break ProducerLoop
				}
			}
		}(machine)
	}

	for _, machine := range machines {
		wg.Add(1)
		go func(machine string) {
			defer wg.Done()
		ProducerLoop:
			for {
				time.Sleep(time.Second * time.Duration(rand.Int63n(20)))

				status := &Status{
					UUID:      uuid.Must(uuid.NewV4()),
					MachineId: machine,
					Status:    status[rand.Intn(len(status))],
					Since:     time.Now(),
				}

				message := &sarama.ProducerMessage{Topic: "Status", Value: status}
				select {
				case producer.Input() <- message:
					enqueued++
				case <-signals:
					producer.AsyncClose() // Trigger a shutdown of the producer.
					break ProducerLoop
				}
			}
		}(machine)
	}

	wg.Wait()

	log.Printf("Successfully produced: %d; errors: %d\n", successes, errors)
}

type Status struct {
	UUID      uuid.UUID
	MachineId string
	Status    string
	Since     time.Time
}

func (this *Status) Encode() ([]byte, error) {
	return json.Marshal(this)
}

func (this *Status) Length() int {
	bytes, _ := this.Encode()
	return len(bytes)
}

type Product struct {
	UUID      uuid.UUID
	MachineId string
	Ok        bool
	Produced  time.Time
}

func (this *Product) Encode() ([]byte, error) {
	return json.Marshal(this)
}

func (this *Product) Length() int {
	bytes, _ := this.Encode()
	return len(bytes)
}
