package controllers

import (
    "forum/database"
    "forum/models"
    "forum/utils"
    "net/http"

    "github.com/gin-gonic/gin"
)

func GetAnnouncement(c *gin.Context) {
    var announcement models.Announcement
    result := database.DB.First(&announcement)
    
    if result.Error != nil {
        c.JSON(http.StatusOK, gin.H{
            "content":      "",
            "content_html": "",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "content":      announcement.Content,
        "content_html": announcement.ContentHTML,
    })
}

func UpdateAnnouncement(c *gin.Context) {
    var input struct {
        Content string `json:"content" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    contentHTML := utils.MarkdownToHTML(input.Content)
    
    var announcement models.Announcement
    result := database.DB.First(&announcement)
    if result.Error != nil {
        announcement.Content = input.Content
        announcement.ContentHTML = contentHTML
        database.DB.Create(&announcement)
    } else {
        announcement.Content = input.Content
        announcement.ContentHTML = contentHTML
        database.DB.Save(&announcement)
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}