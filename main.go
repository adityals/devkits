package main

import (
	"devkits/internal/mitm"
	"devkits/internal/typegen"
	"embed"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	app := application.New(application.Options{
		Name:        "DevKits",
		Description: "Tool kit for development",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	configDir := filepath.Join(home, ".config", "devkits")
	for _, s := range mitm.Services(app.Event, configDir) {
		app.RegisterService(s)
	}

	for _, s := range typegen.Services() {
		app.RegisterService(s)
	}

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "DevKits",
		// Window sized to the golden ratio (1000 / 618 ≈ 1.618).
		Width:            1000,
		Height:           618,
		Mac:              application.MacWindow{},
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})

	// Wails desktop ignores signals, so Ctrl+C would exit without ServiceShutdown and leave the system proxy on.
	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig,
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		<-sig
		// Stop first so a second Ctrl+C still force-kills a hung shutdown.
		signal.Stop(sig)
		app.Quit()
	}()

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
