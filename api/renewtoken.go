package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khandyan95/simplebank/token"
)

func (server *Server) renewToken(ctx *gin.Context) {
	req := struct {
		RefeshToken string `json:"refresh_token_key" binding:"required"`
	}{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	payload, err := server.Maker.ValidateToken(req.RefeshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorMessage(err))
		return
	}

	session, err := server.Store.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorMessage(token.ErrInvalidToken))
		return
	}

	if req.RefeshToken != session.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, errorMessage(token.ErrInvalidToken))
		return
	}

	if payload.Username != session.Username {
		ctx.JSON(http.StatusUnauthorized, errorMessage(token.ErrInvalidToken))
		return
	}

	if !payload.ExpiresAt.Time.Equal(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, errorMessage(token.ErrInvalidToken))
		return
	}

	if session.ExpiresAt.Before(time.Now()) {
		ctx.JSON(http.StatusUnauthorized, errorMessage(jwt.ErrTokenExpired))
		return
	}

	token, aPayload, err := server.Maker.CreateToken(payload.Username, server.Config.TokenExpiryDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, struct {
		AuthToken     string    `json:"authorization_key"`
		AuthExpiresAt time.Time `json:"authorization_expires_at"`
	}{
		AuthToken:     token,
		AuthExpiresAt: aPayload.ExpiresAt.Time,
	})
}
