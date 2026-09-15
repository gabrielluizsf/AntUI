package dev.antui;

/**
 * The one door from Java back into Go.
 *
 * <p>Several parts of Android answer through a listener object — a dialog
 * button, a broadcast, a change of network, a fingerprint. Each wants a
 * different Java type, and native code cannot make one. So the classes in
 * this file are those types, they hold nothing but a number, and every one of
 * them ends up here.
 *
 * <p>The number is a token: Go hands one out when it registers a handler and
 * looks the handler up by it when the call comes back. Java never knows what
 * it means, which is the point — the alternative is a Java field holding a Go
 * pointer, and nothing on either side could keep that honest.
 */
final class Native {

    /**
     * Called for every listener this file implements. The method name says
     * which one was called; the arguments are whatever it was given.
     *
     * <p>It runs on whatever thread Android chose — the main one for a dialog,
     * a binder thread for a broadcast, a background one for a sensor — and it
     * runs synchronously, because some of these have to return a value.
     */
    static native Object invoke(long token, String method, Object[] args);

    private Native() {
    }
}
