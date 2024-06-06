package service

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"github.com/cbroglie/mustache"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/spf13/viper"
)

type LarkConfig struct {
	AppId     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
	ChatId    string `mapstructure:"chat_id"`
}

var configData *LarkConfig

func LoadConfig(fs embed.FS) error {
	data, err := fs.ReadFile("config/config.yml")
	if err != nil {
		return err
	}
	viper.SetConfigType("yml")
	err = viper.ReadConfig(bytes.NewReader(data))
	if err != nil {
		return err
	}
	err = viper.UnmarshalKey("lark", &configData)
	if err != nil {
		return err
	}
	return nil
}

var templates map[string]string

func LoadTemplates(fs embed.FS) error {
	templates = make(map[string]string)
	files, err := fs.ReadDir("template")
	if err != nil {
		return err
	}
	for _, file := range files {
		fileName := file.Name()
		data, _ := fs.ReadFile("template" + "/" + fileName)
		templates[fileName] = string(data)
	}
	return nil
}

func BuildUpdateCard(
	oldVersion string,
	newVersion string,
) (string, error) {
	// use mustache to render card using template/update-card.mustache
	params := map[string]interface{}{
		"oldVersion": oldVersion,
		"newVersion": newVersion,
	}

	template := templates["update-card.mustache"]

	cardStr, err := mustache.Render(template, params)
	if err != nil {
		return "", err
	}
	return cardStr, nil
}

func getLarkConfig() (*LarkConfig, error) {
	if configData == nil {
		return nil, errors.New("config not loaded")
	}
	return configData, nil
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
