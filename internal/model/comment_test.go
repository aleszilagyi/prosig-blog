package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestComment_String(t *testing.T) {
	created := time.Date(2025, 10, 21, 12, 34, 56, 0, time.UTC)
	comment := Comment{
		ID:        1,
		PostID:    42,
		Content:   "Great post!",
		CreatedAt: created,
	}

	expected := fmt.Sprintf("{id: %d, post_id: %d, created_at: %s, content: %s}",
		comment.ID, comment.PostID, created.Format(time.RFC3339), comment.Content)

	assert.Equal(t, expected, comment.String())
}

func TestComment_ToCommentResponse(t *testing.T) {
	created := time.Date(2025, 10, 21, 12, 34, 56, 0, time.UTC)
	comment := Comment{
		ID:        1,
		Content:   "Great post!",
		CreatedAt: created,
	}

	resp := comment.ToCommentResponse()

	assert.Equal(t, comment.ID, resp.ID)
	assert.Equal(t, comment.Content, resp.Content)
	assert.Equal(t, created.Format(time.RFC3339), resp.CreatedAt)
}
