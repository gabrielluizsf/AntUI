//go:build android

// cgo cannot call a C function pointer, and every JNI call is one: JNIEnv is
// a pointer to a table of them, and the way to call FindClass is
// (*env)->FindClass(env, ...). So each call this library makes needs a real
// function to stand in front of it, and this file is that — nothing but
// one-line wrappers, in the same order as <jni.h>.
#ifndef ANTUI_ANDROID_JNI_WRAP_H
#define ANTUI_ANDROID_JNI_WRAP_H

#include <jni.h>
#include <stdlib.h>

// The virtual machine, and this thread's view of it.
static inline jint mw_GetEnv(JavaVM *vm, JNIEnv **env, jint version) {
	return (*vm)->GetEnv(vm, (void **)env, version);
}
static inline jint mw_AttachCurrentThread(JavaVM *vm, JNIEnv **env) {
	return (*vm)->AttachCurrentThread(vm, env, NULL);
}
static inline jint mw_DetachCurrentThread(JavaVM *vm) {
	return (*vm)->DetachCurrentThread(vm);
}

// Classes.
static inline jclass mw_FindClass(JNIEnv *env, const char *name) {
	return (*env)->FindClass(env, name);
}
static inline jclass mw_GetObjectClass(JNIEnv *env, jobject o) {
	return (*env)->GetObjectClass(env, o);
}
static inline jboolean mw_IsInstanceOf(JNIEnv *env, jobject o, jclass c) {
	return (*env)->IsInstanceOf(env, o, c);
}
static inline jboolean mw_IsSameObject(JNIEnv *env, jobject a, jobject b) {
	return (*env)->IsSameObject(env, a, b);
}

// References.
static inline jobject mw_NewGlobalRef(JNIEnv *env, jobject o) {
	return (*env)->NewGlobalRef(env, o);
}
static inline void mw_DeleteGlobalRef(JNIEnv *env, jobject o) {
	(*env)->DeleteGlobalRef(env, o);
}
static inline jobject mw_NewWeakGlobalRef(JNIEnv *env, jobject o) {
	return (*env)->NewWeakGlobalRef(env, o);
}
static inline void mw_DeleteWeakGlobalRef(JNIEnv *env, jobject o) {
	(*env)->DeleteWeakGlobalRef(env, (jweak)o);
}
static inline void mw_DeleteLocalRef(JNIEnv *env, jobject o) {
	(*env)->DeleteLocalRef(env, o);
}
static inline jint mw_PushLocalFrame(JNIEnv *env, jint capacity) {
	return (*env)->PushLocalFrame(env, capacity);
}
static inline jobject mw_PopLocalFrame(JNIEnv *env, jobject keep) {
	return (*env)->PopLocalFrame(env, keep);
}

// Members.
static inline jmethodID mw_GetMethodID(JNIEnv *env, jclass c, const char *n, const char *s) {
	return (*env)->GetMethodID(env, c, n, s);
}
static inline jmethodID mw_GetStaticMethodID(JNIEnv *env, jclass c, const char *n, const char *s) {
	return (*env)->GetStaticMethodID(env, c, n, s);
}
static inline jfieldID mw_GetFieldID(JNIEnv *env, jclass c, const char *n, const char *s) {
	return (*env)->GetFieldID(env, c, n, s);
}
static inline jfieldID mw_GetStaticFieldID(JNIEnv *env, jclass c, const char *n, const char *s) {
	return (*env)->GetStaticFieldID(env, c, n, s);
}

