package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/traivok/go-study/db/sqlc"
	"log"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	server.router = gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validateCurrency)
		if err != nil {
			log.Fatal("Could no register currency validator.")
		}
	}

	server.router.POST("/accounts", server.createAccount)
	server.router.GET("/accounts/:id", server.getAccount)
	server.router.GET("/accounts", server.listAccount)

	server.router.POST("/transfers", server.createTransfer)

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
