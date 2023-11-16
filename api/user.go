package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "picpay_simplificado/db/sqlc"
	"time"

	"github.com/gin-gonic/gin"
)

type getUserRequest struct {
	Username string `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type createUserRequest struct {
	Username       string `json:"username" binding:"required"`
	FullName       string `json:"full_name" binding:"required"`
	CpfCnpj        string `json:"cpf_cnpj" binding:"required"`
	Email          string `json:"email" binding:"required"`
	HashedPassword string `json:"hashed_password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		FullName:       req.FullName,
		CpfCnpj:        req.CpfCnpj,
		Email:          req.Email,
		HashedPassword: req.HashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type deleteUserRequest struct {
	Username string `uri:"id" binding:"required"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetUser(ctx, req.Username)

	if err != nil && err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = server.store.DeleteUser(ctx, req.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse("User deleted"))
}

type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req listUsersRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type updateUserRequest struct {
	Username       string       `json:"username" binding:"required"`
	HashedPassword string       `json:"hashed_password"`
	Email          string       `json:"email"`
	IsMerchant     sql.NullBool `json:"is_merchant"`
}

func (req *updateUserRequest) ValidateUpdateUserResquet() error {
	if req.Email == "" && req.HashedPassword == "" && !req.IsMerchant.Valid {
		x := errors.New("no props received")
		return x
	}
	return nil
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := req.ValidateUpdateUserResquet(); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var passwordValue string
	var emailValue string
	var isMerchantValue sql.NullBool

	if req.Email == "" {
		emailValue = user.Email
	} else {
		emailValue = req.Email
	}

	if req.HashedPassword == "" {
		passwordValue = user.HashedPassword
		user.PasswordChangedAt = time.Now()
	} else {
		passwordValue = req.HashedPassword
	}

	if !req.IsMerchant.Valid {
		isMerchantValue = user.IsMerchant
	} else {
		isMerchantValue = req.IsMerchant
	}

	arg := db.UpdateUserParams{
		Username:       req.Username,
		HashedPassword: passwordValue,
		Email:          emailValue,
		IsMerchant:     isMerchantValue,
	}

	user, err = server.store.UpdateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