// Instance calls.
static inline void mw_CallVoid(JNIEnv *env, jobject o, jmethodID m, jvalue *a) {
	(*env)->CallVoidMethodA(env, o, m, a);
}
static inline jobject mw_CallObject(JNIEnv *env, jobject o, jmethodID m, jvalue *a) {
	return (*env)->CallObjectMethodA(env, o, m, a);
}
static inline jboolean mw_CallBoolean(JNIEnv *env, jobject o, jmethodID m, jvalue *a) {
	return (*env)->CallBooleanMethodA(env, o, m, a);
}
static inline jint mw_CallInt(JNIEnv *env, jobject o, jmethodID m, jvalue *a) {
	return (*env)->CallIntMethodA(env, o, m, a);
}
static inline jlong mw_CallLong(JNIEnv *env, jobject o, jmethodID m, jvalue *a) {
	return (*env)->CallLongMethodA(env, o, m, a);
}
static inline jfloat mw_CallFloat(JNIEnv *env, jobject o, jmethodID m, jvalue *a) {
	return (*env)->CallFloatMethodA(env, o, m, a);
}
static inline jdouble mw_CallDouble(JNIEnv *env, jobject o, jmethodID m, jvalue *a) {
	return (*env)->CallDoubleMethodA(env, o, m, a);
}

// Static calls.
static inline void mw_CallStaticVoid(JNIEnv *env, jclass c, jmethodID m, jvalue *a) {
	(*env)->CallStaticVoidMethodA(env, c, m, a);
}
static inline jobject mw_CallStaticObject(JNIEnv *env, jclass c, jmethodID m, jvalue *a) {
	return (*env)->CallStaticObjectMethodA(env, c, m, a);
}
static inline jboolean mw_CallStaticBoolean(JNIEnv *env, jclass c, jmethodID m, jvalue *a) {
	return (*env)->CallStaticBooleanMethodA(env, c, m, a);
}
static inline jint mw_CallStaticInt(JNIEnv *env, jclass c, jmethodID m, jvalue *a) {
	return (*env)->CallStaticIntMethodA(env, c, m, a);
}
static inline jlong mw_CallStaticLong(JNIEnv *env, jclass c, jmethodID m, jvalue *a) {
	return (*env)->CallStaticLongMethodA(env, c, m, a);
}
static inline jfloat mw_CallStaticFloat(JNIEnv *env, jclass c, jmethodID m, jvalue *a) {
	return (*env)->CallStaticFloatMethodA(env, c, m, a);
}
static inline jdouble mw_CallStaticDouble(JNIEnv *env, jclass c, jmethodID m, jvalue *a) {
	return (*env)->CallStaticDoubleMethodA(env, c, m, a);
}
static inline jobject mw_NewObject(JNIEnv *env, jclass c, jmethodID m, jvalue *a) {
	return (*env)->NewObjectA(env, c, m, a);
}

// Fields.
static inline jobject mw_GetObjectField(JNIEnv *env, jobject o, jfieldID f) {
	return (*env)->GetObjectField(env, o, f);
}
static inline jint mw_GetIntField(JNIEnv *env, jobject o, jfieldID f) {
	return (*env)->GetIntField(env, o, f);
}
static inline jlong mw_GetLongField(JNIEnv *env, jobject o, jfieldID f) {
	return (*env)->GetLongField(env, o, f);
}
static inline jboolean mw_GetBooleanField(JNIEnv *env, jobject o, jfieldID f) {
	return (*env)->GetBooleanField(env, o, f);
}
static inline jfloat mw_GetFloatField(JNIEnv *env, jobject o, jfieldID f) {
	return (*env)->GetFloatField(env, o, f);
}
static inline void mw_SetFloatField(JNIEnv *env, jobject o, jfieldID f, jfloat v) {
	(*env)->SetFloatField(env, o, f, v);
}
static inline void mw_SetObjectField(JNIEnv *env, jobject o, jfieldID f, jobject v) {
	(*env)->SetObjectField(env, o, f, v);
}
static inline void mw_SetIntField(JNIEnv *env, jobject o, jfieldID f, jint v) {
	(*env)->SetIntField(env, o, f, v);
}
static inline jobject mw_GetStaticObjectField(JNIEnv *env, jclass c, jfieldID f) {
	return (*env)->GetStaticObjectField(env, c, f);
}
static inline jint mw_GetStaticIntField(JNIEnv *env, jclass c, jfieldID f) {
	return (*env)->GetStaticIntField(env, c, f);
}
static inline jlong mw_GetStaticLongField(JNIEnv *env, jclass c, jfieldID f) {
	return (*env)->GetStaticLongField(env, c, f);
}

