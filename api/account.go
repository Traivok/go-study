package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/traivok/go-study/db/sqlc"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR" `
}

func (server *Server) createAccount(context *gin.Context) {
	var req createAccountRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(context, arg)

	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(context *gin.Context) {
	var req getAccountRequest

	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(context, req.ID)

	switch err {
	case nil:
		context.JSON(http.StatusOK, account)
	case sql.ErrNoRows:
		context.JSON(http.StatusNotFound, errorResponse(err))
	default:
		context.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return
}

type listAccountRequest struct {
	PageID   int32 `form:"pageId" binding:"required,min=1"`
	PageSize int32 `form:"pageSize" binding:"required,min=8,max=128"`
}

func (server *Server) listAccount(context *gin.Context) {
	var req listAccountRequest

	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(context, args)

	switch err {
	case nil:
		context.JSON(http.StatusOK, accounts)
	case sql.ErrNoRows:
		context.JSON(http.StatusNotFound, errorResponse(err))
	default:
		context.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return
}
