package command

import (
	"errors"

	"github.com/skidy89/openWA/internal/model"
)

const missingProperties = "Missing required options. You must provide both a sender and text to send a fake message."

func runFakeMsg(c *Context) error {
	if c.Chat == nil {
		return errors.New(c.T("no chat open"))
	}
	text := c.Text("text")
	sender := c.Text("sender")
	quoted := c.Text("quoted")
	if text == "" || sender == "" || quoted == "" {
		c.Note(&Note{Title: c.Input, Text: c.T(missingProperties)})
		return nil
	}
	msg := c.Backend.Send(c.Chat.ID, model.Draft{
		Text: text,
		Reply: &model.Message{
			ID:       "3EEBO13891389",
			SenderID: sender + "@s.whatsapp.net",
			Text:     quoted,
			ChatID:   c.Chat.ID,
			Sender:   sender,
		},
	})
	if msg == nil {
		return errors.New(c.T("failed to send message"))
	}
	c.Sent(msg)
	return nil
}
