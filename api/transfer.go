package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/khandyan95/simplebank/db/sqlc"
	"github.com/khandyan95/simplebank/token"
)

type transferRequestParams struct {
	FromAccountId int64   `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64   `json:"to_account_id" binding:"required,min=1"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Currency      string  `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccountTransaction(ctx *gin.Context) {

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

	req := transferRequestParams{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountId, req.Currency)
	if !valid {
		return
	}

	if fromAccount.Owner != payload.Username {
		ctx.JSON(http.StatusUnauthorized, errorMessage(fmt.Errorf("user not authorized to account")))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountId, req.Currency)
	if !valid {
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

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {

	account, err := server.Store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(fmt.Errorf("account with ID %v not found", accountId)))
			return db.Account{}, false
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return db.Account{}, false
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusBadRequest,
			errorMessage(fmt.Errorf("account ID %v currency %v do not match with request currency %v", accountId, account.Currency, currency)))
		return account, false
	}

	return account, true

}
