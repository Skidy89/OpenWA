package command

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/skidy89/openWA/internal/auto"
)

// errNoAuto is what /schedule and /afk say where nothing can send for
// you (some tests).
var errNoAuto = errors.New("This can't be done here.")

func runSchedule(c *Context) error {
	if c.Auto == nil {
		return errors.New(c.T(errNoAuto.Error()))
	}
	at, err := ParseWhen(c.Text("when"), c.Now)
	if err != nil {
		return errors.New(sentence(errors.New(c.T(err.Error()))))
	}
	text := strings.TrimSpace(c.Text("message"))
	d := c.Draft(text)
	j := auto.Job{Chat: c.Chat.ID, At: at, Text: d.Text, Shown: text, Mentions: d.Mentions,
		MentionAll: d.MentionAll, MentionAdmins: d.MentionAdmins}
	if c.Reply != nil {
		j.Reply = c.Reply.ID
	}
	jobNote(c, c.Auto.Schedule(j))
	return nil
}

// previewSchedule says when the message would go.
func previewSchedule(in *Input) (string, bool) {
	when := in.Text("when")
	if when == "" {
		return "", false
	}
	now := in.Now
	if now.IsZero() {
		now = time.Now()
	}
	t, err := ParseWhen(when, now)
	switch {
	case errors.Is(err, errUnfinished):
		return "", false
	case err != nil:
		return sentence(err), false
	}
	return "Sends " + DescribeWhen(t, now), true
}

// jobNote shows a scheduled message, with buttons to send it now or
// cancel it.
func jobNote(c *Context, j auto.Job) {
	n := &Note{Title: c.Input, Text: c.T("Sends ") + DescribeWhen(j.At, c.Now) + ":\n" + j.Shown}
	gone := func() { n.Text, n.Buttons = c.T("It was sent already:\n")+j.Shown, nil }
	n.Buttons = []Button{
		{Label: c.T("Send now"), Run: func() {
			m := c.Auto.SendNow(j.ID)
			if m == nil {
				gone()
				return
			}
			c.Sent(m)
			n.Text, n.Buttons = c.T("Sent:\n")+j.Shown, nil
		}},
		{Label: c.T("Cancel"), Danger: true, Run: func() {
			if !c.Auto.Cancel(j.ID) {
				gone()
				return
			}
			n.Text, n.Buttons = c.T("Cancelled:\n")+j.Shown, nil
		}},
	}
	c.Note(n)
}

// maxListed is how many scheduled messages /scheduled shows.
const maxListed = 10

func runScheduled(c *Context) error {
	if c.Auto == nil {
		return errors.New(c.T(errNoAuto.Error()))
	}
	jobs := c.Auto.Jobs(c.Chat.ID)
	others := len(c.Auto.Jobs("")) - len(jobs)
	elsewhere := ""
	if others > 0 {
		elsewhere = plural(c.Locale, others, "message") + " " + c.T(map[bool]string{true: "is", false: "are"}[others == 1]) +
			" " + c.T("scheduled in other chats.")
	}
	if len(jobs) == 0 {
		c.Note(&Note{Title: c.Input, Text: strings.TrimSpace(c.T("Nothing is scheduled in this chat. ") + elsewhere)})
		return nil
	}
	for _, j := range jobs[:min(len(jobs), maxListed)] {
		jobNote(c, j)
	}
	if more := len(jobs) - maxListed; more > 0 {
		elsewhere = strings.TrimSpace(c.T("And ") + strconv.Itoa(more) + c.T(" more. ") + elsewhere)
	}
	if elsewhere != "" {
		c.Note(&Note{Title: c.Input, Text: elsewhere})
	}
	return nil
}

func runAFK(c *Context) error {
	if c.Auto == nil {
		return errNoAuto
	}
	c.Auto.SetAway(strings.TrimSpace(c.Text("reason")))
	n := &Note{Title: c.Input, Text: c.T("You're AFK until you send a message. Whoever messages you, or mentions you ") +
		c.T("in a group, gets this reply once:\n") + auto.AwayPlain(c.Auto.Away(), c.Now)}
	n.Buttons = []Button{{Label: c.T("I'm back"), Run: func() {
		n.Buttons = nil
		if c.Auto.Away() == nil {
			n.Text = c.T("You were back already.")
			return
		}
		n.Text = auto.BackText(c.Auto.Back())
	}}}
	c.Note(n)
	return nil
}

// sentence is an error as a sentence: capitalized, with a full stop.
func sentence(err error) string {
	s := err.Error()
	return strings.ToUpper(s[:1]) + s[1:] + "."
}
