package ultralight

/*
#cgo CFLAGS: -I./SDK/include
#cgo LDFLAGS: -L./SDK/bin -lUltralight -lUltralightCore -lWebCore -lAppCore -Wl,-rpath,./SDK/bin
#include <AppCore/CAPI.h>
#include <stdlib.h>

extern void go_app_update_cb(void *);
extern void go_win_resize_cb(void *, unsigned int, unsigned int);
extern void go_win_close_cb(void *);

static inline void set_app_update_callback(ULApp app, void *data) {
        if (data == NULL) {
            ulAppSetUpdateCallback(app, NULL, NULL);
        } else {
            ulAppSetUpdateCallback(app, go_app_update_cb, data);
        }
}

static inline void set_win_resize_callback(ULWindow win, void *data) {
        if (data == NULL) {
            ulWindowSetResizeCallback(win, NULL, NULL);
        } else {
            ulWindowSetResizeCallback(win, go_win_resize_cb, data);
        }
}

static inline void set_win_close_callback(ULWindow win, void *data) {
        if (data == NULL) {
            ulWindowSetCloseCallback(win, NULL, NULL);
        } else {
            ulWindowSetCloseCallback(win, go_win_close_cb, data);
        }
}
*/
import "C"
import "unsafe"

type App struct {
	app C.ULApp

	onUpdate func()
}

type Window struct {
	win C.ULWindow
	ovl C.ULOverlay

	onResize func(width, height uint)
	onClose  func()
}

type View struct {
	view C.ULView
}

///
/// Create the App singleton.
///
/// Note: You should only create one of these per application lifetime.
///
func NewApp() *App {
	return &App{app: C.ulCreateApp(C.ulCreateConfig())}
}

///
/// Destroy the App instance.
///
func (app *App) Destroy() {
	C.ulDestroyApp(app.app)
	app.app = nil
}

///
/// Get the main window.
///
func (app *App) Window() *Window {
	return &Window{win: C.ulAppGetWindow(app.app)}
}

///
/// Whether or not the App is running.
///
func (app *App) IsRunning() bool {
	return bool(C.ulAppIsRunning(app.app))
}

///
/// Set a callback for whenever the App updates. You should update all app
/// logic here.
///
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

///
/// Run the main loop.
///
func (app *App) Run() {
	C.ulAppRun(app.app)
}

///
/// Quit the application.
///
func (app *App) Quit() {
	C.ulAppQuit(app.app)
}

var callbackData = map[unsafe.Pointer]interface{}{}

///
/// Create a new Window.
///
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

///
/// Destroy a Window.
///
func (win *Window) Destroy() {
	C.ulDestroyOverlay(win.ovl)
	C.ulDestroyWindow(win.win)
	win.OnResize(nil)
	win.ovl = nil
	win.win = nil
}

///
/// Close a window.
///
func (win *Window) Close() {
	C.ulWindowClose(win.win)
}

///
/// Set the window title.
///
func (win *Window) SetTitle(title string) {
	t := C.CString(title)
	C.ulWindowSetTitle(win.win, t)
	C.free(unsafe.Pointer(t))
}

func (win *Window) Resize(width, height uint) {
	C.ulOverlayResize(win.ovl, C.uint(width), C.uint(height))
}

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

///
/// Get the underlying View.
///
func (win *Window) View() *View {
	return &View{view: C.ulOverlayGetView(win.ovl)}
}

func (view *View) LoadHTML(html string) {
	s := C.CString(html)
	defer C.free(unsafe.Pointer(s))

	C.ulViewLoadHTML(view.view, C.ulCreateString(s))
}

func (view *View) LoadURL(url string) {
	s := C.CString(url)
	defer C.free(unsafe.Pointer(s))

	C.ulViewLoadURL(view.view, C.ulCreateString(s))
}

//export go_app_update_cb
func go_app_update_cb(user_data unsafe.Pointer) {
	app := callbackData[user_data].(*App)
	if app != nil {
		app.onUpdate()
	}
}

//export go_win_resize_cb
func go_win_resize_cb(user_data unsafe.Pointer, width, height C.uint) {
	win := callbackData[user_data].(*Window)
	if win != nil {
		win.onResize(uint(width), uint(height))
	}
}

//export go_win_close_cb
func go_win_close_cb(user_data unsafe.Pointer) {
	win := callbackData[user_data].(*Window)
	if win != nil {
		win.onClose()
	}
}
