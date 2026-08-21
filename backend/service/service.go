package service

import (
	"forum/repository"
)

// ServiceContainer 服务容器
type ServiceContainer struct {
	UserRepo                 *repository.UserRepository
	ArticleRepo              *repository.ArticleRepository
	PersonalNotificationRepo *repository.PersonalNotificationRepository
	PermissionGroupRepo      *repository.PermissionGroupRepository
	SystemLogRepo            *repository.SystemLogRepository
}

// Container 全局服务容器
var Container *ServiceContainer

// InitServices 初始化所有服务
func InitServices() {
	Container = &ServiceContainer{
		UserRepo:                 repository.NewUserRepository(),
		ArticleRepo:              repository.NewArticleRepository(),
		PersonalNotificationRepo: repository.NewPersonalNotificationRepository(),
		PermissionGroupRepo:      repository.NewPermissionGroupRepository(),
		SystemLogRepo:            repository.NewSystemLogRepository(),
	}
}
