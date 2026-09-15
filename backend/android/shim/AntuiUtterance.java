package dev.antui;

import android.speech.tts.UtteranceProgressListener;

/**
 * Says when a piece of speech started, finished or failed.
 *
 * <p>UtteranceProgressListener is an abstract class, so it cannot be a proxy
 * — which is the whole reason this file exists. Everything else the speech
 * engine wants is an interface and needs no Java at all.
 *
 * <p>The utterance id is passed straight through: it is what the caller gave
 * when it asked for the speech, and is how one answer is told from another
 * when several are queued.
 */
final class AntuiUtterance extends UtteranceProgressListener {

    private final long token;

    private AntuiUtterance(long token) {
        this.token = token;
    }

    static Object create(long token) {
        return new AntuiUtterance(token);
    }

    @Override
    public void onStart(String id) {
        Native.invoke(token, "onStart", new Object[]{id});
    }

    @Override
    public void onDone(String id) {
        Native.invoke(token, "onDone", new Object[]{id});
    }

    @Override
    @Deprecated
    public void onError(String id) {
        Native.invoke(token, "onError", new Object[]{id});
    }

    @Override
    public void onError(String id, int code) {
        Native.invoke(token, "onError", new Object[]{id, Integer.valueOf(code)});
    }
}
