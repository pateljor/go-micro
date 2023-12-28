package main

import (
	"fmt"
	"log"
	"math"
	"time"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil{
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ!")
	
	// start listening for messages

	// create consumer

	// watch the queue and consume messages from the queue

}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection


	// dont continue until rabbit is readyw
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			fmt.Println("RabbitMQ Not Ready...")
			counts++
		} else{
			connection = c
			break
		}
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil

}