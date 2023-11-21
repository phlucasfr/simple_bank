package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "picpay_simplificado/db/sqlc"
	"picpay_simplificado/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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
	Username string `json:"username" binding:"required,alphanum"`
	FullName string `json:"full_name" binding:"required"`
	CpfCnpj  string `json:"cpf_cnpj" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required, min=6"`
}

type createUserResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	CpfCnpj  string `json:"cpf_cnpj"`
	Email    string `json:"email"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		FullName:       req.FullName,
		CpfCnpj:        req.CpfCnpj,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		CpfCnpj:  user.CpfCnpj,
		Email:    user.Email,
	}

	ctx.JSON(http.StatusOK, rsp)
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
