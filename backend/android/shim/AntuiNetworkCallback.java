package dev.antui;

import android.net.ConnectivityManager;
import android.net.Network;
import android.net.NetworkCapabilities;

/**
 * Watches the network and tells Go when it changes.
 *
 * <p>NetworkCallback is an abstract class, so this cannot be a proxy. The
 * three methods below are the ones worth reporting: a network appeared, the
 * one being used went away, and what it can do changed — which is how a
 * connection going from metered to not is noticed.
 */
final class AntuiNetworkCallback extends ConnectivityManager.NetworkCallback {

    private final long token;

    AntuiNetworkCallback(long token) {
        this.token = token;
    }

    static Object create(long token) {
        return new AntuiNetworkCallback(token);
    }

    @Override
    public void onAvailable(Network network) {
        Native.invoke(token, "onAvailable", new Object[]{network});
    }

    @Override
    public void onLost(Network network) {
        Native.invoke(token, "onLost", new Object[]{network});
    }

    @Override
    public void onCapabilitiesChanged(Network network, NetworkCapabilities capabilities) {
        Native.invoke(token, "onCapabilitiesChanged", new Object[]{network, capabilities});
    }
}
