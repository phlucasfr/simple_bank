package api

import (
	db "picpay_simplificado/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP for requests
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	//wallets
	router.POST("/wallets", server.createWallet)
	router.GET("/wallets/:id", server.getWallet)
	router.GET("/wallets", server.listWallets)
	router.PUT("/wallets", server.updateWallet)
	router.DELETE("/wallets/:id", server.deleteWallet)

	//users
	router.GET("/users/:id", server.getUser)
	router.POST("/users", server.createUser)
	router.DELETE("/users/:id", server.deleteUser)
	router.GET("/users", server.listUsers)
	router.PUT("/users", server.updateUser)

	//entries
	router.GET("/entries/:id", server.getEntry)
	router.POST("/entries", server.createEntry)
	router.GET("/entries", server.listEntries)

	//transfer
	router.POST("/transfers", server.createTransfer)

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

func successResponse(msg string) gin.H {
	return gin.H{"msg": msg}
}
