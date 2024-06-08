package main

import (
	"database/sql"
	"fmt"

	"github.com/GabrielBrotas/eda-events/internal/database"
	"github.com/GabrielBrotas/eda-events/internal/event"
	"github.com/GabrielBrotas/eda-events/internal/event/handler"
	createaccount "github.com/GabrielBrotas/eda-events/internal/usecase/create_account"
	"github.com/GabrielBrotas/eda-events/internal/usecase/create_client"
	"github.com/GabrielBrotas/eda-events/internal/usecase/create_transaction"
	"github.com/GabrielBrotas/eda-events/internal/usecase/web"
	"github.com/GabrielBrotas/eda-events/internal/usecase/web/webserver"
	"github.com/GabrielBrotas/eda-events/pkg/events"
	"github.com/GabrielBrotas/eda-events/pkg/kafka"
	"github.com/GabrielBrotas/eda-events/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Kafka setup
	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet-group",
	}

	kafkaProducer, err := kafka.NewProducer(&configMap)

	if err != nil {
		panic(err)
	}

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewBalanceUpdatedKafkaHandler(kafkaProducer))

	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	uowInstance := uow.NewUow(db)

	uowInstance.Register("ClientDB", func(tx *sql.Tx) interface{} {
		return database.NewClientDB(tx)
	})

	uowInstance.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(tx)
	})

	uowInstance.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(tx)
	})

	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uowInstance, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)
	createClientUseCase := create_client.NewCreateClientUseCase(uowInstance)
	createAccountUseCase := createaccount.NewCreateAccountUseCase(uowInstance)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running")
	webserver.Start()
}
