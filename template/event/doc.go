// Package event is what the templates agree on: the components that exist,
// the things they can do, and the sound bank that hears about them. It has no
// imports of its own, so the visual and the audible halves of a template can
// both speak it without leaning on each other.
//
// The flow is one-way. A component draws itself, and when the user does
// something to it the component reports a single Event: which component it
// was and what was done. The template's sound bank — the second half of the
// template — listens for those events and plays the matching sound. Anything
// can be a sound bank, so a template can be made to say whatever its author
// wants it to say without touching the way it looks.
package event
