package ultralight

/*
#cgo CFLAGS: -I./SDK/include
#cgo LDFLAGS: -L./SDK/bin -lUltralight -lUltralightCore -lWebCore -lAppCore -Wl,-rpath,./SDK/bin
#include <AppCore/CAPI.h>
#include <stdlib.h>

extern void appUpdateCallback(void *);
extern void winResizeCallback(void *, unsigned int, unsigned int);
extern void winCloseCallback(void *);

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
*/
import "C"
import "unsafe"

// App is the main application object
type App struct {
	app C.ULApp

	onUpdate func()
}

// Window is an application window
type Window struct {
	win C.ULWindow
	ovl C.ULOverlay

	onResize func(width, height uint)
	onClose  func()
}

// View is the window "content"
type View struct {
	view C.ULView
}

// JSContext
type JSContext struct {
	ctx C.JSContextRef
}

// NewApp creates the App singleton.
//
// Note: You should only create one of these per application lifetime.
func NewApp() *App {
	return &App{app: C.ulCreateApp(C.ulCreateConfig())}
}

// Destroy destroys the App instance.
func (app *App) Destroy() {
	C.ulDestroyApp(app.app)
	app.app = nil
}

// Window gets the main application window.
func (app *App) Window() *Window {
	return &Window{win: C.ulAppGetWindow(app.app)}
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

// NewWindow create a new window and sets it as the main application window.
func (app *App) NewWindow(width, height uint, fullscreen bool, title string) *Window {
	win := &Window{win: C.ulCreateWindow(C.ulAppGetMainMonitor(app.app),
		C.uint(width), C.uint(height),
		C.bool(fullscreen),
		C.kWindowFlags_Titled|C.kWindowFlags_Resizable|C.kWindowFlags_Maximizable)}

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
	C.ulDestroyOverlay(win.ovl)
	C.ulDestroyWindow(win.win)
	win.OnResize(nil)
	win.ovl = nil
	win.win = nil
}

// Close closes the window.
func (win *Window) Close() {
	C.ulWindowClose(win.win)
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
	return &View{view: C.ulOverlayGetView(win.ovl)}
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
