package rest

import (
	"encoding/json"
	"net/http"

	"github.com/serverlesscloud/bian-go/domains"
)

// Server represents the REST API server
type Server struct {
	handlers *Handlers
	mux      *http.ServeMux
}

// NewServer creates a new REST server with all routes configured
func NewServer(
	accountService domains.AccountService,
	transactionService domains.TransactionService,
	balanceService domains.BalanceService,
	consentService domains.ConsentService,
) *Server {
	handlers := NewHandlers(accountService, transactionService, balanceService, consentService)
	
	mux := http.NewServeMux()
	server := &Server{
		handlers: handlers,
		mux:      mux,
	}
	
	server.setupRoutes()
	return server
}

// setupRoutes configures all REST API routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.mux.HandleFunc("/health", s.healthCheck)
	
	// Account endpoints
	s.mux.HandleFunc("/accounts/", s.routeAccountRequests)
	
	// Transaction endpoints
	s.mux.HandleFunc("/transactions/", s.handlers.GetTransaction)
	
	// Consent endpoints
	s.mux.HandleFunc("/consents/", s.routeConsentRequests)
}

// routeAccountRequests routes account-related requests based on path
func (s *Server) routeAccountRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		WriteErrorResponse(w, ErrorCodeInvalidInput, "Method not allowed", "Only GET requests are supported", http.StatusMethodNotAllowed)
		return
	}
	
	path := r.URL.Path
	
	// Check for /accounts/{id}/balance
	if len(path) > 10 && path[len(path)-8:] == "/balance" {
		s.handlers.GetAccountBalance(w, r)
		return
	}
	
	// Check for /accounts/{id}/balances (all balance types)
	if len(path) > 11 && path[len(path)-9:] == "/balances" {
		s.handlers.GetBalances(w, r)
		return
	}
	
	// Check for /accounts/{id}/transactions
	if len(path) > 15 && path[len(path)-13:] == "/transactions" {
		s.handlers.GetAccountTransactions(w, r)
		return
	}
	
	// Default to account retrieval /accounts/{id}
	s.handlers.GetAccount(w, r)
}

// routeConsentRequests routes consent-related requests based on path
func (s *Server) routeConsentRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		WriteErrorResponse(w, ErrorCodeInvalidInput, "Method not allowed", "Only GET requests are supported", http.StatusMethodNotAllowed)
		return
	}
	
	path := r.URL.Path
	
	// Check for /consents/{id}/status
	if len(path) > 9 && path[len(path)-7:] == "/status" {
		s.handlers.GetConsentStatus(w, r)
		return
	}
	
	// Default to consent retrieval /consents/{id}
	s.handlers.GetConsent(w, r)
}

// healthCheck provides a simple health check endpoint
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		WriteErrorResponse(w, ErrorCodeInvalidInput, "Method not allowed", "Only GET requests are supported", http.StatusMethodNotAllowed)
		return
	}
	
	response := map[string]interface{}{
		"status":  "healthy",
		"service": "bian-go",
		"version": "1.0.0",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler returns the HTTP handler with middleware applied
func (s *Server) Handler() http.Handler {
	var handler http.Handler = s.mux
	
	// Apply middleware in reverse order (last applied = first executed)
	handler = ContentTypeMiddleware(handler)
	handler = CORSMiddleware([]string{"*"})(handler) // Allow all origins for development
	handler = ErrorRecoveryMiddleware(handler)
	handler = RequestLoggingMiddleware(handler)
	
	return handler
}