package handler

import (
	"backend/internal/dto"
	"backend/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserLoginHandler struct {
	userLoginService service.UserLoginInterface
}

func NewUserLoginHandler(userLoginService service.UserLoginInterface) *UserLoginHandler {
	return &UserLoginHandler{userLoginService: userLoginService}
}

func (h *UserLoginHandler) LoginUser(c *gin.Context) {

	var req dto.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userLoginService.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "user": user})
}

func (h *UserLoginHandler) LogoutUser(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	splitToken := strings.Split(authHeader, " ")

	if len(splitToken) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	token := splitToken[1]

	if err := h.userLoginService.LogoutUser(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}
