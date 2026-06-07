package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 创建举报
func CreateReport(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		TargetType  string `json:"target_type" binding:"required"`
		TargetID    uint   `json:"target_id" binding:"required"`
		Reason      string `json:"reason" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	report := models.Report{
		ReporterID:  userID,
		TargetType:  req.TargetType,
		TargetID:    req.TargetID,
		Reason:      req.Reason,
		Description: req.Description,
		Status:      "pending",
	}

	if err := database.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建举报失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "举报成功", "report": report})
}

// 获取举报列表（管理员）
func GetReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.Report{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	var reports []models.Report

	query.Count(&total)
	query.Preload("Reporter").Preload("Handler").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reports)

	c.JSON(http.StatusOK, gin.H{
		"reports":     reports,
		"total":       total,
		"page":        page,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// 处理举报（管理员）
func HandleReport(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var req struct {
		Status     string `json:"status" binding:"required"`
		HandleNote string `json:"handle_note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var report models.Report
	if result := database.DB.First(&report, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "举报不存在"})
		return
	}

	now := time.Now()
	report.Status = req.Status
	report.HandlerID = &userID
	report.HandleNote = req.HandleNote
	report.HandledAt = &now

	if err := database.DB.Save(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "处理失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "处理成功", "report": report})
}

// 获取举报详情（管理员）
func GetReport(c *gin.Context) {
	id := c.Param("id")

	var report models.Report
	if result := database.DB.Preload("Reporter").Preload("Handler").First(&report, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "举报不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}
