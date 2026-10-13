//go:build android

package app

import "sync"

// Listen registers a handler and gives back the token that identifies it, and
// the function that takes it off again.
//
// **The release function has to be called.** A handler that is never released
// is a Go closure that lives for the life of the process, and — worse — a
// Java object still registered with the platform, which will go on calling it.
func Listen(h Handler) (token int64, release func()) {
	token = lastToken.Add(1)
	handlerMu.Lock()
	handlers[token] = h
	handlerMu.Unlock()
	var once sync.Once
	return token, func() {
		once.Do(func() {
			handlerMu.Lock()
			delete(handlers, token)
			handlerMu.Unlock()
		})
	}
}
