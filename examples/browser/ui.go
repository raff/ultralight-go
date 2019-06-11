package main

import (
	"fmt"

	"github.com/raff/ultralight-go"
)

const (
	UI_HEIGHT = 79
)

type UI struct {
	win  *ultralight.Window
	ovl  *ultralight.Overlay
	tabs map[int]*Tab

	activeTabId  int
	tabIdCounter int

	tabWidth  uint
	tabHeight uint

	updateBack    *ultralight.JSObject
	updateForward *ultralight.JSObject
	updateLoading *ultralight.JSObject
	updateURL     *ultralight.JSObject
	addTab        *ultralight.JSObject
	updateTab     *ultralight.JSObject
	closeTab      *ultralight.JSObject
}

func NewUI(win *ultralight.Window) *UI {
	// re-use the main window overlay

	ovl := win.Overlay(0)
	ovl.Resize(win.Width(), UI_HEIGHT)

	ovl.View().OnConsoleMessage(func(source ultralight.MessageSource, level ultralight.MessageLevel,
		message string, line uint, col uint, sourceId string) {
		fmt.Printf("CONSOLE source=%v level=%v id=%q line=%c col=%v %v\n",
			source, level, sourceId, line, col, message)
	})

	ui := &UI{win: win, ovl: ovl, tabWidth: win.Width(), tabHeight: win.Height() - UI_HEIGHT, tabs: map[int]*Tab{}}

	view := ovl.View()

	/*
		view.OnBeginLoading(func() {
			fmt.Println("begin loading")
		})

		view.OnFinishLoading(func() {
			view := ovl.View()
			fmt.Println("finish loading", view.URL())
		})
	*/

	view.OnDOMReady(func() {
		jscontext := view.JSContext()
		globalObject := jscontext.GlobalObject()

		globalObject.SetPropertyValue("OnBack", ui.OnBack)
		globalObject.SetPropertyValue("OnForward", ui.OnForward)
		globalObject.SetPropertyValue("OnRefresh", ui.OnRefresh)
		globalObject.SetPropertyValue("OnStop", ui.OnStop)
		//globalObject.SetPropertyValue("OnToggleTools", ui.OnToggleTools)
		globalObject.SetPropertyValue("OnRequestNewTab", ui.OnRequestNewTab)
		globalObject.SetPropertyValue("OnRequestTabClose", ui.OnRequestTabClose)
		globalObject.SetPropertyValue("OnActiveTabChange", ui.OnActiveTabChange)
		globalObject.SetPropertyValue("OnRequestChangeURL", ui.OnRequestChangeURL)

		ui.updateBack = globalObject.Property("updateBack").Object()
		ui.updateForward = globalObject.Property("updateForward").Object()
		ui.updateLoading = globalObject.Property("updateLoading").Object()
		ui.updateURL = globalObject.Property("updateURL").Object()
		ui.addTab = globalObject.Property("addTab").Object()
		ui.updateTab = globalObject.Property("updateTab").Object()
		ui.closeTab = globalObject.Property("closeTab").Object()

		ui.CreateNewTab()
	})

	win.OnResize(func(width, height uint) {

		ui.tabWidth = win.Width()
		ui.tabHeight = win.Height() - UI_HEIGHT
		if ui.tabHeight < 1 {
			ui.tabHeight = 1
		}

		ovl.Resize(ui.tabWidth, UI_HEIGHT)

		for _, tab := range ui.tabs {
			tab.Resize(ui.tabWidth, ui.tabHeight)
		}
	})

	ovl.View().LoadURL("file:///assets/ui.html")
	return ui
}

func (ui *UI) activeTab() *Tab {
	return ui.tabs[ui.activeTabId]
}

func (ui *UI) removeTab(i int) {
	if _, ok := ui.tabs[i]; ok {
		ui.tabs[i].Destroy()
		delete(ui.tabs, i)
	}
}

func (ui *UI) OnBack(f, this *ultralight.JSObject, args ...*ultralight.JSValue) *ultralight.JSValue {
	if ui.activeTab() != nil {
		ui.activeTab().View().GoBack()
	}
	return nil
}

