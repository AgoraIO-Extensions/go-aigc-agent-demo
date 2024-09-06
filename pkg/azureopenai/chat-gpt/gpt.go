package chat_gpt

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

var chatGPTObj *ChatGPT

func Inst() *ChatGPT {
	return chatGPTObj
}

type Config struct {
	Key      string
	Model    string
	EndPoint string
}

func NewConfig(key, endpoint string) *Config {
	return &Config{
		Key:      key,
		EndPoint: endpoint,
	}
}

type ChatGPT struct {
	conf   *Config
	client *azopenai.Client
}

func InitChatGPT(cfg *Config) error {
	client, err := azopenai.NewClientWithKeyCredential(cfg.EndPoint, azcore.NewKeyCredential(cfg.Key), nil)
	if err != nil {
		return err
	}
	chatGPTObj = &ChatGPT{
		client: client,
		conf:   cfg,
	}

	return nil
}

type Msg struct {
	Role    string
	Content string
}

func (m Msg) conversion() azopenai.ChatRequestMessageClassification {
	switch m.Role {
	case "system":
		return &azopenai.ChatRequestSystemMessage{Content: to.Ptr(m.Content)}
	case "user":
		return &azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(m.Content)}
	case "assistant":
		return &azopenai.ChatRequestAssistantMessage{Content: to.Ptr(m.Content)}
	}
	return nil
}

func (gpt *ChatGPT) StreamAsk(msgs []Msg, modelName string) (azopenai.GetChatCompletionsStreamResponse, error) {
	var messages []azopenai.ChatRequestMessageClassification
	for _, m := range msgs {
		messages = append(messages, m.conversion())
	}

	return gpt.getChatCompletionsStream(messages, modelName)
}

func (gpt *ChatGPT) getChatCompletionsStream(messages []azopenai.ChatRequestMessageClassification, modelName string) (azopenai.GetChatCompletionsStreamResponse, error) {
	return gpt.client.GetChatCompletionsStream(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		N:              to.Ptr[int32](1), // 表示返回的答案个数（对于一个问题而言，可以返回多个答案以供选择）
		DeploymentName: to.Ptr(modelName),
	}, nil)
}
