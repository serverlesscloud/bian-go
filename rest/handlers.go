package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/serverlesscloud/bian-go/domains"
)

// Handlers contains all REST endpoint handlers
type Handlers struct {
	accountService     domains.AccountService
	transactionService domains.TransactionService
	balanceService     domains.BalanceService
	consentService     domains.ConsentService
}

// NewHandlers creates a new handlers instance
func NewHandlers(
	accountService domains.AccountService,
	transactionService domains.TransactionService,
	balanceService domains.BalanceService,
	consentService domains.ConsentService,
) *Handlers {
	return &Handlers{
		accountService:     accountService,
		transactionService: transactionService,
		balanceService:     balanceService,
		consentService:     consentService,
	}
}

// Account handlers

// GetAccount handles GET /accounts/{id}
func (h *Handlers) GetAccount(w http.ResponseWriter, r *http.Request) {
	accountID := strings.TrimPrefix(r.URL.Path, "/accounts/")
	if accountID == "" || accountID == "/accounts/" {
		WriteInvalidInputError(w, "Account ID is required")
		return
	}
	
	account, err := h.accountService.RetrieveCurrentAccount(r.Context(), accountID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteNotFoundError(w, "account", accountID)
			return
		}
		WriteInternalError(w, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// GetAccountBalance handles GET /accounts/{id}/balance
func (h *Handlers) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from path like /accounts/{id}/balance
	path := strings.TrimPrefix(r.URL.Path, "/accounts/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] == "" {
		WriteInvalidInputError(w, "Account ID is required")
		return
	}
	accountID := parts[0]
	
	balance, err := h.accountService.RetrieveCurrentAccountBalance(r.Context(), accountID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteNotFoundError(w, "account", accountID)
			return
		}
		WriteInternalError(w, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}

// Transaction handlers

// GetTransaction handles GET /transactions/{id}
func (h *Handlers) GetTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := strings.TrimPrefix(r.URL.Path, "/transactions/")
	if transactionID == "" || transactionID == "/transactions/" {
		WriteInvalidInputError(w, "Transaction ID is required")
		return
	}
	
	transaction, err := h.transactionService.RetrievePaymentTransaction(r.Context(), transactionID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteNotFoundError(w, "transaction", transactionID)
			return
		}
		WriteInternalError(w, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// GetAccountTransactions handles GET /accounts/{id}/transactions
func (h *Handlers) GetAccountTransactions(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from path like /accounts/{id}/transactions
	path := strings.TrimPrefix(r.URL.Path, "/accounts/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] == "" {
		WriteInvalidInputError(w, "Account ID is required")
		return
	}
	accountID := parts[0]
	
	// Parse query parameters
	opts := domains.HistoryOptions{}
	
	if fromDateStr := r.URL.Query().Get("fromDate"); fromDateStr != "" {
		fromDate, err := time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			WriteInvalidInputError(w, "Invalid fromDate format, use YYYY-MM-DD")
			return
		}
		opts.FromDate = &fromDate
	}
	
	if toDateStr := r.URL.Query().Get("toDate"); toDateStr != "" {
		toDate, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			WriteInvalidInputError(w, "Invalid toDate format, use YYYY-MM-DD")
			return
		}
		opts.ToDate = &toDate
	}
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			WriteInvalidInputError(w, "Invalid limit, must be a positive integer")
			return
		}
		if limit > 500 {
			WriteInvalidInputError(w, "Limit cannot exceed 500")
			return
		}
		opts.Limit = limit
	}
	
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			WriteInvalidInputError(w, "Invalid offset, must be a non-negative integer")
			return
		}
		opts.Offset = offset
	}
	
	transactions, err := h.transactionService.RetrievePaymentTransactionHistory(r.Context(), accountID, opts)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteNotFoundError(w, "account", accountID)
			return
		}
		WriteInternalError(w, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// Balance handlers

// GetBalances handles GET /accounts/{id}/balances (all balance types)
func (h *Handlers) GetBalances(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from path like /accounts/{id}/balances
	path := strings.TrimPrefix(r.URL.Path, "/accounts/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] == "" {
		WriteInvalidInputError(w, "Account ID is required")
		return
	}
	accountID := parts[0]
	
	balances, err := h.balanceService.RetrieveAccountBalance(r.Context(), accountID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteNotFoundError(w, "account", accountID)
			return
		}
		WriteInternalError(w, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balances)
}

// Consent handlers

// GetConsent handles GET /consents/{id}
func (h *Handlers) GetConsent(w http.ResponseWriter, r *http.Request) {
	consentID := strings.TrimPrefix(r.URL.Path, "/consents/")
	if consentID == "" || consentID == "/consents/" {
		WriteInvalidInputError(w, "Consent ID is required")
		return
	}
	
	consent, err := h.consentService.RetrieveConsent(r.Context(), consentID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteNotFoundError(w, "consent", consentID)
			return
		}
		WriteInternalError(w, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consent)
}

// GetConsentStatus handles GET /consents/{id}/status
func (h *Handlers) GetConsentStatus(w http.ResponseWriter, r *http.Request) {
	// Extract consent ID from path like /consents/{id}/status
	path := strings.TrimPrefix(r.URL.Path, "/consents/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] == "" {
		WriteInvalidInputError(w, "Consent ID is required")
		return
	}
	consentID := parts[0]
	
	status, err := h.consentService.RetrieveConsentStatus(r.Context(), consentID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteNotFoundError(w, "consent", consentID)
			return
		}
		WriteInternalError(w, err)
		return
	}
	
	response := map[string]interface{}{
		"status": status,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}