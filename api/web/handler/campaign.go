package handler

import (
	"api/campaign"
	"api/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	userService user.Service
}

func NewCampaignHandler(campaignService campaign.Service, userService user.Service) *campaignHandler {
	return &campaignHandler{campaignService, userService}
}

func (h *campaignHandler) Index (c *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
	}
	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}

func (h *campaignHandler) New (c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	input := campaign.FormCreateCampaignInput{}
	input.Users = users
	c.HTML(http.StatusOK, "campaign_new.html", input)
}

func (h *campaignHandler) Create (c *gin.Context) {
	var input campaign.FormCreateCampaignInput

	err := c.ShouldBind(&input)
	if err != nil {
		users, e := h.userService.GetAll()
		if e != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}
		input.Users = users
		input.Error = err

		c.HTML(http.StatusOK, "campaign_new.html", input)
	}
	user, err := h.userService.GetUserByID(input.UserID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	createCampaignInput := campaign.CreateCampaignInput{}
	createCampaignInput.Name = input.Name
	createCampaignInput.ShortDescription = input.ShortDescription
	createCampaignInput.Description = input.Description
	createCampaignInput.GoalAmount = input.GoalAmount
	createCampaignInput.Perks = input.Perks
	createCampaignInput.User = user

	_, err = h.campaignService.CreateCampaign(createCampaignInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/campaigns")
}