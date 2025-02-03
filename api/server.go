package api

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/khandyan95/simplebank/db/sqlc"
	"github.com/khandyan95/simplebank/token"
	"github.com/khandyan95/simplebank/util"
	_ "github.com/lib/pq"
)

type Server struct {
	Store  *db.Store
	Router *gin.Engine
	Maker  token.Maker
	Config *util.Config
}

func NewServer() (*Server, error) {

	// Load the config file
	config, err := util.LoadConfing(".")
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("cannot load config file")
	}

	// Open the DB connection
	dbConn, err := sql.Open(config.DBDriver, config.DataSource)
	if err != nil {
		return nil, fmt.Errorf("error in connecting to DB")
	}

	// Create the Maker struct to generate tokens
	maker, err := token.NewJWTMaker(config.SecretKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		Store:  db.NewStore(dbConn),
		Config: &config,
		Maker:  maker,
		Router: gin.Default(),
	}

	// Get Validator
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, fmt.Errorf("failed to load validator")
	}

	// Register custom validator
	if err := v.RegisterValidation("currency", currencyValidator); err != nil {
		return nil, fmt.Errorf("failed to register currency validator")
	}

	// Resister Routes/Handlers
	registerRoutes(server)

	return server, nil
}

func registerRoutes(server *Server) {

	// No auth required
	server.Router.POST("/user", server.createUser)
	server.Router.POST("/user/login", server.loginUser)
	server.Router.POST("/token/renewtoken", server.renewToken)

	//Auth required
	group := server.Router.Group("/")
	group.Use(server.authUser())

	group.POST("/account", server.createAccount)
	group.GET("/account/:id", server.getAccount)
	group.GET("/account", server.listAccounts)
	group.POST("/account/accounttransfer", server.createAccountTransaction)
}

func errorMessage(err error) gin.H {
	return gin.H{"message": err.Error()}
}

func (server *Server) Start() error {
	return server.Router.Run(server.Config.ServerAddress)
}
