package factory

import "github.com/henriqdev/gateway-go/domain/repository"

type RepositoryFactory interface {
	CreateTransactionRepository() repository.TransactionRepository
}
