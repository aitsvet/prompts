package prompts

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddNamedMigration("1_change.go", upChangeRequest, downChangeRequest)
	goose.AddNamedMigration("2_get.go", upGetDataResponse, downGetDataResponse)
}

func upChangeRequest(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`CREATE TABLE change_requests (
        transaction_id bigint NOT NULL,
        account_id bigint NOT NULL,
        amount bigint NOT NULL);`)
	return err
}

func downChangeRequest(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec("DROP TABLE change_requests;")
	return err
}

func upGetDataResponse(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`CREATE TABLE get_data_responses (
        balance bigint NOT NULL,
        transaction_id bigint NOT NULL,
        account_id bigint PRIMARY KEY);`)
	return err
}

func downGetDataResponse(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec("DROP TABLE get_data_responses;")
	return err
}
