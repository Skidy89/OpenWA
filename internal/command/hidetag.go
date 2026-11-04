package command

import (
	"errors"

	"github.com/skidy89/openWA/internal/model"
)

func runHideTag(c *Context) error {
	if c.Chat == nil {
		return errors.New(c.T("no chat open"))
	}
	text := c.Text("text")
	if text == "" {
		return errors.New(c.T("no text provided"))
	}
	if !c.Info.IsGroup {
		c.Note(&Note{Title: c.Input, Text: c.Locale.Text("group only")})
		return nil
	}
	if len(c.Info.Members) == 0 {
		return errors.New(c.T("no members in group"))
	}
	mentions := []string{}
	for _, m := range c.Info.Members {
		mentions = append(mentions, m.ID)
	}
	msg := c.Backend.Send(c.Chat.ID, model.Draft{
		Text:     text,
		Mentions: mentions,
		Reply:    c.Reply,
	})
	if msg == nil {
		return errors.New(c.T("failed to send message"))
	}
	c.Sent(msg)
	return nil
}
