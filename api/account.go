package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/khandyan95/simplebank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD INR"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	accRequest := createAccountRequest{}
	if err := ctx.ShouldBindJSON(&accRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	account, err := server.Store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    accRequest.Owner,
		Name:     accRequest.Name,
		Currency: accRequest.Currency,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type accountIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {

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

	ctx.JSON(http.StatusOK, account)
}

type accountQueryParams struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {

	accQueryRequest := accountQueryParams{}
	if err := ctx.ShouldBindQuery(&accQueryRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	accounts, err := server.Store.ListAccounts(ctx, db.ListAccountsParams{
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
