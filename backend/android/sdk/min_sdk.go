package sdk

// MinSDK is the oldest Android this library supports: API 21, Android 5.0.
// It is not an arbitrary floor — it is the oldest the current NDK will build
// for, and below it there is no compiler to use.
const MinSDK = 21
