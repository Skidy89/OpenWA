package wa

import (
	"testing"

	"github.com/skidy89/openWA/internal/model"
)

func TestPreparePrivateDraftWithoutReply(t *testing.T) {
	b := &Backend{}
	d := model.Draft{Text: "hello"}

	got := b.preparePrivateDraft(d)
	if got.Text != d.Text || got.Reply != nil {
		t.Fatalf("draft changed unexpectedly: %#v", got)
	}
}
