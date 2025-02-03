package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/khandyan95/simplebank/db/sqlc"
	"github.com/khandyan95/simplebank/util"
	"github.com/lib/pq"
)

func (server *Server) loginUser(ctx *gin.Context) {

	req := struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	user, err := server.Store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(fmt.Errorf("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	if err := util.VerifyPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorMessage(err))
		return
	}

	token, payload, err := server.Maker.CreateToken(user.Username, server.Config.TokenExpiryDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	refreshToken, rPayload, err := server.Maker.CreateToken(user.Username, server.Config.RefreshTokenExpiryDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	session, err := server.Store.CreateSession(ctx, db.CreateSessionParams{
		ID:           rPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    rPayload.ExpiresAt.Time,
	})

	if err != nil {
		if pqerror, ok := err.(*pq.Error); ok {
			switch pqerror.Code.Name() {
			case "foreign_key_violation":
				ctx.JSON(http.StatusForbidden, errorMessage(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, struct {
		Username             string    `json:"username"`
		FullName             string    `json:"full_name"`
		Email                string    `json:"email"`
		PasswordChangedAt    time.Time `json:"password_changed_at"`
		CreatedAt            time.Time `json:"created_at"`
		AuthToken            string    `json:"authorization_key"`
		AuthExpiresAt        time.Time `json:"authorization_expires_at"`
		RefeshToken          string    `json:"refresh_token_key"`
		RefeshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	}{
		Username:             user.Username,
		FullName:             user.FullName,
		Email:                user.Email,
		PasswordChangedAt:    user.PasswordChangedAt,
		CreatedAt:            user.CreatedAt,
		AuthToken:            token,
		AuthExpiresAt:        payload.ExpiresAt.Time,
		RefeshToken:          session.RefreshToken,
		RefeshTokenExpiresAt: session.ExpiresAt,
	})

}
