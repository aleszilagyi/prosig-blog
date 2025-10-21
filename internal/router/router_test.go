package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aleszilagyi/prosig-blog/config"
	app_err "github.com/aleszilagyi/prosig-blog/internal/error"
	"github.com/aleszilagyi/prosig-blog/internal/handler"
	"github.com/aleszilagyi/prosig-blog/internal/repository/mocks"
	"github.com/aleszilagyi/prosig-blog/internal/request"
	"github.com/aleszilagyi/prosig-blog/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	config.LoadConfig()
	code := m.Run()
	os.Exit(code)
}

func TestSetupRouter_CreateBlogPost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBlogRepository(ctrl)
	h := handler.NewBlogHandler(mockRepo)
	r := SetupRouter(h)

	t.Run("success", func(t *testing.T) {
		reqBody := request.CreateBlogPostRequest{
			Title:   "My Title",
			Content: "My Content",
		}
		body, _ := json.Marshal(reqBody)

		mockRepo.EXPECT().
			CreatePost(gomock.Any(), reqBody.Title, reqBody.Content).
			Return(123, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.EqualValues(t, 123, resp["post_id"])
	})

	t.Run("malformed JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer([]byte("{invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := request.CreateBlogPostRequest{
			Title:   "",
			Content: "some content",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("repository error", func(t *testing.T) {
		reqBody := request.CreateBlogPostRequest{
			Title:   "Title",
			Content: "Content",
		}
		body, _ := json.Marshal(reqBody)

		mockRepo.EXPECT().
			CreatePost(gomock.Any(), reqBody.Title, reqBody.Content).
			Return(0, errors.New("db error"))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAddCommentRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBlogRepository(ctrl)
	h := handler.NewBlogHandler(mockRepo)
	r := SetupRouter(h)

	t.Run("success - adds a comment", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"comment_content": "Nice post!",
		}
		body, _ := json.Marshal(reqBody)

		mockRepo.
			EXPECT().
			AddComment(gomock.Any(), 1, "Nice post!").
			Return(10, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/posts/1/comments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Contains(t, resp.Body.String(), `"comment_id":10`)
	})

	t.Run("error - malformed json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/posts/1/comments", bytes.NewBufferString(`invalid_json`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("error - invalid post id", func(t *testing.T) {
		reqBody := map[string]interface{}{"comment_content": "Hello"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/posts/invalid/comments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), `"invalid post id"`)
	})

	t.Run("error - invalid request input", func(t *testing.T) {
		// Missing content should fail validation
		reqBody := map[string]interface{}{"comment_content": ""}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/posts/1/comments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("error - repository fails", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"comment_content": "Good job",
		}
		body, _ := json.Marshal(reqBody)

		mockRepo.
			EXPECT().
			AddComment(gomock.Any(), 1, "Good job").
			Return(0, app_err.ErrInternalServer)

		req := httptest.NewRequest(http.MethodPost, "/api/posts/1/comments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, resp.Body.String(), app_err.ErrInternalServer.Error())
	})
}

func TestGetPostWithComments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBlogRepository(ctrl)
	h := handler.NewBlogHandler(mockRepo)
	r := SetupRouter(h)

	t.Run("GetPostWithComments - success", func(t *testing.T) {
		mockPost := &response.PostWithCommentsResponse{
			ID:        1,
			Title:     "Post 1",
			Content:   "Content 1",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			Comments: []*response.CommentResponse{
				{ID: 1, Content: "Nice", CreatedAt: time.Now().Format(time.RFC3339)},
			},
		}

		mockRepo.
			EXPECT().
			GetPostWithComments(gomock.Any(), 1).
			Return(mockPost, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/posts/1", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var body map[string]response.PostWithCommentsResponse
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)
		assert.Equal(t, mockPost.ID, body["post"].ID)
	})

	t.Run("GetPostWithComments - invalid post id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/posts/abc", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "invalid post id")
	})

	t.Run("GetPostWithComments - repo error", func(t *testing.T) {
		mockRepo.
			EXPECT().
			GetPostWithComments(gomock.Any(), 2).
			Return(nil, app_err.ErrInternalServer)

		req := httptest.NewRequest(http.MethodGet, "/api/posts/2", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, resp.Body.String(), app_err.ErrInternalServer.Error())
	})
}

func TestGetAllPostsWithCommentCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBlogRepository(ctrl)
	h := handler.NewBlogHandler(mockRepo)
	r := SetupRouter(h)

	t.Run("GetAllPostsWithCommentCount - success", func(t *testing.T) {
		mockPosts := []*response.PostWithCommentCountResponse{
			{ID: 1, Title: "Post 1", CommentCount: 2},
			{ID: 2, Title: "Post 2", CommentCount: 0},
		}

		mockRepo.
			EXPECT().
			GetAllPostsWithCommentCount(gomock.Any()).
			Return(mockPosts, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var body map[string][]response.PostWithCommentCountResponse
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)
		assert.Len(t, body["posts"], 2)
	})

	t.Run("GetAllPostsWithCommentCount - repo error", func(t *testing.T) {
		mockRepo.
			EXPECT().
			GetAllPostsWithCommentCount(gomock.Any()).
			Return(nil, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, resp.Body.String(), "internal server error")
	})
}
