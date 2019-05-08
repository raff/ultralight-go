package ultralight

/*
#cgo CFLAGS: -I./SDK/include
#cgo LDFLAGS: -L./SDK/bin -lUltralight -lUltralightCore -lWebCore -lAppCore -Wl,-rpath,./SDK/bin
#include <AppCore/CAPI.h>
#include <stdlib.h>

extern void appUpdateCallback(void *);
extern void winResizeCallback(void *, unsigned int, unsigned int);
extern void winCloseCallback(void *);
extern void viewBeginLoadingCallback(void *, ULView);
extern void viewFinishLoadingCallback(void *, ULView);
extern void viewUpdateHistoryCallback(void *, ULView);
extern void viewDOMReadyCallback(void *, ULView);

static inline void set_app_update_callback(ULApp app, void *data) {
        if (data == NULL) {
            ulAppSetUpdateCallback(app, NULL, NULL);
        } else {
            ulAppSetUpdateCallback(app, appUpdateCallback, data);
        }
}

static inline void set_win_resize_callback(ULWindow win, void *data) {
        if (data == NULL) {
            ulWindowSetResizeCallback(win, NULL, NULL);
        } else {
            ulWindowSetResizeCallback(win, winResizeCallback, data);
        }
}

static inline void set_win_close_callback(ULWindow win, void *data) {
        if (data == NULL) {
            ulWindowSetCloseCallback(win, NULL, NULL);
        } else {
            ulWindowSetCloseCallback(win, winCloseCallback, data);
        }
}

static inline void set_view_begin_loading_callback(ULView view, void *data) {
        if (data == NULL) {
            ulViewSetBeginLoadingCallback(view, NULL, NULL);
        } else {
            ulViewSetBeginLoadingCallback(view, viewBeginLoadingCallback, data);
        }
}

static inline void set_view_finish_loading_callback(ULView view, void *data) {
        if (data == NULL) {
            ulViewSetFinishLoadingCallback(view, NULL, NULL);
        } else {
            ulViewSetFinishLoadingCallback(view, viewFinishLoadingCallback, data);
        }
}

static inline void set_view_update_history_callback(ULView view, void *data) {
        if (data == NULL) {
            ulViewSetUpdateHistoryCallback(view, NULL, NULL);
        } else {
            ulViewSetUpdateHistoryCallback(view, viewUpdateHistoryCallback, data);
        }
}

static inline void set_view_dom_ready_callback(ULView view, void *data) {
        if (data == NULL) {
            ulViewSetDOMReadyCallback(view, NULL, NULL);
        } else {
            ulViewSetDOMReadyCallback(view, viewDOMReadyCallback, data);
        }
}
*/
import "C"
import "unsafe"

// App is the main application object
type App struct {
	app     C.ULApp
	windows map[C.ULWindow]*Window

	onUpdate func()
}

// Window is an application window
type Window struct {
	win C.ULWindow
	ovl C.ULOverlay

	app  *App
	view *View

	onResize func(width, height uint)
	onClose  func()
}

// View is the window "content"
type View struct {
	view C.ULView

	onBeginLoading  func()
	onFinishLoading func()
	onUpdateHistory func()
	onDOMReady      func()
}

// JSContext
type JSContext struct {
	ctx C.JSContextRef
}

// NewApp creates the App singleton.
//
// Note: You should only create one of these per application lifetime.
func NewApp() *App {
	return &App{app: C.ulCreateApp(C.ulCreateConfig()), windows: map[C.ULWindow]*Window{}}
}

// Destroy destroys the App instance.
func (app *App) Destroy() {
	C.ulDestroyApp(app.app)
	app.app = nil
	app.windows = nil
}

// Window gets the main application window.
func (app *App) Window() *Window {
	ulwin := C.ulAppGetWindow(app.app)
	if win, ok := app.windows[ulwin]; ok {
		return win
	}

	win := &Window{win: ulwin}
	app.windows[ulwin] = win
	return win
}

// IsRunning checks whether or not the App is running.
func (app *App) IsRunning() bool {
	return bool(C.ulAppIsRunning(app.app))
}

