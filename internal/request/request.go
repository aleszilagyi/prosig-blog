package request

import (
	"errors"

	app_err "github.com/aleszilagyi/prosig-blog/internal/error"
)

type CreateBlogPostRequest struct {
	Title   string `json:"title"`
	Content string `json:"post_content"`
}

type AddCommentRequest struct {
	Content string `json:"comment_content"`
}

func ValidateCreateBlogPost(req *CreateBlogPostRequest) error {
	var err error
	if req.Content == "" {
		return errors.Join(app_err.ErrInvalidInput, errors.New("post content cannot be empty"))
	}

	if req.Title == "" {
		return errors.Join(app_err.ErrInvalidInput, errors.New("post title cannot be empty"))
	}

	return err
}

func ValidateAddComment(req *AddCommentRequest) error {
	if req.Content == "" {
		return errors.Join(app_err.ErrInvalidInput, errors.New("comment content cannot be empty"))
	}

	return nil
}
