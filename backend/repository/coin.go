package repository

import (
	"forum/models"
)

type CoinRecordRepository struct {
	*BaseRepository[models.CoinRecord]
}

func NewCoinRecordRepository() *CoinRecordRepository {
	return &CoinRecordRepository{
		BaseRepository: NewBaseRepository[models.CoinRecord](),
	}
}

func (r *CoinRecordRepository) GetUserCoinRecords(userID uint, page, pageSize int) ([]models.CoinRecord, int64, error) {
	var records []models.CoinRecord
	var total int64

	query := r.db.Model(&models.CoinRecord{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&records).Error

	return records, total, err
}

func (r *CoinRecordRepository) GetTotalCoins(userID uint) (int, error) {
	var total int
	err := r.db.Model(&models.CoinRecord{}).Where("user_id = ?", userID).Select("SUM(amount)").Scan(&total).Error
	return total, err
}