package api

import (
	db "picpay_simplificado/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP for requests
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/wallets", server.createWallet)
	router.GET("/wallets/:id", server.getWallet)
	router.GET("/wallets", server.listWallets)

	//add routes to router
	server.router = router
	return server
}

// StartServer runs the HTTP server on a specif address
func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
