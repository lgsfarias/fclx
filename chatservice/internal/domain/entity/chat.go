package entity

import (
	"errors"

	"github.com/google/uuid"
)

type ChatConfig struct {
	Model            *Model
	Temperature      float32  // 0.0 to 1.0
	TopP             float32  // 0.0 to 1.0 - the cumulative probability for top-k filtering
	N                int      // number of messages to generate
	Stop             []string // stop generating when the given text is generated
	MaxTokens        int      // maximum number of tokens to generate
	PresencePenalty  float32  // -2.0 to 2.0
	FrequencyPenalty float32  // -2.0 to 2.0
}

type Chat struct {
	ID                   string
	UserID               string
	InitialSystemMessage *Message
	Messages             []*Message
	ErasedMessages       []*Message
	Status               string
	TokenUsage           int
	Config               *ChatConfig
}

func NewChat(userID string, initialSystemMessage *Message, config *ChatConfig) (*Chat, error) {
	chat := &Chat{
		ID:                   uuid.New().String(),
		UserID:               userID,
		InitialSystemMessage: initialSystemMessage,
		Status:               "active",
		Config:               config,
		TokenUsage:           0,
	}
	chat.AddMessage(initialSystemMessage)
	if err := chat.Validate(); err != nil {
		return nil, err
	}
	return chat, nil
}

func (c *Chat) Validate() error {
	if c.UserID == "" {
		return errors.New("user_id is empty")
	}
	if c.Status != "active" && c.Status != "ended" {
		return errors.New("invalid status")
	}
	if c.Config.Temperature < 0.0 || c.Config.Temperature > 1.0 {
		return errors.New("invalid temperature")
	}
	// TODO: validate other fields
	return nil
}

func (c *Chat) AddMessage(msg *Message) error {
	if c.Status == "ended" {
		return errors.New("chat is ended")
	}
	for {
		if c.Config.Model.GetMaxTokens() >= msg.GetQtdTokens()+c.TokenUsage {
			c.Messages = append(c.Messages, msg)
			c.RefreshTokenUsage()
			break
		}
		c.ErasedMessages = append(c.ErasedMessages, c.Messages[0])
		c.Messages = c.Messages[1:]
	}
	return nil
}

func (c *Chat) RefreshTokenUsage() {
	c.TokenUsage = 0
	for _, msg := range c.Messages {
		c.TokenUsage += msg.GetQtdTokens()
	}
}

func (c *Chat) GetMessages() []*Message {
	return c.Messages
}

func (c *Chat) CountMessages() int {
	return len(c.Messages)
}

func (c *Chat) EndChat() {
	c.Status = "ended"
}
