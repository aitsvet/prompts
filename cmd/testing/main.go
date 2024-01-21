package main

import (
	"context"
	"testing"
	"time"
	"google.golang.org/grpc"
)

func TestTransferService(client TransferClient, ctx context.Context, t *testing.T) {
	cases := []struct{fromAccountId, toAccountId, amount int64}{
		{1, 2, 500},
	}
	for _, c := range cases {
		resp, err := client.TransferMoney(ctx, &TransferRequest{FromAccountId: c.fromAccountId, ToAccountId: c.toAccountId, Amount: c.amount})
		if err != nil {
			t.Fatalf("could not transfer money: %v", err)
		}
		if resp.GetSuccess() == false {
			t.Errorf("expected true; got false")
		}

		// check if balance has changed correctly
		fromAccountData, err := client.GetAccountData(ctx, &GetDataRequest{AccountId: c.fromAccountId})
		if err != nil {
			t.Fatalf("could not get from account data: %v", err)
		}
		if fromAccountData.GetBalance()+c.amount != 0 { // assuming initial balance is 1000 for account with id 1
			println(fromAccountData.GetBalance(), c.amount)
			t.Errorf("expected from account balance to be %f; got %f", 1000-c.amount, fromAccountData.GetBalance())
		}

		// check if transaction was recorded correctly
		if len(fromAccountData.GetTransactions()) != 1 {
			t.Errorf("expected 1 transaction; got %d", len(fromAccountData.GetTransactions()))
		}
		if fromAccountData.GetTransactions()[0].GetAmount()+c.amount != 0{
			println(fromAccountData.GetTransactions()[0].GetAmount(), c.amount)
			t.Errorf("expected transaction amount to be %f; got %f", c.amount, fromAccountData.GetTransactions()[0].GetAmount())
		}

		// check if balance has changed correctly for the recipient account
		toAccountData, err := client.GetAccountData(ctx, &GetDataRequest{AccountId: c.toAccountId})
		if err != nil {
			t.Fatalf("could not get to account data: %v", err)
		}
		if toAccountData.GetBalance()-c.amount != 0 { // assuming initial balance is 1000 for account with id 2
			t.Errorf("expected from account balance to be %f; got %f", 1000+c.amount, toAccountData.GetBalance())
		}
	}
}

func main() {
	t := &testing.T{}
	conn, err := grpc.Dial("localhost:5002", grpc.WithInsecure()) // replace with your port
	if err != nil {
		t.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()
	client := NewTransferClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	TestTransferService(client, ctx, t) // replace t with your testing instance
}
