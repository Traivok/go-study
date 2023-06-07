package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/traivok/go-study/db/sqlc"
	"net/http"
)

type transferRequest struct {
	FromAccountID int64  `json:"fromAccountID" binding:"required,min=1"`
	ToAccountID   int64  `json:"toAccountID" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(context *gin.Context) {
	var req transferRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(context, req.ToAccountID, req.Currency) ||
		!server.validAccount(context, req.FromAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	account, err := server.store.TransferTx(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, account)
}

func (server *Server) validAccount(context *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(context, accountID)

	switch err {
	case sql.ErrNoRows:
		context.JSON(http.StatusNotFound, errorResponse(err))
		return false
	default:
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	case nil:
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch.\n Account curency: %s. Requested currency %s.", account.ID, account.Currency, currency)
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
