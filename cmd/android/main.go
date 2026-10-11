// Command android starts the Android build of OpenWA.
//
// This entry point intentionally uses the demo backend for now. It validates
// the Gio Android lifecycle and the mobile UI without pulling desktop-only
// concerns into the first Android build.
package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"gioui.org/app"

	"github.com/skidy89/openWA/internal/accounts"
	"github.com/skidy89/openWA/internal/model"
	"github.com/skidy89/openWA/internal/ui"
	"github.com/skidy89/openWA/internal/wa"
)

func main() {
	debug.SetGCPercent(50)

	dataDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	dataDir = filepath.Join(dataDir, "WazzapClients")

	if err := os.MkdirAll(dataDir, 0700); err != nil {
		log.Fatal(err)
	}

	list := accounts.Load(dataDir)

	backend, err := wa.Open(list.Path(list.Active), false, true)
	if err != nil {
		log.Fatal(err)
	}

	options := ui.Options{
		Mobile: true,
		Window: []app.Option{
			app.Title("OpenWA"),
		},
		NotifyDir: filepath.Join(dataDir, "notifications"),
		Accounts:  list,
		Open: func(dir string) (model.Backend, error) {
			return wa.Open(dir, false, true)
		},
	}

	go func() {
		if err := ui.Run(backend, options); err != nil {
			log.Print(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}
