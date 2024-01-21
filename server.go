package prompts

import (
	"context"
	"database/sql"
)

type Server struct {
	DB *sql.DB
	UnimplementedBalanceServer
}

func (s *Server) ChangeBalance(ctx context.Context, in *ChangeRequest) (*ChangeResponse, error) {
	// Begin a new transaction
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return &ChangeResponse{Success: false}, err
	}

	// Rollback the transaction if we get an error at any point
	defer tx.Rollback()

	// Insert the transaction into the database
	_, err = tx.Exec("INSERT INTO change_requests(transaction_id, account_id, amount) VALUES ($1, $2, $3)", in.TransactionId, in.AccountId, in.Amount)
	if err != nil {
		return &ChangeResponse{Success: false}, err
	}

	// Perform UPSERT operation to update the balance in the database
	_, err = tx.Exec("INSERT INTO get_data_responses (balance, transaction_id, account_id) VALUES ($1, $2, $3) ON CONFLICT (account_id) DO UPDATE SET balance = get_data_responses.balance + EXCLUDED.balance, transaction_id = EXCLUDED.transaction_id", in.Amount, in.TransactionId, in.AccountId)
	if err != nil {
		return &ChangeResponse{Success: false}, err
	}

	// Commit the transaction if there are no errors
	err = tx.Commit()
	if err != nil {
		return &ChangeResponse{Success: false}, err
	}

	return &ChangeResponse{Success: true}, nil
}

func (s *Server) GetAccountData(ctx context.Context, in *GetDataRequest) (*GetDataResponse, error) {
	// Query the balance from the database
	var balance int64
	err := s.DB.QueryRow("SELECT balance FROM get_data_responses WHERE account_id = $1", in.AccountId).Scan(&balance)
	if err != nil {
		return nil, err
	}
	// Query all transactions for the account from the database
	rows, err := s.DB.Query("SELECT transaction_id, amount FROM change_requests WHERE account_id = $1 ORDER BY transaction_id", in.AccountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []*GetDataResponse_Transaction
	for rows.Next() {
		t := &GetDataResponse_Transaction{}
		err = rows.Scan(&t.Id, &t.Amount)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	// Return the balance and transactions
	return &GetDataResponse{Balance: balance, Transactions: transactions}, nil
}
