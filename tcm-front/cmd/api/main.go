package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

const webPort = "8080"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	// 连接rabbitMQ
	rabbitConn, err := rabbitConnect()
	if err != nil {
		log.Println(err)
	}
	defer rabbitConn.Close()

	// 实例化app
	app := Config{
		rabbitConn,
	}

	// 创建http服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Router(),
	}

	// start the server
	log.Printf("Starting front service on port %s\n", webPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
