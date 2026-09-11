package main

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// applyStyles registra el CSS de Omarchy para screen.
// Estilo: card angosto, fondo sólido, borde cuadrado, texto gris claro,
// iconos a la izquierda, sin descripciones, submenús jerárquicos.
func applyStyles(screen *gdk.Screen, theme Theme) {
	if screen == nil {
		return
	}
	provider, err := gtk.CssProviderNew()
	if err != nil {
		logf("CssProviderNew err: %v", err)
		return
	}
	if err := provider.LoadFromData(buildCSS(theme)); err != nil {
		logf("LoadFromData err: %v", err)
		return
	}
	gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func buildCSS(t Theme) string {
	bg := FallbackColor(t.Background, "#1a1b26")
	fg := FallbackColor(t.Foreground, "#a9b1d6")
	accent := FallbackColor(t.Accent, "#7aa2f7")
	muted := FallbackColor(t.Muted, "#414868")

	selectedBg := hexToRGBA(fg, 0.08)
	selectedBorder := hexToRGBA(fg, 0.25)

	var b strings.Builder
	fmt.Fprintf(&b, `
window#menu-window {
    background-color: %s;
}
#card {
    background-color: %s;
    color: %s;
    border: 1px solid %s;
    border-radius: 0;
    padding: 0;
}

#header {
    background-color: %s;
    color: %s;
    font-family: "JetBrainsMono Nerd Font", monospace;
    font-weight: bold;
    padding: 12px 14px 10px 14px;
    border-bottom: 1px solid %s;
}
#header-title {
    color: %s;
    font-size: 15px;
    font-family: "JetBrainsMono Nerd Font", monospace;
}
#header-crumb {
    color: %s;
    font-size: 11px;
    font-family: "JetBrainsMono Nerd Font", monospace;
}

#search {
    background-color: %s;
    padding: 8px 12px 6px 12px;
    border-bottom: 1px solid %s;
}
#search entry {
    background-color: %s;
    color: %s;
    border: 1px solid %s;
    border-radius: 0;
    padding: 6px 10px;
    font-family: "JetBrainsMono Nerd Font", monospace;
    font-size: 13px;
}
#search entry:focus {
    border-color: %s;
}
#search prompt {
    color: %s;
    font-weight: bold;
    padding-right: 6px;
}

#list-scroll {
    background-color: %s;
}
#list {
    background-color: %s;
}
#list list { background-color: %s; }
#list list row { padding: 0; border-radius: 0; }

#row {
    background-color: %s;
    padding: 0;
}
#row-box {
    padding: 8px 12px 8px 6px;
}
#row-icon {
    color: %s;
    font-family: "JetBrainsMono Nerd Font", monospace;
    font-size: 15px;
    min-width: 32px;
}
#row-label {
    color: %s;
    font-family: "JetBrainsMono Nerd Font", monospace;
    font-size: 13px;
    font-weight: normal;
}
#row-chevron {
    color: %s;
    font-family: "JetBrainsMono Nerd Font", monospace;
    font-size: 13px;
    min-width: 14px;
}
#row.is-submenu {
    background-color: %s;
}
#row:selected {
    background-color: %s;
    border-left: 2px solid %s;
}
#row:selected #row-label,
#row:selected #row-icon {
    color: %s;
}
#row:selected #row-chevron {
    color: %s;
}
#row:hover:not(:selected) {
    background-color: %s;
}

#footer {
    background-color: %s;
    color: %s;
    padding: 6px 12px;
    border-top: 1px solid %s;
    font-family: "JetBrainsMono Nerd Font", monospace;
    font-size: 11px;
}
#footer key {
    color: %s;
    background-color: %s;
    padding: 1px 6px;
    font-weight: bold;
}
`,
		bg,
		bg, fg, fg,
		bg, fg, muted,
		fg,
		muted,
		bg, muted,
		bg, fg, muted,
		accent,
		fg,
		bg,
		bg,
		bg,
		bg,
		fg,
		fg,
		muted,
		bg,
		selectedBg, selectedBorder,
		fg,
		accent,
		selectedBg,
		bg, muted, muted,
		accent, muted,
	)
	return b.String()
}
