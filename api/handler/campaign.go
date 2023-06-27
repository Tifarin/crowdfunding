package handler

import (
	"api/campaign"
	"api/helper"
	"api/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
 	service campaign.Service

 }

 func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
 	return &campaignHandler{campaignService}
}

func (h *campaignHandler) GetCampaigns (c *gin.Context) {
	//query untuk di bagian endpoint user_id=10
	userID, _ := strconv.Atoi(c.Query("user_id"))
	campaigns, err := h.service.GetCampaigns(userID)

	if err != nil {
		response := helper.APIResponse("error to get campaigns", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("list of campaigns", http.StatusOK,"success",campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign (c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("failed to get detail campaign", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("failed to get detail campaign", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("campaign detail", http.StatusOK,"success",campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign (c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage :=  gin.H{"errors": errors}
		response := helper.APIResponse("failed to create campaign", http.StatusUnprocessableEntity,"error",errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("failed to create campaign", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("campaign detail", http.StatusOK,"success",campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign (c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("failed to update campaign", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage :=  gin.H{"errors": errors}
		response := helper.APIResponse("failed to update campaign", http.StatusUnprocessableEntity,"error",errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	response := helper.APIResponse("update campaign", http.StatusOK,"success",campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}