package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/service/mocks"
)

func TestArticleController_ListPreviews(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("successful request with default page size", func(t *testing.T) {
		// Setup
		mockService := new(mocks.ArticleService)
		controller := NewArticleController(mockService)

		// Mock data
		articleID := uuid.New()
		authorID := uuid.New()
		now := time.Now()
		previews := []model.ArticlePreview{
			{
				ID:                articleID,
				Slug:              "test-article",
				Title:             "Test Article",
				Description:       "Test Description",
				AuthorID:          authorID,
				AuthorDisplayName: "Test Author",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		}

		// Mock expectations
		mockService.On("ListPreviews", mock.Anything, 50, mock.Anything).Return(previews, "next-token", nil)

		// Create test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/articles?pageSize=50", nil)

		// Execute
		controller.ListPreviews(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		// TODO: Parse and verify response body
	})

	t.Run("successful request with custom page size", func(t *testing.T) {
		// Setup
		mockService := new(mocks.ArticleService)
		controller := NewArticleController(mockService)

		// Mock data
		articleID := uuid.New()
		authorID := uuid.New()
		now := time.Now()
		previews := []model.ArticlePreview{
			{
				ID:                articleID,
				Slug:              "test-article",
				Title:             "Test Article",
				Description:       "Test Description",
				AuthorID:          authorID,
				AuthorDisplayName: "Test Author",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		}

		// Mock expectations
		mockService.On("ListPreviews", mock.Anything, 20, mock.Anything).Return(previews, "next-token", nil)

		// Create test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/articles?pageSize=20", nil)

		// Execute
		controller.ListPreviews(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid page size", func(t *testing.T) {
		// Setup
		mockService := new(mocks.ArticleService)
		controller := NewArticleController(mockService)

		// Create test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/articles?pageSize=invalid", nil)

		// Execute
		controller.ListPreviews(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "ListPreviews")
	})

	t.Run("service error", func(t *testing.T) {
		// Setup
		mockService := new(mocks.ArticleService)
		controller := NewArticleController(mockService)

		// Mock expectations
		mockService.On("ListPreviews", mock.Anything, 50, mock.Anything).Return(
			[]model.ArticlePreview{}, "", errors.New("service error"))

		// Create test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/articles?pageSize=50", nil)

		// Execute
		controller.ListPreviews(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
