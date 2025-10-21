package handler

import (
	"errors"
	"net/http"
	"strconv"

	app_err "github.com/aleszilagyi/prosig-blog/internal/error"
	log "github.com/aleszilagyi/prosig-blog/internal/logger"
	"github.com/aleszilagyi/prosig-blog/internal/repository"
	"github.com/aleszilagyi/prosig-blog/internal/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BlogHandler interface {
	CreateBlogPost(ctx *gin.Context)
	AddComment(ctx *gin.Context)
	GetPostWithComments(ctx *gin.Context)
	GetAllPostsWithCommentCount(ctx *gin.Context)
}

type blogHandler struct {
	repo repository.BlogRepository
}

func NewBlogHandler(repo repository.BlogRepository) BlogHandler {
	return &blogHandler{
		repo: repo,
	}
}

func (b *blogHandler) CreateBlogPost(ctx *gin.Context) {
	logger := log.GetLogger()
	req := &request.CreateBlogPostRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Error("[HandlerCreateBlogPost] malformed json", zap.Error(err),
			zap.Int("http_status", http.StatusBadRequest),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "malformed json",
		})
		return
	}

	if err := request.ValidateCreateBlogPost(req); err != nil {
		status, msg := defineHTTPErrorStatus(err)
		ctx.JSON(status, gin.H{
			"error": msg,
		})
		logger.Error("[HandlerCreateBlogPost] invalid request input", zap.Error(err),
			zap.Int("http_status", status),
		)
		return
	}

	logger = logger.With(zap.String("title", req.Title))
	postID, err := b.repo.CreatePost(ctx.Request.Context(), req.Title, req.Content)
	if err != nil {
		status, msg := defineHTTPErrorStatus(err)
		ctx.JSON(status, gin.H{
			"error": msg,
		})
		logger.Error("[HandlerCreateBlogPost] failed to create blog post", zap.Error(err),
			zap.Int("http_status", status),
		)
		return
	}

	data := map[string]interface{}{
		"post_id": postID,
	}

	ctx.JSON(http.StatusCreated, data)
}

func (b *blogHandler) AddComment(ctx *gin.Context) {
	logger := log.GetLogger()
	req := &request.AddCommentRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Error("[HandlerCreateBlogPost] malformed json", zap.Error(err),
			zap.Int("http_status", http.StatusBadRequest),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "malformed json",
		})
		return
	}

	postID, err := getPostIDFromParams(ctx)
	if err != nil {
		status := http.StatusBadRequest
		logger.Error("[HandlerAddComment] invalid post id", zap.Error(err),
			zap.Int("http_status", status),
		)
		ctx.JSON(status, gin.H{
			"error": "invalid post id",
		})
		return
	}

	logger = logger.With(zap.Int("post_id", postID))
	if err := request.ValidateAddComment(req); err != nil {
		status, msg := defineHTTPErrorStatus(err)
		logger.Error("[HandlerAddComment] invalid request input", zap.Error(err),
			zap.Int("http_status", status),
		)
		ctx.JSON(status, gin.H{
			"error": msg,
		})
		return
	}

	commentID, err := b.repo.AddComment(ctx.Request.Context(), postID, req.Content)
	if err != nil {
		status, msg := defineHTTPErrorStatus(err)
		logger.Error("[HandlerAddComment] failed to create comment", zap.Error(err),
			zap.Int("http_status", status),
		)
		ctx.JSON(status, gin.H{
			"error": msg,
		})
		return
	}

	data := map[string]interface{}{
		"comment_id": commentID,
	}

	ctx.JSON(http.StatusCreated, data)
}

func (b *blogHandler) GetPostWithComments(ctx *gin.Context) {
	logger := log.GetLogger()
	postID, err := getPostIDFromParams(ctx)
	if err != nil {
		status := http.StatusBadRequest
		logger.Error("[HandlerGetPostWithComments] invalid post id", zap.Error(err),
			zap.Int("http_status", status),
		)
		ctx.JSON(status, gin.H{
			"error": "invalid post id",
		})
		return
	}

	logger = logger.With(zap.Int("post_id", postID))
	post, err := b.repo.GetPostWithComments(ctx.Request.Context(), postID)
	if err != nil {
		status, msg := defineHTTPErrorStatus(err)
		logger.Error("[HandlerGetPostWithComments] failed to get post with comments", zap.Error(err),
			zap.Int("http_status", status),
		)
		ctx.JSON(status, gin.H{
			"error": msg,
		})
		return
	}

	data := map[string]interface{}{
		"post": post,
	}

	ctx.JSON(http.StatusOK, data)
}

func (b *blogHandler) GetAllPostsWithCommentCount(ctx *gin.Context) {
	logger := log.GetLogger()
	posts, err := b.repo.GetAllPostsWithCommentCount(ctx.Request.Context())
	if err != nil {
		logger.Error("[HandlerGetAllPostsWithCommentCount] failed to get all posts", zap.Error(err))
		status, msg := defineHTTPErrorStatus(err)
		ctx.JSON(status, gin.H{
			"error": msg,
		})
		return
	}

	data := map[string]interface{}{
		"posts": posts,
	}

	ctx.JSON(http.StatusOK, data)
}

func defineHTTPErrorStatus(err error) (httpStatus int, message string) {
	switch {
	case errors.Is(err, app_err.ErrInvalidInput):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, app_err.ErrNotFound):
		return http.StatusNotFound, err.Error()
	default:
		return http.StatusInternalServerError, app_err.ErrInternalServer.Error()
	}
}

func getPostIDFromParams(ctx *gin.Context) (int, error) {
	paramID := ctx.Param("id")
	postID, err := strconv.Atoi(paramID)
	if err != nil {
		return 0, err
	}
	return postID, nil
}
