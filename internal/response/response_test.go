package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapResponse(t *testing.T) {
	t.Run("wraps map successfully", func(t *testing.T) {
		data := map[string]interface{}{
			"key": "value",
			"num": 123,
		}

		resp := WrapResponse(data)

		assert.NotNil(t, resp, "ResponseDataWrapper should not be nil")
		assert.Equal(t, data, resp.Data, "Wrapped data should match input")
		assert.Equal(t, "value", resp.Data["key"])
		assert.Equal(t, 123, resp.Data["num"])
	})

	t.Run("wraps empty map", func(t *testing.T) {
		data := map[string]interface{}{}
		resp := WrapResponse(data)

		assert.NotNil(t, resp)
		assert.Empty(t, resp.Data)
	})

	t.Run("wraps map with posts slice", func(t *testing.T) {
		posts := []*PostWithCommentCountResponse{
			{
				ID:           1,
				Title:        "First Post",
				Content:      "Content 1",
				CreatedAt:    "2025-10-21T00:00:00Z",
				UpdatedAt:    "2025-10-21T00:10:00Z",
				CommentCount: 3,
			},
			{
				ID:           2,
				Title:        "Second Post",
				Content:      "Content 2",
				CreatedAt:    "2025-10-21T01:00:00Z",
				UpdatedAt:    "2025-10-21T01:10:00Z",
				CommentCount: 0,
			},
		}

		data := map[string]interface{}{
			"posts": posts,
		}

		resp := WrapResponse(data)

		assert.NotNil(t, resp)
		assert.Contains(t, resp.Data, "posts")
		assert.Len(t, resp.Data["posts"], 2)

		gotPosts, ok := resp.Data["posts"].([]*PostWithCommentCountResponse)
		assert.True(t, ok, "Expected 'posts' to be of type []*PostWithCommentCountResponse")
		assert.Equal(t, "First Post", gotPosts[0].Title)
		assert.Equal(t, 3, gotPosts[0].CommentCount)
	})
}
