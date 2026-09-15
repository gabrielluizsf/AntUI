# What R8 must not touch in the AntUI shim.
#
# Nothing here is called from Java. Every class, every method and every field
# is reached from native code — by name, at run time, through FindClass and
# GetMethodID — and a shrinker cannot see any of that. To R8 this whole
# package is dead code with a few entry points the manifest happens to name.
#
# The failure is not a build error. R8 removes or renames something, the
# build succeeds, the app installs, and the first time the code path runs it
# throws NoSuchMethodError from inside a JNI call with no indication of why.
#
# AntUI's own pipeline never runs R8 — antuiapk compiles the shim with javac
# and dexes it with d8, and neither shrinks anything — so these rules are for
# an app that merges the shim into a Gradle build, where R8 is on by default
# for every release build.

# The shim itself, whole. Not just the classes: the methods and fields too,
# because it is the members that native code looks up.
-keep class dev.antui.** { *; }

# Any native method anywhere, under the name it was declared with. These are
# attached with RegisterNatives from C, which matches on the name and the
# signature, so a renamed method is an unattached one.
-keepclasseswithmembernames,includedescriptorclasses class * {
    native <methods>;
}

# The activity is named in the manifest, which R8 normally reads and honours
# on its own. It is listed anyway because a manifest is a build input and
# this file is not: an app that merges the shim without merging the manifest
# would otherwise lose it.
-keep public class dev.antui.AntuiActivity { *; }

# The file provider is named in the manifest too, and is instantiated by the
# platform rather than by the app.
-keep public class dev.antui.AntuiFileProvider { *; }
