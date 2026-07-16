package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
)

type ReportService struct{}

var Report = &ReportService{}

func (s *ReportService) CreateReport(reporterID uint, targetType string, targetID uint, reason, description string) error {
	if targetID == 0 {
		return utils.NewError("必须指定举报目标", 400)
	}

	report := models.Report{
		ReporterID:  reporterID,
		TargetType:  targetType,
		TargetID:    targetID,
		Reason:      reason,
		Description: description,
		Status:      "pending",
	}

	if result := database.DB.Create(&report); result.Error != nil {
		return utils.NewError("创建举报失败", 500)
	}
	return nil
}

func (s *ReportService) GetReports(page, pageSize int) ([]models.Report, int, error) {
	var reports []models.Report
	var total int64

	query := database.DB.Model(&models.Report{}).Where("status = ?", "pending")
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Reporter").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&reports).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return reports, totalPages, err
}

func (s *ReportService) GetReport(reportID uint) (*models.Report, error) {
	var report models.Report
	err := database.DB.Preload("Reporter").First(&report, reportID).Error
	if err != nil {
		return nil, utils.NewError("举报不存在", 404)
	}
	return &report, nil
}

func (s *ReportService) HandleReport(reportID uint, status, handleNote string) error {
	var report models.Report
	if result := database.DB.First(&report, reportID); result.Error != nil {
		return utils.NewError("举报不存在", 404)
	}

	report.Status = status
	report.HandleNote = handleNote
	database.DB.Save(&report)

	return nil
}
