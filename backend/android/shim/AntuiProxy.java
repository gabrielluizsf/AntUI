package dev.antui;

import java.lang.reflect.InvocationHandler;
import java.lang.reflect.Method;
import java.lang.reflect.Proxy;

/**
 * One Java object that can be any *interface* Android asks for.
 *
 * <p>java.lang.reflect.Proxy makes a class at run time that implements the
 * interfaces it is given and sends every call to one handler. That covers
 * everything Android delivers through an interface — a dialog's click
 * listener, a sensor's listener, an executor — with no Java written per
 * listener and no dex rebuilt to add one.
 *
 * <p>It does not cover the abstract *classes*: a Proxy cannot extend one, so
 * BroadcastReceiver and NetworkCallback are written out below.
 */
final class AntuiProxy implements InvocationHandler {

    private static final Object[] NONE = new Object[0];

    private final long token;

    private AntuiProxy(long token) {
        this.token = token;
    }

    /**
     * Makes a proxy for the named interfaces. Called from Go, which knows the
     * names and cannot make a Class array itself.
     */
    static Object create(long token, String[] interfaceNames) throws ClassNotFoundException {
        ClassLoader loader = AntuiProxy.class.getClassLoader();
        Class<?>[] interfaces = new Class<?>[interfaceNames.length];
        for (int i = 0; i < interfaceNames.length; i++) {
            interfaces[i] = Class.forName(interfaceNames[i], false, loader);
        }
        return Proxy.newProxyInstance(loader, interfaces, new AntuiProxy(token));
    }

    @Override
    public Object invoke(Object proxy, Method method, Object[] args) {
        if (args == null) {
            args = NONE;
        }
        String name = method.getName();

        // Object's own three methods are answered here rather than in Go.
        // Every proxy inherits them, and a proxy put in a collection would
        // otherwise call across for a hash code.
        if (args.length == 0 && name.equals("hashCode")) {
            return System.identityHashCode(proxy);
        }
        if (args.length == 1 && name.equals("equals")) {
            return proxy == args[0];
        }
        if (args.length == 0 && name.equals("toString")) {
            return "AntuiProxy@" + Long.toHexString(token);
        }

        Object result = Native.invoke(token, name, args);
        if (result == null) {
            // A method returning a primitive cannot be answered with null:
            // the proxy unboxes what comes back and throws on a null. A
            // handler that does not care what it returns should not have to
            // know that.
            return zeroOf(method.getReturnType());
        }
        return result;
    }

    private static Object zeroOf(Class<?> type) {
        if (!type.isPrimitive() || type == void.class) {
            return null;
        }
        if (type == boolean.class) {
            return Boolean.FALSE;
        }
        if (type == char.class) {
            return Character.valueOf('\0');
        }
        if (type == byte.class) {
            return Byte.valueOf((byte) 0);
        }
        if (type == short.class) {
            return Short.valueOf((short) 0);
        }
        if (type == int.class) {
            return Integer.valueOf(0);
        }
        if (type == long.class) {
            return Long.valueOf(0L);
        }
        if (type == float.class) {
            return Float.valueOf(0f);
        }
        return Double.valueOf(0d);
    }
}
