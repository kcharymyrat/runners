package repositories

import (
	"context"
	"database/sql"
)

func BeginTransaction(runnersRepository *RunnersRepository, resResultsRepository *ResultsRepository) error {
	ctx := context.Background()
	transaction, err := resResultsRepository.dbHandler.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	runnersRepository.transaction = transaction
	resResultsRepository.transaction = transaction
	return nil
}

func RollbackTransaction(runnersRepository *RunnersRepository, resResultsRepository *ResultsRepository) error {
	transaction := runnersRepository.transaction
	runnersRepository.transaction = nil
	resResultsRepository.transaction = nil
	return transaction.Rollback()
}

func CommitTransaction(runnersRepository *RunnersRepository, resResultsRepository *ResultsRepository) error {
	transaction := runnersRepository.transaction
	runnersRepository.transaction = nil
	resResultsRepository.transaction = nil
	return transaction.Commit()
}
