// Cocoa backend for AntUI, reached from Go through cgo.
//
// macOS is the one platform where the window system cannot be spoken to over
// a socket or through a DLL: AppKit is Objective-C, so this is the one place
// the library needs a C toolchain. Everything else — Linux and Windows —
// stays pure Go.

#ifndef ANTUI_DARWIN_H
#define ANTUI_DARWIN_H

#include <stdint.h>

// The event kinds handed back to Go, matching canvas.EventType.
enum {
    ANTUI_D_NONE = 0,
    ANTUI_D_CLOSE,
    ANTUI_D_RESIZE,
    ANTUI_D_KEY_DOWN,
    ANTUI_D_KEY_UP,
    ANTUI_D_TEXT,
    ANTUI_D_MOUSE_DOWN,
    ANTUI_D_MOUSE_UP,
    ANTUI_D_MOUSE_MOVE,
    ANTUI_D_MOUSE_WHEEL,
    ANTUI_D_FOCUS,
    ANTUI_D_EXPOSE,
    // Files dragged onto the window. The paths do not travel in the event —
    // it holds numbers only — so `count` comes back in `wheel` and Go reads
    // the paths with antui_d_drop_path.
    ANTUI_D_DROP_FILES
};

// One decoded event. Go copies what it needs out of this and never keeps a
// pointer into it.
typedef struct {
    int32_t  type;
    int32_t  key;        // a ANTUI_KEY_* code, already mapped
    int32_t  button;
    int32_t  mods;
    int32_t  repeat;
    int32_t  x, y;
    int32_t  wheel;
    int32_t  width, height;
    int32_t  focused;
    uint32_t codepoint;
} antui_d_event;

typedef struct antui_d_window antui_d_window;

// Creates and shows the window. Returns NULL on failure.
antui_d_window *antui_d_open(const char *title, int width, int height);
void            antui_d_close(antui_d_window *w);

// Drains the queue into `out`, returning how many events were written.
int antui_d_pump(antui_d_window *w, antui_d_event *out, int max_events);

// Copies `rows` rows of the frame in and puts them on screen. The pixels are
// copied rather than kept, so no Go pointer ever outlives the call.
void antui_d_present(antui_d_window *w, const uint32_t *pixels,
                     int width, int height,
                     int dirty_x, int dirty_y, int dirty_w, int dirty_h);

void antui_d_set_title(antui_d_window *w, const char *title);
int  antui_d_set_fullscreen(antui_d_window *w, int on);
int  antui_d_display_size(antui_d_window *w, int *width, int *height);
int  antui_d_display_refresh(antui_d_window *w);

// The size limits. Zero on any of the four means no bound on that side, and
// an aspect of 0 keeps no ratio. `fixed` and `no_maximize` are the two
// style-mask bits: whether the window can be resized at all, and whether the
// green button zooms it.
void antui_d_set_limits(antui_d_window *w,
                        int min_width, int min_height,
                        int max_width, int max_height,
                        double aspect, int fixed, int no_maximize);

// Asks for a new drawable size. Returns 0 when there is no window.
int antui_d_set_size(antui_d_window *w, int width, int height);

// The files of the last drop. The window holds them until the next one, so
// the paths are good for as long as it takes Go to read them out.
int         antui_d_drop_count(antui_d_window *w);
const char *antui_d_drop_path(antui_d_window *w, int index);

// The application's icon: the picture the Dock shows while it runs, given as
// straight ARGB pixels, a row at a time from the top.
void antui_d_set_icon(antui_d_window *w, const uint32_t *pixels,
                      int width, int height);

// The clipboard. The text is a fresh string the caller frees; the file paths
// belong to the pasteboard's own list, which is read again — and so replaced
// — by the next call to antui_d_clipboard_count.
char       *antui_d_clipboard_text(void);
void        antui_d_set_clipboard(const char *text);
int         antui_d_clipboard_count(void);
const char *antui_d_clipboard_path(int index);

#endif /* ANTUI_DARWIN_H */