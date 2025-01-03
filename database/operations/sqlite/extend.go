package sqliteops

import (
	"context"
	"database/sql"
)

func (this *Queries) IsNoRows(err error) bool {
	return err == sql.ErrNoRows
}

type txBeginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type txFinisher interface {
	Commit() error
	Rollback() error
}

type TxQueries struct {
	*Queries
}

func (this *TxQueries) Commit() error {
	return this.db.(txFinisher).Commit()
}

func (this *TxQueries) Rollback() error {
	return this.db.(txFinisher).Rollback()
}

func (q *Queries) BeginTx(ctx context.Context) (*TxQueries, error) {

	tx, err := q.db.(txBeginner).BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &TxQueries{Queries: New(nil).WithTx(tx)}, nil
}
