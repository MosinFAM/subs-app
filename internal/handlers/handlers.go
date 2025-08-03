package handlers

import (
	"net/http"

	"github.com/MosinFAM/subs-app/internal/models"
	"github.com/MosinFAM/subs-app/internal/repo"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repo repo.Repository
}

// @Summary Create a new subscription
// @Description Create a new subscription for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body models.Subscription true "Subscription data"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Could not create subscription"
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	var s models.Subscription
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	sub, err := h.Repo.CreateSubscription(s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create subscription"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// @Summary List all subscriptions for a user
// @Description Returns all subscriptions for the specified user
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User UUID"
// @Success 200 {array} models.Subscription
// @Failure 400 {object} map[string]string "user_id required"
// @Failure 500 {object} map[string]string "Could not fetch subscriptions"
// @Router /subscriptions [get]
func (h *Handler) ListSubscriptions(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}
	subs, err := h.Repo.ListSubscriptions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch subscriptions"})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// @Summary Get subscription by ID
// @Description Returns the subscription with the specified ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 404 {object} map[string]string "Not found"
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	id := c.Param("id")
	sub, err := h.Repo.GetSubscriptionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// @Summary Update a subscription
// @Description Updates an existing subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param input body models.Subscription true "Updated subscription data"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Update failed"
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	var s models.Subscription
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	s.ID = id
	sub, err := h.Repo.UpdateSubscription(s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// @Summary Delete a subscription
// @Description Deletes the subscription with the specified ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string "Delete failed"
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")
	if err := h.Repo.DeleteSubscription(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Calculate total cost of subscriptions
// @Description Calculates the total subscription cost over a given period, optionally filtered by user ID and service name
// @Tags subscriptions
// @Produce json
// @Param from query string true "Start date in MM-YYYY format"
// @Param to query string true "End date in MM-YYYY format"
// @Param user_id query string false "Filter by user ID"
// @Param service_name query string false "Filter by service name"
// @Success 200 {object} map[string]int "Total cost"
// @Failure 400 {object} map[string]string "Invalid query"
// @Failure 500 {object} map[string]string "Could not calculate total"
// @Router /subscriptions/summary [get]
func (h *Handler) SumSubscriptions(c *gin.Context) {
	var f models.SubscriptionSumRequest
	if err := c.ShouldBindQuery(&f); err != nil || f.From == "" || f.To == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query"})
		return
	}
	sum, err := h.Repo.SumSubscriptions(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate total"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": sum})
}
