// Code generated by SalGen. DO NOT EDIT.
package testdata

import (
	"context"
	"database/sql"
	"github.com/go-gad/sal"
	"github.com/go-gad/sal/looker/testdata/foo-bar"
	"github.com/pkg/errors"
)

type SalStore struct {
	Store
	handler  sal.QueryHandler
	parent   sal.QueryHandler
	ctrl     *sal.Controller
	txOpened bool
}

func NewStore(h sal.QueryHandler, options ...sal.ClientOption) *SalStore {
	s := &SalStore{
		handler:  h,
		ctrl:     sal.NewController(options...),
		txOpened: false,
	}

	return s
}

func (s *SalStore) BeginTx(ctx context.Context, opts *sql.TxOptions) (Store, error) {
	dbConn, ok := s.handler.(sal.TransactionBegin)
	if !ok {
		return nil, errors.New("handler doesn't satisfy the interface TransactionBegin")
	}
	var (
		err error
		tx  *sql.Tx
	)

	ctx = context.WithValue(ctx, sal.ContextKeyTxOpened, s.txOpened)
	ctx = context.WithValue(ctx, sal.ContextKeyOperationType, "Begin")
	ctx = context.WithValue(ctx, sal.ContextKeyMethodName, "BeginTx")

	for _, fn := range s.ctrl.BeforeQuery {
		var fnz sal.FinalizerFunc
		ctx, fnz = fn(ctx, "BEGIN", nil)
		if fnz != nil {
			defer func() { fnz(ctx, err) }()
		}
	}

	tx, err = dbConn.BeginTx(ctx, opts)
	if err != nil {
		err = errors.Wrap(err, "failed to start tx")
		return nil, err
	}

	newClient := &SalStore{
		handler:  tx,
		parent:   s.handler,
		ctrl:     s.ctrl,
		txOpened: true,
	}

	return newClient, nil
}

func (s *SalStore) Tx() sal.Transaction {
	if tx, ok := s.handler.(sal.SqlTx); ok {
		return sal.NewWrappedTx(tx, s.ctrl)
	}
	return nil
}
func (s *SalStore) UpdateAuthor(ctx context.Context, req *foo.Body) error {
	var (
		err      error
		rawQuery = req.Query()
		reqMap   = make(sal.RowMap)
	)

	ctx = context.WithValue(ctx, sal.ContextKeyTxOpened, s.txOpened)
	ctx = context.WithValue(ctx, sal.ContextKeyOperationType, "Exec")
	ctx = context.WithValue(ctx, sal.ContextKeyMethodName, "UpdateAuthor")

	pgQuery, args := sal.ProcessQueryAndArgs(rawQuery, reqMap)

	stmt, err := s.ctrl.PrepareStmt(ctx, s.parent, s.handler, pgQuery)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, fn := range s.ctrl.BeforeQuery {
		var fnz sal.FinalizerFunc
		ctx, fnz = fn(ctx, rawQuery, req)
		if fnz != nil {
			defer func() { fnz(ctx, err) }()
		}
	}

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute Exec")
	}

	return nil
}

// compile time checks
var _ Store = &SalStore{}
