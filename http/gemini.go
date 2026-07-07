package http

import (
	"bytes"
	"context"
	"fmt"

	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"

	gogpt "github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
)

// RequestGeminiRobot 请求 Gemini 机器人
func RequestGeminiRobot(msg string, svcCtx *svc.ServiceContext) (string, error) {
	// 使用 OpenAI 库配置 Gemini 的 API Token 
	cfg := gogpt.DefaultConfig(svcCtx.Config.Gemini.APIToken)
	
	// Gemini 的 OpenAI 兼容 BaseURL 通常是 "https://generativelanguage.googleapis.com/v1beta/openai/"
	// 或者通过你的中转/代理地址从 svcCtx.Config.Gemini.APIUrl 传入
	cfg.BaseURL = svcCtx.Config.Gemini.APIUrl
	
	c := gogpt.NewClientWithConfig(cfg)
	ctx := context.Background()
	msgs := ""
	
	prompt := svcCtx.Config.Gemini.Prompt
	if svcCtx.Config.Gemini.Limit {
		prompt += fmt.Sprintf(" 尽可能的在%v个字内回答", svcCtx.Config.DanmuLen)
	}
	
	req := gogpt.ChatCompletionRequest{
		// 这里的 Model 可以从配置读取，例如 "gemini-2.5-flash" 或 "gemini-1.5-flash"
		Model: svcCtx.Config.Gemini.Model, 
		Messages: []gogpt.ChatCompletionMessage{
			{
				Role:    gogpt.ChatMessageRoleSystem, // Gemini 习惯将 prompt 作为 System 角色
				Content: prompt,
			},
			{
				Role:    gogpt.ChatMessageRoleUser,
				Content: msg,
			},
		},
	}
	
	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	
	logx.Infof("Gemini 本次开销：%v tokens", resp.Usage.TotalTokens)
	for _, v := range resp.Choices {
		data := []byte(v.Message.Content)
		if bytes.HasPrefix(data, []byte{239, 188, 159}) {
			data = bytes.TrimPrefix(data, []byte{239, 188, 159})
		}
		data = bytes.ReplaceAll(data, []byte{10, 10}, []byte{})
		msgs += string(data)
	}
	return msgs, nil
}
