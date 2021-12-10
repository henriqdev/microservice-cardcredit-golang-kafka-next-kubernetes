package kafka

import (
	"testing"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/henriqdev/gateway-go/adapter/presenter/transaction"
	"github.com/henriqdev/gateway-go/domain/entity"
	"github.com/henriqdev/gateway-go/usercase/process_transaction"
	"github.com/stretchr/testify/assert"
)

func TestPRoducerPublish(t *testing.T) {
	expectedOutput := process_transaction.TransactionDtoOutput{
		ID:           "1",
		Status:       entity.REJECTED,
		ErrorMessage: "you dont have limit for this transation",
	}
	// outputJson, _ := json.Marshal(expectedOutput)
	configMap := ckafka.ConfigMap{
		"test.mock.num.brokers": 3,
	}
	producer := NewKafkaProducer(&configMap, transaction.NewTransactionKafkaPresenter())
	err := producer.Publish(expectedOutput, []byte("1"), "test")
	assert.Nil(t, err)
}
