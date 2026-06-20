package controllers

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserStatistics 获取用户统计数据
func GetUserStatistics(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var userStats models.UserStatistics
	if err := database.DB.Where("user_id = ?", userID).First(&userStats).Error; err != nil {
		// 如果用户没有统计记录，创建初始统计
		userStats = models.UserStatistics{
			UserID: userID,
		}
		database.DB.Create(&userStats)
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    userStats,
	})
}

// GetDailyStatistics 获取每日统计数据
func GetDailyStatistics(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	var dailyStats []models.DailyStatistics
	database.DB.Where("date >= ? AND date <= ?", startDate, endDate).
		Order("date desc").
		Find(&dailyStats)

	c.JSON(200, gin.H{
		"success": true,
		"data":    dailyStats,
	})
}

// GetSystemOverview 获取系统概览数据
func GetSystemOverview(c *gin.Context) {
	var overview models.SystemOverview
	if err := database.DB.First(&overview).Error; err != nil {
		// 如果没有概览记录，创建初始数据
		overview = models.SystemOverview{}
		database.DB.Create(&overview)
	}

	// 更新实时数据
	updateSystemOverview(&overview)

	c.JSON(200, gin.H{
		"success": true,
		"data":    overview,
	})
}

// updateSystemOverview 更新系统概览数据
func updateSystemOverview(overview *models.SystemOverview) {
	var count int64

	// 统计总用户数
	database.DB.Model(&models.User{}).Count(&count)
	overview.TotalUsers = int(count)

	// 统计总文章数
	database.DB.Model(&models.Article{}).Where("status != ?", "deleted").Count(&count)
	overview.TotalArticles = int(count)

	// 统计总评论数
	database.DB.Model(&models.Comment{}).Count(&count)
	overview.TotalComments = int(count)

	// 统计总分类数
	database.DB.Model(&models.Category{}).Count(&count)
	overview.TotalCategories = int(count)

	// 统计当前在线用户数
	database.DB.Model(&models.User{}).Where("online_status = ?", "online").Count(&count)
	overview.OnlineUsers = int(count)

	// 统计今日活跃用户数
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.UserActivity{}).Where("date = ?", today).Count(&count)
	overview.TodayActiveUsers = int(count)

	overview.UpdatedAt = time.Now()
	database.DB.Save(overview)
}

// GetArticleStatistics 获取文章统计数据
func GetArticleStatistics(c *gin.Context) {
	articleID := c.Param("id")

	var articleStats models.ArticleStatistics
	if err := database.DB.Where("article_id = ?", articleID).First(&articleStats).Error; err != nil {
		// 如果文章没有统计记录，创建初始统计
		articleStats = models.ArticleStatistics{
			ArticleID: utils.Atouint(articleID),
		}
		database.DB.Create(&articleStats)
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    articleStats,
	})
}

// GetUserActivity 获取用户活跃度数据
func GetUserActivity(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	var activities []models.UserActivity
	database.DB.Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Order("date desc").
		Find(&activities)

	c.JSON(200, gin.H{
		"success": true,
		"data":    activities,
	})
}

// UpdateUserActivity 更新用户活跃度
func UpdateUserActivity(userID uint, activityType string) error {
	today := time.Now().Format("2006-01-02")

	var activity models.UserActivity
	if err := database.DB.Where("user_id = ? AND date = ?", userID, today).First(&activity).Error; err != nil {
		// 创建新的活跃度记录
		activity = models.UserActivity{
			UserID: userID,
			Date:   today,
		}
		database.DB.Create(&activity)
	}

	// 根据活动类型更新计数
	switch activityType {
	case "login":
		activity.LoginCount++
	case "post":
		activity.PostCount++
	case "comment":
		activity.CommentCount++
	case "like":
		activity.LikeCount++
	case "view":
		activity.ViewCount++
	}

	// 计算活跃度分数
	activity.ActiveScore = calculateActiveScore(activity)

	return database.DB.Save(&activity).Error
}

