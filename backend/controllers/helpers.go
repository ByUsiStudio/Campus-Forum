package controllers

import (
	"forum/database"
	"forum/models"
	"strconv"
)

func getUserTitles(userID uint) []models.Title {
	var userTitles []models.UserTitle
	database.DB.Where("user_id = ?", userID).Preload("Title").Find(&userTitles)

	var titles []models.Title
	for _, ut := range userTitles {
		if ut.Title.IsActive {
			titles = append(titles, ut.Title)
		}
	}
	return titles
}

func maskAnonymousUser(article *models.Article, isOwner bool) {
	if article.IsAnonymous && !isOwner {
		article.User = models.User{
			ID:          0,
			Username:    "anonymous",
			DisplayName: "匿名用户",
			Avatar:      "",
		}
	}
}

func maskAnonymousUsers(articles *[]models.Article, currentUserID uint) {
	for i := range *articles {
		isOwner := (*articles)[i].UserID == currentUserID
		maskAnonymousUser(&(*articles)[i], isOwner)
	}
}

func parseUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}
