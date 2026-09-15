// Package notify posts notifications.
//
// Two things about Android's notifications catch everyone once. A
// notification must belong to a channel that has been created, or nothing
// appears and nothing fails; and from API 33 it needs the
// POST_NOTIFICATIONS permission, and without it, again, nothing appears and
// nothing fails. [CreateChannel] and [Allowed] are how to be sure.
package notify
