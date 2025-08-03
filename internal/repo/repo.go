package repo

import "github.com/MosinFAM/subs-app/internal/models"

// go install go.uber.org/mock/mockgen@latest
//
//go:generate mockgen -source=repo.go -destination=repo_mock.go -package=repo Repository
type Repository interface {
	CreateSubscription(s models.Subscription) (models.Subscription, error)
	ListSubscriptions(userID string) ([]models.Subscription, error)
	SumSubscriptions(filter models.SubscriptionSumRequest) (int, error)
	GetSubscriptionByID(id string) (models.Subscription, error)
	UpdateSubscription(s models.Subscription) (models.Subscription, error)
	DeleteSubscription(id string) error
}
