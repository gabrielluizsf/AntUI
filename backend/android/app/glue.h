//go:build android

// The callback table Android fills in. Assigning the function pointers has
// to happen in C, because cgo will not let a file that exports functions
// also define one — so the exports live in activity.go and the one line of
// wiring lives here.
#ifndef ANTUI_ANDROID_APP_GLUE_H
#define ANTUI_ANDROID_APP_GLUE_H

#include <android/native_activity.h>

#include <android/looper.h>

// antui_set_callbacks points every callback in the activity at the Go
// function that handles it. It is called once, from
// ANativeActivity_onCreate, on the UI thread.
void antui_set_callbacks(ANativeActivity *activity);

// The identifier the work pipe is attached to the UI thread's looper with.
#define ANTUI_UI_IDENT 0x6d77

// The addresses of the three Go functions the Java shim calls.
//
// They exist because Android's NativeActivity **dlopens** this library
// rather than loading it through System.loadLibrary — so the machine never
// associates it with the app's class loader, and never finds a native method
// by its name. The symbols are in the library; the lookup does not reach
// them. So they are registered by hand, and their addresses have to come
// from C, because cgo will not hand back the address of a function it
// exported.
void *antui_native_permission_result(void);
void *antui_native_activity_result(void);
void *antui_native_new_intent(void);
void *antui_native_invoke(void);

// antui_watch_ui_pipe asks the UI thread's looper to call back whenever
// something is written to fd. It is how work is handed to that thread
// without a line of Java: Android's own event loop is already there, and a
// pipe is something it knows how to wait on.
int antui_watch_ui_pipe(ALooper *looper, int fd);

#endif
