// Cocoa backend for AntUI. Compiled by cgo only on macOS.

#import <Cocoa/Cocoa.h>
#import <CoreVideo/CoreVideo.h>
#include <string.h>
#include <stdlib.h>

#include "backend_darwin.h"

#define ANTUI_D_MAX_QUEUE 512

// Helper function to safely execute UI code on the Apple main thread.
static void antui_d_on_main(void (^block)(void)) {
    if ([NSThread isMainThread]) {
        block();
    } else {
        dispatch_sync(dispatch_get_main_queue(), block);
    }
}

// Key codes are hardware positions on macOS, not characters, so this table is
// the layout-independent map from where a key is to which key it is.
static int antui_d_key(unsigned short code)
{
    switch (code) {
    case 0x00: return 'A'; case 0x0B: return 'B'; case 0x08: return 'C';
    case 0x02: return 'D'; case 0x0E: return 'E'; case 0x03: return 'F';
    case 0x05: return 'G'; case 0x04: return 'H'; case 0x22: return 'I';
    case 0x26: return 'J'; case 0x28: return 'K'; case 0x25: return 'L';
    case 0x2E: return 'M'; case 0x2D: return 'N'; case 0x1F: return 'O';
    case 0x23: return 'P'; case 0x0C: return 'Q'; case 0x0F: return 'R';
    case 0x01: return 'S'; case 0x11: return 'T'; case 0x20: return 'U';
    case 0x09: return 'V'; case 0x0D: return 'W'; case 0x07: return 'X';
    case 0x10: return 'Y'; case 0x06: return 'Z';
    case 0x1D: return '0'; case 0x12: return '1'; case 0x13: return '2';
    case 0x14: return '3'; case 0x15: return '4'; case 0x17: return '5';
    case 0x16: return '6'; case 0x1A: return '7'; case 0x1C: return '8';
    case 0x19: return '9';
    case 0x31: return ' ';
    case 0x1B: return '-';  case 0x18: return '=';  case 0x2B: return ',';
    case 0x2F: return '.';  case 0x2C: return '/';  case 0x29: return ';';
    case 0x21: return '[';  case 0x1E: return ']';  case 0x2A: return '\\';
    case 0x32: return '`';  case 0x27: return '\'';

    // The special keys, numbered from 256 up exactly as AntUI's Key is.
    case 0x35: return 256;        // escape
    case 0x24: case 0x4C: return 257;  // enter, and the keypad's
    case 0x30: return 258;        // tab
    case 0x33: return 259;        // backspace
    case 0x75: return 261;        // delete (forward)
    case 0x7C: return 262;        // right
    case 0x7B: return 263;        // left
    case 0x7D: return 264;        // down
    case 0x7E: return 265;        // up
    case 0x74: return 266;        // page up
    case 0x79: return 267;        // page down
    case 0x73: return 268;        // home
    case 0x77: return 269;        // end
    case 0x39: return 270;        // caps lock
    case 0x7A: return 271;        // F1
    case 0x78: return 272;        // F2
    case 0x63: return 273;        // F3
    case 0x76: return 274;        // F4
    case 0x60: return 275;        // F5
    case 0x61: return 276;        // F6
    case 0x62: return 277;        // F7
    case 0x64: return 278;        // F8
    case 0x65: return 279;        // F9
    case 0x6D: return 280;        // F10
    case 0x67: return 281;        // F11
    case 0x6F: return 282;        // F12
    case 0x38: return 283;        // left shift
    case 0x3B: return 284;        // left control
    case 0x3A: return 285;        // left alt (option)
    case 0x37: return 286;        // left super (command)
    case 0x3C: return 287;        // right shift
    case 0x3E: return 288;        // right control
    case 0x3D: return 289;        // right alt
    case 0x36: return 290;        // right super
    default:   return 0;
    }
}

static int antui_d_mods(NSEventModifierFlags flags)
{
    int mods = 0;
    if (flags & NSEventModifierFlagShift)   mods |= 1;
    if (flags & NSEventModifierFlagControl) mods |= 2;
    if (flags & NSEventModifierFlagOption)  mods |= 4;
    if (flags & NSEventModifierFlagCommand) mods |= 8;
    return mods;
}

