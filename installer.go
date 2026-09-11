package main

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// maxResults acota las filas construidas por búsqueda. Sin esto, una query
// vacía intentaría crear 13k widgets GTK y congelaría la UI.
const maxResults = 200

type installerUI struct {
	theme Theme
	all   []*Package

	win     *gtk.Window
	list    *gtk.ListBox
	entry   *gtk.SearchEntry
	title   *gtk.Label
	status  *gtk.Label

	// visible[i] corresponde a la fila i del ListBox. No se puede usar un
	// map con clave *gtk.ListBoxRow: gotk3 crea un wrapper Go nuevo en cada
	// GetRowAtIndex/GetSelectedRow y el lookup fallaría.
	visible []*Package
}

func createInstallerWindow(theme Theme, pkgs []*Package) *gtk.Window {
	ui := &installerUI{theme: theme, all: pkgs}
	ui.build()
	return ui.win
}

func (ui *installerUI) build() {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		logf("WindowNew err: %v", err)
		return
	}
	win.SetTitle("uConsole installer")
	win.SetDefaultSize(420, 560)
	win.SetResizable(false)
	win.SetDecorated(false)
	win.SetKeepAbove(true)
	win.SetSkipTaskbarHint(true)
	win.SetSkipPagerHint(true)
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DIALOG)
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetName("menu-window")
	ui.win = win

	card, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	card.SetName("card")
	win.Add(card)

	// header
	header, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	header.SetName("header")
	ui.title, _ = gtk.LabelNew("uConsole installer")
	ui.title.SetName("header-title")
	ui.title.SetHAlign(gtk.ALIGN_START)
	ui.status, _ = gtk.LabelNew("")
	ui.status.SetName("header-crumb")
	ui.status.SetHAlign(gtk.ALIGN_START)
	header.Add(ui.title)
	header.Add(ui.status)
	card.Add(header)

	// search
	searchBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	searchBox.SetName("search")
	prompt, _ := gtk.LabelNew("❯")
	prompt.SetName("prompt")
	searchBox.Add(prompt)
	ui.entry, _ = gtk.SearchEntryNew()
	ui.entry.SetPlaceholderText("buscar paquete")
	searchBox.Add(ui.entry)
	card.Add(searchBox)

	// list
	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	scroll.SetName("list-scroll")
	scroll.SetSizeRequest(420, 400)
	ui.list, _ = gtk.ListBoxNew()
	ui.list.SetName("list")
	ui.list.SetSelectionMode(gtk.SELECTION_SINGLE)
	scroll.Add(ui.list)
	card.Add(scroll)

	card.Add(ui.buildFooter())

	win.Connect("destroy", func() {
		logf("window destroy")
		gtk.MainQuit()
	})
	ui.entry.Connect("search-changed", func() { ui.refresh() })
	ui.entry.Connect("activate", func(e *gtk.SearchEntry) { ui.activateSelected() })

	win.Connect("key-press-event", func(w *gtk.Window, ev *gdk.Event) bool {
		key := gdk.EventKeyNewFromEvent(ev)
		if key == nil {
			return false
		}
		if key.KeyVal() == gdk.KEY_Escape {
			w.Close()
			return true
		}
		return false
	})

	ui.entry.Connect("key-press-event", func(e *gtk.SearchEntry, ev *gdk.Event) bool {
		key := gdk.EventKeyNewFromEvent(ev)
		if key == nil {
			return false
		}
		switch key.KeyVal() {
		case gdk.KEY_Down:
			moveSelection(ui.list, 1)
			return true
		case gdk.KEY_Up:
			moveSelection(ui.list, -1)
			return true
		case gdk.KEY_Page_Down:
			moveSelection(ui.list, 8)
			return true
		case gdk.KEY_Page_Up:
			moveSelection(ui.list, -8)
			return true
		}
		return false
	})

	ui.list.Connect("row-activated", func(l *gtk.ListBox, row *gtk.ListBoxRow) {
		ui.activateRow(row)
	})
	ui.list.Connect("key-press-event", func(l *gtk.ListBox, ev *gdk.Event) bool {
		key := gdk.EventKeyNewFromEvent(ev)
		if key == nil {
			return false
		}
		switch key.KeyVal() {
		case gdk.KEY_Return, gdk.KEY_KP_Enter:
			ui.activateSelected()
			return true
		case gdk.KEY_Escape:
			ui.win.Close()
			return true
		}
		return false
	})

	ui.refresh()
	win.ShowAll()
	ui.entry.GrabFocus()
	logf("installer window shown")
}

