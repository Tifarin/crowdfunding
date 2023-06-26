package campaign

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
//tidak bisa langsung digunakan karena kembaliannya bukan slice campaign
func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter {
		ID: campaign.ID,
		UserID: campaign.UserID,
		Name: campaign.Name,
		ShortDescription: campaign.ShortDescription,
		GoalAmount: campaign.GoalAmount,
		CurrrentAmount: campaign.CurrentAmount,
		Slug: campaign.Slug,
		ImageURL: "",
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