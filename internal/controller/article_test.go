package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/controller"
	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/service/mocks"
)

func TestArticleController(t *testing.T) {
	suite.Run(t, new(ArticleControllerTestSuite))
}

type ArticleControllerTestSuite struct {
	suite.Suite
	mockService *mocks.ArticleService
	controller  *controller.ArticleController
	router      *gin.Engine
}

func (s *ArticleControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *ArticleControllerTestSuite) BeforeTest(suiteName, testName string) {
	s.mockService = new(mocks.ArticleService)
	s.controller = controller.NewArticleController(s.mockService)

	s.router = gin.Default()
	s.controller.RegisterRoutes(s.router)
}

func (s *ArticleControllerTestSuite) TestListPreviews_Success() {
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
	s.mockService.On("ListPreviews", mock.Anything, 50, mock.Anything).Return(previews, "next-token", nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/article-previews?pageSize=50", nil)
	s.router.ServeHTTP(w, req)

	s.Require().Equal(http.StatusOK, w.Code)
	s.mockService.AssertExpectations(s.T())
}
