package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/aleszilagyi/prosig-blog/internal/response"
)

type Post struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Comments  []Comment
}

func (p *Post) String() string {
	createdAt := p.CreatedAt.Format(time.RFC3339)
	updatedAt := p.UpdatedAt.Format(time.RFC3339)
	comments := make([]string, len(p.Comments))
	for idx, comment := range p.Comments {
		comments[idx] = comment.String()
	}
	return fmt.Sprintf("{id: %d, title: %s, created_at: %s, updated_at: %s, content: %s, comments:[%s]}",
		p.ID,
		p.Title,
		createdAt,
		updatedAt,
		p.Content,
		strings.Join(comments, ", "),
	)
}

func (p *Post) ToPostWithCommentCount(commentsCount int) *response.PostWithCommentCountResponse {
	return &response.PostWithCommentCountResponse{
		ID:           p.ID,
		Title:        p.Title,
		Content:      p.Content,
		CreatedAt:    p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    p.UpdatedAt.Format(time.RFC3339),
		CommentCount: commentsCount,
	}
}

func (p *Post) ToPostWithComments() *response.PostWithCommentsResponse {
	comments := make([]*response.CommentResponse, len(p.Comments))
	for idx, comment := range p.Comments {
		comments[idx] = comment.ToCommentResponse()
	}
	return &response.PostWithCommentsResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
		Comments:  comments,
	}
}
