package api

import (
	"database/sql"
	"fmt"
	"net/http"
	db "picpay_simplificado/db/sqlc"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromWalletID int64  `json:"from_wallet_id" binding:"required,min=1"`
	ToWalletID   int64  `json:"to_wallet_id" binding:"required,min=1"`
	Amount       int64  `json:"amount" binding:"required,gt=0"`
	Currency     string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validateWallet(ctx, req.FromWalletID, req.Currency) {
		return
	}

	if !server.validateWallet(ctx, req.ToWalletID, req.Currency) {
		return
	}

	arg := db.TrasferTxParms{
		FromWalletID: req.FromWalletID,
		ToWalletID:   req.ToWalletID,
		Amount:       req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateWallet(ctx *gin.Context, walletID int64, currency string) bool {
	wallet, err := server.store.GetWallet(ctx, walletID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if wallet.Currency != currency {
		err = fmt.Errorf("wallet [%d] currency mismarch: %s vs %s", walletID, wallet.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
