package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khandyan95/simplebank/util"
)

func (server *Server) loginUser(ctx *gin.Context) {

	req := struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := server.Store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(fmt.Errorf("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := util.VerifyPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	token, err := server.Maker.CreateToken(user.Username, server.Config.TokenExpiryDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, struct {
		Username          string    `json:"username"`
		FullName          string    `json:"full_name"`
		Email             string    `json:"email"`
		PasswordChangedAt time.Time `json:"password_changed_at"`
		CreatedAt         time.Time `json:"created_at"`
		AuthToken         string    `json:"authorization_key"`
	}{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		AuthToken:         token,
	})

}
