package handlers

import (
	"context"
	"fmt"
	"go_auth_api_gateway/api/http"
	"go_auth_api_gateway/config"

	"go_auth_api_gateway/genproto/auth_service"

	"go_auth_api_gateway/pkg/jwt"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterUser godoc
// @ID register_user
// @Router /register-user [POST]
// @Summary Register User
// @Description Register User
// @Tags User
// @Accept json
// @Produce json
// @Param user body auth_service.CreateUserRequest true "CreateUserRequestBody"
// @Success 201 {object} http.Response{data=auth_service.UserWithAuth} "User data"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) RegisterUser(c *gin.Context) {
	var user auth_service.CreateUserRequest

	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	if len(user.GetPassword()) < 6 {
		h.handleResponse(c, http.BadRequest, "password must be at least 6 characters")
		return
	}

	usr, _ := h.strg.User().GetByUsername(context.Background(), user.GetUsername())
	fmt.Println("err", usr)
	if usr.Id != "" {
		h.handleResponse(c, http.GRPCError, "user already exists")
		return
	}

	resp, err := h.services.UserService().CreateUser(
		c.Request.Context(),
		&user,
	)

	if err != nil {
		h.handleResponse(c, http.GRPCError, err.Error())
		return
	}

	m := map[interface{}]interface{}{
		"sub": resp.Id,
	}
	accessToken, refreshTokenk, err := jwt.GenJWT(m, config.SigningKey)
	if err != nil {
		h.handleResponse(c, http.GRPCError, err.Error())
		return
	}

	h.handleResponse(c, http.Created, &auth_service.UserWithAuth{
		Id:           resp.GetId(),
		Phone:        resp.GetPhone(),
		FirstName:    resp.GetFirstName(),
		LastName:     resp.GetLastName(),
		Username:     resp.GetUsername(),
		AccessToken:  accessToken,
		RefreshToken: refreshTokenk,
	})
}

// LoginUser godoc
// @ID login_user
// @Router /login-user [POST]
// @Summary Login User
// @Description Login User
// @Tags User
// @Accept json
// @Produce json
// @Param username query string true "username"
// @Param password query string true "password"
// @Success 201 {object} http.Response{data=auth_service.GetByCredentialsRequest} "User data"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) LoginUser(c *gin.Context) {

	resp, err := h.services.UserService().GetByCredentials(
		c.Request.Context(),
		&auth_service.GetByCredentialsRequest{
			Username: c.Query("username"),
			Password: c.Query("password"),
		},
	)

	if err != nil {
		h.handleResponse(c, http.GRPCError, err.Error())
		return
	}

	m := map[interface{}]interface{}{
		"sub": resp.Id,
	}
	accessToken, refreshTokenk, err := jwt.GenJWT(m, config.SigningKey)
	if err != nil {
		h.handleResponse(c, http.GRPCError, err.Error())
		return
	}

	h.handleResponse(c, http.Created, &auth_service.UserWithAuth{
		Id:           resp.GetId(),
		Phone:        resp.GetPhone(),
		FirstName:    resp.GetFirstName(),
		LastName:     resp.GetLastName(),
		Username:     resp.GetUsername(),
		AccessToken:  accessToken,
		RefreshToken: refreshTokenk,
	})
}

// UpdateUser godoc
// @ID update_user
// @Security ApiKeyAuth
// @Router /user [PUT]
// @Summary Update User
// @Description Update User
// @Tags User
// @Accept json
// @Produce json
// @Param user body auth_service.UpdateUserRequest true "UpdateUserRequestBody"
// @Success 200 {object} http.Response{data=auth_service.User} "User data"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) UpdateUser(c *gin.Context) {
	var user auth_service.UpdateUserRequest

	token, err := ExtractToken(c.GetHeader("Authorization"))
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	claims, err := ExtractClaims(token, string(config.SigningKey))
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	err = c.ShouldBindJSON(&user)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	user.Id = claims["sub"].(string)

	resp, err := h.services.UserService().UpdateUser(
		c.Request.Context(),
		&user,
	)

	if err != nil {
		h.handleResponse(c, http.GRPCError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, resp)
}