// ---------------------------------------------------------------------------

@class AntuiView;

struct antui_d_window {
    NSWindow  *window;
    AntuiView *view;
    id         delegate;

    // The frame lives here, in memory this file owns. Go copies into it and
    // never hands a Go pointer across, which is what cgo requires and what
    // keeps the garbage collector from moving a buffer out from under AppKit.
    uint32_t *pixels;
    int       width, height;

    antui_d_event queue[ANTUI_D_MAX_QUEUE];
    int           queued;
    int           fullscreen;

    // The paths of the last drop, kept until the next one: an event carries
    // numbers, and a list of names is not a number.
    NSArray<NSString *> *dropped;
};

static void antui_d_push(antui_d_window *w, antui_d_event ev)
{
    if (!w || w->queued >= ANTUI_D_MAX_QUEUE) return;
    w->queue[w->queued++] = ev;
}

// ---------------------------------------------------------------------------

@interface AntuiView : NSView <NSDraggingDestination>
@property (nonatomic, assign) antui_d_window *owner;
@end

@implementation AntuiView

// The frame arrives already composited and opaque, so telling AppKit that
// spares it blending the whole window against whatever is behind it.
- (BOOL)isOpaque { return YES; }
- (BOOL)acceptsFirstResponder { return YES; }

// --- files dragged onto the window ----------------------------------------

- (NSDragOperation)draggingEntered:(id<NSDraggingInfo>)sender
{
    return NSDragOperationCopy;
}

- (NSDragOperation)draggingUpdated:(id<NSDraggingInfo>)sender
{
    return NSDragOperationCopy;
}

- (BOOL)prepareForDragOperation:(id<NSDraggingInfo>)sender { return YES; }

- (BOOL)performDragOperation:(id<NSDraggingInfo>)sender
{
    antui_d_window *w = self.owner;
    if (!w) return NO;

    NSArray<NSURL *> *urls = [[sender draggingPasteboard]
        readObjectsForClasses:@[[NSURL class]]
                      options:@{NSPasteboardURLReadingFileURLsOnlyKey: @YES}];
    if (urls.count == 0) return NO;

    NSMutableArray<NSString *> *paths = [NSMutableArray arrayWithCapacity:urls.count];
    for (NSURL *url in urls) {
        if (url.path) [paths addObject:url.path];
    }
    if (paths.count == 0) return NO;

    NSArray<NSString *> *previous = w->dropped;
    w->dropped = [paths copy];
    [previous release];

    NSPoint at = [self convertPoint:[sender draggingLocation] fromView:nil];
    NSRect bounds = [self bounds];

    antui_d_event ev;
    memset(&ev, 0, sizeof ev);
    ev.type  = ANTUI_D_DROP_FILES;
    ev.x     = (int32_t)at.x;
    ev.y     = (int32_t)(bounds.size.height - at.y);
    ev.wheel = (int32_t)paths.count;
    antui_d_push(w, ev);
    return YES;
}

- (void)drawRect:(NSRect)dirty
{
    antui_d_window *w = self.owner;
    if (!w || !w->pixels || w->width <= 0 || w->height <= 0) return;

    CGContextRef ctx = [[NSGraphicsContext currentContext] CGContext];
    if (!ctx) return;

    CGColorSpaceRef space = CGColorSpaceCreateDeviceRGB();
    CGDataProviderRef provider = CGDataProviderCreateWithData(
        NULL, w->pixels, (size_t)w->width * (size_t)w->height * 4, NULL);
    CGImageRef image = CGImageCreate(
        (size_t)w->width, (size_t)w->height, 8, 32, (size_t)w->width * 4, space,
        kCGBitmapByteOrder32Little | kCGImageAlphaNoneSkipFirst,
        provider, NULL, false, kCGRenderingIntentDefault);

    if (image) {
        NSRect bounds = [self bounds];
        CGContextSaveGState(ctx);
        CGContextSetInterpolationQuality(ctx, kCGInterpolationNone);
        CGContextTranslateCTM(ctx, 0, bounds.size.height);
        CGContextScaleCTM(ctx, 1.0, -1.0);
        CGContextDrawImage(ctx, CGRectMake(0, 0, bounds.size.width, bounds.size.height),
                           image);
        CGContextRestoreGState(ctx);
        CGImageRelease(image);
    }
    CGDataProviderRelease(provider);
    CGColorSpaceRelease(space);

    antui_d_event ev;
    memset(&ev, 0, sizeof ev);
    ev.type = ANTUI_D_EXPOSE;
    antui_d_push(w, ev);
}
@end

