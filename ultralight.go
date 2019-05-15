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
extern void viewConsoleMessageCallback(void* user_data, ULView caller,
                                       ULMessageSource source, ULMessageLevel level,
                                       ULString message, unsigned int line_number,
                                       unsigned int column_number,
                                       ULString source_id);

extern JSValueRef objFunctionCallback(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject,
                                      size_t argumentCount, JSValueRef *arguments, JSValueRef* exception);

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

static inline void set_view_console_message_callback(ULView view, void *data) {
        if (data == NULL) {
            ulViewSetAddConsoleMessageCallback(view, NULL, NULL);
        } else {
            ulViewSetAddConsoleMessageCallback(view, viewConsoleMessageCallback, data);
        }
}

static inline JSObjectRef make_function_callback(JSContextRef ctx, JSStringRef name) {
        return JSObjectMakeFunctionWithCallback(ctx, name, (JSObjectCallAsFunctionCallback)objFunctionCallback);
}
*/
import "C"
import "unsafe"
import "unicode/utf16"
import "unicode/utf8"
import "reflect"
import "bytes"

import "log"

type JSType int

const (
	JSTypeUndefined = JSType(C.kJSTypeUndefined)
	JSTypeNull      = JSType(C.kJSTypeNull)
	JSTypeBoolean   = JSType(C.kJSTypeBoolean)
	JSTypeNumber    = JSType(C.kJSTypeNumber)
	JSTypeString    = JSType(C.kJSTypeString)
	JSTypeObject    = JSType(C.kJSTypeObject)
)

type MessageSource int

const (
	MessageSourceXML            = MessageSource(C.kMessageSource_XML)
	MessageSourceJS             = MessageSource(C.kMessageSource_JS)
	MessageSourceNetwork        = MessageSource(C.kMessageSource_Network)
	MessageSourceConsoleAPI     = MessageSource(C.kMessageSource_ConsoleAPI)
	MessageSourceStorage        = MessageSource(C.kMessageSource_Storage)
	MessageSourceAppCache       = MessageSource(C.kMessageSource_AppCache)
	MessageSourceRendering      = MessageSource(C.kMessageSource_Rendering)
	MessageSourceCSS            = MessageSource(C.kMessageSource_CSS)
	MessageSourceSecurity       = MessageSource(C.kMessageSource_Security)
	MessageSourceContentBlocker = MessageSource(C.kMessageSource_ContentBlocker)
	MessageSourceOther          = MessageSource(C.kMessageSource_Other)
)

type MessageLevel int

const (
	MessageLevelLog     = MessageLevel(C.kMessageLevel_Log)
	MessageLevelWarning = MessageLevel(C.kMessageLevel_Warning)
	MessageLevelError   = MessageLevel(C.kMessageLevel_Error)
	MessageLevelDebug   = MessageLevel(C.kMessageLevel_Debug)
	MessageLevelInfo    = MessageLevel(C.kMessageLevel_Info)
)

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

	onBeginLoading   func()
	onFinishLoading  func()
	onUpdateHistory  func()
	onDOMReady       func()
	onConsoleMessage func(MessageSource, MessageLevel, string, uint, uint, string)
}

// JSContext
type JSContext struct {
	ctx C.JSContextRef
}

// JSGlobalContext
type JSGlobalContext struct {
	ctx C.JSGlobalContextRef
}

// JSValue
type JSValue struct {
	val C.JSValueRef
	ctx C.JSContextRef
}

// JSObject
type JSObject struct {
	obj C.JSObjectRef
	ctx C.JSContextRef
}

func decodeUTF16(p *C.ushort, l C.ulong) string {
	var u []uint16
	sl := (*reflect.SliceHeader)((unsafe.Pointer(&u)))
	sl.Cap = int(l)
	sl.Len = int(l)
	sl.Data = uintptr(unsafe.Pointer(p))

	runes := utf16.Decode(u)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)

	for _, r := range runes {
		n := utf8.EncodeRune(b8buf, r)
		ret.Write(b8buf[:n])
	}

	return ret.String()
}

