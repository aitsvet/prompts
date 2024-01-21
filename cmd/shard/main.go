package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	prompts "github.com/aitsvet/prompts" // replace with actual path to server.go file

	_ "github.com/lib/pq" // postgres driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" // import reflection for gRPC Server
	"github.com/pressly/goose/v3"
)

func main() {
	// Get the database URL and server address from environment variables
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("missing DB_URL in environment variables")
	}
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		log.Fatalf("missing SERVER_ADDR in environment variables")
	}

	// Open the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Wait for DB availability
	for i := 0; i < 12; i++ { // trying 1 minute by default, can be changed with env variable if needed
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("failed to ping database: %v", err)
		time.Sleep(5 * time.Second)
	}

	// Apply migrations
	err = goose.Up(db, ".") // replace '.' with actual path to directory containing migrations
	if err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s) // enable reflection in gRPC server
	prompts.RegisterBalanceServer(s, &prompts.Server{DB: db})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
