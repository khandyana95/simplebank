package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authPayloadKey    = "authpayload"
	authSupportedType = "bearer"
)

func (server *Server) authUser() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		req := struct {
			AuthToken string `header:"authorization" binding:"required"`
		}{}

		if err := ctx.ShouldBindHeader(&req); err != nil {
			authError := fmt.Errorf("no authorization key provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorMessage(authError))
			return
		}

		fields := strings.Fields(req.AuthToken)
		if len(fields) < 2 {
			authError := fmt.Errorf("invalid authorization key provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorMessage(authError))
			return
		}

		if strings.ToLower(fields[0]) != authSupportedType {
			authError := fmt.Errorf("unsupported authorization key provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorMessage(authError))
			return
		}

		payload, err := server.Maker.ValidateToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorMessage(err))
			return
		}

		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}
