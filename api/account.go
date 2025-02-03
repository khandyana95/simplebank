package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/khandyan95/simplebank/db/sqlc"
	"github.com/khandyan95/simplebank/token"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Name     string `json:"name" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	authPayload, ok := ctx.Get(authPayloadKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorMessage(token.ErrInvalidToken))
		return
	}

	payload, ok := authPayload.(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorMessage(token.ErrInvalidToken))
		return
	}

	accRequest := createAccountRequest{}
	if err := ctx.ShouldBindJSON(&accRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	account, err := server.Store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    payload.Username,
		Name:     accRequest.Name,
		Currency: accRequest.Currency,
	})

	if err != nil {
		if pgerror, ok := err.(*pq.Error); ok {
			switch pgerror.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorMessage(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type accountIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {

	authPayload, ok := ctx.Get(authPayloadKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorMessage(token.ErrInvalidToken))
		return
	}

	payload, ok := authPayload.(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorMessage(token.ErrInvalidToken))
		return
	}

	accIdRequest := accountIdRequest{}
	if err := ctx.ShouldBindUri(&accIdRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	account, err := server.Store.GetAccount(ctx, accIdRequest.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(fmt.Errorf("account not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	if payload.Username != account.Owner {
		ctx.JSON(http.StatusUnauthorized, errorMessage(fmt.Errorf("user not authorized to account")))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type accountQueryParams struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {

	authPayload, ok := ctx.Get(authPayloadKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorMessage(token.ErrInvalidToken))
		return
	}

	payload, ok := authPayload.(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorMessage(token.ErrInvalidToken))
		return
	}

	accQueryRequest := accountQueryParams{}
	if err := ctx.ShouldBindQuery(&accQueryRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	accounts, err := server.Store.ListAccounts(ctx, db.ListAccountsParams{
		Owner:  payload.Username,
		Limit:  accQueryRequest.PageSize,
		Offset: (accQueryRequest.PageID - 1) * accQueryRequest.PageSize,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(fmt.Errorf("accounts not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
