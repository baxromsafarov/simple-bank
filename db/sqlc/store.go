package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

// Store предоставляет все функции для выполнения запросов к БД и транзакций.
type Store struct {
	*Queries
	conn *pgx.Conn
}

// NewStore создает новый объект Store.
func NewStore(conn *pgx.Conn) *Store {
	return &Store{
		conn:    conn,
		Queries: New(conn),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

// TransferTxParams содержит параметры для операции перевода денег.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult содержит результат операции перевода.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Шаг 1: Создаем запись о переводе (transfer)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Шаг 2: Создаем запись о движении средств (entry) для отправителя
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // Сумма отрицательная, так как это списание
		})
		if err != nil {
			return err
		}

		// Шаг 3: Создаем entry для получателя
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // Сумма положительная
		})
		if err != nil {
			return err
		}

		// Шаг 4: Обновляем балансы счетов
		// Чтобы избежать deadlock, обновляем счета в порядке их ID
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:      arg.FromAccountID,
				Balance: -arg.Amount,
			})
			if err != nil {
				return err
			}
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:      arg.ToAccountID,
				Balance: arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:      arg.ToAccountID,
				Balance: arg.Amount,
			})
			if err != nil {
				return err
			}
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:      arg.FromAccountID,
				Balance: -arg.Amount,
			})
			if err != nil {
				return err
			}
		}

		return nil

	})
	return result, err
}
