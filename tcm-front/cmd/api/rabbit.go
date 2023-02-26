package main

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"tcm-web/tcm-signup/event"
	"time"
)

// rabbitConnect 获取rabbitmq connection
func rabbitConnect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			return connection, nil
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
}

func (app *Config) pushToQueue(info *SignRequest) error {
	// 新建一个患者消息生产者
	emitter := event.NewEventEmitter(app.Rabbit, "signup")

	// 序列化
	j, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		return err
	}

	// 放入队列
	err = emitter.Push(string(j))
	if err != nil {
		return err
	}
	return nil
}
