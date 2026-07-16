package controllers

import (
	"forum/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLeaderboard(c *gin.Context) {
	leaderboard, err := service.Leaderboard.GetLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
}

func GetWeeklyLeaderboard(c *gin.Context) {
	leaderboard, err := service.Leaderboard.GetWeeklyLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"weekly_leaderboard": leaderboard})
}

func GetMonthlyLeaderboard(c *gin.Context) {
	leaderboard, err := service.Leaderboard.GetMonthlyLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"monthly_leaderboard": leaderboard})
}