func decodeULString(s C.ULString) string {
	l := C.ulStringGetLength(s)
	if l == 0 {
		return ""
	}

	data := C.ulStringGetData(s)
	return decodeUTF16(data, l)
}

func decodeJSString(s C.JSStringRef) string {
	l := C.JSStringGetLength(s)
	if l == 0 {
		return ""
	}

	data := C.JSStringGetCharactersPtr(s)
	return decodeUTF16(data, l)
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

	win.SetTitle(title)

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

// Whether or not an overlay has keyboard focus.
func (win *Window) HasFocus() bool {
	return bool(C.ulOverlayHasFocus(win.ovl))
}

// Grant this overlay exclusive keyboard focus.
func (win *Window) Focus() {
	C.ulOverlayFocus(win.ovl)
}

// IsFullscreen checks whether or not a window is fullscreen.
func (win *Window) IsFullscreen() bool {
	return bool(C.ulWindowIsFullscreen(win.win))
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
	uls := C.ulCreateString(s)

	defer func() {
		C.ulDestroyString(uls)
		C.free(unsafe.Pointer(s))
	}()

	C.ulViewLoadHTML(view.view, uls)
}

// LoadURL loads a URL into main frame
func (view *View) LoadURL(url string) {
	s := C.CString(url)
	uls := C.ulCreateString(s)

	defer func() {
		C.ulDestroyString(uls)
		C.free(unsafe.Pointer(s))
	}()

	C.ulViewLoadURL(view.view, uls)
}

// URL returns the current URL.
func (view *View) URL() string {
	return decodeULString(C.ulViewGetURL(view.view))
}

// Title returns the current title.
func (view *View) Title() string {
	s := C.ulViewGetTitle(view.view)
	l := C.ulStringGetLength(s)
	if l == 0 {
		return ""
	}

	data := C.ulStringGetData(s)
	return decodeUTF16(data, l)
}

// IsLoading Checks if main frame is loading.
func (view *View) IsLoading() bool {
	return bool(C.ulViewIsLoading(view.view))
}

// JSContext gets the page's JSContext for use with JavaScriptCore API
func (view *View) JSContext() *JSContext {
	return &JSContext{ctx: C.ulViewGetJSContext(view.view)}
}

// EvaluateScript evaluates a raw string of JavaScript and return result
func (view *View) EvaluateScript(script string) *JSValue {
	s := C.CString(script)
	uls := C.ulCreateString(s)

	defer func() {
		C.ulDestroyString(uls)
		C.free(unsafe.Pointer(s))
	}()

	return &JSValue{val: C.ulViewEvaluateScript(view.view, uls), ctx: C.ulViewGetJSContext(view.view)}
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

// Set callback for when a message is added to the console (useful for
// JavaScript / network errors and debugging)
func (view *View) OnConsoleMessage(cb func(source MessageSource, level MessageLevel,
	message string, line uint, col uint, sourceID string)) {
	view.onConsoleMessage = cb
	p := unsafe.Pointer(view.view)

	if cb == nil {
		callbackData[p] = nil
		C.set_view_console_message_callback(view.view, nil)
	} else {
		callbackData[p] = view
		C.set_view_console_message_callback(view.view, p)
	}
}

// Returns a JavaScript value's type.
func (v *JSValue) Type() JSType {
	return JSType(C.JSValueGetType(v.ctx, v.val))
}

// Tests whether a JavaScript value's type is the undefined type.
func (v *JSValue) IsUndefined() bool {
	return bool(C.JSValueIsUndefined(v.ctx, v.val))
}

// Tests whether a JavaScript value's type is the null type.
func (v *JSValue) IsNull() bool {
	return bool(C.JSValueIsNull(v.ctx, v.val))
}

// Tests whether a JavaScript value's type is the boolean type.
func (v *JSValue) IsBoolean() bool {
	return bool(C.JSValueIsBoolean(v.ctx, v.val))
}

// Tests whether a JavaScript value's type is the number type.
func (v *JSValue) IsNumber() bool {
	return bool(C.JSValueIsNumber(v.ctx, v.val))
}

// Tests whether a JavaScript value's type is the string type.
func (v *JSValue) IsString() bool {
	return bool(C.JSValueIsString(v.ctx, v.val))
}

// Tests whether a JavaScript value's type is the object type.
func (v *JSValue) IsObject() bool {
	return bool(C.JSValueIsObject(v.ctx, v.val))
}

// Tests whether a JavaScript value is an array.
func (v *JSValue) IsArray() bool {
	return bool(C.JSValueIsArray(v.ctx, v.val))
}

// Tests whether a JavaScript value is a date.
func (v *JSValue) IsDate() bool {
	return bool(C.JSValueIsDate(v.ctx, v.val))
}

func (v *JSValue) IsFunction() bool {
	if !v.IsObject() {
		return false
	}

	return v.Object().IsFunction()
}

// Converts a JavaScript value to boolean and returns the resulting boolean.
func (v *JSValue) Boolean() bool {
	return bool(C.JSValueToBoolean(v.ctx, v.val))
}

// Converts a JavaScript value to number and returns the resulting number.
func (v *JSValue) Number() float64 {
	return float64(C.JSValueToNumber(v.ctx, v.val, nil))
}

// Converts a JavaScript value to string and copies the result into a JavaScript string.
func (v *JSValue) String() string {
	js := C.JSValueToStringCopy(v.ctx, v.val, nil)
	if js == nil {
		return ""
	}

	return decodeJSString(js)
}

// Converts a JavaScript value to object and returns the resulting object.
func (v *JSValue) Object() *JSObject {
	o := C.JSValueToObject(v.ctx, v.val, nil)
	if o == nil {
		return nil
	}

	return &JSObject{ctx: v.ctx, obj: o}
}

// Creates a JavaScript value of the undefined type.
func (ctx *JSContext) Undefined() JSValue {
	return JSValue{ctx: ctx.ctx, val: C.JSValueMakeUndefined(ctx.ctx)}
}

// Creates a JavaScript value of the null type.
func (ctx *JSContext) Null() JSValue {
	return JSValue{ctx: ctx.ctx, val: C.JSValueMakeNull(ctx.ctx)}
}

// Creates a JavaScript value of the boolean type.
func (ctx *JSContext) Boolean(v bool) JSValue {
	return JSValue{ctx: ctx.ctx, val: C.JSValueMakeBoolean(ctx.ctx, C.bool(v))}
}

// Creates a JavaScript value of the number type.
func (ctx *JSContext) Number(v float64) JSValue {
	return JSValue{ctx: ctx.ctx, val: C.JSValueMakeNumber(ctx.ctx, C.double(v))}
}

// Creates a JavaScript value of the string type.
func (ctx *JSContext) String(v string) JSValue {
	s := C.CString(v)
	js := C.JSStringCreateWithUTF8CString(s)
	defer C.free(unsafe.Pointer(s))

	return JSValue{ctx: ctx.ctx, val: C.JSValueMakeString(ctx.ctx, js)}
}

func (ctx *JSContext) JSValue(v interface{}) JSValue {
	if v == nil {
		return ctx.Null()
	}

	switch t := v.(type) {
	case JSValue:
		return t

	case bool:
		return ctx.Boolean(t)

	case string:
		return ctx.String(t)

	case float64:
		return ctx.Number(t)

	case float32:
		return ctx.Number(float64(t))

	case int:
		return ctx.Number(float64(t))

	case int8:
		return ctx.Number(float64(t))

	case int16:
		return ctx.Number(float64(t))

	case int32:
		return ctx.Number(float64(t))

	case int64:
		return ctx.Number(float64(t))

	default:
		log.Fatalf("cannot convert %#T to JSValue", t)
	}

	return ctx.Undefined() // not reached
}

func makeJSString(v string) C.JSStringRef {
	s := C.CString(v)
	defer C.free(unsafe.Pointer(s))
	return C.JSStringCreateWithUTF8CString(s)
}

type FunctionCallback func(function, this *JSObject, args ...*JSValue) *JSValue

// Convenience method for creating a JavaScript function with a given callback as its implementation.
func (ctx *JSContext) FunctionCallback(name string, cb FunctionCallback) *JSValue {
	obj := C.make_function_callback(ctx.ctx, makeJSString(name))
	p := unsafe.Pointer(obj)
	callbackData[p] = cb
	return &JSValue{ctx: ctx.ctx, val: C.JSValueRef(obj)}
}

// Gets the global object of a JavaScript execution context.
func (ctx *JSContext) GlobalObject() *JSObject {
	return &JSObject{ctx: ctx.ctx, obj: C.JSContextGetGlobalObject(ctx.ctx)}
}

// Gets the global object of a JavaScript execution context.
func (ctx *JSContext) GlobalContext() JSGlobalContext {
	return JSGlobalContext{ctx: C.JSContextGetGlobalContext(ctx.ctx)}
}

// Tests whether an object can be called as a function.
func (o *JSObject) IsFunction() bool {
	return bool(C.JSObjectIsFunction(o.ctx, o.obj))
}

// Calls an object as a function.
func (o *JSObject) Call(this *JSObject, args ...interface{}) *JSValue {
	var thisObj C.JSObjectRef
	var jargs *C.JSValueRef

	if this != nil {
		thisObj = this.obj
	}

	nargs := len(args)
	if nargs > 0 {
		jargs = (*C.JSValueRef)(C.malloc(C.size_t(nargs) * C.size_t(unsafe.Sizeof(uintptr(0)))))
		defer C.free(unsafe.Pointer(jargs))

		var ja []C.JSValueRef
		sl := (*reflect.SliceHeader)(unsafe.Pointer(&ja))
		sl.Cap = nargs
		sl.Len = nargs
		sl.Data = uintptr(unsafe.Pointer(jargs))

		ctx := &JSContext{ctx: o.ctx}

		for i, v := range args {
			ja[i] = ctx.JSValue(v).val
		}
	}

	ret := C.JSObjectCallAsFunction(o.ctx, o.obj, thisObj, C.size_t(nargs), jargs, nil)
	return &JSValue{ctx: o.ctx, val: ret}
}

// Sets a property on an object.
func (o *JSObject) SetProperty(name string, value *JSValue) {
	C.JSObjectSetProperty(o.ctx, o.obj, makeJSString(name), value.val, 0, nil)
}

// Gets a property from an object.
func (o *JSObject) Property(name string) *JSValue {
	return &JSValue{ctx: o.ctx, val: C.JSObjectGetProperty(o.ctx, o.obj, makeJSString(name), nil)}
}

// Gets the names of an object's enumerable properties.
func (o *JSObject) PropertyNames() []string {
	anames := C.JSObjectCopyPropertyNames(o.ctx, o.obj)
	nnames := int(C.JSPropertyNameArrayGetCount(anames))
	if nnames == 0 {
		return nil
	}

	names := make([]string, nnames)
	for i := 0; i < len(names); i++ {
		n := C.JSPropertyNameArrayGetNameAtIndex(anames, C.ulong(i))
		names[i] = decodeJSString(n)
	}

	return names
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

//export viewConsoleMessageCallback
func viewConsoleMessageCallback(userData unsafe.Pointer, caller C.ULView,
	source C.ULMessageSource, level C.ULMessageLevel,
	message C.ULString, line, col C.uint,
	sourceId C.ULString) {
	view := callbackData[userData].(*View)
	if view != nil {
		view.onConsoleMessage(
			MessageSource(source),
			MessageLevel(level),
			decodeULString(message),
			uint(line), uint(col),
			decodeULString(sourceId))
	}
}

//export objFunctionCallback
func objFunctionCallback(ctx C.JSContextRef, function C.JSObjectRef, this C.JSObjectRef,
	nargs C.size_t, args *C.JSValueRef, exc *C.JSValueRef) C.JSValueRef {

	if data := callbackData[unsafe.Pointer(function)]; data != nil {
		if cb, ok := data.(FunctionCallback); ok {
			f := &JSObject{ctx: ctx, obj: function}
			fthis := &JSObject{ctx: ctx, obj: this}
			fargs := make([]*JSValue, nargs)

			if int(nargs) > 0 {
				var ja []C.JSValueRef
				sl := (*reflect.SliceHeader)(unsafe.Pointer(&ja))
				sl.Cap = int(nargs)
				sl.Len = int(nargs)
				sl.Data = uintptr(unsafe.Pointer(args))

				for i, v := range ja {
					fargs[i] = &JSValue{ctx: ctx, val: v}
				}
			}

			// FunctionCallback func(function, this *JSObject, args ...*JSValue) *JSValue
			if ret := cb(f, fthis, fargs...); ret != nil {
				return ret.val
			}
		} else {
			log.Printf("expected FunctionCallback got %#v\n", data)
		}
	}

	return C.JSValueMakeNull(ctx)
}

type Config struct {
	cfg C.ULConfig
}

type configOption func(c *Config)

func EnableImages(enabled bool) configOption {
	return func(c *Config) {
		c.EnableImages(enabled)
	}
}

func EnableJavascript(enabled bool) configOption {
	return func(c *Config) {
		c.EnableImages(enabled)
	}
}

func UseBGRA(enabled bool) configOption {
	return func(c *Config) {
		c.UseBGRAForOffscreenRendering(enabled)
	}
}

func DeviceScaleHint(value float64) configOption {
	return func(c *Config) {
		c.DeviceScaleHint(value)
	}
}

func FontFamilyStandard(fontName string) configOption {
	return func(c *Config) {
		c.FontFamilyStandard(fontName)
	}
}

func FontFamilyFixed(fontName string) configOption {
	return func(c *Config) {
		c.FontFamilyFixed(fontName)
	}
}

func FontFamilySerif(fontName string) configOption {
	return func(c *Config) {
		c.FontFamilySerif(fontName)
	}
}

func FontFamilySansSerif(fontName string) configOption {
	return func(c *Config) {
		c.FontFamilySansSerif(fontName)
	}
}

func UserAgent(agent string) configOption {
	return func(c *Config) {
		c.UserAgent(agent)
	}
}

func UserStylesheet(css string) configOption {
	return func(c *Config) {
		c.UserStylesheet(css)
	}
}

// Create config with default values (see <Ultralight/platform/Config.h>).
func NewConfig(options ...configOption) *Config {
	c := &Config{cfg: C.ulCreateConfig()}

	for _, opt := range options {
		opt(c)
	}

	return c
}

// Destroy config.
func (c *Config) Destroy() {
	C.ulDestroyConfig(c.cfg)
	c.cfg = nil
}

// Set whether images should be enabled (Default = True)
func (c *Config) EnableImages(enabled bool) {
	C.ulConfigSetEnableImages(c.cfg, C.bool(enabled))
}

// Set whether JavaScript should be eanbled (Default = True)
func (c *Config) EnableJavascript(enabled bool) {
	C.ulConfigSetEnableJavaScript(c.cfg, C.bool(enabled))
}

// Set whether we should use BGRA byte order (instead of RGBA) for View
// bitmaps. (Default = False)
func (c *Config) UseBGRAForOffscreenRendering(enabled bool) {
	C.ulConfigSetUseBGRAForOffscreenRendering(c.cfg, C.bool(enabled))
}

// Set the amount that the application DPI has been scaled, used for
// scaling device coordinates to pixels and oversampling raster shapes.
// (Default = 1.0)
func (c *Config) DeviceScaleHint(value float64) {
	C.ulConfigSetDeviceScaleHint(c.cfg, C.double(value))
}

// Set default font-family to use (Default = Times New Roman)
func (c *Config) FontFamilyStandard(fontName string) {
	s := C.CString(fontName)
	uls := C.ulCreateString(s)

	defer func() {
		C.ulDestroyString(uls)
		C.free(unsafe.Pointer(s))
	}()

	C.ulConfigSetFontFamilyStandard(c.cfg, uls)
}

// Set default font-family to use for fixed fonts, eg <pre> and <code>.
// (Default = Courier New)
func (c *Config) FontFamilyFixed(fontName string) {
	s := C.CString(fontName)
	uls := C.ulCreateString(s)

	defer func() {
		C.ulDestroyString(uls)
		C.free(unsafe.Pointer(s))
	}()

	C.ulConfigSetFontFamilyFixed(c.cfg, uls)
}

// Set default font-family to use for serif fonts. (Default = Times New Roman)
func (c *Config) FontFamilySerif(fontName string) {
	s := C.CString(fontName)
	uls := C.ulCreateString(s)

	defer func() {
		C.ulDestroyString(uls)
		C.free(unsafe.Pointer(s))
	}()

	C.ulConfigSetFontFamilySerif(c.cfg, uls)
}

// Set default font-family to use for sans-serif fonts. (Default = Arial)
func (c *Config) FontFamilySansSerif(fontName string) {
	s := C.CString(fontName)
	uls := C.ulCreateString(s)

	defer func() {
		C.ulDestroyString(uls)
		C.free(unsafe.Pointer(s))
	}()

	C.ulConfigSetFontFamilySansSerif(c.cfg, uls)
}

// Set user agent string. (See <Ultralight/platform/Config.h> for the default)
func (c *Config) UserAgent(agent string) {
	s := C.CString(agent)
	uls := C.ulCreateString(s)

	defer func() {
		C.ulDestroyString(uls)
		C.free(unsafe.Pointer(s))
	}()

	C.ulConfigSetUserAgent(c.cfg, uls)
}

// Set user stylesheet (CSS). (Default = Empty)
func (c *Config) UserStylesheet(css string) {
	s := C.CString(css)
	uls := C.ulCreateString(s)

	defer func() {
		C.ulDestroyString(uls)
		C.free(unsafe.Pointer(s))
	}()

	C.ulConfigSetUserStylesheet(c.cfg, uls)
}

type Renderer struct {
	rnd C.ULRenderer
}

// Create renderer (create this only once per application lifetime).
func NewRenderer(c *Config) *Renderer {
	return &Renderer{rnd: C.ulCreateRenderer(c.cfg)}
}

// Destroy renderer.
func (r *Renderer) Destroy() {
	C.ulDestroyRenderer(r.rnd)
	r.rnd = nil
}

// Update timers and dispatch internal callbacks (JavaScript and network)
func (r *Renderer) Update() {
	C.ulUpdate(r.rnd)
}

// Render all active Views to their respective bitmaps.
func (r *Renderer) Render() {
	C.ulRender(r.rnd)
}

// Create a View with certain size (in device coordinates).
func (r *Renderer) NewView(width, height uint, transparent bool) *View {
	return &View{view: C.ulCreateView(r.rnd, C.uint(width), C.uint(height), C.bool(transparent))}
}

// Destroy a View.
func (v *View) Destroy() {
	v.OnBeginLoading(nil)
	v.OnFinishLoading(nil)
	v.OnUpdateHistory(nil)
	v.OnDOMReady(nil)
	v.OnConsoleMessage(nil)
	C.ulDestroyView(v.view)
	v.view = nil
}

// Get bitmap (will reset the dirty flag).
//func (v *View) Bitmap() image.Image {
//    bitmap := C.ulViewGetBitmap(v.view)
//}

// Write bitmap to a PNG on disk.
func (v *View) WriteToPNG(filename string) bool {
	path := C.CString(filename)
	defer C.free(unsafe.Pointer(path))

	return bool(C.ulBitmapWritePNG(C.ulViewGetBitmap(v.view), path))
}
