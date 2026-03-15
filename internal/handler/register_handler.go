package handler 

import (
	"backend/internal/dto"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRegisterHandler struct {
	userRegisterService service.UserRegisterInterface
}

func NewUserRegisterHandler(userRegisterService service.UserRegisterInterface) *UserRegisterHandler {
	return &UserRegisterHandler{userRegisterService: userRegisterService}
}

func (h *UserRegisterHandler) RegisterUser(c *gin.Context) {

	var req dto.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRegisterService.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user": user})
}