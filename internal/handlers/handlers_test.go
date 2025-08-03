package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MosinFAM/subs-app/internal/models"
	"github.com/MosinFAM/subs-app/internal/repo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getTestContext(method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, w
}

func getTestContextWithQuery(method, path string, query string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, path+"?"+query, nil)
	c.Request = req

	return c, w
}

func TestHandler_CreateSubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockRepository(ctrl)
	h := &Handler{Repo: mockRepo}

	validSub := models.Subscription{
		ServiceName: "Netflix",
		Price:       1299,
		UserID:      "user-123",
		StartDate:   "01-2024",
	}

	tests := []struct {
		name       string
		reqBody    interface{}
		mockSetup  func()
		wantStatus int
	}{
		{
			name:    "success",
			reqBody: validSub,
			mockSetup: func() {
				mockRepo.EXPECT().CreateSubscription(gomock.AssignableToTypeOf(models.Subscription{})).
					Return(validSub, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "bad request invalid json",
			reqBody:    nil,
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "internal error",
			reqBody: validSub,
			mockSetup: func() {
				mockRepo.EXPECT().CreateSubscription(gomock.AssignableToTypeOf(models.Subscription{})).
					Return(models.Subscription{}, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.reqBody != nil {
				body, _ = json.Marshal(tt.reqBody)
			} else {
				body = []byte(`{invalid json`)
			}

			tt.mockSetup()
			c, w := getTestContext("POST", "/subscriptions", body)
			h.CreateSubscription(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var resp models.Subscription
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, validSub.ServiceName, resp.ServiceName)
			}
		})
	}
}

func TestHandler_ListSubscriptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockRepository(ctrl)
	h := &Handler{Repo: mockRepo}

	subs := []models.Subscription{
		{
			ID:          "sub1",
			ServiceName: "Netflix",
			Price:       1299,
			UserID:      "user-123",
			StartDate:   "01-2024",
		},
	}

	tests := []struct {
		name       string
		query      string
		mockSetup  func()
		wantStatus int
		wantLen    int
	}{
		{
			name:  "success",
			query: "user_id=user-123",
			mockSetup: func() {
				mockRepo.EXPECT().ListSubscriptions("user-123").
					Return(subs, nil)
			},
			wantStatus: http.StatusOK,
			wantLen:    len(subs),
		},
		{
			name:       "bad request missing user_id",
			query:      "",
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "internal error",
			query: "user_id=user-123",
			mockSetup: func() {
				mockRepo.EXPECT().ListSubscriptions("user-123").
					Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			c, w := getTestContextWithQuery("GET", "/subscriptions", tt.query)
			h.ListSubscriptions(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var resp []models.Subscription
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Len(t, resp, tt.wantLen)
			}
		})
	}
}

func TestHandler_GetSubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockRepository(ctrl)
	h := &Handler{Repo: mockRepo}

	sub := models.Subscription{
		ID:          "sub1",
		ServiceName: "Netflix",
		Price:       1299,
		UserID:      "user-123",
		StartDate:   "01-2024",
	}

	tests := []struct {
		name       string
		paramID    string
		mockSetup  func()
		wantStatus int
	}{
		{
			name:    "success",
			paramID: "sub1",
			mockSetup: func() {
				mockRepo.EXPECT().GetSubscriptionByID("sub1").
					Return(sub, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "not found",
			paramID: "missing",
			mockSetup: func() {
				mockRepo.EXPECT().GetSubscriptionByID("missing").
					Return(models.Subscription{}, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			c, w := getTestContext("GET", "/subscriptions/"+tt.paramID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.paramID}}
			h.GetSubscription(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestHandler_UpdateSubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockRepository(ctrl)
	h := &Handler{Repo: mockRepo}

	validSub := models.Subscription{
		ServiceName: "Netflix",
		Price:       1299,
		UserID:      "user-123",
		StartDate:   "01-2024",
	}

	tests := []struct {
		name       string
		paramID    string
		reqBody    interface{}
		mockSetup  func()
		wantStatus int
	}{
		{
			name:    "success",
			paramID: "sub1",
			reqBody: validSub,
			mockSetup: func() {
				mockRepo.EXPECT().UpdateSubscription(gomock.AssignableToTypeOf(models.Subscription{})).
					Return(validSub, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "bad request invalid json",
			paramID:    "sub1",
			reqBody:    nil,
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "internal error",
			paramID: "sub1",
			reqBody: validSub,
			mockSetup: func() {
				mockRepo.EXPECT().UpdateSubscription(gomock.AssignableToTypeOf(models.Subscription{})).
					Return(models.Subscription{}, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.reqBody != nil {
				body, _ = json.Marshal(tt.reqBody)
			} else {
				body = []byte(`{invalid json`)
			}

			tt.mockSetup()
			c, w := getTestContext("PUT", "/subscriptions/"+tt.paramID, body)
			c.Params = gin.Params{{Key: "id", Value: tt.paramID}}
			h.UpdateSubscription(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestHandler_DeleteSubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockRepository(ctrl)
	h := &Handler{Repo: mockRepo}

	tests := []struct {
		name       string
		paramID    string
		mockSetup  func()
		wantStatus int
	}{
		{
			name:    "success",
			paramID: "sub1",
			mockSetup: func() {
				mockRepo.EXPECT().DeleteSubscription("sub1").Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:    "internal error",
			paramID: "sub1",
			mockSetup: func() {
				mockRepo.EXPECT().DeleteSubscription("sub1").Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			c, w := getTestContext("DELETE", "/subscriptions/"+tt.paramID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.paramID}}
			h.DeleteSubscription(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusNoContent {
				assert.Equal(t, 0, w.Body.Len())
			}
		})
	}
}

func TestHandler_SumSubscriptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockRepository(ctrl)
	h := &Handler{Repo: mockRepo}

	tests := []struct {
		name       string
		query      string
		mockSetup  func()
		wantStatus int
		wantTotal  int
	}{
		{
			name:  "success",
			query: "from=01-2024&to=12-2024",
			mockSetup: func() {
				mockRepo.EXPECT().SumSubscriptions(gomock.AssignableToTypeOf(models.SubscriptionSumRequest{})).
					Return(10000, nil)
			},
			wantStatus: http.StatusOK,
			wantTotal:  10000,
		},
		{
			name:       "bad request invalid query",
			query:      "to=12-2024",
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "internal error",
			query: "from=01-2024&to=12-2024",
			mockSetup: func() {
				mockRepo.EXPECT().SumSubscriptions(gomock.AssignableToTypeOf(models.SubscriptionSumRequest{})).
					Return(0, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			c, w := getTestContextWithQuery("GET", "/subscriptions/summary", tt.query)
			h.SumSubscriptions(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var resp map[string]int
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTotal, resp["total"])
			}
		})
	}
}
