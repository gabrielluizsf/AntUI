package dev.antui;

import android.hardware.biometrics.BiometricPrompt;

/**
 * The answer to a fingerprint or a face.
 *
 * <p>BiometricPrompt.AuthenticationCallback is abstract and arrived in API 28.
 * A class that mentions it is only loaded when it is first used, so this one
 * costs nothing on a device too old to have it — as long as nothing touches
 * it there, which is what the version check in Go is for.
 */
final class AntuiAuth extends BiometricPrompt.AuthenticationCallback {

    private final long token;

    AntuiAuth(long token) {
        this.token = token;
    }

    static Object create(long token) {
        return new AntuiAuth(token);
    }

    @Override
    public void onAuthenticationSucceeded(BiometricPrompt.AuthenticationResult result) {
        Native.invoke(token, "onAuthenticationSucceeded", new Object[0]);
    }

    @Override
    public void onAuthenticationFailed() {
        // Not an error: a finger that was not recognised. The prompt stays up
        // and the user tries again, so this is news rather than an answer.
        Native.invoke(token, "onAuthenticationFailed", new Object[0]);
    }

    @Override
    public void onAuthenticationError(int code, CharSequence message) {
        Native.invoke(token, "onAuthenticationError",
                new Object[]{Integer.valueOf(code), String.valueOf(message)});
    }
}
