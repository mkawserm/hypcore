package views

import (
	"github.com/golang/glog"
	core2 "github.com/mkawserm/hypcore/package/core"
	"net/http"
	"regexp"
)

type DynamicView struct {
	Context *core2.HContext
}

func (dView *DynamicView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("PATH: " + r.URL.Path)

	// check for auth
	if dView.Context.HasAuth() {
		uid := ""
		ok := false

		h := httpGetHeader(r.Header, core2.HeaderAuthorizationCanonical)
		if h == "" {
			httpBadRequest(w, []byte("Oops! No Authorization header found !!!"))
			return

		} else {
			uid, ok = dView.Context.Auth.GetUID([]byte(h), dView.Context.AuthBearer)
		}

		if ok {
			if uid == "" {
				httpBadRequest(w, []byte("Oops! No UID found from AuthInterface !!!"))
				return
			}

		} else { // Failed to authorize. not ok
			httpBadRequest(w, []byte("Oops! Invalid Authorization data !!!"))
			return
		}
	}

	glog.Infoln("Processing Middleware in the DynamicView.")
	for _, mi := range dView.Context.MiddlewareList {
		next := mi.ServeHTTP(dView.Context, r, w)
		if next == false {
			return
		}
	}
	glog.Infoln("Middleware processing complete.")

	glog.Infoln("Dynamic route dispatch started.")
	for _, route := range dView.Context.RouteList {
		var rc = regexp.MustCompile(route.Pattern)
		if rc.Match([]byte(r.URL.Path)) {
			route.HttpHandlerObject.ServeHTTP(dView.Context, r, w)
			return
		}
	}

	glog.Infoln("Dynamic route dispatch failed to find a route.")
	glog.Errorln("Showing 404 Http Error")

	httpNotFound(w, []byte("Oops!!!"))
}
