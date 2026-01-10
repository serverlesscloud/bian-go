package main

import (
	"log"

	"github.com/serverlesscloud/bian-go/providers/mock"
	"github.com/serverlesscloud/bian-go/server"
)

func main() {
	// Initialize mock provider with sample data
	provider := mock.NewProvider()
	
	// Create server configuration
	config := server.DefaultConfig()
	
	// Create unified server with both REST and GraphQL APIs
	srv := server.NewServer(
		provider, // AccountService
		provider, // TransactionService
		provider, // BalanceService
		provider, // ConsentService
		config,
	)
	
	// Start server (blocks until shutdown)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}