// ---------------------------------------------------------------------------

@interface AntuiDelegate : NSObject <NSWindowDelegate>
@property (nonatomic, assign) antui_d_window *owner;
@end

@implementation AntuiDelegate

- (BOOL)windowShouldClose:(id)sender
{
    (void)sender;
    antui_d_event ev;
    memset(&ev, 0, sizeof ev);
    ev.type = ANTUI_D_CLOSE;
    antui_d_push(self.owner, ev);
    return NO;
}

- (void)windowDidResize:(NSNotification *)notification
{
    (void)notification;
    antui_d_window *w = self.owner;
    if (!w) return;

    NSRect bounds = [w->view bounds];
    int width  = (int)(bounds.size.width + 0.5);
    int height = (int)(bounds.size.height + 0.5);
    if (width <= 0 || height <= 0 || (width == w->width && height == w->height)) return;

    uint32_t *grown = (uint32_t *)calloc((size_t)width * (size_t)height, 4);
    if (!grown) return;
    free(w->pixels);
    w->pixels = grown;
    w->width  = width;
    w->height = height;

    antui_d_event ev;
    memset(&ev, 0, sizeof ev);
    ev.type = ANTUI_D_RESIZE;
    ev.width = width;
    ev.height = height;
    antui_d_push(w, ev);
}

- (void)windowDidBecomeKey:(NSNotification *)notification
{
    (void)notification;
    antui_d_event ev;
    memset(&ev, 0, sizeof ev);
    ev.type = ANTUI_D_FOCUS;
    ev.focused = 1;
    antui_d_push(self.owner, ev);
}

- (void)windowDidResignKey:(NSNotification *)notification
{
    (void)notification;
    antui_d_event ev;
    memset(&ev, 0, sizeof ev);
    ev.type = ANTUI_D_FOCUS;
    ev.focused = 0;
    antui_d_push(self.owner, ev);
}
@end

// ---------------------------------------------------------------------------

