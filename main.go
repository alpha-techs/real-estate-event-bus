package main

import (
	"embed"
	"event-bus/routes"
	"event-bus/service"
	"github.com/gin-gonic/gin"
)

type UpdateBeRequest struct {
	Version string `json:"version"`
}

type GhActionBeReleaseRequest struct {
	Tag string `json:"tag"`
}

//go:embed config/config.yml
var configFs embed.FS

//go:embed template/*
var templatesFs embed.FS

func main() {
	_ = service.LoadConfig(configFs)

	_ = service.LoadTemplates(templatesFs)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	/**
	 * 更新后端服务
	 */
	r.POST("/update-be", func(c *gin.Context) {
		var json UpdateBeRequest

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Call the UpdateBe function from the service package
		err := service.UpdateBe(
			json.Version,
		)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"version": json.Version,
		})
	})

	/**
	 * 获取后端服务版本
	 */
	r.GET("/get-be-version", func(c *gin.Context) {
		version, err := service.GetBeVersion()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"version": version,
		})
	})

	/**
	 * Github Action触发，后端服务新版本发布
	 */
	r.POST("/gh-action/be-release", func(c *gin.Context) {
		// 获取当前版本
		oldVersion, err := service.GetBeVersion()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		// 获取新版本
		var json GhActionBeReleaseRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		newVersion := json.Tag
		cardContent, err := service.BuildUpdateCard(oldVersion, newVersion)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		_, err = service.SendCardToChat(cardContent)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"oldVersion": oldVersion,
			"newVersion": newVersion,
		})
	})

	/**
	 * 监听Lark卡片点击事件
	 */
	routes.RegisterLarkCardRoutes(r)

	/**
	 * 监听Lark事件
	 */
	routes.RegisterLarkEventRoutes(r)

	err := r.Run(":9200")
	if err != nil {
		return
	}
}