// OnUpdate sets a callback for whenever the App updates.
// You should update all app logic here.
func (app *App) OnUpdate(cb func()) {
	app.onUpdate = cb
	p := unsafe.Pointer(app.app)

	if cb == nil {
		callbackData[p] = nil
		C.set_app_update_callback(app.app, nil)
	} else {
		callbackData[p] = app
		C.set_app_update_callback(app.app, p)
	}
}

// Run runs the main loop.
func (app *App) Run() {
	C.ulAppRun(app.app)
}

// Quit the application.
func (app *App) Quit() {
	C.ulAppQuit(app.app)
}

var callbackData = map[unsafe.Pointer]interface{}{}
var reverseCallbback = map[interface{}]unsafe.Pointer{}

// NewWindow create a new window and sets it as the main application window.
func (app *App) NewWindow(width, height uint, fullscreen bool, title string) *Window {
	win := &Window{win: C.ulCreateWindow(C.ulAppGetMainMonitor(app.app),
		C.uint(width), C.uint(height),
		C.bool(fullscreen),
		C.kWindowFlags_Titled|C.kWindowFlags_Resizable|C.kWindowFlags_Maximizable),
		app: app}

	if title != "" {
		t := C.CString(title)
		C.ulWindowSetTitle(win.win, t)
		C.free(unsafe.Pointer(t))
	}

	C.ulAppSetWindow(app.app, win.win)

	win.ovl = C.ulCreateOverlay(win.win, C.ulWindowGetWidth(win.win), C.ulWindowGetHeight(win.win), 0, 0)
	return win
}

// Destroy destroys the window.
func (win *Window) Destroy() {
	delete(win.app.windows, win.win)
	C.ulDestroyOverlay(win.ovl)
	C.ulDestroyWindow(win.win)
	win.OnResize(nil)
	win.ovl = nil
	win.win = nil
	win.app = nil

}

// Close closes the window.
func (win *Window) Close() {
	C.ulWindowClose(win.win)

	// should this remove the window from win.app ?
}

// SetTitle sets the window title.
func (win *Window) SetTitle(title string) {
	t := C.CString(title)
	C.ulWindowSetTitle(win.win, t)
	C.free(unsafe.Pointer(t))
}

// Resize resizes the window (and underlying View).
// Dimensions should be specified in device coordinates.
func (win *Window) Resize(width, height uint) {
	C.ulOverlayResize(win.ovl, C.uint(width), C.uint(height))
}

// OnResize sets a callback to be notified when a window resizes
// (parameters are passed back in device coordinates).
func (win *Window) OnResize(cb func(width, height uint)) {
	win.onResize = cb
	p := unsafe.Pointer(win.win)

	if cb == nil {
		callbackData[p] = nil
		C.set_win_resize_callback(win.win, nil)
	} else {
		callbackData[p] = win
		C.set_win_resize_callback(win.win, p)
	}
}

// OnClose sets a callback to be notified when a window closes.
func (win *Window) OnClose(cb func()) {
	win.onClose = cb
	p := unsafe.Pointer(win.win)

	if cb == nil {
		callbackData[p] = nil
		C.set_win_close_callback(win.win, nil)
	} else {
		callbackData[p] = win
		C.set_win_close_callback(win.win, p)
	}
}

// View gets the underlying View.
func (win *Window) View() *View {
	if win.view == nil {
		win.view = &View{view: C.ulOverlayGetView(win.ovl)}
	}

	return win.view
}

// LoadHTML loads a raw string of html
func (view *View) LoadHTML(html string) {
	s := C.CString(html)
	defer C.free(unsafe.Pointer(s))

	C.ulViewLoadHTML(view.view, C.ulCreateString(s))
}

// LoadURL loads a URL into main frame
func (view *View) LoadURL(url string) {
	s := C.CString(url)
	defer C.free(unsafe.Pointer(s))

	C.ulViewLoadURL(view.view, C.ulCreateString(s))
}

/*
// URL returns the current URL.
func (view *View) URL() string {
    s := C.ulViewGetURL(view.view)
    if C.ulStringGetLength(s) == 0 {
        return ""
    }

    data := C.ulStringGetData(s)
    return C.GoString(data)
}
*/

// IsLoading Checks if main frame is loading.
func (view *View) IsLoading() bool {
	return bool(C.ulViewIsLoading(view.view))
}

// JSContext gets the page's JSContext for use with JavaScriptCore API
func (view *View) JSContext() *JSContext {
	return &JSContext{ctx: C.ulViewGetJSContext(view.view)}
}

