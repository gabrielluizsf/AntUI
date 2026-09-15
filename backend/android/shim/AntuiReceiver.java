package dev.antui;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;

/**
 * A broadcast receiver that hands what it receives to Go.
 *
 * <p>It has to be a real class rather than a proxy because BroadcastReceiver
 * is abstract, and it has to be public with a no-argument constructor because
 * the system builds one itself when a broadcast arrives at an app that is not
 * running — which is exactly what a notification's action button does.
 *
 * <p>It is used both ways: built by the system for a notification, and built
 * here and registered for a broadcast the app subscribes to. The two differ
 * only in where the token comes from.
 */
public final class AntuiReceiver extends BroadcastReceiver {

    /** The extra the token travels in. */
    public static final String TOKEN = "dev.antui.token";

    /** Zero for one the system built; set for one registered from Go. */
    private final long token;

    /** What the system calls when it builds one itself. */
    public AntuiReceiver() {
        this.token = 0;
    }

    private AntuiReceiver(long token) {
        this.token = token;
    }

    /** Makes one to register by hand, for a broadcast the app subscribes to. */
    static Object create(long token) {
        return new AntuiReceiver(token);
    }

    @Override
    public void onReceive(Context context, Intent intent) {
        // A receiver registered from Go carries its own token. One the system
        // built has none, and takes it from the intent — which is how a
        // notification's action button gets back to the right handler after
        // the process it was posted from is long gone.
        long which = token;
        if (which == 0 && intent != null) {
            which = intent.getLongExtra(TOKEN, 0);
        }
        Native.invoke(which, "onReceive", new Object[]{intent});
    }
}
