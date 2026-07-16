package service

import (
	"forum/repository"
)

type ServiceContainer struct {
	UserRepo                 *repository.UserRepository
	ArticleRepo              *repository.ArticleRepository
	CommentRepo              *repository.CommentRepository
	CommentLikeRepo          *repository.CommentLikeRepository
	CategoryRepo             *repository.CategoryRepository
	LikeRepo                 *repository.LikeRepository
	FavoriteRepo             *repository.FavoriteRepository
	FriendRepo               *repository.FriendRepository
	FriendRequestRepo        *repository.FriendRequestRepository
	NotificationRepo         *repository.NotificationRepository
	UserNotificationRepo     *repository.UserNotificationRepository
	PersonalNotificationRepo *repository.PersonalNotificationRepository
	TopicRepo                *repository.TopicRepository
	TopicFollowRepo          *repository.TopicFollowRepository
	ArticleTopicRepo         *repository.ArticleTopicRepository
	SignInRecordRepo         *repository.SignInRecordRepository
	SignInConfigRepo         *repository.SignInConfigRepository
	ReportRepo               *repository.ReportRepository
	TitleRepo                *repository.TitleRepository
	UserTitleRepo            *repository.UserTitleRepository
	UserLevelRepo            *repository.UserLevelRepository
	ExperienceRecordRepo     *repository.ExperienceRecordRepository
	LevelConfigRepo          *repository.LevelConfigRepository
	AchievementRepo          *repository.AchievementRepository
	UserAchievementRepo      *repository.UserAchievementRepository
	UserStatisticsRepo       *repository.UserStatisticsRepository
	DailyStatisticsRepo      *repository.DailyStatisticsRepository
	ArticleStatisticsRepo    *repository.ArticleStatisticsRepository
	SystemOverviewRepo       *repository.SystemOverviewRepository
	DeletionRequestRepo      *repository.DeletionRequestRepository
	AnnouncementRepo         *repository.AnnouncementRepository
	SiteConfigRepo           *repository.SiteConfigRepository
	SidebarConfigRepo        *repository.SidebarConfigRepository
	LeaderboardRepo          *repository.LeaderboardRepository
	UserBadgeRepo            *repository.UserBadgeRepository
	CoinRecordRepo           *repository.CoinRecordRepository
	PermissionGroupRepo      *repository.PermissionGroupRepository
	SystemLogRepo            *repository.SystemLogRepository
}

var Container *ServiceContainer

func InitServices() {
	Container = &ServiceContainer{
		UserRepo:                 repository.NewUserRepository(),
		ArticleRepo:              repository.NewArticleRepository(),
		CommentRepo:              repository.NewCommentRepository(),
		CommentLikeRepo:          repository.NewCommentLikeRepository(),
		CategoryRepo:             repository.NewCategoryRepository(),
		LikeRepo:                 repository.NewLikeRepository(),
		FavoriteRepo:             repository.NewFavoriteRepository(),
		FriendRepo:               repository.NewFriendRepository(),
		FriendRequestRepo:        repository.NewFriendRequestRepository(),
		NotificationRepo:         repository.NewNotificationRepository(),
		UserNotificationRepo:     repository.NewUserNotificationRepository(),
		PersonalNotificationRepo: repository.NewPersonalNotificationRepository(),
		TopicRepo:                repository.NewTopicRepository(),
		TopicFollowRepo:          repository.NewTopicFollowRepository(),
		ArticleTopicRepo:         repository.NewArticleTopicRepository(),
		SignInRecordRepo:         repository.NewSignInRecordRepository(),
		SignInConfigRepo:         repository.NewSignInConfigRepository(),
		ReportRepo:               repository.NewReportRepository(),
		TitleRepo:                repository.NewTitleRepository(),
		UserTitleRepo:            repository.NewUserTitleRepository(),
		UserLevelRepo:            repository.NewUserLevelRepository(),
		ExperienceRecordRepo:     repository.NewExperienceRecordRepository(),
		LevelConfigRepo:          repository.NewLevelConfigRepository(),
		AchievementRepo:          repository.NewAchievementRepository(),
		UserAchievementRepo:      repository.NewUserAchievementRepository(),
		UserStatisticsRepo:       repository.NewUserStatisticsRepository(),
		DailyStatisticsRepo:      repository.NewDailyStatisticsRepository(),
		ArticleStatisticsRepo:    repository.NewArticleStatisticsRepository(),
		SystemOverviewRepo:       repository.NewSystemOverviewRepository(),
		DeletionRequestRepo:      repository.NewDeletionRequestRepository(),
		AnnouncementRepo:         repository.NewAnnouncementRepository(),
		SiteConfigRepo:           repository.NewSiteConfigRepository(),
		SidebarConfigRepo:        repository.NewSidebarConfigRepository(),
		LeaderboardRepo:          repository.NewLeaderboardRepository(),
		UserBadgeRepo:            repository.NewUserBadgeRepository(),
		CoinRecordRepo:           repository.NewCoinRecordRepository(),
		PermissionGroupRepo:      repository.NewPermissionGroupRepository(),
		SystemLogRepo:            repository.NewSystemLogRepository(),
	}
}