// Strings. The UTF-16 pair is the correct one; the UTF-8 pair speaks Java's
// "modified UTF-8", which is not UTF-8.
static inline jstring mw_NewString(JNIEnv *env, const jchar *u, jsize len) {
	return (*env)->NewString(env, u, len);
}
static inline jsize mw_GetStringLength(JNIEnv *env, jstring s) {
	return (*env)->GetStringLength(env, s);
}
static inline void mw_GetStringRegion(JNIEnv *env, jstring s, jsize start, jsize len, jchar *buf) {
	(*env)->GetStringRegion(env, s, start, len, buf);
}
static inline jstring mw_NewStringUTF(JNIEnv *env, const char *s) {
	return (*env)->NewStringUTF(env, s);
}

// Arrays.
static inline jsize mw_GetArrayLength(JNIEnv *env, jarray a) {
	return (*env)->GetArrayLength(env, a);
}
static inline jintArray mw_NewIntArray(JNIEnv *env, jsize n) {
	return (*env)->NewIntArray(env, n);
}
static inline jbyteArray mw_NewByteArray(JNIEnv *env, jsize n) {
	return (*env)->NewByteArray(env, n);
}
static inline jobjectArray mw_NewObjectArray(JNIEnv *env, jsize n, jclass c, jobject init) {
	return (*env)->NewObjectArray(env, n, c, init);
}
static inline void mw_GetIntArrayRegion(JNIEnv *env, jintArray a, jsize s, jsize n, jint *b) {
	(*env)->GetIntArrayRegion(env, a, s, n, b);
}
static inline void mw_SetIntArrayRegion(JNIEnv *env, jintArray a, jsize s, jsize n, const jint *b) {
	(*env)->SetIntArrayRegion(env, a, s, n, b);
}
static inline void mw_GetByteArrayRegion(JNIEnv *env, jbyteArray a, jsize s, jsize n, jbyte *b) {
	(*env)->GetByteArrayRegion(env, a, s, n, b);
}
static inline void mw_SetByteArrayRegion(JNIEnv *env, jbyteArray a, jsize s, jsize n, const jbyte *b) {
	(*env)->SetByteArrayRegion(env, a, s, n, b);
}
static inline jobject mw_GetObjectArrayElement(JNIEnv *env, jobjectArray a, jsize i) {
	return (*env)->GetObjectArrayElement(env, a, i);
}
static inline void mw_SetObjectArrayElement(JNIEnv *env, jobjectArray a, jsize i, jobject v) {
	(*env)->SetObjectArrayElement(env, a, i, v);
}

// A direct view of an array's storage. Between these two calls the runtime
// may have suspended the collector, so nothing else may be called and
// nothing may block.
static inline void *mw_GetCritical(JNIEnv *env, jarray a) {
	return (*env)->GetPrimitiveArrayCritical(env, a, NULL);
}
static inline void mw_ReleaseCritical(JNIEnv *env, jarray a, void *p, jint mode) {
	(*env)->ReleasePrimitiveArrayCritical(env, a, p, mode);
}

// Exceptions.
static inline jboolean mw_ExceptionCheck(JNIEnv *env) {
	return (*env)->ExceptionCheck(env);
}
static inline jthrowable mw_ExceptionOccurred(JNIEnv *env) {
	return (*env)->ExceptionOccurred(env);
}
static inline void mw_ExceptionClear(JNIEnv *env) {
	(*env)->ExceptionClear(env);
}
static inline void mw_ExceptionDescribe(JNIEnv *env) {
	(*env)->ExceptionDescribe(env);
}
static inline jint mw_ThrowNew(JNIEnv *env, jclass c, const char *msg) {
	return (*env)->ThrowNew(env, c, msg);
}

// Native methods, for the rare case the name-mangling convention will not do.
static inline jint mw_RegisterNatives(JNIEnv *env, jclass c, const JNINativeMethod *m, jint n) {
	return (*env)->RegisterNatives(env, c, m, n);
}
static inline jint mw_UnregisterNatives(JNIEnv *env, jclass c) {
	return (*env)->UnregisterNatives(env, c);
}

#endif
