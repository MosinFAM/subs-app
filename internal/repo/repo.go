package repo

import "github.com/MosinFAM/subs-app/internal/models"

type Repository interface {
	CreateSubscription(s models.Subscription) (models.Subscription, error)
	ListSubscriptions(userID string) ([]models.Subscription, error)
	SumSubscriptions(filter models.SubscriptionSumRequest) (int, error)
	GetSubscriptionByID(id string) (models.Subscription, error)
	UpdateSubscription(s models.Subscription) (models.Subscription, error)
	DeleteSubscription(id string) error
}
