package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MosinFAM/subs-app/internal/models"
	"github.com/google/uuid"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) CreateSubscription(s models.Subscription) (models.Subscription, error) {
	s.ID = uuid.New().String()

	start, err := time.Parse("01-2006", s.StartDate)
	if err != nil {
		return s, err
	}
	var end *time.Time
	if s.EndDate != nil {
		e, err := time.Parse("01-2006", *s.EndDate)
		if err != nil {
			return s, err
		}
		end = &e
	}

	_, err = r.db.Exec(`
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, s.ID, s.ServiceName, s.Price, s.UserID, start, end)

	return s, err
}

func (r *PostgresRepo) ListSubscriptions(userID string) ([]models.Subscription, error) {
	rows, err := r.db.Query(`
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var s models.Subscription
		var start time.Time
		var end *time.Time
		err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &start, &end)
		if err != nil {
			return nil, err
		}
		s.StartDate = start.Format("01-2006")
		if end != nil {
			str := end.Format("01-2006")
			s.EndDate = &str
		}
		subs = append(subs, s)
	}
	return subs, nil
}

func (r *PostgresRepo) SumSubscriptions(filter models.SubscriptionSumRequest) (int, error) {
	from, err := time.Parse("01-2006", filter.From)
	if err != nil {
		return 0, err
	}
	to, err := time.Parse("01-2006", filter.To)
	if err != nil {
		return 0, err
	}

	query := `
		SELECT SUM(price)
		FROM subscriptions
		WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)
	`
	args := []interface{}{to, from}
	idx := 3

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", idx)
		args = append(args, *filter.UserID)
		idx++
	}
	if filter.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name ILIKE $%d", idx)
		args = append(args, "%"+*filter.ServiceName+"%")
	}

	var sum sql.NullInt64
	err = r.db.QueryRow(query, args...).Scan(&sum)
	if err != nil {
		return 0, err
	}
	if !sum.Valid {
		return 0, nil
	}
	return int(sum.Int64), nil
}
