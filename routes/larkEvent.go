package routes

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
)

func RegisterLarkEventRoutes(engine *gin.Engine) {
	engine.POST("/lark/event", handleLarkEvent)
}

type LarkEventRequest struct {
	Schema    string `json:"schema"`
	Challenge string `json:"challenge"`
	Header    struct {
		EventId    string `json:"event_id"`
		EventType  string `json:"event_type"`
		CreateTime int64  `json:"create_time"`
		Token      string `json:"token"`
		AppId      string `json:"app_id"`
		TenantKey  string `json:"tenant_key"`
	} `json:"header"`
	Event struct {
		Operator struct {
			OperatorName string `json:"operator_name"`
			OperatorId   struct {
				UnionId string `json:"union_id"`
				UserId  string `json:"user_id"`
				OpenId  string `json:"open_id"`
			} `json:"operator_id"`
		} `json:"operator"`
		EventKey string `json:"event_key"`
	} `json:"event"`
}

func handleLarkEvent(c *gin.Context) {
	var request LarkEventRequest

	body, _ := io.ReadAll(c.Request.Body)
	println(string(body))

	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	err := c.ShouldBindJSON(&request)
	if err != nil {
		println(err.Error())
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if request.Header.EventType == "url_verification" {
		c.JSON(200, gin.H{
			"challenge": request.Challenge,
		})
		return
	}

	c.JSON(200, gin.H{})
}
