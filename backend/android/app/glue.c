//go:build android

#include "_cgo_export.h"

// The pointers below are Go functions, exported by activity.go. Android
// calls every one of them on the UI thread.
//
// onSaveInstanceState returns a block the framework then frees, so what Go
// hands back has to come from malloc and not from the Go heap.
// Android declares the content rect as const, and cgo has no way to express
// const on a pointer parameter — so the exported Go function takes a plain
// ARect* and this wrapper is where the const is dropped, in one visible
// place rather than by a cast at the assignment.
static void contentRectChanged(ANativeActivity *a, const ARect *r) {
	antuiOnContentRectChanged(a, (ARect *)r);
}

void *antui_native_permission_result(void) {
	return (void *)Java_dev_antui_AntuiActivity_nativeOnPermissionResult;
}

void *antui_native_activity_result(void) {
	return (void *)Java_dev_antui_AntuiActivity_nativeOnActivityResult;
}

void *antui_native_new_intent(void) {
	return (void *)Java_dev_antui_AntuiActivity_nativeOnNewIntent;
}

void *antui_native_invoke(void) {
	return (void *)Java_dev_antui_Native_invoke;
}

int antui_watch_ui_pipe(ALooper *looper, int fd) {
	return ALooper_addFd(looper, fd, ANTUI_UI_IDENT, ALOOPER_EVENT_INPUT,
	                     antuiOnUIWork, NULL);
}

void antui_set_callbacks(ANativeActivity *a) {
	a->callbacks->onStart = antuiOnStart;
	a->callbacks->onResume = antuiOnResume;
	a->callbacks->onPause = antuiOnPause;
	a->callbacks->onStop = antuiOnStop;
	a->callbacks->onDestroy = antuiOnDestroy;
	a->callbacks->onSaveInstanceState = antuiOnSaveInstanceState;

	a->callbacks->onWindowFocusChanged = antuiOnWindowFocusChanged;
	a->callbacks->onNativeWindowCreated = antuiOnNativeWindowCreated;
	a->callbacks->onNativeWindowResized = antuiOnNativeWindowResized;
	a->callbacks->onNativeWindowRedrawNeeded = antuiOnNativeWindowRedrawNeeded;
	a->callbacks->onNativeWindowDestroyed = antuiOnNativeWindowDestroyed;

	a->callbacks->onInputQueueCreated = antuiOnInputQueueCreated;
	a->callbacks->onInputQueueDestroyed = antuiOnInputQueueDestroyed;

	a->callbacks->onContentRectChanged = contentRectChanged;
	a->callbacks->onConfigurationChanged = antuiOnConfigurationChanged;
	a->callbacks->onLowMemory = antuiOnLowMemory;
}
