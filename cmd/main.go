package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/aleszilagyi/prosig-blog/config"
	"github.com/aleszilagyi/prosig-blog/internal/handler"
	log "github.com/aleszilagyi/prosig-blog/internal/logger"
	"github.com/aleszilagyi/prosig-blog/internal/repository"
	"github.com/aleszilagyi/prosig-blog/internal/router"
	"github.com/aleszilagyi/prosig-blog/internal/storage"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer func() {
		if r := recover(); r != nil {
			// Capture stack trace
			stack := debug.Stack()
			logger.Error("unhandled panic recovered",
				zap.Any("recover", r),
				zap.ByteString("stacktrace", stack),
			)
			os.Exit(1)
		}
	}()

	config.LoadConfig()
	logger = log.GetLogger()

	dbConn, err := storage.Connect(config.GetConfigs())
	if err != nil {
		logger.Error("[Setup] failed to start db connection", zap.Error(err))
	}
	defer dbConn.Close()

	repo := repository.NewBlogRepository(dbConn)
	blogHandler := handler.NewBlogHandler(repo)
	server := router.SetupRouter(blogHandler)

	port := config.GetConfigs().AppConfig.Port
	server.Run(fmt.Sprintf(":%d", port))
}
