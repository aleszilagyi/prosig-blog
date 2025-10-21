package response

type ResponseDataWrapper[T any] struct {
	Data T `json:"data"`
}

type CommentResponse struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_ad"`
}

type PostWithCommentCountResponse struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	CommentCount int    `json:"comment_count"`
}

type PostWithCommentsResponse struct {
	ID        int                `json:"id"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
	Comments  []*CommentResponse `json:"comments"`
}

func WrapResponse(data map[string]interface{}) *ResponseDataWrapper[map[string]interface{}] {
	return &ResponseDataWrapper[map[string]interface{}]{
		Data: data,
	}
}