antui_d_window *antui_d_open(const char *title, int width, int height)
{
    __block antui_d_window *w = NULL;
    
    antui_d_on_main(^{
        @autoreleasepool {
            w = (antui_d_window *)calloc(1, sizeof *w);
            if (!w) return;

            w->width  = width;
            w->height = height;
            w->pixels = (uint32_t *)calloc((size_t)width * (size_t)height, 4);
            if (!w->pixels) { free(w); w = NULL; return; }

            [NSApplication sharedApplication];
            [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

            NSRect frame = NSMakeRect(0, 0, width, height);
            NSWindowStyleMask style = NSWindowStyleMaskTitled |
                                      NSWindowStyleMaskClosable |
                                      NSWindowStyleMaskMiniaturizable |
                                      NSWindowStyleMaskResizable;

            w->window = [[NSWindow alloc] initWithContentRect:frame
                                                   styleMask:style
                                                     backing:NSBackingStoreBuffered
                                                       defer:NO];
            if (!w->window) { free(w->pixels); free(w); w = NULL; return; }

            w->view = [[AntuiView alloc] initWithFrame:frame];
            w->view.owner = w;

            AntuiDelegate *delegate = [[AntuiDelegate alloc] init];
            delegate.owner = w;
            w->delegate = delegate;

            [w->window setContentView:w->view];
            [w->window setDelegate:delegate];
            [w->window setTitle:[NSString stringWithUTF8String:title ? title : ""]];
            [w->window setAcceptsMouseMovedEvents:YES];
            [w->view registerForDraggedTypes:@[NSPasteboardTypeFileURL]];
            [w->window makeFirstResponder:w->view];
            [w->window center];
            [w->window makeKeyAndOrderFront:nil];

            [NSApp activateIgnoringOtherApps:YES];
            [NSApp finishLaunching];
        }
    });
    return w;
}

void antui_d_close(antui_d_window *w)
{
    if (!w) return;
    antui_d_on_main(^{
        @autoreleasepool {
            [w->window setDelegate:nil];
            [w->window close];
            [w->dropped release];
            w->dropped = nil;
        }
    });
    free(w->pixels);
    free(w);
}

int antui_d_pump(antui_d_window *w, antui_d_event *out, int max_events)
{
    if (!w || !out || max_events <= 0) return 0;
    __block int count = 0;

    antui_d_on_main(^{
        @autoreleasepool {
            for (;;) {
                NSEvent *event = [NSApp nextEventMatchingMask:NSEventMaskAny
                                                    untilDate:[NSDate distantPast]
                                                       inMode:NSDefaultRunLoopMode
                                                      dequeue:YES];
                if (!event) break;

                antui_d_event ev;
                memset(&ev, 0, sizeof ev);
                ev.mods = antui_d_mods([event modifierFlags]);

                NSEventType type = [event type];
                switch (type) {
                case NSEventTypeMouseMoved:
                case NSEventTypeLeftMouseDragged:
                case NSEventTypeRightMouseDragged:
                case NSEventTypeOtherMouseDragged:
                case NSEventTypeLeftMouseDown:
                case NSEventTypeLeftMouseUp:
                case NSEventTypeRightMouseDown:
                case NSEventTypeRightMouseUp:
                case NSEventTypeOtherMouseDown:
                case NSEventTypeOtherMouseUp:
                case NSEventTypeScrollWheel: {
                    NSPoint point = [event locationInWindow];
                    ev.x = (int)point.x;
                    ev.y = w->height - (int)point.y;
                    break;
                }
                default:
                    break;
                }

                switch (type) {
                case NSEventTypeMouseMoved:
                case NSEventTypeLeftMouseDragged:
                case NSEventTypeRightMouseDragged:
                case NSEventTypeOtherMouseDragged:
                    ev.type = ANTUI_D_MOUSE_MOVE;
                    antui_d_push(w, ev);
                    break;

                case NSEventTypeLeftMouseDown:
                case NSEventTypeRightMouseDown:
                case NSEventTypeOtherMouseDown:
                    ev.type = ANTUI_D_MOUSE_DOWN;
                    ev.button = (type == NSEventTypeLeftMouseDown)  ? 0 :
                                (type == NSEventTypeRightMouseDown) ? 1 : 2;
                    antui_d_push(w, ev);
                    break;

                case NSEventTypeLeftMouseUp:
                case NSEventTypeRightMouseUp:
                case NSEventTypeOtherMouseUp:
                    ev.type = ANTUI_D_MOUSE_UP;
                    ev.button = (type == NSEventTypeLeftMouseUp)  ? 0 :
                                (type == NSEventTypeRightMouseUp) ? 1 : 2;
                    antui_d_push(w, ev);
                    break;

                case NSEventTypeScrollWheel: {
                    double delta = [event deltaY];
                    ev.type = ANTUI_D_MOUSE_WHEEL;
                    ev.wheel = delta > 0 ? 1 : (delta < 0 ? -1 : 0);
                    if (ev.wheel) antui_d_push(w, ev);
                    break;
                }

                case NSEventTypeKeyDown:
                case NSEventTypeKeyUp: {
                    ev.type = (type == NSEventTypeKeyDown) ? ANTUI_D_KEY_DOWN : ANTUI_D_KEY_UP;
                    ev.key = antui_d_key([event keyCode]);
                    ev.repeat = (type == NSEventTypeKeyDown && [event isARepeat]) ? 1 : 0;
                    antui_d_push(w, ev);

                    if (type == NSEventTypeKeyDown && !(ev.mods & 8) && !(ev.mods & 2)) {
                        NSString *characters = [event characters];
                        NSUInteger charCount = [characters length];
                        for (NSUInteger i = 0; i < charCount; ++i) {
                            unichar unit = [characters characterAtIndex:i];
                            if (unit < 32 || unit == 127 || unit >= 0xF700) continue;
                            antui_d_event text;
                            memset(&text, 0, sizeof text);
                            text.type = ANTUI_D_TEXT;
                            text.codepoint = unit;
                            text.mods = ev.mods;
                            antui_d_push(w, text);
                        }
                    }
                    break;
                }
                default:
                    break;
                }

                [NSApp sendEvent:event];
            }

            count = w->queued < max_events ? w->queued : max_events;
            memcpy(out, w->queue, (size_t)count * sizeof(antui_d_event));
            
            if (count < w->queued) {
                memmove(w->queue, w->queue + count,
                        (size_t)(w->queued - count) * sizeof(antui_d_event));
            }
            w->queued -= count;
        }
    });

    return count;
}

void antui_d_present(antui_d_window *w, const uint32_t *pixels,
                     int width, int height,
                     int dirty_x, int dirty_y, int dirty_w, int dirty_h)
{
    if (!w || !pixels || dirty_w <= 0 || dirty_h <= 0) return;
    if (width != w->width || height != w->height) return;

    // Memory copy operations are completely safe off the main thread.
    for (int row = 0; row < dirty_h; ++row) {
        int y = dirty_y + row;
        if (y < 0 || y >= w->height) continue;
        memcpy(w->pixels + (size_t)y * w->width + dirty_x,
               pixels + (size_t)y * width + dirty_x,
               (size_t)dirty_w * 4);
    }

    // Drawing operations must be synchronized.
    antui_d_on_main(^{
        @autoreleasepool {
            [w->view setNeedsDisplay:YES];
            [w->view displayIfNeeded];
        }
    });
}

void antui_d_set_title(antui_d_window *w, const char *title)
{
    if (!w) return;
    antui_d_on_main(^{
        @autoreleasepool {
            [w->window setTitle:[NSString stringWithUTF8String:title ? title : ""]];
        }
    });
}

int antui_d_set_fullscreen(antui_d_window *w, int on)
{
    if (!w) return 0;
    __block int result = 1;
    antui_d_on_main(^{
        @autoreleasepool {
            int already = ([w->window styleMask] & NSWindowStyleMaskFullScreen) != 0;
            if (already != (on != 0)) [w->window toggleFullScreen:nil];
            w->fullscreen = on ? 1 : 0;
        }
    });
    return result;
}

void antui_d_set_limits(antui_d_window *w,
                        int min_width, int min_height,
                        int max_width, int max_height,
                        double aspect, int fixed, int no_maximize)
{
    if (!w || !w->window) return;
    antui_d_on_main(^{
        @autoreleasepool {
            NSSize small = NSMakeSize(min_width > 0 ? min_width : 1,
                                      min_height > 0 ? min_height : 1);
            NSSize large = NSMakeSize(max_width > 0 ? max_width : CGFLOAT_MAX,
                                      max_height > 0 ? max_height : CGFLOAT_MAX);
            [w->window setContentMinSize:small];
            [w->window setContentMaxSize:large];

            if (aspect > 0.0) {
                [w->window setContentAspectRatio:NSMakeSize(aspect, 1.0)];
            } else {
                [w->window setContentResizeIncrements:NSMakeSize(1.0, 1.0)];
            }

            NSWindowStyleMask mask = [w->window styleMask];
            if (!(mask & NSWindowStyleMaskFullScreen)) {
                if (fixed) {
                    mask &= ~NSWindowStyleMaskResizable;
                } else {
                    mask |= NSWindowStyleMaskResizable;
                }
                [w->window setStyleMask:mask];

                NSButton *zoom = [w->window standardWindowButton:NSWindowZoomButton];
                [zoom setEnabled:(fixed || no_maximize) ? NO : YES];
            }
        }
    });
}

int antui_d_set_size(antui_d_window *w, int width, int height)
{
    if (!w || !w->window || width < 1 || height < 1) return 0;
    antui_d_on_main(^{
        @autoreleasepool {
            [w->window setContentSize:NSMakeSize(width, height)];
        }
    });
    return 1;
}

int antui_d_display_size(antui_d_window *w, int *width, int *height)
{
    __block int success = 0;
    antui_d_on_main(^{
        @autoreleasepool {
            NSScreen *screen = w && w->window ? [w->window screen] : [NSScreen mainScreen];
            if (!screen) screen = [NSScreen mainScreen];
            if (screen) {
                NSRect frame = [screen frame];
                if (width)  *width  = (int)frame.size.width;
                if (height) *height = (int)frame.size.height;
                success = 1;
            }
        }
    });
    return success;
}

void antui_d_set_icon(antui_d_window *w, const uint32_t *pixels,
                      int width, int height)
{
    (void)w;
    if (!pixels || width <= 0 || height <= 0) return;

    antui_d_on_main(^{
        @autoreleasepool {
            NSBitmapImageRep *rep = [[NSBitmapImageRep alloc]
                initWithBitmapDataPlanes:NULL
                              pixelsWide:width
                              pixelsHigh:height
                           bitsPerSample:8
                         samplesPerPixel:4
                                hasAlpha:YES
                                isPlanar:NO
                          colorSpaceName:NSDeviceRGBColorSpace
                            bitmapFormat:NSBitmapFormatAlphaFirst
                             bytesPerRow:width * 4
                            bitsPerPixel:32];
            if (!rep) return;

            memcpy([rep bitmapData], pixels, (size_t)width * (size_t)height * 4);

            NSImage *icon = [[NSImage alloc] initWithSize:NSMakeSize(width, height)];
            [icon addRepresentation:rep];
            [NSApp setApplicationIconImage:icon];
        }
    });
}

int antui_d_drop_count(antui_d_window *w)
{
    __block int count = 0;
    antui_d_on_main(^{
        if (w && w->dropped) {
            count = (int)w->dropped.count;
        }
    });
    return count;
}

const char *antui_d_drop_path(antui_d_window *w, int index)
{
    __block const char *path = NULL;
    antui_d_on_main(^{
        if (w && w->dropped && index >= 0 && index < (int)w->dropped.count) {
            path = [w->dropped[index] UTF8String];
        }
    });
    return path;
}

static NSArray<NSString *> *antui_d_clip_files = nil;

char *antui_d_clipboard_text(void)
{
    __block char *result = NULL;
    antui_d_on_main(^{
        @autoreleasepool {
            NSString *text = [[NSPasteboard generalPasteboard]
                stringForType:NSPasteboardTypeString];
            if (text) {
                const char *utf8 = [text UTF8String];
                if (utf8) {
                    result = strdup(utf8);
                }
            }
        }
    });
    return result;
}

void antui_d_set_clipboard(const char *text)
{
    if (!text) return;
    antui_d_on_main(^{
        @autoreleasepool {
            NSPasteboard *board = [NSPasteboard generalPasteboard];
            [board clearContents];
            [board setString:[NSString stringWithUTF8String:text]
                     forType:NSPasteboardTypeString];
        }
    });
}

int antui_d_clipboard_count(void)
{
    __block int count = 0;
    antui_d_on_main(^{
        @autoreleasepool {
            NSArray<NSURL *> *urls = [[NSPasteboard generalPasteboard]
                readObjectsForClasses:@[[NSURL class]]
                              options:@{NSPasteboardURLReadingFileURLsOnlyKey: @YES}];
            NSMutableArray<NSString *> *paths = [NSMutableArray array];
            for (NSURL *url in urls) {
                if (url.path) [paths addObject:url.path];
            }
            NSArray<NSString *> *previous = antui_d_clip_files;
            antui_d_clip_files = [paths copy];
            [previous release];
            
            count = (int)antui_d_clip_files.count;
        }
    });
    return count;
}

const char *antui_d_clipboard_path(int index)
{
    __block const char *path = NULL;
    antui_d_on_main(^{
        if (antui_d_clip_files && index >= 0 && index < (int)antui_d_clip_files.count) {
            path = [antui_d_clip_files[index] UTF8String];
        }
    });
    return path;
}

int antui_d_display_refresh(antui_d_window *w)
{
    __block int refresh = 0;
    antui_d_on_main(^{
        @autoreleasepool {
            NSScreen *screen = w && w->window ? [w->window screen] : [NSScreen mainScreen];
            if (screen) {
                if (@available(macOS 12.0, *)) {
                    NSInteger hz = [screen maximumFramesPerSecond];
                    if (hz > 0 && hz < 1000) {
                        refresh = (int)hz;
                    }
                }
            }
        }
    });
    return refresh;
}