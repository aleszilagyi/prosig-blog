package router

import (
	"github.com/aleszilagyi/prosig-blog/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(handler handler.BlogHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/posts", handler.CreateBlogPost)
		api.POST("/posts/:id/comments", handler.AddComment)
		api.GET("/posts", handler.GetAllPostsWithCommentCount)
		api.GET("/posts/:id", handler.GetPostWithComments)
	}

	return r
}
