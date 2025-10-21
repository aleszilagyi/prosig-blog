package request

import (
	"errors"
	"strings"
	"testing"

	app_err "github.com/aleszilagyi/prosig-blog/internal/error"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateBlogPost(t *testing.T) {
	tests := []struct {
		name      string
		req       *CreateBlogPostRequest
		wantError bool
		errMsg    string
	}{
		{
			name: "valid request",
			req: &CreateBlogPostRequest{
				Title:   "Hello",
				Content: "This is a blog post",
			},
			wantError: false,
		},
		{
			name: "empty content",
			req: &CreateBlogPostRequest{
				Title:   "Hello",
				Content: "",
			},
			wantError: true,
			errMsg:    "post content cannot be empty",
		},
		{
			name: "empty title",
			req: &CreateBlogPostRequest{
				Title:   "",
				Content: "Some content",
			},
			wantError: true,
			errMsg:    "post title cannot be empty",
		},
		{
			name: "empty title and content",
			req: &CreateBlogPostRequest{
				Title:   "",
				Content: "",
			},
			wantError: true,
			errMsg:    "post content cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateBlogPost(tt.req)

			if tt.wantError {
				assert.Error(t, err, "expected an error")
				assert.True(t, errors.Is(err, app_err.ErrInvalidInput), "error should wrap ErrInvalidInput")
				assert.True(t, strings.Contains(err.Error(), tt.errMsg), "error message should contain: %s", tt.errMsg)
			} else {
				assert.NoError(t, err, "expected no error")
			}
		})
	}
}

func TestValidateAddComment(t *testing.T) {
	tests := []struct {
		name      string
		req       *AddCommentRequest
		wantError bool
		errMsg    string
	}{
		{
			name: "valid comment",
			req: &AddCommentRequest{
				Content: "Nice post!",
			},
			wantError: false,
		},
		{
			name: "empty comment",
			req: &AddCommentRequest{
				Content: "",
			},
			wantError: true,
			errMsg:    "comment content cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddComment(tt.req)

			if tt.wantError {
				assert.Error(t, err, "expected an error")
				assert.True(t, errors.Is(err, app_err.ErrInvalidInput), "error should wrap ErrInvalidInput")
				assert.True(t, strings.Contains(err.Error(), tt.errMsg), "error message should contain: %s", tt.errMsg)
			} else {
				assert.NoError(t, err, "expected no error")
			}
		})
	}
}
