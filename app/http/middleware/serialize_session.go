package middleware

import (
	"context"
	"elcomware/app/facades"
	"hash/maphash"
	"time"

	"github.com/goravel/framework/contracts/http"
)

// Fixed stripes bound memory even when clients send arbitrary session cookies.
// The process-random hash also avoids predictable cross-session contention.
var sessionStripes = func() [1024]chan struct{} {
	var stripes [1024]chan struct{}
	for i := range stripes {
		stripes[i] = make(chan struct{}, 1)
	}
	return stripes
}()
var sessionStripeSeed = maphash.MakeSeed()

type SerializeSessionMiddleware struct{}

func (*SerializeSessionMiddleware) Signature() string { return "serialize_session" }
func (*SerializeSessionMiddleware) Handle(ctx http.Context) {
	id := ctx.Request().Cookie(facades.Config().GetString("session.cookie"))
	if id == "" {
		// Without a cookie, StartSession creates an independent new session.
		ctx.Request().Next()
		return
	}
	wait, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()
	release, err := acquireSession(wait, id)
	if err != nil {
		ctx.Response().Header("Retry-After", "1")
		_ = ctx.Response().String(http.StatusServiceUnavailable, "Session is busy. Please try again.").Abort()
		return
	}
	defer release()
	// Must wrap StartSession, including its deferred session persistence after Next.
	ctx.Request().Next()
}

// SerializeSession protects session read-modify-write within one server process.
// Multi-process deployments require distributed locking or transactional carts.
func SerializeSession() http.Middleware { return &SerializeSessionMiddleware{} }

// acquireSession waits only until the request is cancelled or its wait budget expires.
func acquireSession(ctx context.Context, id string) (func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	lock := sessionStripes[maphash.String(sessionStripeSeed, id)%uint64(len(sessionStripes))]
	select {
	case lock <- struct{}{}:
		if err := ctx.Err(); err != nil {
			<-lock
			return nil, err
		}
		return func() { <-lock }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
