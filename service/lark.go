package service

import (
	"context"
	"errors"
	"github.com/cbroglie/mustache"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/spf13/viper"
)

func BuildUpdateCard(
	oldVersion string,
	newVersion string,
) (string, error) {
	print("oldVersion: ", oldVersion)
	print("newVersion: ", newVersion)
	// use mustache to render card using template/update-card.mustache
	params := map[string]interface{}{
		"oldVersion": oldVersion,
		"newVersion": newVersion,
	}

	cardStr, err := mustache.RenderFile("template/update-card.mustache", params)
	if err != nil {
		return "", err
	}
	return cardStr, nil
}

type LarkConfig struct {
	AppId     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
	ChatId    string `mapstructure:"chat_id"`
}

func getLarkConfig() (LarkConfig, error) {
	// read config/config.yml using viper
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		return LarkConfig{}, err
	}
	var larkConfig LarkConfig
	err = viper.UnmarshalKey("lark", &larkConfig)
	if err != nil {
		return LarkConfig{}, err
	}
	return larkConfig, nil
}

func SendCardToChat(
	card string,
) (*larkim.CreateMessageResp, error) {
	// get lark config
	larkConfig, err := getLarkConfig()
	if err != nil {
		return nil, err
	}
	client := lark.NewClient(larkConfig.AppId, larkConfig.AppSecret)
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(larkConfig.ChatId).
			Content(card).
			Build()).
		Build()

	resp, err := client.Im.V1.Message.Create(context.Background(), req)
	if err != nil {
		return nil, err
	}
	if !resp.Success() {
		return nil, errors.New(resp.Msg)
	}
	return resp, nil
}
