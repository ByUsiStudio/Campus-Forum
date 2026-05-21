package controllers

import (
    "forum/database"
    "forum/models"
    "net/http"
    _"strconv"

    "github.com/gin-gonic/gin"
)

func GetDeletionRequests(c *gin.Context) {
    var requests []models.DeletionRequest
    database.DB.Preload("Article").Preload("User").Where("status = ?", "pending").Order("created_at DESC").Find(&requests)
    
    c.JSON(http.StatusOK, gin.H{
        "requests": requests,
    })
}

func ApproveDeletion(c *gin.Context) {
    id := c.Param("id")
    
    var request models.DeletionRequest
    if result := database.DB.Preload("Article").First(&request, id); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "申请不存在"})
        return
    }
    
    // 删除文章
    database.DB.Delete(&request.Article)
    
    // 更新申请状态
    request.Status = "approved"
    database.DB.Save(&request)
    
    c.JSON(http.StatusOK, gin.H{"message": "已批准删除"})
}

func RejectDeletion(c *gin.Context) {
    id := c.Param("id")
    
    var request models.DeletionRequest
    if result := database.DB.First(&request, id); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "申请不存在"})
        return
    }
    
    request.Status = "rejected"
    database.DB.Save(&request)
    
    c.JSON(http.StatusOK, gin.H{"message": "已拒绝删除"})
}