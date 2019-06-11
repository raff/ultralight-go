package main

import (
	"github.com/raff/ultralight-go"
)

type Tab struct {
	ovl          *ultralight.Overlay
	id           int
	readyToClose bool
}

func NewTab(ui *UI, id int, width, height uint, x, y int) *Tab {
	ovl := ui.win.NewOverlay(width, height, x, y)

	view := ovl.View()

	view.OnChangeTitle(func(title string) {
		ui.UpdateTabTitle(id, title)
	})

	view.OnChangeURL(func(url string) {
		ui.UpdateTabURL(id, url)
	})

	/*
		view.OnChangeCursor(func(cursor ultralight.Cursor) {
			if id == ui.activeTabId {
				ui.SetCursor(cursor)
			}
		})
	*/

	view.OnBeginLoading(func() {
		ui.UpdateTabNavigation(id, view.IsLoading(), view.CanGoBack(), view.CanGoForward())
	})

	view.OnFinishLoading(func() {
		ui.UpdateTabNavigation(id, view.IsLoading(), view.CanGoBack(), view.CanGoForward())
	})

	view.OnUpdateHistory(func() {
		ui.UpdateTabNavigation(id, view.IsLoading(), view.CanGoBack(), view.CanGoForward())
	})

	return &Tab{ovl: ovl, id: id}
}

func (tab *Tab) Destroy() {
	tab.ovl.Destroy()
	tab.ovl = nil
}

func (tab *Tab) Show() {
	tab.ovl.Show()
	tab.ovl.Focus()

	/*
	   if tab.inspectorOverlay {
	       tab.inspectorOverlay.Show()
	   }
	*/
}

func (tab *Tab) Hide() {
	tab.ovl.Hide()
	tab.ovl.Unfocus()

	/*
	   if tab.inspectorOverlay {
	       tab.inspectorOverlay.Hide()
	   }
	*/
}

func (tab *Tab) Resize(width, height uint) {
	if height < 1 {
		height = 1
	}
	tab.ovl.Resize(width, height)
}

func (tab *Tab) View() *ultralight.View {
	return tab.ovl.View()
}
