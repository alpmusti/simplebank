package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/alpmusti/simplebank/db/sqlc"
	"github.com/alpmusti/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithError(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, respondWithError(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, respondWithError(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithError(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, respondWithError(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, respondWithError(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("accounts doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, respondWithError(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountRequest struct {
	PageID   int32 `form:"pageId" binding:"required,min=1"`
	PageSize int32 `form:"pageSize" binding:"required,min=5,max=20"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithError(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithError(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
