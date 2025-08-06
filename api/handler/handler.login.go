package handler

import (
	"github.com/gin-gonic/gin"
	type_api "github.com/mrjvadi/BackendPanelVpn/types/type-api"
	"net/http"
)

// Login godoc
// @Summary      User Login
// @Description  User login with username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      type_api.LoginRequest      true  "Login Request"
// @Success      200      {object}  type_api.BaseResponse[type_api.LoginResponse]  "Login successful"
// @Failure      400      {object}  type_api.BaseResponseError               "Invalid request data"
// @Failure      401      {object}  type_api.BaseResponseError               "Unauthorized access"
// @Failure      500      {object}  type_api.BaseResponseError               "Internal server error"
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req type_api.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, type_api.BaseResponseError{
			Code:    400,
			Message: "Invalid request data",
		})
		return
	}

	res := h.service.AuthService.Login(c, req.Username, req.Password)
	if res.IsError() {
		c.JSON(http.StatusInternalServerError, type_api.BaseResponseError{
			Code:    res.Code.Int(),
			Message: res.Message.String(),
		})

		return
	}

	c.JSON(http.StatusOK, type_api.BaseResponse[type_api.LoginResponse]{
		Code:    200,
		Message: "Login successful",
		Data:    res.Data,
	})

	return
}

// Logout godoc
// @Summary      User Logout
// @Description  Revoke the current user's session token
// @Tags         Auth
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer <token>"
// @Success      200            {object}  type_api.BaseResponse[any]  "Logout successful"
// @Failure      400            {object}  type_api.BaseResponseError    "Bad request"
// @Failure      401            {object}  type_api.BaseResponseError    "Unauthorized"
// @Failure      500            {object}  type_api.BaseResponseError    "Internal server error"
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {

	// get token in gin context
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, type_api.BaseResponseError{
			Code:    http.StatusBadRequest,
			Message: "Token is required",
		})
		return
	}

	res := h.service.AuthService.Logout(token.(string))
	if res.IsError() {
		c.JSON(http.StatusInternalServerError, type_api.BaseResponseError{
			Code:    res.Code.Int(),
			Message: res.Message.String(),
		})
		return
	}
	c.JSON(http.StatusOK, type_api.BaseResponse[any]{
		Code:    res.Code.Int(),
		Message: res.Message.String(),
		Data:    nil,
	})
	return
}
