package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPost_String(t *testing.T) {
	created := time.Date(2025, 10, 21, 12, 0, 0, 0, time.UTC)
	updated := created.Add(time.Hour)
	post := Post{
		ID:        1,
		Title:     "Hello",
		Content:   "This is a post",
		CreatedAt: created,
		UpdatedAt: updated,
		Comments: []Comment{
			{ID: 1, PostID: 1, Content: "Nice!", CreatedAt: created},
			{ID: 2, PostID: 1, Content: "Great post", CreatedAt: created},
		},
	}

	output := post.String()

	// Check that output contains expected fields
	assert.Contains(t, output, "id: 1")
	assert.Contains(t, output, "title: Hello")
	assert.Contains(t, output, created.Format(time.RFC3339))
	assert.Contains(t, output, updated.Format(time.RFC3339))
	assert.Contains(t, output, "This is a post")
	assert.Contains(t, output, "post_id: 1")
	assert.Contains(t, output, "content: Great post")
	assert.Contains(t, output, "content: Nice!")
}

func TestPost_ToPostWithCommentCount(t *testing.T) {
	created := time.Date(2025, 10, 21, 12, 0, 0, 0, time.UTC)
	updated := created.Add(time.Hour)
	post := Post{
		ID:        1,
		Title:     "Hello",
		Content:   "This is a post",
		CreatedAt: created,
		UpdatedAt: updated,
	}

	resp := post.ToPostWithCommentCount(3)

	assert.Equal(t, post.ID, resp.ID)
	assert.Equal(t, post.Title, resp.Title)
	assert.Equal(t, post.Content, resp.Content)
	assert.Equal(t, created.Format(time.RFC3339), resp.CreatedAt)
	assert.Equal(t, updated.Format(time.RFC3339), resp.UpdatedAt)
	assert.Equal(t, 3, resp.CommentCount)
}

func TestPost_ToPostWithComments(t *testing.T) {
	created := time.Date(2025, 10, 21, 12, 0, 0, 0, time.UTC)
	updated := created.Add(time.Hour)
	post := Post{
		ID:        1,
		Title:     "Hello",
		Content:   "This is a post",
		CreatedAt: created,
		UpdatedAt: updated,
		Comments: []Comment{
			{ID: 1, PostID: 1, Content: "Nice!"},
			{ID: 2, PostID: 1, Content: "Great post"},
		},
	}

	resp := post.ToPostWithComments()

	assert.Equal(t, post.ID, resp.ID)
	assert.Equal(t, post.Title, resp.Title)
	assert.Equal(t, post.Content, resp.Content)
	assert.Equal(t, created.Format(time.RFC3339), resp.CreatedAt)
	assert.Equal(t, updated.Format(time.RFC3339), resp.UpdatedAt)
	assert.Len(t, resp.Comments, 2)

	// Check nested comments
	assert.Equal(t, 1, resp.Comments[0].ID)
	assert.Equal(t, "Nice!", resp.Comments[0].Content)
	assert.Equal(t, 2, resp.Comments[1].ID)
	assert.Equal(t, "Great post", resp.Comments[1].Content)
}