// calculateActiveScore 计算活跃度分数
func calculateActiveScore(activity models.UserActivity) float64 {
	// 使用加权公式计算活跃度分数
	loginWeight := 1.0
	postWeight := 5.0
	commentWeight := 2.0
	likeWeight := 0.5
	viewWeight := 0.1

	return float64(activity.LoginCount)*loginWeight +
		float64(activity.PostCount)*postWeight +
		float64(activity.CommentCount)*commentWeight +
		float64(activity.LikeCount)*likeWeight +
		float64(activity.ViewCount)*viewWeight
}

// UpdateDailyStatistics 更新每日统计数据
func UpdateDailyStatistics() error {
	today := time.Now().Format("2006-01-02")

	var dailyStats models.DailyStatistics
	if err := database.DB.Where("date = ?", today).First(&dailyStats).Error; err != nil {
		// 创建新的每日统计记录
		dailyStats = models.DailyStatistics{
			Date: today,
		}
		database.DB.Create(&dailyStats)
	}

	// 统计新增用户数
	var newUsers int64
	database.DB.Model(&models.User{}).Where("DATE(created_at) = ?", today).Count(&newUsers)
	dailyStats.NewUsers = int(newUsers)

	// 统计活跃用户数
	var activeUsers int64
	database.DB.Model(&models.UserActivity{}).Where("date = ?", today).Count(&activeUsers)
	dailyStats.ActiveUsers = int(activeUsers)

	// 统计新增文章数
	var newArticles int64
	database.DB.Model(&models.Article{}).Where("DATE(created_at) = ? AND status != ?", today, "deleted").Count(&newArticles)
	dailyStats.NewArticles = int(newArticles)

	// 统计新增评论数
	var newComments int64
	database.DB.Model(&models.Comment{}).Where("DATE(created_at) = ?", today).Count(&newComments)
	dailyStats.NewComments = int(newComments)

	// 统计总浏览量
	var totalViews int64
	database.DB.Model(&models.ViewHistory{}).Where("DATE(created_at) = ?", today).Count(&totalViews)
	dailyStats.TotalViews = int(totalViews)

	// 统计总点赞数
	var totalLikes int64
	database.DB.Model(&models.Like{}).Where("DATE(created_at) = ?", today).Count(&totalLikes)
	dailyStats.TotalLikes = int(totalLikes)

	// 统计总分享数
	var totalShares int64
	database.DB.Model(&models.Article{}).Where("DATE(updated_at) = ? AND share_count > 0", today).Count(&totalShares)
	dailyStats.TotalShares = int(totalShares)

	// 统计总签到数
	var totalSignIns int64
	database.DB.Model(&models.SignInRecord{}).Where("DATE(sign_in_at) = ?", today).Count(&totalSignIns)
	dailyStats.TotalSignIns = int(totalSignIns)

	// 统计峰值在线人数
	var peakOnline int64
	database.DB.Model(&models.User{}).Where("online_status = ?", "online").Count(&peakOnline)
	dailyStats.PeakOnline = int(peakOnline)

	return database.DB.Save(&dailyStats).Error
}

// GetStatisticsDashboard 获取统计仪表板数据（管理员）
func GetStatisticsDashboard(c *gin.Context) {
	// 获取系统概览
	var overview models.SystemOverview
	if err := database.DB.First(&overview).Error; err != nil {
		overview = models.SystemOverview{}
		database.DB.Create(&overview)
	}
	updateSystemOverview(&overview)

	// 获取最近7天的每日统计
	var recentStats []models.DailyStatistics
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	database.DB.Where("date >= ? AND date <= ?", startDate, endDate).
		Order("date asc").
		Find(&recentStats)

	// 获取热门文章统计
	var hotArticles []models.Article
	database.DB.Where("status = ?", "published").
		Order("view_count desc").
		Limit(10).
		Find(&hotArticles)

	// 获取活跃用户统计
	var activeUsers []models.User
	database.DB.Where("online_status = ?", "online").
		Order("last_active_at desc").
		Limit(10).
		Find(&activeUsers)

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"overview":     overview,
			"recent_stats": recentStats,
			"hot_articles": hotArticles,
			"active_users": activeUsers,
		},
	})
}
