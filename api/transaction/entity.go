package transaction

import (
	"api/campaign"
	"api/user"
	"time"

	"github.com/leekchan/accounting"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	Campaign   campaign.Campaign
	User       user.User
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
func (t Transaction) AmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp ", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(t.Amount)
}
