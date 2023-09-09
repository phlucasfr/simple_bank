package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "picpay_simplificado/db/sqlc"

	"github.com/gin-gonic/gin"
)

type getUserRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Id)

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
	FullName string `json:"full_name" binding:"required"`
	CpfCnpj  string `json:"cpf_cnpj" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		FullName: req.FullName,
		CpfCnpj:  req.CpfCnpj,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type deleteUserRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetUser(ctx, req.Id)

	if err != nil && err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = server.store.DeleteUser(ctx, req.Id)

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
	ID         int64        `json:"id" binding:"required,min=1"`
	Password   string       `json:"password"`
	Email      string       `json:"email"`
	IsMerchant sql.NullBool `json:"is_merchant"`
}

func (req *updateUserRequest) ValidateUpdateUserResquet() error {
	if req.Email == "" && req.Password == "" && !req.IsMerchant.Valid {
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

	user, err := server.store.GetUser(ctx, req.ID)

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

	if req.Password == "" {
		passwordValue = user.Password
	} else {
		passwordValue = req.Password
	}

	if !req.IsMerchant.Valid {
		isMerchantValue = user.IsMerchant
	} else {
		isMerchantValue = req.IsMerchant
	}

	arg := db.UpdateUserParams{
		ID:         req.ID,
		Password:   passwordValue,
		Email:      emailValue,
		IsMerchant: isMerchantValue,
	}

	user, err = server.store.UpdateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
