package model

import (
	"fmt"
	"time"

	"github.com/aleszilagyi/prosig-blog/internal/response"
)

type Comment struct {
	ID        int
	PostID    int
	Content   string
	CreatedAt time.Time
}

func (c *Comment) String() string {
	createdAt := c.CreatedAt.Format(time.RFC3339)
	return fmt.Sprintf("{id: %d, post_id: %d, created_at: %s, content: %s}",
		c.ID,
		c.PostID,
		createdAt,
		c.Content,
	)
}

func (c *Comment) ToCommentResponse() *response.CommentResponse {
	return &response.CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
	}
}