// CanGoBack checks if can navigate backwards in history
func (view *View) CanGoBack() bool {
	return bool(C.ulViewCanGoBack(view.view))
}

// CanGoForward checks if can navigate forwards in history
func (view *View) CanGoForward() bool {
	return bool(C.ulViewCanGoForward(view.view))
}

// GoBack navigates backwards in history
func (view *View) GoBack() {
	C.ulViewGoBack(view.view)
}

// GoForward navigates forwards in history
func (view *View) GoForward() {
	C.ulViewGoForward(view.view)
}

// GoToHistoryOffset navigates to arbitrary offset in history
func (view *View) GoToHistoryOffset(offset int) {
	C.ulViewGoToHistoryOffset(view.view, C.int(offset))
}

// Reload reloads the current page
func (view *View) Reload() {
	C.ulViewReload(view.view)
}

// Stop stops all page loads
func (view *View) Stop() {
	C.ulViewStop(view.view)
}

// Set callback for when the page begins loading new URL into main frame
func (view *View) OnBeginLoading(cb func()) {
	view.onBeginLoading = cb
	p := unsafe.Pointer(view.view)

	if cb == nil {
		callbackData[p] = nil
		C.set_view_begin_loading_callback(view.view, nil)
	} else {
		callbackData[p] = view
		C.set_view_begin_loading_callback(view.view, p)
	}
}

// Set callback for when the page finishes loading new URL into main frame
func (view *View) OnFinishLoading(cb func()) {
	view.onFinishLoading = cb
	p := unsafe.Pointer(view.view)

	if cb == nil {
		callbackData[p] = nil
		C.set_view_finish_loading_callback(view.view, nil)
	} else {
		callbackData[p] = view
		C.set_view_finish_loading_callback(view.view, p)
	}
}

// Set callback for when the history (back/forward state) is modified
func (view *View) OnUpdateHistory(cb func()) {
	view.onUpdateHistory = cb
	p := unsafe.Pointer(view.view)

	if cb == nil {
		callbackData[p] = nil
		C.set_view_update_history_callback(view.view, nil)
	} else {
		callbackData[p] = view
		C.set_view_update_history_callback(view.view, p)
	}
}

// Set callback for when all JavaScript has been parsed and the document is
// ready. This is the best time to make initial JavaScript calls to your page.
func (view *View) OnDOMReady(cb func()) {
	view.onDOMReady = cb
	p := unsafe.Pointer(view.view)

	if cb == nil {
		callbackData[p] = nil
		C.set_view_dom_ready_callback(view.view, nil)
	} else {
		callbackData[p] = view
		C.set_view_dom_ready_callback(view.view, p)
	}
}

//export appUpdateCallback
func appUpdateCallback(userData unsafe.Pointer) {
	app := callbackData[userData].(*App)
	if app != nil {
		app.onUpdate()
	}
}

//export winResizeCallback
func winResizeCallback(userData unsafe.Pointer, width, height C.uint) {
	win := callbackData[userData].(*Window)
	if win != nil {
		win.onResize(uint(width), uint(height))
	}
}

//export winCloseCallback
func winCloseCallback(userData unsafe.Pointer) {
	win := callbackData[userData].(*Window)
	if win != nil {
		win.onClose()
	}
}

//export viewBeginLoadingCallback
func viewBeginLoadingCallback(userData unsafe.Pointer, caller C.ULView) {
	view := callbackData[userData].(*View)
	if view != nil {
		view.onBeginLoading()
	}
}

//export viewFinishLoadingCallback
func viewFinishLoadingCallback(userData unsafe.Pointer, caller C.ULView) {
	view := callbackData[userData].(*View)
	if view != nil {
		view.onFinishLoading()
	}
}

//export viewUpdateHistoryCallback
func viewUpdateHistoryCallback(userData unsafe.Pointer, caller C.ULView) {
	view := callbackData[userData].(*View)
	if view != nil {
		view.onUpdateHistory()
	}
}

//export viewDOMReadyCallback
func viewDOMReadyCallback(userData unsafe.Pointer, caller C.ULView) {
	view := callbackData[userData].(*View)
	if view != nil {
		view.onDOMReady()
	}
}
