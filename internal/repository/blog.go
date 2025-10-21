//go:generate mockgen -destination ./mocks/mock_blog.go -package mocks github.com/aleszilagyi/prosig-blog/internal/repository BlogRepository
package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	app_err "github.com/aleszilagyi/prosig-blog/internal/error"
	log "github.com/aleszilagyi/prosig-blog/internal/logger"
	"github.com/aleszilagyi/prosig-blog/internal/model"
	"github.com/aleszilagyi/prosig-blog/internal/response"
	"go.uber.org/zap"
)

const (
	queryAllPostsWithCommentCount = `
		SELECT b.id, b.title, b.created_at, b.updated_at, COUNT(c.id) AS comment_count
		FROM blog_posts b
		LEFT JOIN comments c ON c.blog_post_id = b.id
		GROUP BY b.id
		ORDER BY b.created_at DESC
	`

	queryPostWithComments = `
		SELECT 
			b.id AS post_id,
			b.title,
			b.content,
			b.created_at AS post_created_at,
			b.updated_at AS post_updated_at,
			c.id AS comment_id,
			c.content AS comment_content,
			c.created_at AS comment_created_at
		FROM blog_posts b
		LEFT JOIN comments c ON c.blog_post_id = b.id
		WHERE b.id = $1
		ORDER BY c.created_at DESC
	`

	queryCreatePost = `
		INSERT INTO blog_posts (title, content)
		VALUES ($1, $2)
		RETURNING id
	`

	queryAddComment = `
		INSERT INTO comments (blog_post_id, content)
		VALUES ($1, $2)
		RETURNING id
	`
)

type BlogRepository interface {
	GetAllPostsWithCommentCount(ctx context.Context) ([]*response.PostWithCommentCountResponse, error)
	GetPostWithComments(ctx context.Context, id int) (*response.PostWithCommentsResponse, error)
	CreatePost(ctx context.Context, title, content string) (int, error)
	AddComment(ctx context.Context, blogPostID int, content string) (int, error)
}

type blogRepository struct {
	db *sql.DB
}

func NewBlogRepository(db *sql.DB) BlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) GetAllPostsWithCommentCount(ctx context.Context) ([]*response.PostWithCommentCountResponse, error) {
	logger := log.GetLogger()

	rows, err := r.db.QueryContext(ctx, queryAllPostsWithCommentCount)
	if err != nil {
		logger.Error("[RepoGetAllPostsWithCommentCount] failed to query all posts", zap.Error(err))
		return nil, errors.Join(app_err.ErrInternalServer, err)
	}
	defer rows.Close()

	var posts []*response.PostWithCommentCountResponse
	for rows.Next() {
		var id int
		var title string
		var createdAt, updatedAt time.Time
		var count int
		if err := rows.Scan(&id, &title, &createdAt, &updatedAt, &count); err != nil {
			logger.Error("[RepoGetAllPostsWithCommentCount] failed to scan post", zap.Error(err))
			// Return all available posts, do not block
			continue
		}

		post := &model.Post{
			ID:        id,
			Title:     title,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		posts = append(posts, post.ToPostWithCommentCount(count))
	}

	if err := rows.Err(); err != nil {
		logger.Error("[RepoGetAllPostsWithCommentCount] row iteration error", zap.Error(err))
		return nil, errors.Join(app_err.ErrInternalServer, err)
	}

	return posts, nil
}

func (r *blogRepository) GetPostWithComments(ctx context.Context, requestPostID int) (*response.PostWithCommentsResponse, error) {
	logger := log.GetLogger().With(zap.Int("post_id", requestPostID))
	rows, err := r.db.QueryContext(ctx, queryPostWithComments, requestPostID)
	if err != nil {
		logger.Error("[RepoGetPostWithComments] failed to query post", zap.Error(err))
		return nil, errors.Join(app_err.ErrInternalServer, err)
	}
	defer rows.Close()

	var post *model.Post
	for rows.Next() {
		var (
			postID         *int
			title          *string
			content        *string
			postCreatedAt  *time.Time
			postUpdatedAt  *time.Time
			commentID      *int
			commentContent *string
			commentCreated *time.Time
		)

		err := rows.Scan(
			&postID,
			&title,
			&content,
			&postCreatedAt,
			&postUpdatedAt,
			&commentID,
			&commentContent,
			&commentCreated,
		)
		if err != nil {
			logger.Error("[RepoGetPostWithComments] failed to scan post", zap.Error(err))
			return nil, errors.Join(app_err.ErrInternalServer, err)
		}

		// Initialize post only once
		if post == nil {
			post = &model.Post{
				ID:        *postID,
				Title:     *title,
				Content:   *content,
				CreatedAt: *postCreatedAt,
				UpdatedAt: *postUpdatedAt,
			}
		}

		// Add comments if present
		if commentID != nil && commentContent != nil {
			comment := model.Comment{
				ID:        *commentID,
				PostID:    *postID,
				Content:   *commentContent,
				CreatedAt: *commentCreated,
			}
			post.Comments = append(post.Comments, comment)
		}
	}

	if err := rows.Err(); err != nil {
		logger.Error("[RepoGetPostWithComments] row iteration error", zap.Error(err))
		return nil, errors.Join(app_err.ErrInternalServer, err)
	}

	if post == nil {
		logger.Info("[RepoGetPostWithComments] could not find the post")
		return nil, app_err.ErrNotFound
	}

	return post.ToPostWithComments(), nil
}

func (r *blogRepository) CreatePost(ctx context.Context, title, content string) (int, error) {
	logger := log.GetLogger().With(zap.String("post_title", title))
	var id int
	err := r.db.QueryRowContext(ctx, queryCreatePost, title, content).Scan(&id)
	if err != nil {
		logger.Error("[RepoCreatePost] could not persist the post", zap.Error(err))
		return 0, errors.Join(app_err.ErrInternalServer, err)
	}
	logger.Info("[RepoCreatePost] post created", zap.Int("post_id", id))
	return id, nil
}

func (r *blogRepository) AddComment(ctx context.Context, blogPostID int, content string) (int, error) {
	logger := log.GetLogger().With(zap.Int("post_id", blogPostID))
	var id int
	err := r.db.QueryRowContext(ctx, queryAddComment, blogPostID, content).Scan(&id)
	if err != nil {
		logger.Error("[RepoAddComment] could not add the comment to blog post", zap.Error(err))
		return 0, errors.Join(app_err.ErrInternalServer, err)
	}
	logger.Info("[RepoAddComment] comment added to blog post", zap.Int("comment_id", id))
	return id, nil
}
