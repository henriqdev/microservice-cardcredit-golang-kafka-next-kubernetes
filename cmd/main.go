package main

import (
	"database/sql"
	"encoding/json"
	"log"

	// "os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/henriqdev/gateway-go/adapter/broker/kafka"
	"github.com/henriqdev/gateway-go/adapter/factory"
	"github.com/henriqdev/gateway-go/adapter/presenter/transaction"
	"github.com/henriqdev/gateway-go/usercase/process_transaction"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	repositoryFactory := factory.NewRepositoryDatabaseFactory(db)
	repository := repositoryFactory.CreateTransactionRepository()
	configMapProducer := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		// "security.protocol": os.Getenv("SECURITY_PROTOCOL"),
		// "sasl.mechanisms":   os.Getenv("SASL_MECHANISMS"),
		// "sasl.username":     os.Getenv("SASL_USERNAME"),
		// "sasl.password":     os.Getenv("SASL_PASSWORD"),
	}
	kafkaPresenter := transaction.NewTransactionKafkaPresenter()
	producer := kafka.NewKafkaProducer(configMapProducer, kafkaPresenter)

	var msgChan = make(chan *ckafka.Message)
	configMapConsumer := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		// "security.protocol": os.Getenv("SECURITY_PROTOCOL"),
		// "sasl.mechanisms":   os.Getenv("SASL_MECHANISMS"),
		// "sasl.username":     os.Getenv("SASL_USERNAME"),
		// "sasl.password":     os.Getenv("SASL_PASSWORD"),
		"client.id": "goapp",
		"group.id":  "goapp",
	}
	topics := []string{"transactions"}
	consumer := kafka.NewConsumer(configMapConsumer, topics)
	go consumer.Consume(msgChan)

	usecase := process_transaction.NewProcessTransaction(repository, producer, "transactions_result")

	for msg := range msgChan {
		var input process_transaction.TransactionDtoInput
		json.Unmarshal(msg.Value, &input)
		usecase.Execute(input)
	}
}
