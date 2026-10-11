package template

// loadingEnabledByDefault says whether a template starts with its loading
// screen on. It is on for the program the developer ships — a screen the
// app has to ask to be without is the screen that catches the app which
// forgot to say Ready — and the test suite turns it off in an init, because
// the suite draws into canvases it then looks at and a screen painted over
// everything would be in the way of every one of those pictures. The tests
// that are about the loading screen turn it back on around their own cases.
var loadingEnabledByDefault = true
