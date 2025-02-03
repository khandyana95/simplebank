package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/khandyan95/simplebank/db/sqlc"
	"github.com/khandyan95/simplebank/util"
	"github.com/lib/pq"
)

func (server *Server) createUser(ctx *gin.Context) {

	req := struct {
		Username string `json:"username" binding:"required,alphanum"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	user, err := server.Store.CreateUser(ctx, db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	})

	if err != nil {
		if pqerror, ok := err.(*pq.Error); ok {
			switch pqerror.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorMessage(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return

	}

	ctx.JSON(http.StatusOK, struct {
		Username  string    `json:"username"`
		FullName  string    `json:"full_name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

// func (server *Server) updateUser(ctx *gin.Context) {

// 	var req = struct {
// 		Username string `json:"username" binding:"required"`
// 		Password string `json:"password" binding:"required"`
// 	}{}

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorMessage(err))
// 		return
// 	}

// 	hashedPassword, err := util.HashPassword(req.Password)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
// 		return
// 	}

// 	user, err := server.Store.UpdateUser(ctx, db.UpdateUserParams{
// 		Username:          req.Username,
// 		HashedPassword:    hashedPassword,
// 		PasswordChangedAt: time.Now(),
// 	})

// 	if err != nil {
// 		if pqerror, ok := err.(*pq.Error); ok {
// 			ctx.JSON(http.StatusInternalServerError, errorMessage(pqerror))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
// 		return

// 	}

// 	ctx.JSON(http.StatusOK, struct {
// 		Username          string    `json:"username"`
// 		PasswordChangedAt time.Time `json:"created_at"`
// 	}{
// 		Username:          user.Username,
// 		PasswordChangedAt: user.PasswordChangedAt,
// 	})

// }
