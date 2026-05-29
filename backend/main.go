package main

import (
	"encoding/json"
	"forum/controllers"
	"forum/database"
	"forum/middleware"
	"forum/utils"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Config struct {
	Port      string         `json:"port"`
	JWTSecret string         `json:"jwt_secret"`
	WebDAV    WebDAVConfig   `json:"webdav"`
	Database  DatabaseConfig `json:"database"`
}

type WebDAVConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

var config Config

func main() {
	// 加载配置
	loadConfig()

	// 初始化JWT
	controllers.InitJWT(config.JWTSecret)

	// 初始化WebDAV
	utils.InitWebDAV(config.WebDAV.URL, config.WebDAV.Username, config.WebDAV.Password)

	// 初始化数据库
	dsn := config.Database.Username + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + config.Database.Port + ")/" + config.Database.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	database.InitDB(dsn)

	// 自动迁移和初始化
	database.AutoMigrate()
	database.CheckAndInitAdmin()

	// 设置路由
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 设置最大请求体大小为 20G（适用于视频上传）
	r.MaxMultipartMemory = 20480 << 20 // 20G

	// WebDAV代理路由
	r.Any("/proxy/webdav/*path", utils.ProxyWebDAVHandler)

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由
	api := r.Group("/api")
	{
		// 公开路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.POST("/init-admin", controllers.InitAdmin)
			auth.GET("/check-init", controllers.CheckInit)
		}

		// 需要认证的路由
		protected := api.Group("/")
		protected.Use(middleware.Auth(config.JWTSecret))
		{
			protected.GET("/profile", controllers.GetProfile)
			protected.PUT("/profile", controllers.UpdateProfile)
			protected.POST("/upload/avatar", controllers.UploadAvatar)
			protected.POST("/upload/image", controllers.UploadImage)
			protected.POST("/upload/video", controllers.UploadVideo)

			// 文章相关
			protected.POST("/articles", controllers.CreateArticle)
			protected.PUT("/articles/:id", controllers.UpdateArticle)
			protected.DELETE("/articles/:id", controllers.DeleteArticle)
			protected.POST("/articles/:id/like", controllers.LikeArticle)
			protected.DELETE("/articles/:id/like", controllers.UnlikeArticle)

			// 评论相关
			protected.POST("/articles/:id/comments", controllers.CreateComment)
			protected.DELETE("/comments/:id", controllers.DeleteComment)
			protected.POST("/comments/:id/like", controllers.LikeComment)
			protected.DELETE("/comments/:id/like", controllers.UnlikeComment)

			// 分区相关
			protected.POST("/categories", middleware.AdminOnly(), controllers.CreateCategory)
			protected.PUT("/categories/:id", middleware.AdminOnly(), controllers.UpdateCategory)
			protected.DELETE("/categories/:id", middleware.AdminOnly(), controllers.DeleteCategory)

			// 侧边栏配置
			protected.PUT("/sidebar-config", middleware.AdminOnly(), controllers.UpdateSidebarConfig)

			// 删除审核
			protected.GET("/deletion-requests", middleware.AdminOnly(), controllers.GetDeletionRequests)
			protected.POST("/deletion-requests/:id/approve", middleware.AdminOnly(), controllers.ApproveDeletion)
			protected.POST("/deletion-requests/:id/reject", middleware.AdminOnly(), controllers.RejectDeletion)

			// 公告
			protected.PUT("/announcement", middleware.AdminOnly(), controllers.UpdateAnnouncement)

			// 网站配置
			protected.PUT("/site-config", middleware.AdminOnly(), controllers.UpdateSiteConfig)
			protected.POST("/site-config/test-smtp", middleware.AdminOnly(), controllers.TestSMTPConfig)

			// 通知相关
			protected.GET("/notifications", controllers.GetNotifications)
			protected.GET("/notifications/unread-count", controllers.GetUnreadCount)
			protected.POST("/notifications/:id/read", controllers.MarkNotificationRead)
			protected.POST("/notifications/read-all", controllers.MarkAllNotificationsRead)
			protected.POST("/notifications", middleware.AdminOnly(), controllers.CreateNotification)
			protected.GET("/notifications/admin", middleware.AdminOnly(), controllers.GetAdminNotifications)
			protected.DELETE("/notifications/:id", middleware.AdminOnly(), controllers.DeleteNotification)

			// 评论回复通知
			protected.GET("/comment-reply-notifications", controllers.GetCommentReplyNotifications)
			protected.POST("/comment-reply-notifications/:id/read", controllers.MarkCommentReplyNotificationRead)
			protected.POST("/comment-reply-notifications/read-all", controllers.MarkAllCommentReplyNotificationsRead)

			// 关注相关
			protected.POST("/follow/:id", controllers.FollowUser)
			protected.DELETE("/follow/:id", controllers.UnfollowUser)
			protected.GET("/following", controllers.GetFollowingList)
			protected.GET("/followers", controllers.GetFollowerList)
			protected.GET("/follow/status/:id", controllers.CheckFollowStatus)
			protected.GET("/mutual", controllers.GetMutualFriends)

			// 收藏相关
			protected.POST("/articles/:id/favorite", controllers.AddFavorite)
			protected.DELETE("/articles/:id/favorite", controllers.RemoveFavorite)
			protected.GET("/favorites", controllers.GetFavorites)
			protected.GET("/articles/:id/favorite/check", controllers.CheckFavorite)

			// 聊天相关
			protected.GET("/chat/messages/:id", controllers.GetChatMessages)
			protected.POST("/chat/messages/:id", controllers.GetChatMessages) // POST方式获取聊天记录
			protected.GET("/chat/sessions", controllers.GetChatSessions)
			protected.POST("/chat/send", controllers.SendChatMessage)
			protected.GET("/chat/unread-count", controllers.GetChatUnreadCount)

			// 后台管理
			protected.GET("/admin/statistics", middleware.AdminOnly(), controllers.GetStatistics)
			protected.GET("/admin/users", middleware.AdminOnly(), controllers.GetAllUsers)
			protected.PUT("/admin/users/:id/role", middleware.AdminOnly(), controllers.UpdateUserRole)
			protected.POST("/admin/users/:id/ban", middleware.AdminOnly(), controllers.BanUser)
			protected.POST("/admin/users/:id/unban", middleware.AdminOnly(), controllers.UnbanUser)
			protected.DELETE("/admin/users/:id", middleware.AdminOnly(), controllers.DeleteUser)
			protected.GET("/admin/articles", middleware.AdminOnly(), controllers.GetAllArticles)
			protected.PUT("/admin/articles/:id/status", middleware.AdminOnly(), controllers.UpdateArticleStatus)
			protected.GET("/admin/comments", middleware.AdminOnly(), controllers.GetAllComments)
			protected.DELETE("/admin/comments/:id", middleware.AdminOnly(), controllers.DeleteCommentAdmin)

			// 头衔管理
			protected.GET("/titles", controllers.GetAllTitles)
			protected.POST("/titles", middleware.AdminOnly(), controllers.CreateTitle)
			protected.PUT("/titles/:id", middleware.AdminOnly(), controllers.UpdateTitle)
			protected.DELETE("/titles/:id", middleware.AdminOnly(), controllers.DeleteTitle)
			protected.POST("/titles/grant", middleware.AdminOnly(), controllers.GrantTitle)
			protected.POST("/titles/revoke", middleware.AdminOnly(), controllers.RevokeTitle)
			protected.GET("/users/:id/titles", controllers.GetUserTitles)
		}

		// 公开路由
		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/:id", controllers.GetArticle)
		api.GET("/categories", controllers.GetCategories)
		api.GET("/announcement", controllers.GetAnnouncement)
		api.GET("/sidebar-config", controllers.GetSidebarConfig)
		api.GET("/site-config", controllers.GetSiteConfig)

		// 用户公开信息
		api.GET("/users/:id", controllers.GetUserByID)
		api.GET("/users/:id/articles", controllers.GetUserArticles)

		// 密码重置
		api.POST("/password/reset-code", controllers.SendResetCode)
		api.POST("/password/reset", controllers.ResetPassword)
	}

	// WebSocket路由
	r.GET("/ws/chat", controllers.WebSocketHandler)

	utils.Section("服务启动")
	utils.Info("服务器地址: http://0.0.0.0:%s", config.Port)
	utils.Info("API 文档: http://0.0.0.0:%s/swagger/index.html", config.Port)
	utils.Section("")

	r.Run(":" + config.Port)
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		utils.Error("无法加载配置文件: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		utils.Error("解析配置文件失败: %v", err)
		os.Exit(1)
	}
	utils.Success("配置文件加载成功")
}
