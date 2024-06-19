package routes

import (
	"bytes"
	"event-bus/service"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"io"
)

func RegisterLarkCardRoutes(engine *gin.Engine) {
	engine.POST("/lark/card", handleLarkCard)
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
				Command string `json:"command"`
				Params  string `json:"params"`
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

func handleLarkCard(c *gin.Context) {
	var request LarkCardEvent

	body, _ := io.ReadAll(c.Request.Body)
	println(string(body))

	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	if err := c.ShouldBindJSON(&request); err != nil {
		println(err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if request.Type == "url_verification" {
		c.JSON(200, gin.H{"challenge": request.Challenge})
		return
	}

	command := request.Event.Action.Value.Command
	params := request.Event.Action.Value.Params

	if command == "update-be" {
		type UpdateBeParams []string
		var updateBeParams UpdateBeParams
		err := sonic.Unmarshal([]byte(params), &updateBeParams)

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		currentVersion, err := service.GetBeVersion()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		newVersion := updateBeParams[0]
		messageId := request.Event.Context.OpenMessageID

		// update be in background and return immediately
		go func() {
			err := service.UpdateBe(newVersion)
			if err != nil {
				// 更新卡片，提示用户更新失败
				println(newVersion, " 更新失败")
			} else {
				// 更新卡片，提示用户更新成功
				successCard, err := service.BuildUpdateSuccessCard(currentVersion, newVersion)
				if err != nil {
					println(err.Error())
				}
				_, err = service.UpdateCard(successCard, messageId)
				if err != nil {
					println(err.Error())
				}
			}
		}()

		c.JSON(200, gin.H{})
		return
	}

	// 无法识别的命令
	c.JSON(200, gin.H{})
}
