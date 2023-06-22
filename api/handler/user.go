package handler

import (
	"api/helper"
	"api/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *userHandler {
	return &userHandler{service}
}
func (h *userHandler) RegisterHandler(c *gin.Context) {
	// tangkap input dari user
	// mapping input dari user ke struct RegisterUserInput
	// struct diatas kita passing sebagai parameter service

	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	NewUser, err := h.service.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := user.FormatUser(NewUser, "tokeninitokeninitoken")
	response := helper.APIResponse("account has been created", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}
