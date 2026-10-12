package command

import (
	"errors"
)

func runReplicateMessage(c *Context) error {
	if c.Chat == nil {
		return errors.New("Open a chat first.")
	}
	if c.Reply == nil {
		return errors.New("Reply to a message, then use /replicate.")
	}
	n := busy(c, c.T("Replicating message…"))
	chatID, reply := c.Chat.ID, c.Reply
	text := c.Text("text")
	if text == "" {
		c.note()
		return errors.New("No text to replicate.")
	}
	d := c.Draft(text)
	c.Do(func() func() {
		err := c.Backend.ReplicateMessage(chatID, d, reply)
		return func() {
			n.Busy = false
			if err != nil {
				fail(n, err.Error())
				return
			}
			n.Text = c.T("Message replicated.")
		}
	})
	return nil
}
