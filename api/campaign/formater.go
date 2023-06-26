package campaign

import "strings"

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrrentAmount   int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

type CampaignDetailFormatter struct {
	ID               int                   `json:"id"`
	Name             string                `json:"name"`
	ShortDescription string                `json:"short_description"`
	Description      string                `json:"description"`
	ImageURL         string                `json:"image_url"`
	GoalAmount       int                   `json:"goal_amount"`
	CurrrentAmount   int                   `json:"current_amount"`
	UserID           int                   `json:"user_id"`
	Slug             string                `json:"slug"`
	Perks            []string              `json:"perks"`
	User             CampaignUserFormatter `json:"user"`
	Images []CampaignImageFormatter `json:"images"`
}

type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignImageFormatter struct {
	ImageURL string `json:"image_url"`
	IsPrimary bool `json:"is_primary"`
}
// tidak bisa langsung digunakan karena kembaliannya bukan slice campaign
func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		GoalAmount:       campaign.GoalAmount,
		CurrrentAmount:   campaign.CurrentAmount,
		Slug:             campaign.Slug,
		ImageURL:         "",
	}
	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	if len(campaigns) == 0 {
		return []CampaignFormatter{}
	}
	var campaignsFormatter []CampaignFormatter

	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}
	return campaignsFormatter
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	campaignDetailFormatter := CampaignDetailFormatter{}
	campaignDetailFormatter.ID = campaign.ID
	campaignDetailFormatter.Name = campaign.Name
	campaignDetailFormatter.ShortDescription = campaign.ShortDescription
	campaignDetailFormatter.Description = campaign.Description
	campaignDetailFormatter.GoalAmount = campaign.GoalAmount
	campaignDetailFormatter.CurrrentAmount = campaign.CurrentAmount
	campaignDetailFormatter.Slug = campaign.Slug
	campaignDetailFormatter.UserID = campaign.UserID
	campaignDetailFormatter.ImageURL = ""
	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}
	var perks []string
	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk))
	}
	campaignDetailFormatter.Perks = perks

	user := campaign.User
	campaignUserFormatter := CampaignUserFormatter{}
	campaignUserFormatter.Name = user.Name
	campaignUserFormatter.ImageURL = user.AvatarFileName

	campaignDetailFormatter.User = campaignUserFormatter
	
	images := []CampaignImageFormatter{}

	for _, image := range campaign.CampaignImages {
		caampaignImageFormatter:= CampaignImageFormatter{}
		caampaignImageFormatter.ImageURL = image.FileName
		isPrimary :=  false
		if image.IsPrimary == 1 {
			isPrimary = true
		}
		caampaignImageFormatter.IsPrimary = isPrimary

		images = append(images, caampaignImageFormatter)
	}
	campaignDetailFormatter.Images = images

	return campaignDetailFormatter
}