func (ui *UI) OnForward(f, this *ultralight.JSObject, args ...*ultralight.JSValue) *ultralight.JSValue {
	if ui.activeTab() != nil {
		ui.activeTab().View().GoForward()
	}
	return nil
}

func (ui *UI) OnRefresh(f, this *ultralight.JSObject, args ...*ultralight.JSValue) *ultralight.JSValue {
	if ui.activeTab() != nil {
		ui.activeTab().View().Reload()
	}
	return nil
}

func (ui *UI) OnStop(f, this *ultralight.JSObject, args ...*ultralight.JSValue) *ultralight.JSValue {
	if ui.activeTab() != nil {
		ui.activeTab().View().Stop()
	}
	return nil
}

func (ui *UI) OnRequestNewTab(f, this *ultralight.JSObject, args ...*ultralight.JSValue) *ultralight.JSValue {
	ui.CreateNewTab()
	return nil
}

func (ui *UI) OnRequestTabClose(f, this *ultralight.JSObject, args ...*ultralight.JSValue) *ultralight.JSValue {

	if len(args) == 1 {
		id := int(args[0].Number())

		tab := ui.tabs[id]
		if tab == nil {
			return nil
		}

		if len(ui.tabs) == 1 {
			app.Quit()
		}

		if id != ui.activeTabId {
			ui.removeTab(id)
		} else {
			tab.readyToClose = true
		}

		ui.closeTab.Call(nil, id)
	}

	return nil
}

func (ui *UI) OnActiveTabChange(f, this *ultralight.JSObject, args ...*ultralight.JSValue) *ultralight.JSValue {
	if len(args) == 1 {
		id := int(args[0].Number())

		if id == ui.activeTabId {
			return nil
		}

		if ui.tabs[id] == nil {
			return nil
		}

		ui.activeTab().Hide()
		if ui.activeTab().readyToClose {
			ui.removeTab(ui.activeTabId)
		}

		ui.activeTabId = id
		ui.activeTab().Show()

		view := ui.activeTab().View()

		ui.SetLoading(view.IsLoading())
		ui.SetCanGoBack(view.CanGoBack())
		ui.SetCanGoForward(view.CanGoBack())
		ui.SetURL(view.URL())
	}

	return nil
}

func (ui *UI) OnRequestChangeURL(f, this *ultralight.JSObject, args ...*ultralight.JSValue) *ultralight.JSValue {
	if len(args) == 1 {
		url := args[0].String()

		if len(ui.tabs) != 0 {
			ui.activeTab().View().LoadURL(url)
		}
	}

	return nil
}

func (ui *UI) CreateNewTab() {
	id := ui.tabIdCounter
	ui.tabIdCounter += 1

	tab := NewTab(ui, id, ui.tabWidth, ui.tabHeight, 0, UI_HEIGHT)
	ui.tabs[id] = tab

	tab.View().LoadURL("file:///assets/new_tab_page.html")
	ui.addTab.Call(nil, id, "New tab", "")
}

func (ui *UI) UpdateTabTitle(id int, title string) {
	ui.updateTab.Call(nil, id, title, "")
}

func (ui *UI) UpdateTabURL(id int, url string) {
	if id == ui.activeTabId && len(ui.tabs) > 0 {
		ui.SetURL(url)
	}
}

func (ui *UI) UpdateTabNavigation(id int, isLoading, canGoBack, canGoForward bool) {
	if id == ui.activeTabId && len(ui.tabs) > 0 {
		ui.SetLoading(isLoading)
		ui.SetCanGoBack(canGoBack)
		ui.SetCanGoForward(canGoForward)
	}
}

func (ui *UI) SetLoading(isLoading bool) {
	ui.updateLoading.Call(nil, isLoading)
}

func (ui *UI) SetCanGoBack(canGoBack bool) {
	ui.updateBack.Call(nil, canGoBack)
}

func (ui *UI) SetCanGoForward(canGoForward bool) {
	ui.updateForward.Call(nil, canGoForward)
}

func (ui *UI) SetURL(url string) {
	ui.updateURL.Call(nil, url)
}

//func (ui *UI) SetCursor(cursor ultralight.Cursor) { }
