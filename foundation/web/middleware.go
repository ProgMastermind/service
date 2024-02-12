package web

// Middleware is a function designed to run some code before and/or after
// another Handler. It is desinged to remove boilerplate or other concerns not
// direct to any given Handler

type Middleware func(Handler) Handler

// wrapMiddleware creates a new Handler by wrapping middleware around final handler.
// The middleware's Handlers will be executed by requests in the order they
// are provided

func wrapMiddleware(mw []Middleware, handler Handler) Handler {

	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backward ensures that the
	// first middleware of the slice is the first to be executed by requests

	for i := len(mw) - 1; i >= 0; i-- {
		mvFunc := mw[i]
		if mvFunc != nil {
			handler = mvFunc(handler)
		}

	}

	return handler
}
