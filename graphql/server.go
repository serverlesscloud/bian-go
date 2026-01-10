package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/serverlesscloud/bian-go/domains"
	"github.com/serverlesscloud/bian-go/graphql/generated"
)

// Server represents the GraphQL server
type Server struct {
	resolver *Resolver
	handler  http.Handler
}

// NewServer creates a new GraphQL server
func NewServer(
	accountService domains.AccountService,
	transactionService domains.TransactionService,
	balanceService domains.BalanceService,
	consentService domains.ConsentService,
) *Server {
	resolver := NewResolver(accountService, transactionService, balanceService, consentService)
	
	config := generated.Config{Resolvers: resolver}
	schema := generated.NewExecutableSchema(config)
	
	srv := handler.NewDefaultServer(schema)
	
	return &Server{
		resolver: resolver,
		handler:  srv,
	}
}

// Handler returns the GraphQL HTTP handler
func (s *Server) Handler() http.Handler {
	return s.handler
}

// PlaygroundHandler returns the GraphQL Playground handler for development
func (s *Server) PlaygroundHandler() http.Handler {
	return playground.Handler("GraphQL Playground", "/graphql")
}