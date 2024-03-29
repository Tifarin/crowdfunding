package transaction

import (
	"api/campaign"
	"api/payment"
	"errors"
	"strconv"
)

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessPayment(input TransactionNotificationInput) error
	GetAllTransaction() ([]Transaction, error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService payment.Service
}

func NewService(repository Repository, campaignRepository campaign.Repository,paymentService payment.Service) *service {
	return &service{repository, campaignRepository,paymentService}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}
	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("not an owner of the campaign")
	}
	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (s *service) GetTransactionByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"

	newTransaction, err :=  s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}
	paymentTransaction := payment.Transaction {
		ID: newTransaction.ID,
		Amount: newTransaction.Amount,
	}
	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction,nil
	}
	newTransaction.PaymentURL = paymentURL
	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction,nil
	}
	return newTransaction, nil
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID)
	transaction, err := s.repository.GetByID(transaction_id)
	if err != nil {
		return err
	}
	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" || input.TransactionStatus == "success" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updateTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updateTransaction.CampaignID)
	if err != nil {
		return  err
	}
	if updateTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updateTransaction.Amount
		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}
	return nil

}

func (s *service) GetAllTransaction() ([]Transaction, error) {
	transactions, err := s.repository.FindAll()
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}