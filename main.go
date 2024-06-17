package main

import (
	"embed"
	"event-bus/service"
	"github.com/gin-gonic/gin"
)

type UpdateBeRequest struct {
	Version string `json:"version"`
}

type GhActionBeReleaseRequest struct {
	Tag string `json:"tag"`
}

type LarkCardEvent struct {
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
	Token     string `json:"token"`
	Schema    string `json:"schema"`
	Header    struct {
		EventId    string `json:"event_id"`
		Token      string `json:"token"`
		CreateTime string `json:"create_time"`
		EventType  string `json:"event_type"`
		TenantKey  string `json:"tenant_key"`
		AppId      string `json:"app_id"`
	} `json:"header"`
	Event struct {
		Operator struct {
			TenantKey string `json:"tenant_key"`
			UserID    string `json:"user_id"`
			OpenID    string `json:"open_id"`
		} `json:"operator"`
		Token  string `json:"token"`
		Action struct {
			Value struct {
				Command string   `json:"command"`
				Params  []string `json:"params"`
			} `json:"value"`
			Tag    string `json:"tag"`
			Option string `json:"option"`
		} `json:"action"`
		Host         string `json:"host"`
		DeliveryType string `json:"delivery_type"`
		Context      struct {
			Url           string `json:"url"`
			PreviewToken  string `json:"preview_token"`
			OpenMessageID string `json:"open_message_id"`
			OpenChatID    string `json:"open_chat_id"`
		} `json:"context"`
	}
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
	r.POST("/lark/card", func(c *gin.Context) {
		var json LarkCardEvent

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if json.Type == "url_verification" {
			c.JSON(200, gin.H{"challenge": json.Challenge})
			return
		}

		print(c.Request.Body)

		command := json.Event.Action.Value.Command
		params := json.Event.Action.Value.Params

		if command == "update-be" {
			version := params[0]
			// update be in background and return immediately
			go func() {
				err := service.UpdateBe(version)
				if err != nil {
					// 更新卡片，提示用户更新失败
					print("更新失败")
				} else {
					// 更新卡片，提示用户更新成功
					print("更新成功")
				}
			}()

			c.JSON(200, gin.H{})
			return
		}

		// 无法识别的命令
		c.JSON(200, gin.H{})
	})

	err := r.Run(":9100")
	if err != nil {
		return
	}
}
