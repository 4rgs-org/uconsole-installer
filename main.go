package main

import (
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	defer recoverPanic()
	setupLog()
	logf("main: starting uconsole-installer")

	if os.Getenv("WAYLAND_DISPLAY") == "" {
		os.Setenv("WAYLAND_DISPLAY", "wayland-1")
	}
	if os.Getenv("XDG_RUNTIME_DIR") == "" {
		os.Setenv("XDG_RUNTIME_DIR", "/run/user/1000")
	}

	gtk.Init(&os.Args)

	screen, err := gdk.ScreenGetDefault()
	if err != nil || screen == nil {
		logf("screen err: %v", err)
		os.Exit(1)
	}

	theme := LoadTheme()
	logf("theme: %s", theme.Name)
	applyStyles(screen, theme)

	pkgs, err := LoadPackages()
	if err != nil {
		logf("LoadPackages err: %v", err)
		showError("No se pudo leer la lista de paquetes.\n\n" + err.Error())
		os.Exit(1)
	}

	win := createInstallerWindow(theme, pkgs)
	if win == nil {
		logf("createInstallerWindow failed")
		os.Exit(1)
	}
	gtk.Main()
	logf("gtk.Main returned")
}
