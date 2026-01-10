package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serverlesscloud/bian-go/domains"
	"github.com/serverlesscloud/bian-go/graphql"
	"github.com/serverlesscloud/bian-go/rest"
)

// Server represents the unified HTTP server with both REST and GraphQL APIs
type Server struct {
	httpServer *http.Server
	port       string
}

// Config holds server configuration
type Config struct {
	Port                string
	EnablePlayground    bool
	AllowedCORSOrigins  []string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
}

// DefaultConfig returns default server configuration
func DefaultConfig() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		EnablePlayground:   getEnv("ENABLE_PLAYGROUND", "true") == "true",
		AllowedCORSOrigins: []string{"*"}, // Allow all origins for development
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		IdleTimeout:        60 * time.Second,
	}
}

// NewServer creates a new unified server with both REST and GraphQL endpoints
func NewServer(
	accountService domains.AccountService,
	transactionService domains.TransactionService,
	balanceService domains.BalanceService,
	consentService domains.ConsentService,
	config *Config,
) *Server {
	if config == nil {
		config = DefaultConfig()
	}
	
	// Create REST server
	restServer := rest.NewServer(accountService, transactionService, balanceService, consentService)
	
	// Create GraphQL server
	graphqlServer := graphql.NewServer(accountService, transactionService, balanceService, consentService)
	
	// Create main HTTP mux
	mux := http.NewServeMux()
	
	// Mount REST API endpoints
	mux.Handle("/accounts/", restServer.Handler())
	mux.Handle("/transactions/", restServer.Handler())
	mux.Handle("/consents/", restServer.Handler())
	mux.Handle("/health", restServer.Handler())
	
	// Mount GraphQL endpoint
	mux.Handle("/graphql", graphqlServer.Handler())
	
	// Mount GraphQL Playground (development only)
	if config.EnablePlayground {
		mux.Handle("/playground", graphqlServer.PlaygroundHandler())
	}
	
	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      mux,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}
	
	return &Server{
		httpServer: httpServer,
		port:       config.Port,
	}
}

// Start starts the server and blocks until shutdown
func (s *Server) Start() error {
	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", s.port)
		log.Printf("📊 REST API: http://localhost:%s/", s.port)
		log.Printf("🎮 GraphQL API: http://localhost:%s/graphql", s.port)
		log.Printf("🛝 GraphQL Playground: http://localhost:%s/playground", s.port)
		log.Printf("❤️  Health Check: http://localhost:%s/health", s.port)
		
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("🛑 Server shutting down...")
	
	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Shutdown server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	
	log.Println("✅ Server shutdown complete")
	return nil
}

// Stop gracefully stops the server
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	return s.httpServer.Shutdown(ctx)
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}