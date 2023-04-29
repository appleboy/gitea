package gpt

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/util"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
)

func CreateSimpleChatCompletion(ctx context.Context, client *openai.Client, content string) (resp openai.ChatCompletionResponse, err error) {
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		MaxTokens:   500,
		Temperature: 0.7,
		TopP:        1,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}
	return client.CreateChatCompletion(ctx, req)
}

func NewOpenAIClient(authToken string, proxyUrl string) (*openai.Client, error) {
	c := openai.DefaultConfig(authToken)
	c.HTTPClient.Timeout = 60 * time.Second

	if proxyUrl != "" {
		u, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(u.Scheme, "socks5") {
			dialer, err := proxy.FromURL(u, proxy.Direct)
			if err != nil {
				return nil, err
			}
			contextDialer := dialer.(proxy.ContextDialer)
			c.HTTPClient.Transport = &http.Transport{
				DialContext: contextDialer.DialContext,
			}
		} else {
			c.HTTPClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(u),
			}
		}
	}
	return openai.NewClientWithConfig(c), nil
}

func NewOpenAIClientFromSetting() (*openai.Client, error) {
	authToken := setting.CfgProvider.Section("gpt").Key("openai_auth_token").String()
	proxyUrl := setting.CfgProvider.Section("gpt").Key("proxy").String()
	log.Info("authToken: %s", authToken)
	log.Info("proxyUrl: %s", proxyUrl)
	if authToken == "" {
		return nil, util.ErrNotExist
	}
	return NewOpenAIClient(authToken, proxyUrl)
}
