package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/thanaphat2005/simple-bank/db/sqlc"
	"github.com/thanaphat2005/simple-bank/token"
	"github.com/thanaphat2005/simple-bank/utils"
)

type Server struct {
	config utils.Configuration
	store  db.Store
	maker  token.Maker
	router *gin.Engine
}

func NewServer(store db.Store, configuration utils.Configuration, tokenMaker token.Maker) (*Server, error) {
	server := &Server{store: store, maker: tokenMaker, config: configuration}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRouter()

	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	router.POST("users/login", server.loginUser)
	router.POST("users", server.createUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.maker))

	authRoutes.POST("accounts", server.createAccount)
	authRoutes.GET("accounts/:id", server.getAccount)
	authRoutes.GET("accounts", server.listAccount)
	authRoutes.POST("transfers", server.createTransfer)
	server.router = router
}

func (server *Server) Start(address string) error {
	err := server.router.Run(address)
	if err != nil {
		return err
	}
	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
