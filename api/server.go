package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/khandyan95/simplebank/db/sqlc"
)

type Server struct {
	Store  *db.Store
	Router *gin.Engine
}

func NewServer(store *db.Store) *Server {

	server := &Server{}
	server.Store = store

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", currencyValidator)
	}

	router.POST("/user", server.createUser)
	// router.POST("/user/changepassword", server.updateUser)

	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/account", server.listAccounts)
	router.POST("/account/accounttransfer", server.createAccountTransaction)

	server.Router = router

	return server
}

func errorMessage(err error) gin.H {
	return gin.H{"message": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
