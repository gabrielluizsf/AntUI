// Notifications for the Cocoa backend. Compiled by cgo only on macOS.

#import <Cocoa/Cocoa.h>
#import <UserNotifications/UserNotifications.h>

#include "backend_darwin.h"

// A notification posted while the app is frontmost is dropped by default,
// and a program that notifies is usually frontmost; the delegate's one job
// is to say "show it anyway".
API_AVAILABLE(macos(10.14))
@interface AntuiNotificationDelegate : NSObject <UNUserNotificationCenterDelegate>
@end

@implementation AntuiNotificationDelegate
- (void)userNotificationCenter:(UNUserNotificationCenter *)center
       willPresentNotification:(UNNotification *)notification
         withCompletionHandler:(void (^)(UNNotificationPresentationOptions))completionHandler
{
    if (@available(macOS 11.0, *)) {
        completionHandler(UNNotificationPresentationOptionBanner |
                          UNNotificationPresentationOptionList |
                          UNNotificationPresentationOptionSound);
    } else {
        completionHandler(UNNotificationPresentationOptionAlert |
                          UNNotificationPresentationOptionSound);
    }
}
@end

int antui_d_send_notification(const char *title, const char *body)
{
    if (!title || !body) return 0;
    if (@available(macOS 10.14, *)) {
        // Without a bundle identifier currentNotificationCenter throws
        // rather than returning nil, so it must not be asked at all.
        if ([[NSBundle mainBundle] bundleIdentifier] == nil) return -1;
    } else {
        return -1;
    }

    // initWithUTF8String answers nil for bytes that are not UTF-8.
    NSString *t = [[NSString alloc] initWithUTF8String:title];
    NSString *b = [[NSString alloc] initWithUTF8String:body];
    if (!t || !b) {
        [t release];
        [b release];
        return 0;
    }

    if (@available(macOS 10.14, *)) {
        UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
        static AntuiNotificationDelegate *delegate;
        static dispatch_once_t once;
        dispatch_once(&once, ^{
            delegate = [[AntuiNotificationDelegate alloc] init];
            center.delegate = delegate;
        });
        // Asking again once the user has answered returns that answer at
        // once without a prompt, so every post asks rather than keeping the
        // answer — the user can change it in System Settings at any time.
        // The block is copied by the call, which retains t and b for it.
        [center requestAuthorizationWithOptions:(UNAuthorizationOptionAlert | UNAuthorizationOptionSound)
                              completionHandler:^(BOOL granted, NSError *error) {
            if (!granted) return;
            UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
            content.title = t;
            content.body = b;
            content.sound = [UNNotificationSound defaultSound];
            UNNotificationRequest *request =
                [UNNotificationRequest requestWithIdentifier:[[NSUUID UUID] UUIDString]
                                                     content:content
                                                     trigger:nil];
            [center addNotificationRequest:request withCompletionHandler:nil];
            [content release];
        }];
    }
    [t release];
    [b release];
    return 1;
}