func (ui *installerUI) buildFooter() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	box.SetName("footer")
	parts := []struct{ key, action string }{
		{"↑↓", "nav"},
		{"↵", "install"},
		{"Esc", "close"},
	}
	for i, p := range parts {
		if i > 0 {
			sep, _ := gtk.LabelNew("│")
			box.Add(sep)
		}
		k, _ := gtk.LabelNew("")
		k.SetMarkup(fmt.Sprintf(
			`<tt><span foreground="%s" background="%s"> %s </span></tt>`,
			ui.theme.Background, ui.theme.Accent, p.key,
		))
		box.Add(k)
		t, _ := gtk.LabelNew(p.action)
		box.Add(t)
	}
	return box
}

func (ui *installerUI) refresh() {
	query, _ := ui.entry.GetText()

	ui.visible = ui.visible[:0]
	for {
		c := ui.list.GetRowAtIndex(0)
		if c == nil {
			break
		}
		ui.list.Remove(c)
	}

	hits := Search(ui.all, query, maxResults)
	for _, p := range hits {
		ui.addRow(p)
		ui.visible = append(ui.visible, p)
	}

	total := len(ui.all)
	if len(hits) >= maxResults {
		ui.status.SetText(fmt.Sprintf("%d+ de %d paquetes", maxResults, total))
	} else {
		ui.status.SetText(fmt.Sprintf("%d de %d paquetes", len(hits), total))
	}

	if len(ui.visible) > 0 {
		if row := ui.list.GetRowAtIndex(0); row != nil {
			ui.list.SelectRow(row)
		}
	}
}

func (ui *installerUI) addRow(p *Package) {
	row, _ := gtk.ListBoxRowNew()
	row.SetName("row")

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	box.SetName("row-box")

	icon := " "
	if p.Installed {
		icon = " "
	}
	ic, _ := gtk.LabelNew(icon)
	ic.SetName("row-icon")
	ic.SetHAlign(gtk.ALIGN_CENTER)
	ic.SetSizeRequest(32, -1)
	box.Add(ic)

	name, _ := gtk.LabelNew(p.Name)
	name.SetName("row-label")
	name.SetHAlign(gtk.ALIGN_START)
	name.SetEllipsize(3)
	box.Add(name)

	repo, _ := gtk.LabelNew(p.Repo)
	repo.SetName("row-chevron")
	box.PackEnd(repo, false, false, 0)

	row.Add(box)
	ui.list.Add(row)
	// Las filas creadas después del ShowAll() inicial nacen ocultas.
	row.ShowAll()
}

func (ui *installerUI) activateSelected() {
	if row := ui.list.GetSelectedRow(); row != nil {
		ui.activateRow(row)
	}
}

func (ui *installerUI) activateRow(row *gtk.ListBoxRow) {
	if row == nil {
		return
	}
	idx := row.GetIndex()
	if idx < 0 || idx >= len(ui.visible) {
		logf("activateRow: idx %d fuera de rango (%d)", idx, len(ui.visible))
		return
	}
	p := ui.visible[idx]

	if p.Installed {
		logf("remove %s", p.Name)
		executeCommand("foot -e sudo pacman -R " + p.Name)
	} else {
		logf("install %s", p.Name)
		executeCommand("foot -e sudo pacman -S " + p.Name)
	}
	ui.win.Close()
}

var _ = strings.TrimSpace
