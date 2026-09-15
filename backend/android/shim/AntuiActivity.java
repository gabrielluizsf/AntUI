package dev.antui;

import android.app.NativeActivity;
import android.content.Intent;
import android.os.Bundle;

/**
 * The whole of the Java in this library.
 *
 * <p>Android hands a native app a NativeActivity and a set of C callbacks,
 * and that covers the window, the input and the lifecycle. It does not cover
 * the three things below, because each of them is delivered by *overriding a
 * method* on the activity — and there is no way to override a Java method
 * from native code. So there is one subclass, it overrides those three, and
 * each one hands its arguments straight to Go.
 *
 * <p>What this class must not do is anything else. Every line here is a line
 * that has to be compiled, dexed, shipped and kept in step with the Go it
 * calls; logic that could live in Go and does not is a liability. It holds no
 * state, makes no decisions, and has no methods of its own.
 *
 * <p>The native methods are found by name rather than registered: a Go
 * function exported as Java_dev_antui_AntuiActivity_nativeOnPermissionResult
 * is what the machine looks for the first time the method is called.
 */
public class AntuiActivity extends NativeActivity {

    /**
     * The activity, for Java that needs a Context and has no other way to
     * one. Go does not use it — the platform already hands the native side a
     * reference to this object, and a second path to the same thing is a
     * second thing to keep in step.
     */
    private static AntuiActivity instance;

    public static AntuiActivity get() {
        return instance;
    }

    @Override
    protected void onCreate(Bundle state) {
        // Before super, which is what loads the native library and calls
        // into it. Anything that ran from there and wanted the instance
        // would otherwise find null.
        instance = this;
        super.onCreate(state);
    }

    @Override
    protected void onDestroy() {
        if (instance == this) {
            instance = null;
        }
        super.onDestroy();
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, String[] permissions,
                                           int[] results) {
        nativeOnPermissionResult(requestCode, permissions, results);
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        nativeOnActivityResult(requestCode, resultCode, data);
    }

    @Override
    protected void onNewIntent(Intent intent) {
        super.onNewIntent(intent);
        // So that getIntent() returns the one that just arrived rather than
        // the one the activity was started with, which is what every caller
        // expects and is not the default.
        setIntent(intent);
        nativeOnNewIntent(intent);
    }

    private static native void nativeOnPermissionResult(int requestCode,
                                                        String[] permissions,
                                                        int[] results);

    private static native void nativeOnActivityResult(int requestCode,
                                                      int resultCode,
                                                      Intent data);

    private static native void nativeOnNewIntent(Intent intent);
}
