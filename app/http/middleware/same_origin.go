package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"net/url"
	"strings"
)

// ValidMutationOrigin uses browser-controlled metadata to protect session mutations.
// Requests without browser metadata and cookies may be used by non-browser clients.
func ValidMutationOrigin(method, host, scheme, origin, referer, fetchSite, cookie string) bool {
	if method == "GET" || method == "HEAD" || method == "OPTIONS" {
		return true
	}
	if fetchSite != "" && fetchSite != "same-origin" && fetchSite != "none" {
		return false
	}
	same := func(raw string) bool {
		u, err := url.Parse(raw)
		return err == nil && u.Scheme == scheme && strings.EqualFold(u.Host, host) && u.User == nil
	}
	if origin != "" {
		return same(origin)
	}
	if referer != "" {
		return same(referer)
	}
	if fetchSite == "same-origin" {
		return true
	}
	return cookie == "" && fetchSite == ""
}

type SameOriginMiddleware struct{}

func (m *SameOriginMiddleware) Signature() string { return "same_origin" }
func (m *SameOriginMiddleware) Handle(ctx http.Context) {
	req := ctx.Request()
	scheme := "http"
	if req.Origin().TLS != nil {
		scheme = "https"
	}
	if !ValidMutationOrigin(req.Method(), req.Host(), scheme, req.Header("Origin"), req.Header("Referer"), req.Header("Sec-Fetch-Site"), req.Header("Cookie")) {
		_ = ctx.Response().String(http.StatusForbidden, "This request must come from the same site.").Abort()
		return
	}
	req.Next()
}
func SameOrigin() http.Middleware { return &SameOriginMiddleware{} }
