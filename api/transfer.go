package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/khandyan95/simplebank/db/sqlc"
)

type transferRequestParams struct {
	FromAccountId int64   `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64   `json:"to_account_id" binding:"required,min=1"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Currency      string  `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccountTransaction(ctx *gin.Context) {

	req := transferRequestParams{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountId, req.Currency) {
		return
	}

	if !server.validAccount(ctx, req.ToAccountId, req.Currency) {
		return
	}

	txnResult, err := server.Store.CreateTransaction(ctx, db.CreateTransactionParams{
		FromAccountId: req.FromAccountId,
		ToAccountId:   req.ToAccountId,
		Amount:        req.Amount,
		Currency:      req.Currency,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, txnResult)
}

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) bool {

	accout, err := server.Store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(fmt.Errorf("account with ID %v not found", accountId)))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return false
	}

	if accout.Currency != currency {
		ctx.JSON(http.StatusBadRequest,
			errorMessage(fmt.Errorf("account ID %v currency %v do not match with request currency %v", accountId, accout.Currency, currency)))
		return false
	}

	return true

}
