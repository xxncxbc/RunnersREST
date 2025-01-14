package repositories

import "database/sql"

type RunnersRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewRunnersRepository(dbHandler *sql.DB, transaction *sql.Tx) *RunnersRepository {
	return &RunnersRepository{dbHandler: dbHandler, transaction: transaction}
}
