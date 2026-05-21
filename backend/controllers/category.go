package controllers

import (
    "forum/database"
    "forum/models"
    "net/http"
    _"strconv"

    "github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
    var input struct {
        Name        string `json:"name" binding:"required"`
        Description string `json:"description"`
        SortOrder   int    `json:"sort_order"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    category := models.Category{
        Name:        input.Name,
        Description: input.Description,
        SortOrder:   input.SortOrder,
    }
    
    if result := database.DB.Create(&category); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "创建分区失败"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message":  "创建成功",
        "category": category,
    })
}

func UpdateCategory(c *gin.Context) {
    id := c.Param("id")
    var category models.Category
    
    if result := database.DB.First(&category, id); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "分区不存在"})
        return
    }
    
    var input struct {
        Name        string `json:"name"`
        Description string `json:"description"`
        SortOrder   int    `json:"sort_order"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if input.Name != "" {
        category.Name = input.Name
    }
    if input.Description != "" {
        category.Description = input.Description
    }
    category.SortOrder = input.SortOrder
    
    database.DB.Save(&category)
    
    c.JSON(http.StatusOK, gin.H{
        "message":  "更新成功",
        "category": category,
    })
}

func DeleteCategory(c *gin.Context) {
    id := c.Param("id")
    
    // 检查是否有文章使用此分区
    var count int64
    database.DB.Model(&models.Article{}).Where("category_id = ?", id).Count(&count)
    if count > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "该分区下还有文章，无法删除"})
        return
    }
    
    database.DB.Delete(&models.Category{}, id)
    c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetCategories(c *gin.Context) {
    var categories []models.Category
    database.DB.Order("sort_order ASC, id ASC").Find(&categories)
    c.JSON(http.StatusOK, gin.H{"categories": categories})
}