package repository

import (
	"forum/database"
	"gorm.io/gorm"
)

// Repository 通用数据访问接口
type Repository[T any] interface {
	Create(entity *T) error
	Update(entity *T) error
	Delete(id uint) error
	GetByID(id uint) (*T, error)
	GetAll(conditions ...interface{}) ([]T, error)
	GetPage(page, pageSize int, conditions ...interface{}) ([]T, int64, error)
	Count(conditions ...interface{}) (int64, error)
	First(conditions ...interface{}) (*T, error)
}

// BaseRepository 通用数据访问实现
type BaseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository 创建通用 Repository
func NewBaseRepository[T any]() *BaseRepository[T] {
	return &BaseRepository[T]{
		db: database.DB,
	}
}

// Create 创建记录
func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// Update 更新记录
func (r *BaseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete 删除记录
func (r *BaseRepository[T]) Delete(id uint) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}

// GetByID 根据ID获取记录
func (r *BaseRepository[T]) GetByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetAll 获取所有记录
func (r *BaseRepository[T]) GetAll(conditions ...interface{}) ([]T, error) {
	var entities []T
	query := r.db
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.Find(&entities).Error
	return entities, err
}

// GetPage 分页查询
func (r *BaseRepository[T]) GetPage(page, pageSize int, conditions ...interface{}) ([]T, int64, error) {
	var entities []T
	var total int64
	
	query := r.db.Model(new(T))
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	
	query.Count(&total)
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&entities).Error
	
	return entities, total, err
}

// Count 统计记录数
func (r *BaseRepository[T]) Count(conditions ...interface{}) (int64, error) {
	var count int64
	query := r.db.Model(new(T))
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.Count(&count).Error
	return count, err
}

// First 获取第一条记录
func (r *BaseRepository[T]) First(conditions ...interface{}) (*T, error) {
	var entity T
	query := r.db
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// DB 获取原始数据库连接
func (r *BaseRepository[T]) DB() *gorm.DB {
	return r.db
}

// Where 添加查询条件
func (r *BaseRepository[T]) Where(query interface{}, args ...interface{}) *gorm.DB {
	return r.db.Where(query, args...)
}

// Preload 预加载关联
func (r *BaseRepository[T]) Preload(relation string) *gorm.DB {
	return r.db.Preload(relation)
}

// Order 排序
func (r *BaseRepository[T]) Order(value interface{}) *gorm.DB {
	return r.db.Order(value)
}