// DeleteUser godoc
// @ID delete_user
// @Security ApiKeyAuth
// @Router /user [DELETE]
// @Summary Delete User
// @Description Get User
// @Tags User
// @Accept json
// @Produce json
// @Success 200  {object} http.Response{data=string} "Succes"
// @Response 400 {object} http.Response{data=string} "Invalid Argument"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) DeleteUser(c *gin.Context) {

	token, err := ExtractToken(c.GetHeader("Authorization"))
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	claims, err := ExtractClaims(token, string(config.SigningKey))
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	userID := claims["sub"].(string)

	_, err = h.services.UserService().DeleteUser(
		c.Request.Context(),
		&auth_service.UserPrimaryKey{
			Id: userID,
		},
	)
	if err != nil {
		// check if the error has a code attribute
		if err1, ok := status.FromError(err); ok {
			// access the error code using the .Code attribute
			errorCode := err1.Code()
			if errorCode == codes.InvalidArgument {
				h.handleResponse(c, http.BadRequest, "User was already deleted")
				return
			}
			h.handleResponse(c, http.GRPCError, err.Error())
			return
		}
		h.handleResponse(c, http.OK, "Deleted successfully")
	}
}

// UpdateUser godoc
// @ID reset_password
// @Security ApiKeyAuth
// @Router /user/reset-password [PUT]
// @Summary Update User
// @Description Reset Password
// @Tags User
// @Accept json
// @Produce json
// @Param password query string true "password"
// @Success 200 {object} http.Response{data=auth_service.LoginResponse} "User data"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) ResetPassword(c *gin.Context) {

	token, err := ExtractToken(c.GetHeader("Authorization"))
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	claims, err := ExtractClaims(token, string(config.SigningKey))
	if err != nil {
		fmt.Println("IN here 1")
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	UserId := claims["sub"].(string)

	_, err = h.services.UserService().ResetPassword(
		c.Request.Context(),
		&auth_service.ResetPasswordRequest{
			UserId:   UserId,
			Password: c.Query("password"),
		},
	)
	if err != nil {
		h.handleResponse(c, http.GRPCError, err.Error())
		return
	}

	if err != nil {
		h.handleResponse(c, http.GRPCError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, "Password updated successfully")
}

// // GetUserList godoc
// // @ID get_user_list
// // @Router /user [GET]
// // @Summary Get User List
// // @Description  Get User List
// // @Tags User
// // @Accept json
// // @Produce json
// // @Param offset query integer false "offset"
// // @Param limit query integer false "limit"
// // @Param search query string false "search"
// // @Param client-platform-id query string false "client-platform-id"
// // @Param client-type-id query string false "client-type-id"
// // @Param project-id query string false "project-id"
// // @Success 200 {object} http.Response{data=auth_service.GetUserListResponse} "GetUserListResponseBody"
// // @Response 400 {object} http.Response{data=string} "Invalid Argument"
// // @Failure 500 {object} http.Response{data=string} "Server Error"
// func (h *Handler) GetUserList(c *gin.Context) {
// 	offset, err := h.getOffsetParam(c)
// 	if err != nil {
// 		h.handleResponse(c, http.InvalidArgument, err.Error())
// 		return
// 	}

// 	limit, err := h.getLimitParam(c)
// 	if err != nil {
// 		h.handleResponse(c, http.InvalidArgument, err.Error())
// 		return
// 	}

// 	resp, err := h.services.UserService().GetUserList(
// 		c.Request.Context(),
// 		&auth_service.GetUserListRequest{
// 			Limit:  int32(limit),
// 			Offset: int32(offset),
// 			Search: c.Query("search"),
// 		},
// 	)

// 	if err != nil {
// 		h.handleResponse(c, http.GRPCError, err.Error())
// 		return
// 	}

// 	h.handleResponse(c, http.OK, resp)
// }

// // GetUserByID godoc
// // @ID get_user_by_id
// // @Router /user/{user-id} [GET]
// // @Summary Get User By ID
// // @Description Get User By ID
// // @Tags User
// // @Accept json
// // @Produce json
// // @Param user-id path string true "user-id"
// // @Success 200 {object} http.Response{data=auth_service.User} "UserBody"
// // @Response 400 {object} http.Response{data=string} "Invalid Argument"
// // @Failure 500 {object} http.Response{data=string} "Server Error"
// func (h *Handler) GetUserByID(c *gin.Context) {
// 	userID := c.Param("user-id")

// 	if !util.IsValidUUID(userID) {
// 		h.handleResponse(c, http.InvalidArgument, "user id is an invalid uuid")
// 		return
// 	}

// 	resp, err := h.services.UserService().GetUserByID(
// 		c.Request.Context(),
// 		&auth_service.UserPrimaryKey{
// 			Id: userID,
// 		},
// 	)

// 	if err != nil {
// 		h.handleResponse(c, http.GRPCError, err.Error())
// 		return
// 	}

// 	h.handleResponse(c, http.OK, resp)
// }
