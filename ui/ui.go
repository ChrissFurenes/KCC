package ui

import "github.com/rivo/tview"

type UIview struct {
	Uiapp       *tview.Application
	ConfigList  *tview.List
	InfoData    *tview.TextView
	CommandList *tview.TextView
	Grid        *tview.Grid
}
