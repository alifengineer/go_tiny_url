package handlers

import (
	"encoding/json"
	"go_auth_api_gateway/api/http"
	"go_auth_api_gateway/genproto/auth_service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthMiddleware(c *gin.Context) {
	var result auth_service.HasAccessModel
	if ok := h.hasAccess(c, &result); !ok {
		c.Abort()
		return
	}

	c.Set("auth", &result)

	c.Next()
}

func (h *Handler) hasAccess(c *gin.Context, result *auth_service.HasAccessModel) bool {

	bearerToken := c.GetHeader("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) != 2 || strArr[0] != "Bearer" {
		h.handleResponse(c, http.Forbidden, "token error: wrong format")
		return false
	}
	accessToken := strArr[1]

	claims, err := ExtractClaims(accessToken, h.cfg.SecretKey)

	if err != nil {
		h.handleResponse(c, http.Forbidden, "no access")
		return false
	}

	var tm time.Time
	switch exp := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(exp), 0)
	case json.Number:
		v, _ := exp.Int64()
		tm = time.Unix(v, 0)
	}

	if tm.Unix() < time.Now().Unix() {
		h.handleResponse(c, http.Forbidden, "token expired")
		return false
	}

	userId := claims["sub"]

	_, err = h.services.UserService().GetUserByID(c.Request.Context(), &auth_service.UserPrimaryKey{Id: userId.(string)})
	if err != nil {
		h.handleResponse(c, http.Forbidden, err.Error())
		return false
	}

	result.UserId = userId.(string)

	return true
}
